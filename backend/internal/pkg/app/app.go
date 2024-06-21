package app

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"kinogo/config"
	"kinogo/internal/app/endpoint/grpcMovies"
	"kinogo/internal/app/endpoint/restAuth"
	"kinogo/internal/app/services/auth"
	movies "kinogo/internal/app/services/movies"
	"kinogo/pkg/cache"
	"kinogo/pkg/db"
	"kinogo/pkg/logger"
	pbMovies "kinogo/pkg/movies_v1"
	"log"
	"net"
	"time"
)

type App struct {
	movies *movies.Service
	auth   *auth.Service

	server *grpc.Server
	router *gin.Engine
}

func New() (*App, error) {
	// инициализируем конфиг, логгер и кэш
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Ошибка при попытке спарсить .env файл в структуру: %v", err)
	}

	logger.Init(cfg.LoggerLevel)

	a := &App{}

	NewGRPC(a)
	NewREST(a, cfg.Auth)

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

func NewGRPC(a *App) {
	a.server = grpc.NewServer()

	// обьявляем сервисы
	a.movies = movies.New()

	// регистрируем эндпоинты
	serviceMovies := &grpcMovies.Endpoint{
		Movies: a.movies,
	}
	pbMovies.RegisterMoviesV1Server(a.server, serviceMovies)
}

func NewREST(a *App, cfgAuth config.Auth) {
	a.router = gin.Default()

	cfgCors := cors.DefaultConfig()
	cfgCors.AllowOrigins = []string{"http://localhost:5173", "http://127.0.0.1:5173", "http://localhost", "http://127.0.0.1"}
	cfgCors.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"}
	cfgCors.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	cfgCors.AllowCredentials = true
	cfgCors.MaxAge = 24 * time.Hour

	a.router.Use(cors.New(cfgCors))

	// регистрируем сервисы
	a.auth = auth.New()

	// регистрируем эндпоинты
	serviceAuth := &restAuth.Endpoint{
		Auth: a.auth,
	}

	// регистрируем маршруты
	a.router.POST("/auth/telegram/callback", serviceAuth.TelegramAuthCallback(cfgAuth.JWTSecret, cfgAuth.BotToken))

	a.router.GET("/stream/:id/:quality/playlist.m3u8", func(c *gin.Context) {
		quality := c.Param("quality")
		id := c.Param("id")
		filename := "media/" + id + "/" + quality + "/playlist.m3u8"

		// Установка заголовков
		c.Header("Content-Type", "application/vnd.apple.mpegurl")
		c.Header("Content-Disposition", "attachment; filename=playlist.m3u8")

		// Отправка файла
		c.File(filename)
	})

	a.router.GET("/stream/:id/:quality/:file", func(c *gin.Context) {
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
}

func (a *App) RunGRPC() error {
	lis, err := net.Listen("tcp", ":8080")
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

func (a *App) RunREST() error {
	err := a.router.Run(":4000")
	if err != nil {
		return err
	}

	return nil
}

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
