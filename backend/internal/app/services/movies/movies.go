package movies_v1

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/lib/pq"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"kinogo/internal/app/models"
	"kinogo/pkg/cache"
	"kinogo/pkg/db"
	"kinogo/pkg/movies_v1"
	"log"
	"strings"
	"time"
)

type Service struct {
}

func New() *Service {
	return &Service{}
}

func (s Service) GetMoviesService(limit int32, page int32) ([]models.Movies, error) {
	var movies []models.Movies

	// Получение данных из Redis
	moviesJSON, err := cache.Rdb.Get(cache.Ctx, "movies_"+fmt.Sprint(limit)+"_"+fmt.Sprint(page)).Result()
	if err == nil && moviesJSON != "" {
		err = json.Unmarshal([]byte(moviesJSON), &movies)
		if err != nil {
			log.Printf("Ошибка десериализации: %v", err)
		}
		return movies, nil
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
		  SELECT m.id, m.title, m.description, m.releasedate, m.scorekp, m.scoreimdb, m.poster, m.typemovie, array_agg(g.name)
		  FROM movies m
		  JOIN "moviesGenres" mg ON m.id = mg.idmovie
		  JOIN genres g ON mg.idgenre = g.id
		  GROUP BY m.id
		  ORDER BY m.id DESC %s %s`, limitQuery, pageQuery)

	rows, err := db.Conn.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	found := false
	for rows.Next() {
		var id, releaseDate, typeMovie int32
		var title, description, poster string
		var scoreKP, scoreIMDB float64
		var genresArray pq.StringArray
		errScan := rows.Scan(&id, &title, &description, &releaseDate, &scoreKP, &scoreIMDB, &poster, &typeMovie, &genresArray)
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
	err = cache.Rdb.Set(cache.Ctx, "movies_"+fmt.Sprint(limit)+"_"+fmt.Sprint(page), moviesJSONbyte, 1*time.Minute).Err()
	if err != nil {
		log.Printf("Ошибка при сохранении данных в Redis: %v", err)
	}

	return movies, nil
}

func (s Service) GetMovieByIdService(id int32) (models.Movie, error) {
	var movie models.Movie

	// Получение данных из Redis
	movieJSON, err := cache.Rdb.Get(cache.Ctx, "movies_id_"+fmt.Sprint(id)).Result()
	if err == nil && movieJSON != "" {
		err = json.Unmarshal([]byte(movieJSON), &movie)
		if err != nil {
			log.Fatalf("Ошибка десериализации: %v", err)
		}
		return movie, nil
	} else if !errors.Is(err, redis.Nil) {
		// Если ошибка не связана с отсутствием ключа, логируем её
		log.Printf("Ошибка при получении данных из Redis: %v", err)
	}

	// Data not in Redis, get from database
	query := `
       	SELECT m.*, array_agg(DISTINCT g.name) as genres, array_agg(DISTINCT c.name) as countries
        FROM movies m
        JOIN "moviesGenres" mg ON m.id = mg.idmovie
        JOIN genres g ON mg.idgenre = g.id
        JOIN "moviesCountries" mc ON m.id = mc.idmovie
        JOIN countries c ON mc.idcountry = c.id
        WHERE m.id = $1
        GROUP BY m.id`
	rows, err := db.Conn.Query(query, id)
	if err != nil {
		return models.Movie{}, err
	}
	defer rows.Close()

	found := false
	for rows.Next() {
		var movieID, releaseDate, timeMovie, typeMovie int32
		var title, description, poster string
		var scoreKP, scoreIMDB float64
		var countriesArray, genresArray pq.StringArray

		errScan := rows.Scan(&movieID, &title, &description, &releaseDate, &timeMovie, &scoreKP, &scoreIMDB, &poster, &typeMovie, &genresArray, &countriesArray)
		if errScan != nil {
			return models.Movie{}, errScan
		}

		movie = models.Movie{
			Id:          movieID,
			Title:       title,
			Description: description,
			Country:     strings.Join(countriesArray, ", "),
			ReleaseDate: releaseDate,
			TimeMovie:   timeMovie,
			ScoreKP:     scoreKP,
			ScoreIMDB:   scoreIMDB,
			Poster:      poster,
			TypeMovie:   typeMovie,
			Views:       0,
			Likes:       0,
			Dislikes:    0,
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
		log.Printf("Ошибка при сохранении данных в Redis: %v", err)
	}

	return movie, nil
}

func (s Service) GetMoviesByFilterService(filtersMap map[string]interface{}) ([]models.Movies, error) {
	var movies []models.Movies

	stringMap := fmt.Sprint(filtersMap)

	// Получение данных из Redis
	movieJSON, err := cache.Rdb.Get(cache.Ctx, "movies_filter_"+stringMap).Result()
	if err == nil && movieJSON != "" {
		err = json.Unmarshal([]byte(movieJSON), &movies)
		if err != nil {
			log.Fatalf("Ошибка десериализации: %v", err)
		}
		return movies, nil
	} else if !errors.Is(err, redis.Nil) {
		// Если ошибка не связана с отсутствием ключа, логируем её
		log.Printf("Ошибка при получении данных из Redis: %v", err)
	}

	// Data not in Redis, get from database
	var baseQuery strings.Builder
	var args []interface{}
	argPos := 1

	// Начинаем запрос с основной части
	baseQuery.WriteString(`
		SELECT m.id, m.title, m.description, m.releasedate, m.scorekp, m.scoreimdb, m.poster, m.typemovie, array_agg(g.name), COUNT(mv."movieId") AS views
		FROM movies m
		JOIN "moviesGenres" mg ON m.id = mg.idmovie
		JOIN genres g ON mg.idgenre = g.id
		LEFT JOIN "moviesViews" mv ON m.id = mv."movieId"
		WHERE 1=1`)

	// Условия для года выпуска
	yearMin, okMin := filtersMap["yearMin"].(int32)
	yearMax, okMax := filtersMap["yearMax"].(int32)

	if okMin && okMax {
		if yearMin != 0 && yearMax == 0 {
			baseQuery.WriteString(fmt.Sprintf(" AND m.releasedate >= $%d", argPos))
			args = append(args, yearMin)
			argPos += 1
		} else if yearMin == 0 && yearMax != 0 {
			baseQuery.WriteString(fmt.Sprintf(" AND m.releasedate <= $%d", argPos))
			args = append(args, yearMax)
			argPos += 1
		} else if yearMin != 0 && yearMax != 0 {
			baseQuery.WriteString(fmt.Sprintf(" AND m.releasedate BETWEEN $%d AND $%d", argPos, argPos+1))
			args = append(args, yearMin, yearMax)
			argPos += 2
		}
	}

	// Условия для жанров
	if genres, ok := filtersMap["genres"].([]*movies_v1.Genres); ok && len(genres) > 0 {
		genreNames := make([]string, len(genres))
		for i, genre := range genres {
			genreNames[i] = genre.Name
		}

		genrePlaceholders := make([]string, len(genreNames))
		for i := range genreNames {
			genrePlaceholders[i] = fmt.Sprintf("$%d", argPos)
			args = append(args, genreNames[i])
			argPos++
		}

		baseQuery.WriteString(fmt.Sprintf(`
			AND m.id IN (
				SELECT m.id
				FROM movies m
				JOIN "moviesGenres" mg ON m.id = mg.idmovie
				JOIN genres g ON mg.idgenre = g.id
				WHERE g.name IN (%s)
				GROUP BY m.id
			)`, strings.Join(genrePlaceholders, ",")))
	}

	// Условия для типа фильма
	if typeMovie, ok := filtersMap["typeMovie"].(int32); ok && typeMovie != 0 {
		baseQuery.WriteString(fmt.Sprintf(" AND m.typemovie = $%d", argPos))
		args = append(args, typeMovie)
		argPos++
	}

	// Условия для поиска по названию
	if search, ok := filtersMap["search"].(string); ok && search != "" {
		baseQuery.WriteString(fmt.Sprintf(" AND word_similarity(m.title, $%d) > 0.1", argPos))
		args = append(args, search)
		argPos++
	}

	// Условия для лучшего фильма
	orderBy := ""
	if bestMovie, ok := filtersMap["bestMovie"].(bool); ok && bestMovie {
		orderBy = "ORDER BY views DESC"
	}

	// Устанавливаем LIMIT и OFFSET
	limit, limitOk := filtersMap["limit"].(int32)
	page, pageOk := filtersMap["page"].(int32)

	limitOffset := ""
	if limitOk && limit > 0 {
		offset := int32(0)
		if pageOk && page > 0 {
			offset = (page - 1) * limit
		}

		limitOffset = fmt.Sprintf("LIMIT $%d OFFSET $%d", argPos, argPos+1)
		args = append(args, limit, offset)
		argPos += 2
	}

	// Завершаем запрос и возвращаем его вместе с аргументами
	query := fmt.Sprintf("%s GROUP BY m.id %s %s", baseQuery.String(), orderBy, limitOffset)

	rows, err := db.Conn.Query(query, args...)
	if err != nil {
		return []models.Movies{}, err
	}
	defer rows.Close()

	found := false
	for rows.Next() {
		var id, releaseDate, typeMovie, views int32
		var title, description, poster string
		var scoreKP, scoreIMDB float64
		var genresArray pq.StringArray

		errScan := rows.Scan(&id, &title, &description, &releaseDate, &scoreKP, &scoreIMDB, &poster, &typeMovie, &genresArray, &views)
		if errScan != nil {
			return []models.Movies{}, errScan
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
		log.Printf("Ошибка при сохранении данных в Redis: %v", err)
	}

	return movies, nil
}

func (s Service) AddMoviesService(moviesMap map[string]interface{}) (int32, error) {
	var countries []string
	var genres []string

	if countriesInterface, ok := moviesMap["countries"]; ok {
		if countriesSlice, ok := countriesInterface.([]*movies_v1.Countries); ok {
			for _, country := range countriesSlice {
				countries = append(countries, country.Name)
			}
		}
	}

	if genresInterface, ok := moviesMap["genres"]; ok {
		if genresSlice, ok := genresInterface.([]*movies_v1.Genres); ok {
			for _, genre := range genresSlice {
				genres = append(genres, genre.Name)
			}
		}
	}

	tx, err := db.Conn.Begin()
	if err != nil {
		return 0, fmt.Errorf("failed to begin transaction: %v", err)
	}

	// Добавление фильма в таблицу movies
	insertMovieSQL := `INSERT INTO movies (title, description, releasedate, timemovie, scorekp, scoreimdb, poster, typemovie, views, likes, dislikes) 
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) RETURNING id`
	var movieID int
	err = tx.QueryRow(insertMovieSQL, moviesMap["title"], moviesMap["description"], moviesMap["releaseDate"], moviesMap["timeMovie"], moviesMap["scoreKP"], moviesMap["scoreIMDB"], moviesMap["poster"], moviesMap["typeMovie"], 1, 1, 1).Scan(&movieID)
	if err != nil {
		tx.Rollback()
		return 0, fmt.Errorf("failed to insert movie: %v", err)
	}

	// Добавление стран в таблицу moviesCountries
	insertCountriesSQL := `INSERT INTO "moviesCountries" (idmovie, idcountry) VALUES ($1, $2)`
	for _, countryName := range countries {
		// Получение ID страны из таблицы countries
		var countryID int
		err = tx.QueryRow(`SELECT id FROM countries WHERE name = $1`, countryName).Scan(&countryID)
		if err != nil {
			tx.Rollback()
			return 0, fmt.Errorf("failed to get country ID: %v", err)
		}

		// Вставка связи фильма и страны
		_, err = tx.Exec(insertCountriesSQL, movieID, countryID)
		if err != nil {
			tx.Rollback()
			return 0, fmt.Errorf("failed to insert movie country: %v", err)
		}
	}

	// Добавление жанров в таблицу moviesgenres
	insertGenresSQL := `INSERT INTO "moviesGenres" (idmovie, idgenre) VALUES ($1, $2)`
	for _, genreName := range genres {
		// Получение ID жанра из таблицы genres
		var genreID int
		err = tx.QueryRow(`SELECT id FROM genres WHERE name = $1`, genreName).Scan(&genreID)
		if err != nil {
			tx.Rollback()
			return 0, fmt.Errorf("failed to get genre ID: %v", err)
		}

		// Вставка связи фильма и жанра
		_, err = tx.Exec(insertGenresSQL, movieID, genreID)
		if err != nil {
			tx.Rollback()
			return 0, fmt.Errorf("failed to insert movie genre: %v", err)
		}
	}

	// Фиксация транзакции
	if errCommit := tx.Commit(); errCommit != nil {
		return 0, fmt.Errorf("failed to commit transaction: %v", errCommit)
	}

	return int32(movieID), nil
}

func (s Service) DeleteMoviesService(id int32) error {
	tx, err := db.Conn.Begin()
	if err != nil {
		return err
	}

	// Удаление связей фильма и стран
	deleteCountriesSQL := `DELETE FROM "moviesCountries" WHERE idmovie = $1`
	_, err = tx.Exec(deleteCountriesSQL, id)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to delete movie countries: %v", err)
	}

	// Удаление записей из таблицы moviesgenres связанных с movieID
	deleteMoviesGenresSQL := `DELETE FROM "moviesGenres" WHERE idmovie = $1`
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
