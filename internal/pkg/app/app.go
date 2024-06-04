package app

import (
	"github.com/gin-gonic/gin"
	"kinogo/internal/app/endpoint"
	"kinogo/internal/app/middleware"
	"kinogo/internal/app/service"
)

type App struct {
	e *endpoint.Endpoint
	s *service.Service
}

func New() (*App, error) {
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
}
