package app

import (
	"github.com/gin-gonic/gin"
	"kinogo/internal/app/endpoint"
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

	//router.Use(middleware.AuthCheck())

	err := router.Run(":4000")
	if err != nil {
		return a, err
	}

	return a, nil
}
