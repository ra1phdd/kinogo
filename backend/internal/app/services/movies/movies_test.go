package movies_v1

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/lib/pq"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"kinogo/internal/app/models"
	"kinogo/internal/app/services/testutil"
	"kinogo/pkg/cache"
	"kinogo/pkg/db"
	"kinogo/pkg/movies_v1"
	"testing"
)

func TestMovies_GetMovies(t *testing.T) {
	conn, mockDB, mock, mr, rdb := testutil.SetupMocks()
	defer mockDB.Close()
	defer mr.Close()

	db.Conn = conn
	cache.Rdb = rdb

	s := Service{}

	t.Run("Data from Redis", func(t *testing.T) {
		// Подготавливаем данные для Redis
		movies := []models.Movies{
			{Id: 1, Title: "Movie 1", Genres: "Action,Drama"},
			{Id: 2, Title: "Movie 2", Genres: "Comedy"},
		}
		moviesJSON, _ := json.Marshal(movies)
		err := mr.Set("movies_10_1", string(moviesJSON))
		if err != nil {
			t.Fatalf("error mr set")
		}

		result, err := s.GetMoviesService(10, 1)
		assert.NoError(t, err)
		assert.Equal(t, movies, result)
	})

	t.Run("Data from Database", func(t *testing.T) {
		mr.FlushAll()

		rows := sqlmock.NewRows([]string{"id", "title", "description", "releasedate", "scorekp", "scoreimdb", "poster", "typemovie", "genres"}).
			AddRow(1, "Movie 1", "Description 1", 2021, 7.5, 7.8, "poster1.jpg", 1, pq.StringArray{"Action", "Drama"}).
			AddRow(2, "Movie 2", "Description 2", 2022, 8.0, 8.2, "poster2.jpg", 2, pq.StringArray{"Comedy"})

		mock.ExpectQuery("SELECT m.id, m.title, m.description, m.releasedate, m.scorekp, m.scoreimdb, m.poster, m.typemovie, array_agg\\(g.name\\)").
			WillReturnRows(rows)

		result, err := s.GetMoviesService(10, 1)
		assert.NoError(t, err)
		assert.Len(t, result, 2)
		assert.Equal(t, "Movie 1", result[0].Title)
		assert.Equal(t, "Movie 2", result[1].Title)
	})

	t.Run("No Data", func(t *testing.T) {
		mr.FlushAll()

		mock.ExpectQuery("SELECT m.id, m.title, m.description, m.releasedate, m.scorekp, m.scoreimdb, m.poster, m.typemovie, array_agg\\(g.name\\)").
			WillReturnRows(sqlmock.NewRows([]string{"id", "title", "description", "releasedate", "scorekp", "scoreimdb", "poster", "typemovie", "genres"}))

		result, err := s.GetMoviesService(10, 1)
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "нет значений в БД")
	})

	t.Run("Database Error", func(t *testing.T) {
		mr.FlushAll()

		mock.ExpectQuery("SELECT m.id, m.title, m.description, m.releasedate, m.scorekp, m.scoreimdb, m.poster, m.typemovie, array_agg\\(g.name\\)").
			WillReturnError(errors.New("database error"))

		result, err := s.GetMoviesService(10, 1)
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "database error")
	})

	t.Run("Redis Error", func(t *testing.T) {
		// Симулируем ошибку Redis
		cache.Rdb = redis.NewClient(&redis.Options{
			Addr: "localhost:6379", // Неправильный адрес
		})

		rows := sqlmock.NewRows([]string{"id", "title", "description", "releasedate", "scorekp", "scoreimdb", "poster", "typemovie", "genres"}).
			AddRow(1, "Movie 1", "Description 1", 2021, 7.5, 7.8, "poster1.jpg", 1, pq.StringArray{"Action", "Drama"})

		mock.ExpectQuery("SELECT m.id, m.title, m.description, m.releasedate, m.scorekp, m.scoreimdb, m.poster, m.typemovie, array_agg\\(g.name\\)").
			WillReturnRows(rows)

		result, err := s.GetMoviesService(10, 1)
		assert.NoError(t, err) // Функция должна вернуть данные из БД без ошибки
		assert.Len(t, result, 1)
		assert.Equal(t, "Movie 1", result[0].Title)

		// Восстанавливаем правильное подключение к Redis
		cache.Rdb = redis.NewClient(&redis.Options{
			Addr: mr.Addr(),
		})
	})
}

func TestMovies_GetMovieByIdService(t *testing.T) {
	conn, mockDB, mock, mr, rdb := testutil.SetupMocks()
	defer mockDB.Close()
	defer mr.Close()

	db.Conn = conn
	cache.Rdb = rdb

	s := Service{}

	t.Run("Data from Redis", func(t *testing.T) {
		// Подготавливаем данные для Redis
		movie := models.Movie{
			Id:          1,
			Title:       "Movie 1",
			Description: "Description 1",
			Country:     "USA, UK",
			ReleaseDate: 2021,
			TimeMovie:   120,
			ScoreKP:     7.5,
			ScoreIMDB:   7.8,
			Poster:      "poster1.jpg",
			TypeMovie:   1,
			Genres:      "Action, Drama",
		}
		movieJSON, _ := json.Marshal(movie)
		err := mr.Set("movies_id_1", string(movieJSON))
		if err != nil {
			t.Fatalf("error mr set")
		}

		result, err := s.GetMovieByIdService(1)
		assert.NoError(t, err)
		assert.Equal(t, movie, result)
	})

	t.Run("Data from Database", func(t *testing.T) {
		mr.FlushAll()

		rows := sqlmock.NewRows([]string{"id", "title", "description", "releasedate", "timemovie", "scorekp", "scoreimdb", "poster", "typemovie", "genres", "countries"}).
			AddRow(1, "Movie 1", "Description 1", 2021, 120, 7.5, 7.8, "poster1.jpg", 1, pq.StringArray{"Action", "Drama"}, pq.StringArray{"USA", "UK"})

		mock.ExpectQuery("SELECT m.*, array_agg\\(DISTINCT g.name\\) as genres, array_agg\\(DISTINCT c.name\\) as countries").
			WithArgs(int32(1)).
			WillReturnRows(rows)

		result, err := s.GetMovieByIdService(1)
		assert.NoError(t, err)
		assert.Equal(t, int32(1), result.Id)
		assert.Equal(t, "Movie 1", result.Title)
		assert.Equal(t, "Action, Drama", result.Genres)
		assert.Equal(t, "USA, UK", result.Country)
	})

	t.Run("No Data", func(t *testing.T) {
		mr.FlushAll()

		mock.ExpectQuery("SELECT m.*, array_agg\\(DISTINCT g.name\\) as genres, array_agg\\(DISTINCT c.name\\) as countries").
			WithArgs(int32(1)).
			WillReturnRows(sqlmock.NewRows([]string{"id", "title", "description", "releasedate", "timemovie", "scorekp", "scoreimdb", "poster", "typemovie", "genres", "countries"}))

		result, err := s.GetMovieByIdService(1)
		assert.Error(t, err)
		assert.Equal(t, models.Movie{}, result)
		assert.Contains(t, err.Error(), "нет значений в БД")
	})

	t.Run("Database Error", func(t *testing.T) {
		mr.FlushAll()

		mock.ExpectQuery("SELECT m.*, array_agg\\(DISTINCT g.name\\) as genres, array_agg\\(DISTINCT c.name\\) as countries").
			WithArgs(int32(1)).
			WillReturnError(errors.New("database error"))

		result, err := s.GetMovieByIdService(1)
		assert.Error(t, err)
		assert.Equal(t, models.Movie{}, result)

		assert.Contains(t, err.Error(), "database error")
	})

	t.Run("Redis Error", func(t *testing.T) {
		// Симулируем ошибку Redis
		cache.Rdb = redis.NewClient(&redis.Options{
			Addr: "localhost:6379", // Неправильный адрес
		})

		rows := sqlmock.NewRows([]string{"id", "title", "description", "releasedate", "timemovie", "scorekp", "scoreimdb", "poster", "typemovie", "genres", "countries"}).
			AddRow(1, "Movie 1", "Description 1", 2021, 120, 7.5, 7.8, "poster1.jpg", 1, pq.StringArray{"Action", "Drama"}, pq.StringArray{"USA", "UK"})

		mock.ExpectQuery("SELECT m.*, array_agg\\(DISTINCT g.name\\) as genres, array_agg\\(DISTINCT c.name\\) as countries").
			WithArgs(int32(1)).
			WillReturnRows(rows)

		result, err := s.GetMovieByIdService(1)
		assert.NoError(t, err) // Функция должна вернуть данные из БД без ошибки
		assert.Equal(t, int32(1), result.Id)
		assert.Equal(t, "Movie 1", result.Title)

		// Восстанавливаем правильное подключение к Redis
		cache.Rdb = redis.NewClient(&redis.Options{
			Addr: mr.Addr(),
		})
	})
}

func TestService_GetMoviesByFilterService(t *testing.T) {
	conn, mockDB, mock, mr, rdb := testutil.SetupMocks()
	defer mockDB.Close()
	defer mr.Close()

	db.Conn = conn
	cache.Rdb = rdb

	s := Service{}

	t.Run("Data from Redis", func(t *testing.T) {
		// Подготавливаем данные для Redis
		movies := []models.Movies{
			{Id: 1, Title: "Movie 1", Genres: "Action,Drama"},
			{Id: 2, Title: "Movie 2", Genres: "Comedy"},
		}
		moviesJSON, _ := json.Marshal(movies)
		filtersMap := map[string]interface{}{"yearMin": int32(2000), "yearMax": int32(2022)}
		stringMap := fmt.Sprint(filtersMap)
		err := mr.Set("movies_filter_"+stringMap, string(moviesJSON))
		if err != nil {
			t.Fatalf("error mr set")
		}

		result, err := s.GetMoviesByFilterService(filtersMap)
		assert.NoError(t, err)
		assert.Equal(t, movies, result)
	})

	t.Run("Data from Database", func(t *testing.T) {
		mr.FlushAll()

		filtersMap := map[string]interface{}{
			"yearMin":   int32(2000),
			"yearMax":   int32(2022),
			"genres":    []*movies_v1.Genres{{Name: "Action"}},
			"typeMovie": int32(1),
			"search":    "Movie",
			"bestMovie": true,
			"limit":     int32(10),
			"page":      int32(1),
		}

		rows := sqlmock.NewRows([]string{"id", "title", "description", "releasedate", "scorekp", "scoreimdb", "poster", "typemovie", "genres", "views"}).
			AddRow(1, "Movie 1", "Description 1", 2021, 7.5, 7.8, "poster1.jpg", 1, pq.StringArray{"Action", "Drama"}, 100).
			AddRow(2, "Movie 2", "Description 2", 2022, 8.0, 8.2, "poster2.jpg", 1, pq.StringArray{"Action", "Comedy"}, 200)

		mock.ExpectQuery(`SELECT m.id, m.title, m.description, m.releasedate, m.scorekp, m.scoreimdb, m.poster, m.typemovie, array_agg\(g.name\), COUNT\(mv."movieId"\) AS views`).
			WithArgs(2000, 2022, "Action", int32(1), "Movie", int32(10), int32(0)).
			WillReturnRows(rows)

		result, err := s.GetMoviesByFilterService(filtersMap)
		assert.NoError(t, err)
		assert.Len(t, result, 2)
		assert.Equal(t, "Movie 1", result[0].Title)
		assert.Equal(t, "Movie 2", result[1].Title)
	})

	t.Run("No Data", func(t *testing.T) {
		mr.FlushAll()

		filtersMap := map[string]interface{}{"yearMin": int32(2000), "yearMax": int32(2022)}

		mock.ExpectQuery(`SELECT m.id, m.title, m.description, m.releasedate, m.scorekp, m.scoreimdb, m.poster, m.typemovie, array_agg\(g.name\), COUNT\(mv."movieId"\) AS views`).
			WithArgs(2000, 2022).
			WillReturnRows(sqlmock.NewRows([]string{"id", "title", "description", "releasedate", "scorekp", "scoreimdb", "poster", "typemovie", "genres", "views"}))

		result, err := s.GetMoviesByFilterService(filtersMap)
		assert.Error(t, err)
		assert.Empty(t, result)
		assert.Contains(t, err.Error(), "нет значений в БД")
	})

	t.Run("Database Error", func(t *testing.T) {
		mr.FlushAll()

		filtersMap := map[string]interface{}{"yearMin": int32(2000), "yearMax": int32(2022)}

		mock.ExpectQuery(`SELECT m.id, m.title, m.description, m.releasedate, m.scorekp, m.scoreimdb, m.poster, m.typemovie, array_agg\(g.name\), COUNT\(mv."movieId"\) AS views`).
			WithArgs(2000, 2022).
			WillReturnError(errors.New("database error"))

		result, err := s.GetMoviesByFilterService(filtersMap)
		assert.Error(t, err)
		assert.Empty(t, result)
		assert.Contains(t, err.Error(), "database error")
	})

	t.Run("Redis Error", func(t *testing.T) {
		// Симулируем ошибку Redis
		cache.Rdb = redis.NewClient(&redis.Options{
			Addr: "localhost:6379", // Неправильный адрес
		})

		filtersMap := map[string]interface{}{"yearMin": int32(2000), "yearMax": int32(2022)}

		rows := sqlmock.NewRows([]string{"id", "title", "description", "releasedate", "scorekp", "scoreimdb", "poster", "typemovie", "genres", "views"}).
			AddRow(1, "Movie 1", "Description 1", 2021, 7.5, 7.8, "poster1.jpg", 1, pq.StringArray{"Action", "Drama"}, 100)

		mock.ExpectQuery(`SELECT m.id, m.title, m.description, m.releasedate, m.scorekp, m.scoreimdb, m.poster, m.typemovie, array_agg\(g.name\), COUNT\(mv."movieId"\) AS views`).
			WithArgs(2000, 2022).
			WillReturnRows(rows)

		result, err := s.GetMoviesByFilterService(filtersMap)
		assert.NoError(t, err) // Функция должна вернуть данные из БД без ошибки
		assert.Len(t, result, 1)
		assert.Equal(t, "Movie 1", result[0].Title)

		// Восстанавливаем правильное подключение к Redis
		cache.Rdb = redis.NewClient(&redis.Options{
			Addr: mr.Addr(),
		})
	})
}

func TestService_AddMoviesService(t *testing.T) {
	conn, mockDB, mock, _, _ := testutil.SetupMocks()
	defer mockDB.Close()

	db.Conn = conn

	s := Service{}

	t.Run("Successful Addition", func(t *testing.T) {
		moviesMap := map[string]interface{}{
			"title":       "Test Movie",
			"description": "Test Description",
			"releaseDate": 2023,
			"timeMovie":   120,
			"scoreKP":     8.0,
			"scoreIMDB":   7.5,
			"poster":      "test_poster.jpg",
			"typeMovie":   1,
			"countries":   []*movies_v1.Countries{{Name: "USA"}, {Name: "UK"}},
			"genres":      []*movies_v1.Genres{{Name: "Action"}, {Name: "Drama"}},
		}

		mock.ExpectBegin()
		mock.ExpectQuery("INSERT INTO movies").
			WithArgs(moviesMap["title"], moviesMap["description"], moviesMap["releaseDate"], moviesMap["timeMovie"], moviesMap["scoreKP"], moviesMap["scoreIMDB"], moviesMap["poster"], moviesMap["typeMovie"], 1, 1, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

		mock.ExpectQuery("SELECT id FROM countries WHERE name = ?").
			WithArgs("USA").
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
		mock.ExpectExec("INSERT INTO \"moviesCountries\"").
			WithArgs(1, 1).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectQuery("SELECT id FROM countries WHERE name = ?").
			WithArgs("UK").
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(2))
		mock.ExpectExec("INSERT INTO \"moviesCountries\"").
			WithArgs(1, 2).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectQuery("SELECT id FROM genres WHERE name = ?").
			WithArgs("Action").
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
		mock.ExpectExec("INSERT INTO \"moviesGenres\"").
			WithArgs(1, 1).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectQuery("SELECT id FROM genres WHERE name = ?").
			WithArgs("Drama").
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(2))
		mock.ExpectExec("INSERT INTO \"moviesGenres\"").
			WithArgs(1, 2).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectCommit()

		movieID, err := s.AddMoviesService(moviesMap)
		assert.NoError(t, err)
		assert.Equal(t, int32(1), movieID)
	})

	t.Run("Transaction Begin Error", func(t *testing.T) {
		moviesMap := map[string]interface{}{
			"title":       "Test Movie",
			"description": "Test Description",
			"releaseDate": 2023,
			"timeMovie":   120,
			"scoreKP":     8.0,
			"scoreIMDB":   7.5,
			"poster":      "test_poster.jpg",
			"typeMovie":   1,
			"countries":   []*movies_v1.Countries{},
			"genres":      []*movies_v1.Genres{},
		}

		mock.ExpectBegin().WillReturnError(errors.New("begin error"))

		_, err := s.AddMoviesService(moviesMap)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to begin transaction")
	})

	t.Run("Insert Movie Error", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectQuery("INSERT INTO movies").WillReturnError(errors.New("insert error"))
		mock.ExpectRollback()

		_, err := s.AddMoviesService(map[string]interface{}{
			"title":       "Test Movie",
			"description": "Test Description",
			"releaseDate": 2023,
			"timeMovie":   120,
			"scoreKP":     8.0,
			"scoreIMDB":   7.5,
			"poster":      "test_poster.jpg",
			"typeMovie":   1,
			"countries":   []*movies_v1.Countries{},
			"genres":      []*movies_v1.Genres{},
		})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to insert movie")
	})

	t.Run("Get Country ID Error", func(t *testing.T) {
		moviesMap := map[string]interface{}{
			"title":       "Test Movie",
			"description": "Test Description",
			"releaseDate": 2023,
			"timeMovie":   120,
			"scoreKP":     8.0,
			"scoreIMDB":   7.5,
			"poster":      "test_poster.jpg",
			"typeMovie":   1,
			"countries":   []*movies_v1.Countries{{Name: "NonexistentCountry"}},
			"genres":      []*movies_v1.Genres{},
		}

		mock.ExpectBegin()
		mock.ExpectQuery("INSERT INTO movies").
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
		mock.ExpectQuery("SELECT id FROM countries WHERE name = ?").
			WithArgs("NonexistentCountry").
			WillReturnError(sql.ErrNoRows)
		mock.ExpectRollback()

		_, err := s.AddMoviesService(moviesMap)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to get country ID")
	})

	t.Run("Insert Country Error", func(t *testing.T) {
		moviesMap := map[string]interface{}{
			"title":       "Test Movie",
			"description": "Test Description",
			"releaseDate": 2023,
			"timeMovie":   120,
			"scoreKP":     8.0,
			"scoreIMDB":   7.5,
			"poster":      "test_poster.jpg",
			"typeMovie":   1,
			"countries":   []*movies_v1.Countries{{Name: "USA"}},
			"genres":      []*movies_v1.Genres{},
		}

		mock.ExpectBegin()
		mock.ExpectQuery("INSERT INTO movies").
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
		mock.ExpectQuery("SELECT id FROM countries WHERE name = ?").
			WithArgs("USA").
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
		mock.ExpectExec("INSERT INTO \"moviesCountries\"").
			WithArgs(1, 1).
			WillReturnError(errors.New("insert country error"))
		mock.ExpectRollback()

		_, err := s.AddMoviesService(moviesMap)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to insert movie country")
	})

	t.Run("Get Genre ID Error", func(t *testing.T) {
		moviesMap := map[string]interface{}{
			"title":       "Test Movie",
			"description": "Test Description",
			"releaseDate": 2023,
			"timeMovie":   120,
			"scoreKP":     8.0,
			"scoreIMDB":   7.5,
			"poster":      "test_poster.jpg",
			"typeMovie":   1,
			"countries":   []*movies_v1.Countries{},
			"genres":      []*movies_v1.Genres{{Name: "NonexistentGenre"}},
		}

		mock.ExpectBegin()
		mock.ExpectQuery("INSERT INTO movies").
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
		mock.ExpectQuery("SELECT id FROM genres WHERE name = ?").
			WithArgs("NonexistentGenre").
			WillReturnError(sql.ErrNoRows)
		mock.ExpectRollback()

		_, err := s.AddMoviesService(moviesMap)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to get genre ID")
	})

	t.Run("Insert Genre Error", func(t *testing.T) {
		moviesMap := map[string]interface{}{
			"title":       "Test Movie",
			"description": "Test Description",
			"releaseDate": 2023,
			"timeMovie":   120,
			"scoreKP":     8.0,
			"scoreIMDB":   7.5,
			"poster":      "test_poster.jpg",
			"typeMovie":   1,
			"countries":   []*movies_v1.Countries{},
			"genres":      []*movies_v1.Genres{{Name: "Action"}},
		}

		mock.ExpectBegin()
		mock.ExpectQuery("INSERT INTO movies").
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
		mock.ExpectQuery("SELECT id FROM genres WHERE name = ?").
			WithArgs("Action").
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
		mock.ExpectExec("INSERT INTO \"moviesGenres\"").
			WithArgs(1, 1).
			WillReturnError(errors.New("insert genre error"))
		mock.ExpectRollback()

		_, err := s.AddMoviesService(moviesMap)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to insert movie genre")
	})

	t.Run("Commit Error", func(t *testing.T) {
		moviesMap := map[string]interface{}{
			"title":       "Test Movie",
			"description": "Test Description",
			"releaseDate": 2023,
			"timeMovie":   120,
			"scoreKP":     8.0,
			"scoreIMDB":   7.5,
			"poster":      "test_poster.jpg",
			"typeMovie":   1,
			"countries":   []*movies_v1.Countries{},
			"genres":      []*movies_v1.Genres{},
		}

		mock.ExpectBegin()
		mock.ExpectQuery("INSERT INTO movies").
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
		mock.ExpectCommit().WillReturnError(errors.New("commit error"))

		_, err := s.AddMoviesService(moviesMap)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to commit transaction")
	})
}

func TestMovies_DeleteMoviesService(t *testing.T) {
	conn, mockDB, mock, mr, rdb := testutil.SetupMocks()
	defer mockDB.Close()
	defer mr.Close()

	db.Conn = conn
	cache.Rdb = rdb

	s := Service{}

	t.Run("Successful Deletion", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(`DELETE FROM "moviesCountries" WHERE idmovie = \$1`).
			WithArgs(int32(1)).
			WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectExec(`DELETE FROM "moviesGenres" WHERE idmovie = \$1`).
			WithArgs(int32(1)).
			WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectExec(`DELETE FROM movies WHERE id = \$1`).
			WithArgs(int32(1)).
			WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectCommit()

		err := s.DeleteMoviesService(1)
		assert.NoError(t, err)
	})

	t.Run("Begin Transaction Error", func(t *testing.T) {
		mock.ExpectBegin().WillReturnError(errors.New("begin transaction error"))

		err := s.DeleteMoviesService(1)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "begin transaction error")
	})

	t.Run("Delete Countries Error", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(`DELETE FROM "moviesCountries" WHERE idmovie = \$1`).
			WithArgs(int32(1)).
			WillReturnError(errors.New("delete countries error"))
		mock.ExpectRollback()

		err := s.DeleteMoviesService(1)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to delete movie countries: delete countries error")
	})

	t.Run("Delete Genres Error", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(`DELETE FROM "moviesCountries" WHERE idmovie = \$1`).
			WithArgs(int32(1)).
			WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectExec(`DELETE FROM "moviesGenres" WHERE idmovie = \$1`).
			WithArgs(int32(1)).
			WillReturnError(errors.New("delete genres error"))
		mock.ExpectRollback()

		err := s.DeleteMoviesService(1)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "delete genres error")
	})

	t.Run("Delete Movie Error", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(`DELETE FROM "moviesCountries" WHERE idmovie = \$1`).
			WithArgs(int32(1)).
			WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectExec(`DELETE FROM "moviesGenres" WHERE idmovie = \$1`).
			WithArgs(int32(1)).
			WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectExec(`DELETE FROM movies WHERE id = \$1`).
			WithArgs(int32(1)).
			WillReturnError(errors.New("delete movie error"))
		mock.ExpectRollback()

		err := s.DeleteMoviesService(1)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "delete movie error")
	})

	t.Run("Commit Error", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(`DELETE FROM "moviesCountries" WHERE idmovie = \$1`).
			WithArgs(int32(1)).
			WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectExec(`DELETE FROM "moviesGenres" WHERE idmovie = \$1`).
			WithArgs(int32(1)).
			WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectExec(`DELETE FROM movies WHERE id = \$1`).
			WithArgs(int32(1)).
			WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectCommit().WillReturnError(errors.New("commit error"))

		err := s.DeleteMoviesService(1)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "commit error")
	})

	t.Run("No Rows Affected", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(`DELETE FROM "moviesCountries" WHERE idmovie = \$1`).
			WithArgs(int32(1)).
			WillReturnResult(sqlmock.NewResult(0, 0))
		mock.ExpectExec(`DELETE FROM "moviesGenres" WHERE idmovie = \$1`).
			WithArgs(int32(1)).
			WillReturnResult(sqlmock.NewResult(0, 0))
		mock.ExpectExec(`DELETE FROM movies WHERE id = \$1`).
			WithArgs(int32(1)).
			WillReturnResult(sqlmock.NewResult(0, 0))
		mock.ExpectCommit()

		err := s.DeleteMoviesService(1)
		assert.NoError(t, err)
	})
}
