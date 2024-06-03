package endpoint

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"kinogo/internal/app/models"
	"kinogo/pkg/cache"
	"kinogo/pkg/db"
	"kinogo/pkg/logger"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Service interface {
	ParseTemplatesMain(w http.ResponseWriter, allData models.AllData) error
	GetMoviesFromDB(query string, args ...string) ([]models.MovieData, error)
}

type Endpoint struct {
	s Service
}

func New(s Service) *Endpoint {
	return &Endpoint{
		s: s,
	}
}

func (e *Endpoint) GetAllContents(c *gin.Context) {
	var moviesSlice []models.MovieData

	// Try to get data from Redis
	movies, err := cache.Rdb.Get(cache.Ctx, "QueryAllContents").Result()
	if err == nil && movies != "" {
		var cachedMovies []models.MovieData
		if err := json.Unmarshal([]byte(movies), &cachedMovies); err == nil {
			c.JSON(http.StatusOK, cachedMovies)
			return
		}
	}

	// Data not in Redis, get from database
	query := "SELECT movies.*, array_agg(genres.name) AS genres FROM movies JOIN moviesgenres ON movies.id = moviesgenres.idmovie JOIN genres ON moviesgenres.idgenre = genres.id GROUP BY movies.id ORDER BY movies.id DESC"
	moviesSlice, err = e.s.GetMoviesFromDB(query, "")
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	// Save data to Redis
	moviesJSON, err := json.Marshal(moviesSlice)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "впаыдлопыадловьдльоа"})
		return
	}
	err = cache.Rdb.Set(cache.Ctx, "QueryAllContents", moviesJSON, 1*time.Minute).Err()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving contents"})
		return
	}

	// Return contents as JSON
	c.JSON(http.StatusOK, moviesSlice)
}

func (e *Endpoint) GetAllMovies(c *gin.Context) {
	var moviesSlice []models.MovieData

	// Try to get data from Redis
	movies, err := cache.Rdb.Get(cache.Ctx, "QueryAllMovies").Result()
	if err == nil && movies != "" {
		var cachedMovies []models.MovieData
		if err := json.Unmarshal([]byte(movies), &cachedMovies); err == nil {
			c.JSON(http.StatusOK, cachedMovies)
			return
		}
	}

	// Data not in Redis, get from database
	query := "SELECT movies.*, array_agg(genres.name) AS genres FROM movies JOIN moviesgenres ON movies.id = moviesgenres.idmovie JOIN genres ON moviesgenres.idgenre = genres.id WHERE movies.typemovie = 'movie' GROUP BY movies.id ORDER BY movies.id DESC"
	moviesSlice, err = e.s.GetMoviesFromDB(query, "")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	// Save data to Redis
	moviesJSON, err := json.Marshal(moviesSlice)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "впаыдлопыадловьдльоа"})
		return
	}
	err = cache.Rdb.Set(cache.Ctx, "QueryAllMovies", moviesJSON, 1*time.Minute).Err()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving contents"})
		return
	}

	// Return contents as JSON
	c.JSON(http.StatusOK, moviesSlice)
}

func (e *Endpoint) GetAllCartoons(c *gin.Context) {
	var moviesSlice []models.MovieData

	// Try to get data from Redis
	movies, err := cache.Rdb.Get(cache.Ctx, "QueryAllCartoons").Result()
	if err == nil && movies != "" {
		var cachedMovies []models.MovieData
		if err := json.Unmarshal([]byte(movies), &cachedMovies); err == nil {
			c.JSON(http.StatusOK, cachedMovies)
			return
		}
	}

	// Data not in Redis, get from database
	query := "SELECT movies.*, array_agg(genres.name) AS genres FROM movies JOIN moviesgenres ON movies.id = moviesgenres.idmovie JOIN genres ON moviesgenres.idgenre = genres.id WHERE movies.typemovie = 'cartoon' GROUP BY movies.id ORDER BY movies.id DESC"
	moviesSlice, err = e.s.GetMoviesFromDB(query, "")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	// Save data to Redis
	moviesJSON, err := json.Marshal(moviesSlice)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "впаыдлопыадловьдльоа"})
		return
	}
	err = cache.Rdb.Set(cache.Ctx, "QueryAllCartoons", moviesJSON, 1*time.Minute).Err()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving contents"})
		return
	}

	// Return contents as JSON
	c.JSON(http.StatusOK, moviesSlice)
}

func (e *Endpoint) GetAllTelecasts(c *gin.Context) {
	var moviesSlice []models.MovieData

	// Try to get data from Redis
	movies, err := cache.Rdb.Get(cache.Ctx, "QueryAllTelecasts").Result()
	if err == nil && movies != "" {
		var cachedMovies []models.MovieData
		if err := json.Unmarshal([]byte(movies), &cachedMovies); err == nil {
			c.JSON(http.StatusOK, cachedMovies)
			return
		}
	}

	// Data not in Redis, get from database
	query := "SELECT movies.*, array_agg(genres.name) AS genres FROM movies JOIN moviesgenres ON movies.id = moviesgenres.idmovie JOIN genres ON moviesgenres.idgenre = genres.id WHERE movies.typemovie = 'telecast' GROUP BY movies.id ORDER BY movies.id DESC"
	moviesSlice, err = e.s.GetMoviesFromDB(query, "")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	// Save data to Redis
	moviesJSON, err := json.Marshal(moviesSlice)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "впаыдлопыадловьдльоа"})
		return
	}
	err = cache.Rdb.Set(cache.Ctx, "QueryAllTelecasts", moviesJSON, 1*time.Minute).Err()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving contents"})
		return
	}

	// Return contents as JSON
	c.JSON(http.StatusOK, moviesSlice)
}

func (e *Endpoint) GetMovieByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid movie ID"})
		return
	}

	var moviesSlice models.MovieData
	movies, err := cache.Rdb.Get(cache.Ctx, "QueryMovie_"+fmt.Sprint(id)).Result()
	if err == nil && movies != "" {
		var cachedMovies []models.MovieData
		if err := json.Unmarshal([]byte(movies), &cachedMovies); err == nil {
			c.JSON(http.StatusOK, cachedMovies)
			return
		}
	}

	// Data not in Redis, get from database
	query := "SELECT movies.*,array_agg(genres.name) AS genres FROM movies	JOIN moviesgenres ON movies.id = moviesgenres.idmovie JOIN genres ON moviesgenres.idgenre = genres.id WHERE movies.id = $1 GROUP BY movies.id"
	moviesSliceData, err := e.s.GetMoviesFromDB(query, idParam)
	if err != nil {
		if err.Error() != "Нет данных в БД" {
			c.JSON(http.StatusNotFound, gin.H{"error": err})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	moviesSlice = moviesSliceData[0]

	// Save data to Redis
	moviesJSON, err := json.Marshal(moviesSlice)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "впаыдлопыадловьдльоа"})
		return
	}
	err = cache.Rdb.Set(cache.Ctx, "QueryMovie_"+fmt.Sprint(id), moviesJSON, 1*time.Minute).Err()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving contents"})
		return
	}

	// Return contents as JSON
	c.JSON(http.StatusOK, moviesSlice)
}

func (e *Endpoint) FilterMovies(c *gin.Context) {
	arrayGenre := c.QueryArray("genre")
	yearMinStr := c.Query("year__min")
	yearMaxStr := c.Query("year__max")

	logger.Debug("Получение данных фильтра из формы", zap.Any("Жанры", arrayGenre), zap.String("Минимальный год", yearMinStr), zap.String("Максимальный год", yearMaxStr))

	nameGenres := make([]string, 0, len(arrayGenre))
	for _, genre := range arrayGenre {
		if genre == "выбрать" {
			break
		}
		nameGenres = append(nameGenres, genre)
	}

	var idGenres []int
	if len(nameGenres) > 0 {
		query := `SELECT id FROM genres WHERE name IN (?)`
		query, args, err := sqlx.In(query, nameGenres)
		if err != nil {
			logger.Error("ошибка при создании запроса для получения id жанров по названию: ", zap.Error(err), zap.Any("Названия жанров", nameGenres))
			return
		}

		query = db.Conn.Rebind(query)
		err = db.Conn.Select(&idGenres, query, args...)
		if err != nil {
			logger.Error("Ошибка при получении id жанров по названию: ", zap.Error(err), zap.Any("Запрос", query), zap.Any("Аргументы", args), zap.Any("Названия жанров", nameGenres))
			return
		}
	}

	var idGenresStr []string
	for _, id := range idGenres {
		idGenresStr = append(idGenresStr, strconv.Itoa(id))
	}

	var moviesSlice []models.MovieData
	movies, err := cache.Rdb.Get(cache.Ctx, "QueryFilterMovies_"+strings.Join(idGenresStr, "_")+"_"+yearMinStr+"_"+yearMaxStr).Result()
	if err == nil && movies != "" {
		var cachedMovies []models.MovieData
		if err := json.Unmarshal([]byte(movies), &cachedMovies); err == nil {
			c.JSON(http.StatusOK, cachedMovies)
			return
		}
	}

	// Данных нет в Redis, получаем их из базы данных
	var query string
	if len(idGenres) > 0 {
		query = "SELECT movies.*, array_agg(genres.name) AS genres FROM movies JOIN moviesgenres ON movies.id = moviesgenres.idmovie JOIN genres ON moviesgenres.idgenre = genres.id WHERE genres.id = ANY($1) AND (releasedate >= $2 AND releasedate <= $3) GROUP BY movies.id ORDER BY movies.id DESC"
		moviesSlice, err = e.s.GetMoviesFromDB(query, strings.Join(nameGenres, ","), yearMinStr, yearMaxStr)
	} else {
		query = "SELECT movies.*, array_agg(genres.name) AS genres FROM movies JOIN moviesgenres ON movies.id = moviesgenres.idmovie JOIN genres ON moviesgenres.idgenre = genres.id WHERE (releasedate >= $1 AND releasedate <= $2) GROUP BY movies.id ORDER BY movies.id DESC"
		moviesSlice, err = e.s.GetMoviesFromDB(query, yearMinStr, yearMaxStr)
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	// Сохраняем данные в Redis
	moviesJSON, err := json.Marshal(moviesSlice)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving contents"})
		return
	}
	err = cache.Rdb.Set(cache.Ctx, "QueryFilterMovies_"+strings.Join(idGenresStr, "_")+"_"+yearMinStr+"_"+yearMaxStr, moviesJSON, 1*time.Minute).Err()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving contents"})
		return
	}

	c.JSON(http.StatusOK, moviesSlice)
}

func (e *Endpoint) SearchMovies(c *gin.Context) {
	textSearch := strings.ToLower(c.Query("text"))
	logger.Debug("Получение текста поиска из формы", zap.String("textSearch", textSearch))

	var moviesSlice []models.MovieData
	movies, err := cache.Rdb.Get(cache.Ctx, "QuerySearchMovies_"+textSearch).Result()
	if err == nil && movies != "" {
		var cachedMovies []models.MovieData
		if err := json.Unmarshal([]byte(movies), &cachedMovies); err == nil {
			c.JSON(http.StatusOK, cachedMovies)
			return
		}
	}

	// Данных нет в Redis, получаем их из базы данных
	query := "SELECT movies.*, array_agg(genres.name) AS genres FROM movies JOIN moviesgenres ON movies.id = moviesgenres.idmovie JOIN genres ON moviesgenres.idgenre = genres.id\n\t\tWHERE word_similarity(movies.title, $1) > 0.1 GROUP BY movies.id"
	moviesSlice, err = e.s.GetMoviesFromDB(query, textSearch)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	// Сохраняем данные в Redis
	var moviesJSON []byte
	moviesJSON, err = json.Marshal(moviesSlice)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving contents"})
		return
	}
	err = cache.Rdb.Set(cache.Ctx, "QuerySearchMovies_"+textSearch, moviesJSON, 1*time.Minute).Err()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving contents"})
		return
	}

	logger.Debug("Получен контент из БД/кэша", zap.Any("movies", movies))
	c.JSON(http.StatusOK, moviesSlice)
}

func (e *Endpoint) GetBestMovie(c *gin.Context) {
	var moviesSlice models.MovieData
	movies, err := cache.Rdb.Get(cache.Ctx, "QueryBestMovie").Result()
	if err == nil && movies != "" {
		var cachedMovies []models.MovieData
		if err := json.Unmarshal([]byte(movies), &cachedMovies); err == nil {
			c.JSON(http.StatusOK, cachedMovies)
			return
		}
	}

	// Данных нет в Redis, получаем их из базы данных
	query := "SELECT movies.*, array_agg(genres.name) AS genres FROM movies JOIN moviesgenres ON movies.id = moviesgenres.idmovie JOIN genres ON moviesgenres.idgenre = genres.id GROUP BY movies.id ORDER BY movies.views DESC LIMIT 1"
	moviesSliceData, err := e.s.GetMoviesFromDB(query, "")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	moviesSlice = moviesSliceData[0]

	// Сохраняем данные в Redis
	moviesJSON, err := json.Marshal(moviesSlice)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving contents"})
		return
	}
	err = cache.Rdb.Set(cache.Ctx, "QueryBestMovie", moviesJSON, 5*time.Minute).Err()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving contents"})
		return
	}

	c.JSON(http.StatusOK, moviesSlice)
}
