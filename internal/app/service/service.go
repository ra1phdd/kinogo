package service

import (
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"html/template"
	"kinogo/internal/app/models"
	"kinogo/pkg/db"
	"kinogo/pkg/logger"
	"net/http"
	"strconv"
	"strings"
)

type Service struct {
}

func New() *Service {
	return &Service{}
}

func (s *Service) ParseTemplatesMain(w http.ResponseWriter, allData models.AllData) error {
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
		return err
	}

	err = tmpl.Execute(w, allData)
	if err != nil {
		logger.Error("Ошибка выполнения шаблонов", zap.Error(err), zap.Any("allData", allData))
		return err
	}
	return nil
}

func (s *Service) GetMoviesFromDB(query string, args ...string) ([]models.MovieData, error) {
	var moviesSlice []models.MovieData

	var rows *sqlx.Rows
	var err error
	if args[0] == "" {
		rows, err = db.Conn.Queryx(query)
	} else {
		id, err := strconv.Atoi(args[0])
		if err != nil {
			rows, err = db.Conn.Queryx(query, args)
		}
	}
	if err != nil {
		fmt.Println(err)
		return []models.MovieData{}, errors.New("error retrieving contents")
	}
	defer rows.Close()

	found := false

	for rows.Next() {
		found = true

		var MoviesData models.MovieData
		err := rows.Scan(
			&MoviesData.Id,
			&MoviesData.Title,
			&MoviesData.Description,
			&MoviesData.Country,
			&MoviesData.ReleaseDate,
			&MoviesData.TimeMovie,
			&MoviesData.ScoreKP,
			&MoviesData.ScoreIMDB,
			&MoviesData.Poster,
			&MoviesData.TypeMovie,
			&MoviesData.Views,
			&MoviesData.Likes,
			&MoviesData.Dislikes,
			&MoviesData.Genres,
		)
		if err != nil {
			return []models.MovieData{}, errors.New("error retrieving cont243ents")
		}

		MoviesData.Genres = strings.Replace(strings.Trim(MoviesData.Genres, "{}"), ",", ", ", -1)

		MoviesData.ScoreKP, err = strconv.ParseFloat(fmt.Sprintf("%.1f", MoviesData.ScoreKP), 64)
		if err != nil {
			return []models.MovieData{}, errors.New("error retrieving conten12ts")
		}

		MoviesData.ScoreIMDB, err = strconv.ParseFloat(fmt.Sprintf("%.1f", MoviesData.ScoreIMDB), 64)
		if err != nil {
			return []models.MovieData{}, errors.New("error retrieving content32s")
		}

		moviesSlice = append(moviesSlice, MoviesData)
	}

	if !found {
		return []models.MovieData{}, errors.New("нет данных в БД")
	}

	return moviesSlice, nil
}

/*func (s *Service) SearchHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		logger.Error("Ошибка парсинга из формы", zap.Error(err))
		return
	}

	textSearch := strings.ToLower(r.Form.Get("search"))
	logger.Debug("Получение текста поиска из формы", zap.String("textSearch", textSearch))

	movies, err := s.GetSearchMovies(textSearch)
	if err != nil {
		logger.Error("Ошибка получения контента по поиску", zap.Error(err), zap.Any("movies", movies))
		return
	}
	logger.Debug("Получен контент из БД/кэша", zap.Any("movies", movies))

	var sb strings.Builder
	for _, movie := range movies {
		sb.WriteString(fmt.Sprintf("<a href='/id/%d'>%s (%d)</a>", movie.Id, movie.Title, movie.ReleaseDate))
	}

	_, _ = fmt.Fprintln(w, sb.String())
}*/

func (s *Service) ValidateInput(input string) error {
	// Проверка на пустую строку
	if strings.TrimSpace(input) == "" {
		return errors.New("входные данные пустые")
	}

	// Проверка на максимальную длину
	const maxLength = 200
	if len(input) > maxLength {
		return errors.New("входные данные слишком длинные (больше 200 символов)")
	}

	// Проверка на специальные символы
	specialChars := []string{"<", ">", "&", "%"}
	for _, char := range specialChars {
		if strings.Contains(input, char) {
			return errors.New("входные данные содержат специальные символы (<, >, &, %)")
		}
	}

	// Все проверки пройдены, возвращаем nil
	return nil
}
