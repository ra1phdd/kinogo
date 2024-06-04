package endpoint

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"kinogo/internal/app/config"
	"kinogo/internal/app/models"
	"kinogo/pkg/cache"
	"kinogo/pkg/logger"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"
)

func TestMain(m *testing.M) {
	configuration, err := config.NewConfig("/media/ra1ph/FILES/Программирование/kinogo/.env")
	if err != nil {
		panic(err)
	}

	// Инициализация логгепа
	logger.Init(configuration.LoggerLevel)

	// Инициализация кэша Redis
	cache.Init(configuration.RedisAddr, configuration.RedisPort, configuration.RedisPassword)

	os.Exit(m.Run())
}

func TestGetAllContents(t *testing.T) {
	// Create a new Gin router
	router := gin.New()

	// Create a new Endpoint instance
	endpoint := &Endpoint{
		s: &mockService{},
	}

	// Register the GetAllContents handler
	router.GET("/api/v1/contents", endpoint.GetAllContents)

	// Create a new HTTP request
	req, err := http.NewRequest("GET", "/api/v1/contents", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a new response recorder
	recorder := httptest.NewRecorder()

	// Perform the request
	router.ServeHTTP(recorder, req)

	// Check the response status code
	if recorder.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, recorder.Code)
	}

	// Check the response body
	var responseBody []models.MovieData
	err = json.Unmarshal(recorder.Body.Bytes(), &responseBody)
	if err != nil {
		t.Fatal(err)
	}

	// Check the response data
	expectedMovies := []models.MovieData{
		{Id: 1, Title: "Movie 1", Description: "Description 1", Country: "USA", ReleaseDate: 2021, TimeMovie: 120, ScoreKP: 7.5, ScoreIMDB: 8.0, Poster: "poster1.jpg", TypeMovie: "Action", Views: 1000, Likes: 100, Dislikes: 10, Genres: "Action,Adventure"},
		{Id: 2, Title: "Movie 2", Description: "Description 2", Country: "UK", ReleaseDate: 2022, TimeMovie: 90, ScoreKP: 6.8, ScoreIMDB: 7.2, Poster: "poster2.jpg", TypeMovie: "Comedy", Views: 500, Likes: 50, Dislikes: 5, Genres: "Comedy,Drama"},
	}
	if !reflect.DeepEqual(responseBody, expectedMovies) {
		t.Errorf("Expected movies %v, got %v", expectedMovies, responseBody)
	}

	// Check if the data is cached in Redis
	cachedMovies, err := cache.Rdb.Get(cache.Ctx, "QueryAllContents").Result()
	if err != nil {
		t.Errorf("Expected data to be cached in Redis, but got error: %v", err)
	}

	var cachedMoviesData []models.MovieData
	err = json.Unmarshal([]byte(cachedMovies), &cachedMoviesData)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(cachedMoviesData, expectedMovies) {
		t.Errorf("Expected cached movies %v, got %v", expectedMovies, cachedMoviesData)
	}
}

type mockService struct{}

func (m *mockService) ParseTemplatesMain(w http.ResponseWriter, allData models.AllData) error {
	return nil
}

func (m *mockService) GetMoviesFromDB(query string, params models.QueryParams) ([]models.MovieData, error) {
	// Mock the database response
	movies := []models.MovieData{
		{Id: 1, Title: "Movie 1", Description: "Description 1", Country: "USA", ReleaseDate: 2021, TimeMovie: 120, ScoreKP: 7.5, ScoreIMDB: 8.0, Poster: "poster1.jpg", TypeMovie: "Action", Views: 1000, Likes: 100, Dislikes: 10, Genres: "Action,Adventure"},
		{Id: 2, Title: "Movie 2", Description: "Description 2", Country: "UK", ReleaseDate: 2022, TimeMovie: 90, ScoreKP: 6.8, ScoreIMDB: 7.2, Poster: "poster2.jpg", TypeMovie: "Comedy", Views: 500, Likes: 50, Dislikes: 5, Genres: "Comedy,Drama"},
	}
	return movies, nil
}
