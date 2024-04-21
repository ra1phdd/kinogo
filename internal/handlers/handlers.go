package handlers

import (
	"html/template"
	"kinogo/internal/models"
	"kinogo/internal/services"
	"kinogo/pkg/auth"
	"kinogo/pkg/db"
	"kinogo/pkg/logger"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

// Главная страница
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Security-Policy", "frame-ancestors http://127.0.0.1")
	// Вывод 404 на несуществующую страницу
	validPaths := map[string]bool{
		"/":          true,
		"/filter":    true,
		"/search":    true,
		"/films":     true,
		"/cartoons":  true,
		"/telecasts": true,
	}
	if _, ok := validPaths[r.URL.Path]; !ok {
		http.NotFound(w, r)
		return
	}

	verify := auth.VerifyTelegramData(&auth.Auth, auth.BotToken)
	if !verify {
		logger.Warn("Зафиксирована попытка поддельной авторизации")
	} else {
		logger.Info("Пользователь авторизован")
	}

	var streaming bool
	var movies []models.MovieData
	var bestMovie models.MovieData
	var err error
	switch r.URL.Path {
	case "/films":
		movies, err = services.GetAllFilms()
		if err != nil {
			logger.Error("Ошибка при получении всех фильмов из БД/кэша", zap.Error(err), zap.Any("movies", movies))
			return
		}
		streaming = false
	case "/cartoons":
		movies, err = services.GetAllCartoons()
		if err != nil {
			logger.Error("Ошибка при получении всех мультфильмов из БД/кэша", zap.Error(err), zap.Any("movies", movies))
			return
		}
		streaming = false
	case "/telecasts":
		movies, err = services.GetAllTelecasts()
		if err != nil {
			logger.Error("Ошибка при получении всех передач из БД/кэша", zap.Error(err), zap.Any("movies", movies))
			return
		}
		streaming = false
	default:
		movies, err = services.GetAllMovies()
		if err != nil {
			logger.Error("Ошибка при получении всего контента из БД/кэша", zap.Error(err), zap.Any("movies", movies))
			return
		}

		streaming, err = services.IsStreaming()
		if err != nil {
			logger.Error("Ошибка при получении информации о наличии стрима из БД/кэша", zap.Error(err), zap.Any("streaming", streaming))
			return
		}
	}
	logger.Debug("Получен контент из БД/кэша", zap.Any("movies", movies))

	var allData models.AllData

	bestMovie, err = services.GetBestMovie()
	if err != nil {
		logger.Warn("Ошибка при получении популярного фильма из БД/кэша", zap.Error(err), zap.Any("movies", movies))
		allData.GeneralData = models.GeneralData{
			BestMovieAside: false,
		}
	} else {
		allData.GeneralData = models.GeneralData{
			BestMovieAside: true,
		}
		allData.BestMovieData = bestMovie
	}

	allData.GeneralData = models.GeneralData{
		Stream:       streaming,
		Auth:         verify,
		IndexHandler: true,
		SearchAside:  true,
		FilterAside:  true,
	}
	allData.MovieData = append(allData.MovieData, movies...)
	allData.UserData = auth.Auth

	ParseTemplatesMain(w, allData)
}

// Страница фильтра
func FilterIndexHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	arrayGenre := r.Form["genre"]
	yearMinStr := r.FormValue("year__min")
	yearMaxStr := r.FormValue("year__max")

	logger.Debug("Получение данных фильтра из формы", zap.Any("Жанры", arrayGenre), zap.String("Минимальный год", yearMinStr), zap.String("Максимальный год", yearMaxStr))

	yearMin, errMin := strconv.Atoi(yearMinStr)
	if errMin != nil {
		yearMin = 1980
	}
	yearMax, errMax := strconv.Atoi(yearMaxStr)
	if errMax != nil {
		yearMax = time.Now().Year()
	}

	nameGenres := make([]string, 0, len(arrayGenre))
	for _, genre := range arrayGenre {
		if genre == "выбрать" {
			break
		}
		nameGenres = append(nameGenres, genre)
	}

	verify := auth.VerifyTelegramData(&auth.Auth, auth.BotToken)
	if !verify {
		logger.Warn("Зафиксирована попытка поддельной авторизации")
	} else {
		logger.Info("Пользователь авторизован")
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

	movies, err := services.GetFilterMovies(idGenres, yearMin, yearMax)
	if err != nil {
		logger.Error("Ошибка при получении контента по фильтрам", zap.Error(err), zap.Any("movies", movies))
		return
	}
	logger.Debug("Получен контент из БД/кэша", zap.Any("movies", movies))

	boolFilter := len(idGenres) > 0 || yearMin != 1980 || yearMax != time.Now().Year()
	logger.Debug("Фильтр", zap.Bool("boolFilter", boolFilter))

	allData := models.AllData{
		GeneralData: models.GeneralData{
			Auth:          verify,
			FilterHandler: true,
			SearchAside:   true,
			FilterAside:   true,
		},
		FilterData: models.FilterData{
			BoolFilter: boolFilter,
			Genre:      nameGenres,
			YearMin:    yearMin,
			YearMax:    yearMax,
		},
	}
	allData.MovieData = append(allData.MovieData, movies...)
	allData.UserData = auth.Auth

	ParseTemplatesMain(w, allData)
}

func MovieHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	path := r.URL.Path
	parts := strings.Split(path, "/")
	idStr := parts[len(parts)-1]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	verify := auth.VerifyTelegramData(&auth.Auth, auth.BotToken)
	if !verify {
		logger.Warn("Зафиксирована попытка поддельной авторизации")
	} else {
		logger.Info("Пользователь авторизован")
	}

	movie, err := services.GetMovie(id)
	if err != nil {
		logger.Error("Ошибка получения контента по id", zap.Error(err), zap.Any("movie", movie), zap.Int("id", id))
		return
	} else if movie.Id == 0 {
		http.NotFound(w, r)
		return
	}
	logger.Debug("Получен контент из БД/кэша", zap.Any("movie", movie))

	movie.Views, movie.Likes, movie.Dislikes, err = services.GetStatsToDB(id)
	if err != nil {
		logger.Error("Ошибка получении статистики фильма по id", zap.Error(err), zap.Int("id", id), zap.Int64("Просмотры", movie.Views), zap.Int64("Лайки", movie.Likes), zap.Int64("Дизлайки", movie.Dislikes))
	}
	logger.Debug("Статистика фильма", zap.Int("id", id), zap.Int64("Просмотры", movie.Views), zap.Int64("Лайки", movie.Likes), zap.Int64("Дизлайки", movie.Dislikes))

	comments, err := services.GetCommentsFromDB(id)
	if err != nil {
		logger.Error("Ошибка получения комментариев по id", zap.Error(err), zap.Any("comments", comments), zap.Int("id", id))
		return
	}
	logger.Debug("Получены комментарии из БД/кэш", zap.Any("comments", comments), zap.Int("id", id))
	allComm := services.BuildCommentTree(comments)

	var allData models.AllData
	allData.GeneralData = models.GeneralData{
		Auth:         verify,
		MovieHandler: true,
		SearchAside:  true,
		FilterAside:  true,
	}
	allData.MovieData = append(allData.MovieData, movie)
	allData.UserData = auth.Auth
	allData.CommentsData = allComm

	services.HandleView(r, movie.Id)

	ParseTemplatesMain(w, allData)
}

func SearchPageHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	textSearch := strings.ToLower(r.Form.Get("search"))
	logger.Debug("Получение текста поиска из формы", zap.String("textSearch", textSearch))

	verify := auth.VerifyTelegramData(&auth.Auth, auth.BotToken)
	if !verify {
		logger.Warn("Зафиксирована попытка поддельной авторизации")
	} else {
		logger.Info("Пользователь авторизован")
	}

	movies, err := services.GetSearchMovies(textSearch)
	if err != nil {
		logger.Error("Ошибка получения контента по поиску", zap.Error(err), zap.Any("movies", movies))
		return
	}
	logger.Debug("Получен контент из БД/кэша", zap.Any("movies", movies))

	var allData models.AllData
	allData.GeneralData = models.GeneralData{
		Auth:          verify,
		TextSearch:    textSearch,
		SearchHandler: true,
		SearchAside:   true,
		FilterAside:   true,
	}
	allData.MovieData = append(allData.MovieData, movies...)
	allData.UserData = auth.Auth

	ParseTemplatesMain(w, allData)
}

// Легаси
func AdminHandler(w http.ResponseWriter, r *http.Request) {
	ParseTemplatesAdmin(w, models.AllData{})
}

func ParseTemplatesMain(w http.ResponseWriter, allData models.AllData) {
	tmpl, err := template.ParseFiles(
		"web/main/templates/index.html",
		"web/main/templates/twitch.html",
		"web/main/templates/searchaside.html",
		"web/main/templates/filteraside.html",
		"web/main/templates/moviecard.html",
		"web/main/templates/filter.html",
		"web/main/templates/movie.html",
		"web/main/templates/bestmovieaside.html",
		"web/main/templates/comments.html",
	)
	if err != nil {
		logger.Error("Ошибка парсинга шаблонов", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, allData)
	if err != nil {
		logger.Error("Ошибка выполнения шаблонов", zap.Error(err), zap.Any("allData", allData))
		return
	}
}

func ParseTemplatesAdmin(w http.ResponseWriter, allData models.AllData) {
	tmpl, err := template.ParseFiles(
		"web/admin/templates/index.html",
	)
	if err != nil {
		logger.Error("Ошибка парсинга шаблонов", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, allData)
	if err != nil {
		logger.Error("Ошибка выполнения шаблонов", zap.Error(err), zap.Any("allData", allData))
		return
	}
}
