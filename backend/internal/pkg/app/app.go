package app

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"kinogo/config"
	"kinogo/internal/app/endpoint/grpcComments"
	"kinogo/internal/app/endpoint/grpcMetrics"
	"kinogo/internal/app/endpoint/grpcMovies"
	"kinogo/internal/app/endpoint/restAuth"
	"kinogo/internal/app/endpoint/restUpload"
	icAuth "kinogo/internal/app/interceptors/auth"
	"kinogo/internal/app/interceptors/uuid"
	auth "kinogo/internal/app/services/auth"
	comments "kinogo/internal/app/services/comments"
	metrics "kinogo/internal/app/services/metrics"
	movies "kinogo/internal/app/services/movies"
	"kinogo/pkg/cache"
	pbComments "kinogo/pkg/comments_v1"
	"kinogo/pkg/db"
	"kinogo/pkg/logger"
	pbMetrics "kinogo/pkg/metrics_v1"
	pbMovies "kinogo/pkg/movies_v1"
	"log"
	"net"
	"time"
)

type App struct {
	movies   *movies.Service
	comments *comments.Service
	auth     *auth.Service
	metrics  *metrics.Service

	server *grpc.Server
	router *gin.Engine
}

func New() (*App, error) {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Ошибка при попытке спарсить .env файл в структуру: %v", err)
	}

	logger.Init(cfg.LoggerLevel)

	a := &App{}

	NewGRPC(a, cfg.ApiKey)
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

func NewGRPC(a *App, ApiKey string) {
	a.server = grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			icUuid.UUIDCheckerInterceptor,
			icAuth.AuthCheckerInterceptor,
		),
	)

	// обьявляем сервисы
	a.movies = movies.New(ApiKey)
	a.comments = comments.New()
	a.auth = auth.New()
	a.metrics = metrics.New()

	// регистрируем эндпоинты
	serviceMovies := &grpcMovies.Endpoint{
		Movies: a.movies,
	}
	serviceComments := &grpcComments.Endpoint{
		Comments: a.comments,
	}
	serviceMetrics := &grpcMetrics.Endpoint{
		Metrics: a.metrics,
	}
	pbMovies.RegisterMoviesV1Server(a.server, serviceMovies)
	pbComments.RegisterCommentsV1Server(a.server, serviceComments)
	pbMetrics.RegisterMetricsV1Server(a.server, serviceMetrics)
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

	serviceUpload := &restUpload.Endpoint{}

	// регистрируем маршруты
	a.router.POST("/auth/telegram/callback", serviceAuth.TelegramAuthCallback(cfgAuth.JWTSecret, cfgAuth.BotToken))

	a.router.POST("/upload", serviceUpload.UploadFileHandler)

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
