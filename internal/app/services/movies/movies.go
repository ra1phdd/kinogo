package movies_v1

import (
	"encoding/json"
	"fmt"
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

func (s Service) GetMovieByIdService(id int32) (models.Movie, error) {
	var movie models.Movie

	// Получение данных из Redis
	movieJSON, err := cache.Rdb.Get(cache.Ctx, "movies_id_"+fmt.Sprint(id)).Result()
	if err != nil {
		return models.Movie{}, err
	}
	if movieJSON != "" {
		err = json.Unmarshal([]byte(movieJSON), &movie)
		if err != nil {
			log.Fatalf("Ошибка десериализации: %v", err)
		}
		return movie, nil
	}

	// Data not in Redis, get from database
	query := `
        SELECT m.*, array_agg(g.name)
        FROM movies m
        JOIN moviesgenres mg ON m.id = mg.idmovie
        JOIN genres g ON mg.idgenre = g.id
        WHERE m.id = $1 GROUP BY m.id`
	rows, err := db.Conn.Query(query, id)
	if err != nil {
		return models.Movie{}, err
	}
	defer rows.Close()

	found := false
	for rows.Next() {
		var id, releaseDate, timeMovie, views, likes, dislikes, typeMovie int32
		var title, description, country, poster string
		var scoreKP, scoreIMDB float64
		var genresArray []string

		err := rows.Scan(&id, &title, &description, &country, &releaseDate, &timeMovie, &scoreKP, &scoreIMDB, &poster, &typeMovie, &views, &likes, &dislikes)
		if err != nil {
			return models.Movie{}, err
		}

		movie = models.Movie{
			Id:          id,
			Title:       title,
			Description: description,
			Country:     country,
			ReleaseDate: releaseDate,
			TimeMovie:   timeMovie,
			ScoreKP:     scoreKP,
			ScoreIMDB:   scoreIMDB,
			Poster:      poster,
			TypeMovie:   typeMovie,
			Views:       views,
			Likes:       likes,
			Dislikes:    dislikes,
			Genres:      strings.Join(genresArray, ", "),
		}

		found = true
	}

	if !found {
		return models.Movie{}, status.Error(codes.NotFound, "нет значений в БД")
	}

	// Save data to Redis
	movieJSONbyte, err := json.Marshal(movie)
	if err != nil {
		return models.Movie{}, err
	}
	err = cache.Rdb.Set(cache.Ctx, "movies_id_"+fmt.Sprint(id), movieJSONbyte, 60*time.Minute).Err()
	if err != nil {
		return models.Movie{}, err
	}

	return movie, nil
}

func (s Service) GetMoviesByFilterService(filtersMap map[string]interface{}) ([]models.Movies, error) {
	var movies []models.Movies

	stringMap := fmt.Sprint(filtersMap)

	// Получение данных из Redis
	movieJSON, err := cache.Rdb.Get(cache.Ctx, "movies_filter_"+stringMap).Result()
	if err != nil {
		return []models.Movies{}, err
	}
	if movieJSON != "" {
		err = json.Unmarshal([]byte(movieJSON), &movies)
		if err != nil {
			log.Fatalf("Ошибка десериализации: %v", err)
		}
		return movies, nil
	}

	// Data not in Redis, get from database
	var baseQuery string
	var args []interface{}
	argPos := 1

	// Условия для года выпуска
	if yearMin, yearMinOk := filtersMap["yearMin"].(int); yearMinOk {
		if yearMax, yearMaxOk := filtersMap["yearMax"].(int); yearMaxOk {
			baseQuery += fmt.Sprintf(" AND m.releasedate BETWEEN $%d AND $%d", argPos, argPos+1)
			args = append(args, yearMin, yearMax)
			argPos += 2
		}
	}

	genres, genresOk := filtersMap["genres"].([]int)
	if !genresOk {
		return []models.Movies{}, fmt.Errorf("error: value is not a []int")
	}

	// Условия для жанров
	if len(genres) > 0 {
		genrePlaceholders := make([]string, len(genres))
		for i, genre := range genres {
			genrePlaceholders[i] = fmt.Sprintf("$%d", argPos)
			args = append(args, genre)
			argPos++
		}
		baseQuery += fmt.Sprintf(" AND g.id IN (%s)", strings.Join(genrePlaceholders, ","))
	}

	// Условия для типа фильма
	if typeMovie, typeMovieOk := filtersMap["typeMovie"].(string); typeMovieOk && typeMovie != "" {
		baseQuery += fmt.Sprintf(" AND m.typemovie = $%d", argPos)
		args = append(args, typeMovie)
		argPos++
	}

	if search, searchOk := filtersMap["search"].(string); searchOk && search != "" {
		baseQuery += " AND word_similarity(m.title, $1) > 0.1"
		args = append(args, search)
		argPos++
	}

	// Условия для лучшего фильма
	var bestMovieQuery string
	if bestMovie, bestMovieOk := filtersMap["bestMovie"].(bool); bestMovieOk && bestMovie {
		bestMovieQuery = " ORDER BY m.views"
	}

	var offset, page, limit int32
	if pageIntf, ok := filtersMap["page"]; ok {
		if page, ok = pageIntf.(int32); !ok {
			page = 0
		}
	}

	// Извлекаем значение limit и проверяем его тип
	if limitIntf, ok := filtersMap["limit"]; ok {
		if limit, ok = limitIntf.(int32); !ok {
			limit = 0
		}
	}

	// Рассчитываем offset на основе page и limit
	if page > 0 && limit > 0 {
		offset = (page - 1) * limit
	} else {
		offset = 0
		limit = 0
	}

	// Модифицируем строку запроса, чтобы включить LIMIT и OFFSET
	query := fmt.Sprintf(`
		SELECT m.id, m.title, m.description, m.releasedate, m.scorekp, m.scoreimdb, m.poster, m.typemovie, array_agg(g.name)
		FROM movies m
		JOIN moviesgenres mg ON m.id = mg.idmovie
		JOIN genres g ON mg.idgenre = g.id
		WHERE 1=1 %s
		GROUP BY m.id %s
		LIMIT $%d OFFSET $%d`, baseQuery, bestMovieQuery, argPos, argPos+1)

	// Добавляем значения limit и offset в массив аргументов (args)
	args = append(args, limit, offset)

	rows, err := db.Conn.Query(query, args...)
	if err != nil {
		return []models.Movies{}, err
	}
	defer rows.Close()

	found := false
	for rows.Next() {
		var id, releaseDate, typeMovie int32
		var title, description, poster string
		var scoreKP, scoreIMDB float64
		var genresArray []string

		err := rows.Scan(&id, &title, &description, &releaseDate, &scoreKP, &scoreIMDB, &poster, &typeMovie)
		if err != nil {
			return []models.Movies{}, err
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
		return []models.Movies{}, status.Error(codes.NotFound, "нет значений в БД")
	}

	// Save data to Redis
	moviesJSONbyte, err := json.Marshal(movies)
	if err != nil {
		return []models.Movies{}, err
	}
	err = cache.Rdb.Set(cache.Ctx, "movies_filter_"+stringMap, moviesJSONbyte, 5*time.Minute).Err()
	if err != nil {
		return []models.Movies{}, err
	}

	return movies, nil
}

func (s Service) AddMoviesService(moviesMap map[string]interface{}) (int32, error) {
	genres, genresOk := moviesMap["genres"].([]int)
	if !genresOk {
		return 0, fmt.Errorf("error: value is not a []int")
	}

	tx, err := db.Conn.Begin()
	if err != nil {
		return 0, fmt.Errorf("failed to begin transaction: %v", err)
	}

	// Добавление фильма в таблицу movies
	insertMovieSQL := `INSERT INTO movies (title, description, country, releasedate, timemovie, scorekp, scoreimdb, poster, typemovie, views, likes, dislikes)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12) RETURNING id`
	var movieID int
	err = tx.QueryRow(insertMovieSQL, moviesMap["title"], moviesMap["description"], moviesMap["country"], moviesMap["releaseDate"], moviesMap["timeMovie"], moviesMap["scoreKP"], moviesMap["scoreIMDB"], moviesMap["poster"], moviesMap["typeMovie"], moviesMap["views"], moviesMap["likes"], moviesMap["dislikes"]).Scan(&movieID)
	if err != nil {
		tx.Rollback()
		return 0, fmt.Errorf("failed to insert movie: %v", err)
	}

	// Добавление жанров в таблицу moviesgenres
	insertGenresSQL := `INSERT INTO moviesgenres (idmovie, idgenre) VALUES ($1, $2)`
	for _, genreID := range genres {
		_, err = tx.Exec(insertGenresSQL, movieID, genreID)
		if err != nil {
			tx.Rollback()
			return 0, fmt.Errorf("failed to insert movie genre: %v", err)
		}
	}

	// Фиксация транзакции
	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("failed to commit transaction: %v", err)
	}

	return int32(movieID), nil
}

func (s Service) DelMoviesService(id int32) error {
	tx, err := db.Conn.Begin()
	if err != nil {
		return err
	}

	// Удаление записей из таблицы moviesgenres связанных с movieID
	deleteMoviesGenresSQL := `DELETE FROM moviesgenres WHERE idmovie = $1`
	_, err = tx.Exec(deleteMoviesGenresSQL, id)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Удаление фильма из таблицы movies
	deleteMovieSQL := `DELETE FROM movies WHERE id = $1`
	_, err = tx.Exec(deleteMovieSQL, id)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Фиксация транзакции
	if errCommit := tx.Commit(); errCommit != nil {
		return errCommit
	}

	return nil
}

func New() *Service {
	return &Service{}
}
