package comments_v1

import (
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
	"strings"
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
	rows, err := db.Conn.Query(`SELECT id, "userId", content, "createdAt", "updatedAt" FROM comments WHERE "movieId" = $1`, movieID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	found := false
	for rows.Next() {
		var id, userId int32
		var content string
		var createdAt, updatedAt time.Time
		errScan := rows.Scan(&id, &userId, &content, &createdAt, &updatedAt)
		if errScan != nil {
			return nil, errScan
		}

		moviesItem := models.Movies{
			Id:          id,
			Title:       title,
			Description: description,
			ReleaseDate: releaseDate,
			ScoreKP:     scoreKP,
			ScoreIMDB:   scoreIMDB,
			Poster:      poster,
			TypeMovie:   typeMovie,
			Genres:      strings.Join(genresArray, ", "),
		}

		movies = append(movies, moviesItem)

		found = true
	}

	if !found {
		return nil, status.Error(codes.NotFound, "нет значений в БД")
	}

	// Save data to Redis
	moviesJSONbyte, err := json.Marshal(movies)
	if err != nil {
		return nil, err
	}
	err = cache.Rdb.Set(cache.Ctx, "comments_"+fmt.Sprint(movieId)+"_"+fmt.Sprint(limit)+"_"+fmt.Sprint(page), moviesJSONbyte, 1*time.Minute).Err()
	if err != nil {
		log.Printf("Ошибка при сохранении данных в Redis: %v", err)
	}

	return movies, nil
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
