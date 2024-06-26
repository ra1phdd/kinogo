package comments_v1

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"kinogo/internal/app/models"
	"kinogo/pkg/cache"
	"kinogo/pkg/db"
	"log"
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

	commentsJSON, err := cache.Rdb.Get(cache.Ctx, "comments_"+fmt.Sprint(movieId)+"_"+fmt.Sprint(limit)+"_"+fmt.Sprint(page)).Result()
	if err == nil && commentsJSON != "" {
		err = json.Unmarshal([]byte(commentsJSON), &comments)
		if err != nil {
			log.Printf("Ошибка десериализации: %v", err)
		}
		return comments, nil
	} else if !errors.Is(err, redis.Nil) {
		// Если ошибка не связана с отсутствием ключа, логируем её
		log.Printf("Ошибка при получении данных из Redis: %v", err)
	}

	var limitQuery, pageQuery string
	if limit > 0 {
		limitQuery = fmt.Sprintf("LIMIT %d", limit)
	}
	if page > 0 && limit > 0 {
		pageQuery = fmt.Sprintf("OFFSET %d", (page-1)*limit)
	}

	query := fmt.Sprintf(`
		  SELECT c.id, c."parentId", c.text, c."createdAt", c."updatedAt", u.username, u.photourl, u.first_name, u.last_name
		  FROM comments c
		  JOIN users u ON c."userId" = u.id
		  WHERE c."movieId" = $1 ORDER BY c.id DESC %s %s`, limitQuery, pageQuery)

	rows, err := db.Conn.Query(query, movieId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	allComments := make(map[int32]models.Comments)
	var rootComments []models.Comments

	for rows.Next() {
		var id int32
		var parentId sql.NullInt32
		var text string
		var createdAt, updatedAt time.Time
		var username, photoUrl, firstName, lastName string
		errScan := rows.Scan(&id, &parentId, &text, &createdAt, &updatedAt, &username, &photoUrl, &firstName, &lastName)
		if errScan != nil {
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
		return nil, status.Error(codes.NotFound, "нет значений в БД")
	}

	// Построение дерева комментариев
	comments = buildCommentTree(rootComments, allComments)

	commentsJSONbyte, err := json.Marshal(comments)
	if err != nil {
		return nil, err
	}
	err = cache.Rdb.Set(cache.Ctx, "comments_"+fmt.Sprint(movieId)+"_"+fmt.Sprint(limit)+"_"+fmt.Sprint(page), commentsJSONbyte, 1*time.Minute).Err()
	if err != nil {
		log.Printf("Ошибка при сохранении данных в Redis: %v", err)
	}

	return comments, nil
}

func (s Service) AddCommentService(data map[string]interface{}) (int32, error) {
	var parentId *int32 = nil
	if data["parentId"].(int32) != 0 {
		id := data["parentId"].(int32)
		parentId = &id
	}

	var id int32
	err := db.Conn.QueryRow(
		`INSERT INTO comments ("userId", "movieId", "parentId", text, "createdAt", "updatedAt")
         VALUES ($1, $2, $3, $4, $5, $5)
         RETURNING id`, data["userId"], data["movieId"], parentId, data["text"], data["createdAt"]).Scan(&id)
	if err != nil {
		return 0, err
	}

	pattern := "comments_" + fmt.Sprint(data["movieId"].(int32)) + "_*_*"

	// Получение всех ключей, соответствующих шаблону
	keys, err := cache.Rdb.Keys(context.Background(), pattern).Result()
	if err != nil {
		return 0, fmt.Errorf("failed to get keys: %w", err)
	}

	// Удаление всех ключей
	if len(keys) > 0 {
		if err := cache.Rdb.Del(context.Background(), keys...).Err(); err != nil {
			return 0, fmt.Errorf("failed to delete keys: %w", err)
		}
	}

	return id, nil
}

func (s Service) UpdateCommentService(data map[string]interface{}) error {
	_, err := db.Conn.Exec(
		`UPDATE comments
         SET text = $1, "updatedAt" = $2
         WHERE id = $3`, data["text"], data["updatedAt"], data["id"])
	if err != nil {
		return fmt.Errorf("failed to update comment: %w", err)
	}

	var movieId int32
	err = db.Conn.QueryRow(
		`SELECT "movieId" FROM comments WHERE id = $1`, data["id"]).Scan(&movieId)
	if err != nil {
		return fmt.Errorf("failed to get movieId for comment: %w", err)
	}

	pattern := "comments_" + fmt.Sprint(movieId) + "_*_*"

	// Получение всех ключей, соответствующих шаблону
	keys, err := cache.Rdb.Keys(context.Background(), pattern).Result()
	if err != nil {
		return fmt.Errorf("failed to get keys: %w", err)
	}

	// Удаление всех ключей
	if len(keys) > 0 {
		if err := cache.Rdb.Del(context.Background(), keys...).Err(); err != nil {
			return fmt.Errorf("failed to delete keys: %w", err)
		}
	}

	return nil
}

func (s Service) DelCommentService(id int32) error {
	rows, err := db.Conn.Query(
		`SELECT id, "movieId" FROM comments WHERE "parentId" >= $1 ORDER BY id DESC`, id)
	if err != nil {
		return fmt.Errorf("failed to fetch children IDs: %w", err)
	}
	defer rows.Close()

	var movieId int32
	var childIds []int32
	for rows.Next() {
		var childId int32
		if errScan := rows.Scan(&childId, &movieId); errScan != nil {
			return fmt.Errorf("failed to scan children IDs: %w", errScan)
		}
		childIds = append(childIds, childId)
	}

	allIds := append(childIds, id)

	fmt.Println(allIds)

	// Delete all comments in a single transaction for consistency
	tx, err := db.Conn.BeginTx(cache.Ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}

	// Delete comments
	for _, commentId := range allIds {
		_, errExec := tx.Exec(`DELETE FROM comments WHERE id = $1`, commentId)
		if errExec != nil {
			tx.Rollback()
			return fmt.Errorf("failed to delete comment %d: %w", commentId, errExec)
		}
	}

	// Commit the transaction
	if errCommit := tx.Commit(); errCommit != nil {
		return fmt.Errorf("failed to commit transaction: %w", errCommit)
	}

	pattern := "comments_" + fmt.Sprint(movieId) + "_*_*"

	// Получение всех ключей, соответствующих шаблону
	keys, err := cache.Rdb.Keys(context.Background(), pattern).Result()
	if err != nil {
		return fmt.Errorf("failed to get keys: %w", err)
	}

	// Удаление всех ключей
	if len(keys) > 0 {
		if err := cache.Rdb.Del(context.Background(), keys...).Err(); err != nil {
			return fmt.Errorf("failed to delete keys: %w", err)
		}
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
