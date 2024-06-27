package comments_v1

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"kinogo/internal/app/models"
	"kinogo/pkg/cache"
	"kinogo/pkg/db"
	"kinogo/pkg/logger"
	"sort"
	"time"
)

type Service struct {
}

func New() *Service {
	return &Service{}
}

func (s Service) GetCommentsByIdService(movieId int32, limit int32, page int32) ([]models.Comments, error) {
	var comments []models.Comments

	commentsJSON, err := cache.Rdb.Get(cache.Ctx, fmt.Sprintf("comments_%d_%d_%d", movieId, limit, page)).Result()
	if err == nil && commentsJSON != "" {
		if err = json.Unmarshal([]byte(commentsJSON), &comments); err != nil {
			logger.Error("Ошибка десериализации данных из Redis", zap.String("commentsJSON", commentsJSON))
			return nil, err
		}
		return comments, nil
	}
	if err != nil && !errors.Is(err, redis.Nil) {
		logger.Error("Ошибка при получении данных из Redis")
	}

	query := `
		SELECT c.id, c."parentId", c.text, c."createdAt", c."updatedAt", u.username, u.photourl, u.first_name, u.last_name
		FROM comments c
		JOIN users u ON c."userId" = u.id
		WHERE c."movieId" = $1
		ORDER BY c.id DESC`

	params := []interface{}{movieId}
	if limit > 0 {
		query += " LIMIT $2"
		params = append(params, limit)
		if page > 0 {
			query += " OFFSET $3"
			params = append(params, (page-1)*limit)
		}
	}

	rows, err := db.Conn.Query(query, params...)
	if err != nil {
		logger.Error("Ошибка выполнения SQL-запроса", zap.String("query", query), zap.Any("params", params))
		return nil, err
	}
	defer func() {
		if errClose := rows.Close(); errClose != nil {
			logger.Error("Ошибка при закрытии rows", zap.Error(errClose))
		}
	}()

	var allComments = make(map[int32]models.Comments)
	var rootComments []models.Comments

	for rows.Next() {
		var id int32
		var parentId sql.NullInt32
		var text string
		var createdAt, updatedAt time.Time
		var username, photoUrl, firstName, lastName string
		errScan := rows.Scan(&id, &parentId, &text, &createdAt, &updatedAt, &username, &photoUrl, &firstName, &lastName)
		if errScan != nil {
			logger.Error("Ошибка сканирования строки результата запроса")
			return nil, errScan
		}

		comment := models.Comments{
			ID:        id,
			ParentID:  parentId.Int32,
			Text:      text,
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
			User: models.User{
				Username:  username,
				PhotoUrl:  photoUrl,
				FirstName: firstName,
				LastName:  lastName,
			},
		}

		allComments[id] = comment

		if !parentId.Valid || parentId.Int32 == 0 {
			rootComments = append(rootComments, comment)
		}
	}

	if len(allComments) == 0 {
		logger.Warn("Нет значений в БД", zap.Any("comments", allComments))
		return nil, status.Error(codes.NotFound, "Нет значений в БД")
	}

	comments = buildCommentTree(rootComments, allComments)

	commentsJSONbyte, err := json.Marshal(comments)
	if err != nil {
		logger.Error("Ошибка сериализации данных в Redis", zap.Any("comments", comments))
		return nil, err
	}
	err = cache.Rdb.Set(cache.Ctx, fmt.Sprintf("comments_%d_%d_%d", movieId, limit, page), commentsJSONbyte, 5*time.Minute).Err()
	if err != nil {
		logger.Error("Ошибка при сохранении данных в Redis")
	}

	return comments, nil
}

func (s Service) AddCommentService(data map[string]interface{}) (int32, error) {
	// Так как parentId может быть как числом, так и nil
	// то обьявляем его как указатель на int32 и присваиваем nil
	var parentId *int32 = nil
	if data["parentId"].(int32) != 0 {
		id := data["parentId"].(int32)
		parentId = &id
	}

	var id int32
	err := db.Conn.QueryRow(
		`INSERT INTO comments ("userId", "movieId", "parentId", text, "createdAt", "updatedAt")
         VALUES ($1, $2, $3, $4, $5, $5) RETURNING id`, data["userId"], data["movieId"], parentId, data["text"], data["createdAt"]).Scan(&id)
	if err != nil {
		logger.Error("Ошибка выполнения SQL-запроса", zap.Any("comment", data))
		return 0, err
	}

	pattern := fmt.Sprintf("comments_%d_*_*", data["movieId"].(int32))
	err = cache.ClearCacheByPattern(pattern)
	if err != nil {
		logger.Error("Ошибка очистки кеша по паттерну", zap.String("pattern", pattern))
		return 0, err
	}

	return id, nil
}

func (s Service) UpdateCommentService(data map[string]interface{}) error {
	_, err := db.Conn.Exec(`UPDATE comments SET text = $1, "updatedAt" = $2 WHERE id = $3`, data["text"], data["updatedAt"], data["id"])
	if err != nil {
		logger.Error("Ошибка выполнения SQL-запроса", zap.Any("comment", data))
		return err
	}

	var movieId int32
	err = db.Conn.QueryRow(`SELECT "movieId" FROM comments WHERE id = $1`, data["id"]).Scan(&movieId)
	if err != nil {
		logger.Error("Ошибка выполнения SQL-запроса (получение movieId по id комментария)", zap.Any("id", data["id"]))
		return err
	}

	pattern := fmt.Sprintf("comments_%d_*_*", data["movieId"].(int32))
	err = cache.ClearCacheByPattern(pattern)
	if err != nil {
		logger.Error("Ошибка очистки кеша по паттерну", zap.String("pattern", pattern))
		return err
	}

	return nil
}

func (s Service) DelCommentService(id int32) error {
	rows, err := db.Conn.Query(`
		WITH RECURSIVE CommentTree AS (
			SELECT id, "movieId"
			FROM comments
			WHERE id = $1
			UNION ALL
			SELECT c.id, c."movieId"
			FROM comments c
			INNER JOIN CommentTree ct ON ct.id = c."parentId"
		)
		SELECT id, "movieId" FROM CommentTree ORDER BY id DESC`, id)
	if err != nil {
		logger.Error("Ошибка выполнения SQL-запроса", zap.Int32("id", id))
		return err
	}
	defer func() {
		if errClose := rows.Close(); errClose != nil {
			logger.Error("Ошибка при закрытии rows", zap.Error(errClose))
		}
	}()

	var movieId, commentId int32
	var commentIds []int32
	for rows.Next() {
		if errScan := rows.Scan(&commentId, &movieId); errScan != nil {
			logger.Error("Ошибка сканирования строки результата запроса")
			return errScan
		}

		commentIds = append(commentIds, commentId)
	}

	tx, err := db.Conn.BeginTx(cache.Ctx, nil)
	if err != nil {
		logger.Error("Ошибка начала транзации в БД")
		return err
	}

	for _, item := range commentIds {
		_, errExec := tx.Exec(`DELETE FROM comments WHERE id = $1`, item)
		if errExec != nil {
			errRollback := tx.Rollback()
			if errRollback != nil {
				logger.Error("Ошибка отката транзации в БД")
				return errRollback
			}
			logger.Error("Ошибка выполнения SQL-запроса в транзакции", zap.Int32("id", item))
			return errExec
		}
	}

	if errCommit := tx.Commit(); errCommit != nil {
		logger.Error("Ошибка окончания транзации в БД")
		return errCommit
	}

	pattern := fmt.Sprintf("comments_%d_*_*", movieId)
	err = cache.ClearCacheByPattern(pattern)
	if err != nil {
		logger.Error("Ошибка очистки кеша по паттерну", zap.String("pattern", pattern))
		return err
	}

	return nil
}

func buildCommentTree(comments []models.Comments, allComments map[int32]models.Comments) []models.Comments {
	for i, comment := range comments {
		var children []models.Comments
		for _, potentialChild := range allComments {
			if potentialChild.ParentID == comment.ID {
				children = append(children, potentialChild)
			}
		}
		if len(children) > 0 {
			sort.SliceStable(children, func(i, j int) bool {
				return children[i].ID < children[j].ID
			})
			comments[i].Children = buildCommentTree(children, allComments)
		}
	}
	return comments
}
