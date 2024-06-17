package movies_v1

import (
	"encoding/json"
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

func (s Service) GetMoviesService() ([]models.Movies, error) {
	var movies []models.Movies

	// Получение данных из Redis
	moviesJSON, err := cache.Rdb.Get(cache.Ctx, "movies").Result()
	if err != nil {
		return nil, err
	}
	if moviesJSON != "" {
		err = json.Unmarshal([]byte(moviesJSON), &movies)
		if err != nil {
			log.Fatalf("Ошибка десериализации: %v", err)
		}
		return movies, nil
	}

	// Data not in Redis, get from database
	query := "SELECT movies.*, array_agg(genres.name) AS genres FROM movies JOIN moviesgenres ON movies.id = moviesgenres.idmovie JOIN genres ON moviesgenres.idgenre = genres.id GROUP BY movies.id ORDER BY movies.id DESC"
	rows, err := db.Conn.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	found := false
	for rows.Next() {
		var id, releaseDate, typeMovie int32
		var title, description, poster, genres string
		var scoreKP, scoreIMDB float64
		err := rows.Scan(&id, &title, &description, &releaseDate, &scoreKP, &scoreIMDB, &poster, &typeMovie, &genres)
		if err != nil {
			return nil, err
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
			Genres:      genres,
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
	err = cache.Rdb.Set(cache.Ctx, "movies", moviesJSONbyte, 1*time.Minute).Err()
	if err != nil {
		return nil, err
	}

	return movies, nil
}

func (s Service) GetMoviesByIdService(i int32) (models.Movies, error) {
	//TODO implement me
	panic("implement me")
}

func (s Service) GetMoviesByCategoryService(s2 string) ([]models.Movies, error) {
	//TODO implement me
	panic("implement me")
}

func (s Service) AddMoviesService(s5 string, s2 string, s3 string, s4 string) (int32, error) {
	//TODO implement me
	panic("implement me")
}

func (s Service) DelMoviesService(i int32) error {
	//TODO implement me
	panic("implement me")
}

func New() *Service {
	return &Service{}
}
