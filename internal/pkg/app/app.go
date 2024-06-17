package app

import (
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"kinogo/config"
	"kinogo/internal/app/endpoint/grpcMovies"
	"kinogo/internal/app/services/movies"
	"kinogo/pkg/cache"
	"kinogo/pkg/db"
	"kinogo/pkg/logger"
	pbMovies "kinogo/pkg/movies_v1"
	"log"
	"net"
)

type App struct {
	movies *movies.Service

	server *grpc.Server
}

func New() (*App, error) {
	// инициализируем конфиг, логгер и кэш
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Ошибка при попытке спарсить .env файл в структуру: %v", err)
	}

	logger.Init(cfg.LoggerLevel)

	a := &App{}

	a.server = grpc.NewServer()

	// обьявляем сервисы
	a.movies = movies.New()

	// регистрируем эндпоинты
	serviceMovies := &grpcMovies.Endpoint{
		Movies: a.movies,
	}
	pbMovies.RegisterMoviesV1Server(a.server, serviceMovies)

	err = cache.Init(cfg.Redis.RedisAddr+":"+cfg.Redis.RedisPort, cfg.Redis.RedisUsername, cfg.Redis.RedisPassword, cfg.Redis.RedisDBId)
	if err != nil {
		logger.Error("ошибка при инициализации кэша: ", zap.Error(err))
		return nil, err
	}

	err = db.Init(cfg.DB.DBUser, cfg.DB.DBPassword, cfg.DB.DBHost, cfg.DB.DBName)
	if err != nil {
		logger.Fatal("ошибка при инициализации БД: ", zap.Error(err))
		return nil, err
	}

	return a, nil
}

func (a *App) Run() error {
	lis, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		logger.Fatal("Ошибка при открытии listener: ", zap.Error(err))
	}

	err = a.server.Serve(lis)
	if err != nil {
		logger.Fatal("Ошибка при инициализации сервера: ", zap.Error(err))
		return err
	}

	return nil
}

func (a *App) Stop() {
	logger.Info("закрытие gRPC сервера")

	a.server.GracefulStop()
}

/*func New() (*App, error) {
	a := &App{}

	router := gin.Default()
	v1 := router.Group("/api/v1")

	v1.Use(middleware.CORSPolicy())

	a.s = service.New()
	a.e = endpoint.New(a.s)

	v1.GET("/contents", a.e.GetAllContents)
	v1.GET("/movies", a.e.GetAllMovies)
	v1.GET("/cartoons", a.e.GetAllCartoons)
	v1.GET("/telecasts", a.e.GetAllTelecasts)

	// Эндпоинт для получения информации о самом просматриваемом контенте на сайте
	v1.GET("/content/best", a.e.GetBestMovie)

	// Эндпоинт для получения контента по поиску
	v1.GET("/content/search", a.e.SearchMovies)

	// Эндпоинт для получения контента с фильтрами
	v1.POST("/content/filter", a.e.FilterMovies)

	// Эндпоинт для получения информации о конкретном фильме по его идентификатору
	v1.GET("/movies/:id", a.e.GetMovieByID)

	v1.GET("/stream/:id", func(c *gin.Context) {
		id := c.Param("id")
		filename := "media/" + id + "/stream.m3u8"

		// Установка заголовков
		c.Header("Content-Type", "application/vnd.apple.mpegurl")
		c.Header("Content-Disposition", "attachment; filename=stream.m3u8")

		// Отправка файла
		c.File(filename)
	})

	v1.HEAD("/stream/:id", func(c *gin.Context) {
		c.Header("Content-Type", "application/vnd.apple.mpegurl")
		c.Header("Content-Disposition", "attachment; filename=stream.m3u8")

		c.Status(200) // Возвращаем только статус-код 200 OK без тела ответа
	})

	v1.GET("/stream/:id/:quality/stream.m3u8", func(c *gin.Context) {
		quality := c.Param("quality")
		id := c.Param("id")
		filename := "media/" + id + "/stream_" + quality + ".m3u8"

		// Установка заголовков
		c.Header("Content-Type", "application/vnd.apple.mpegurl")
		c.Header("Content-Disposition", "attachment; filename=stream.m3u8")

		// Отправка файла
		c.File(filename)
	})

	v1.GET("/stream/:id/:quality/:file", func(c *gin.Context) {
		file := c.Param("file")
		quality := c.Param("quality")
		id := c.Param("id")
		filename := "media/" + id + "/" + quality + "/" + file

		// Установка заголовков
		c.Header("Content-Type", "application/vnd.apple.mpegurl")
		c.Header("Content-Disposition", "attachment; filename="+file)

		// Отправка файла
		c.File(filename)
	})

	//router.Use(middleware.AuthCheck())

	err := router.Run(":4000")
	if err != nil {
		return a, err
	}

	return a, nil
}*/

/*mux := http.NewServeMux()

// Добавление видео
mux.HandleFunc("/resultmovie", services.ResultMovieHandler)
mux.HandleFunc("/addmovie", services.AddMovieHandler)

mux.HandleFunc("/like", func(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	id, err := strconv.Atoi(r.Form.Get("like"))
	if err != nil {
		logger.Error("Ошибка парсинга ID фильма для постановки лайка")
	}
	logger.Debug("Постановка лайка", zap.Int("id", id))
	services.HandleLike(r, int64(id))
})
mux.HandleFunc("/dislike", func(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	id, err := strconv.Atoi(r.Form.Get("dislike"))
	if err != nil {
		logger.Error("Ошибка парсинга ID фильма для постановки дизлайка")
	}
	logger.Debug("Постановка дизлайка", zap.Int("id", id))
	services.HandleDislike(r, int64(id))
})*/
