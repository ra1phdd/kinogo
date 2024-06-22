package comments_v1

import (
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
	"time"
)

type Service struct {
}

func New() *Service {
	return &Service{}
}

func (s Service) GetCommentsByIdService(movieId int32, limit int32, page int32) ([]models.Comments, error) {
	var comments []models.Comments

	// Получение данных из Redis
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

	// Data not in Redis, get from database
	var limitQuery, pageQuery string
	if limit > 0 {
		limitQuery = fmt.Sprintf("LIMIT %d", limit)
	}
	if page > 0 && limit > 0 {
		pageQuery = fmt.Sprintf("OFFSET %d", (page-1)*limit)
	}

	query := fmt.Sprintf(`
		  SELECT id, "userId", "parentId", text, "createdAt", "updatedAt" FROM comments WHERE "movieId" = $1 %s %s`, limitQuery, pageQuery)

	rows, err := db.Conn.Query(query, movieId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	allComments := make(map[int32]models.Comments)
	var rootComments []models.Comments

	for rows.Next() {
		var id, userId int32
		var parentId sql.NullInt32
		var text string
		var createdAt, updatedAt time.Time
		errScan := rows.Scan(&id, &userId, &parentId, &text, &createdAt, &updatedAt)
		if errScan != nil {
			return nil, errScan
		}

		comment := models.Comments{
			ID:        id,
			UserID:    userId,
			ParentID:  parentId.Int32,
			Text:      text,
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
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

	// Применение пагинации
	if limit > 0 && page > 0 {
		start := (page - 1) * limit
		end := start + limit
		if start >= int32(len(comments)) {
			comments = []models.Comments{}
		} else if end > int32(len(comments)) {
			comments = comments[start:]
		} else {
			comments = comments[start:end]
		}
	}

	// Save data to Redis
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
	//TODO implement me
	panic("implement me")
}

func (s Service) UpdateCommentService(data map[string]interface{}) error {
	//TODO implement me
	panic("implement me")
}

func (s Service) DelCommentService(id int32, parentId int32) error {
	//TODO implement me
	panic("implement me")
}

func buildCommentTree(comments []models.Comments, allComments map[int32]models.Comments) []models.Comments {
	for i, comment := range comments {
		children := []models.Comments{}
		for _, potentialChild := range allComments {
			if potentialChild.ParentID == comment.ID {
				children = append(children, potentialChild)
			}
		}
		if len(children) > 0 {
			comments[i].Children = buildCommentTree(children, allComments)
		}
	}
	return comments
}
