package handlers

import (
	"fmt"
	"html/template"
	"kinogo/internal/models"
	"kinogo/internal/services"
	"kinogo/pkg/db"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
)

// Главная страница
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" && r.URL.Path != "/filter" && r.URL.Path != "/search" && r.URL.Path != "/films" && r.URL.Path != "/cartoons" && r.URL.Path != "/telecasts" {
		http.NotFound(w, r)
		return
	}

	var streaming bool
	var movies []models.MovieData
	var bestMovie models.MovieData
	var err error
	if r.URL.Path == "/films" || r.URL.Path == "/cartoons" || r.URL.Path == "/telecasts" {
		if r.URL.Path == "/films" {
			// Запрос к БД/кэшу списка фильмов
			movies, err = services.GetAllFilms()
			if err != nil {
				fmt.Println("Ошибка:", err)
				return
			}
			fmt.Println(movies)
		} else if r.URL.Path == "/cartoons" {
			// Запрос к БД/кэшу списка фильмов
			movies, err = services.GetAllCartoons()
			if err != nil {
				fmt.Println("Ошибка:", err)
				return
			}
			fmt.Println(movies)
		} else if r.URL.Path == "/telecasts" {
			// Запрос к БД/кэшу списка фильмов
			movies, err = services.GetAllTelecasts()
			if err != nil {
				fmt.Println("Ошибка:", err)
				return
			}
			fmt.Println(movies)
		}
		streaming = false
	} else {
		// Запрос к БД/кэшу списка фильмов
		movies, err = services.GetAllMovies()
		if err != nil {
			fmt.Println("Ошибка:", err)
			return
		}
		fmt.Println(movies)

		// Запрос к API Twitch/кэшу статуса стрима
		streaming, err = services.IsStreaming()
		if err != nil {
			fmt.Println("Ошибка:", err)
			return
		}
		fmt.Println(streaming)
	}

	bestMovie, err = services.GetBestMovie()
	if err != nil {
		fmt.Println("Ошибка:", err)
		return
	}
	fmt.Println(bestMovie)

	var allData models.AllData
	allData.GeneralData = models.GeneralData{
		Stream:         streaming,
		IndexHandler:   true,
		SearchAside:    true,
		FilterAside:    true,
		BestMovieAside: true,
	}
	allData.MovieData = append(allData.MovieData, movies...)
	allData.BestMovieData = bestMovie

	ParseTemplates(w, allData)
}

// Страница фильтра
func FilterIndexHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	arrayGenre := r.Form["genre"]
	yearMinStr := r.FormValue("year__min")
	yearMaxStr := r.FormValue("year__max")

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

	var idGenres []int
	if len(nameGenres) > 0 {
		query := `SELECT id FROM genres WHERE name IN (?)`
		query, args, err := sqlx.In(query, nameGenres)
		if err != nil {
			fmt.Println("ошибка при подготовке запроса: ", err)
			return
		}

		query = db.Conn.Rebind(query)
		err = db.Conn.Select(&idGenres, query, args...)
		if err != nil {
			fmt.Println("ошибка при выполнении запроса: ", err)
			return
		}
	}

	movies, err := services.GetFilterMovies(idGenres, yearMin, yearMax)
	if err != nil {
		fmt.Println("Ошибка:", err)
		return
	}

	boolFilter := len(idGenres) > 0 || yearMin != 1980 || yearMax != time.Now().Year()

	allData := models.AllData{
		GeneralData: models.GeneralData{
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

	ParseTemplates(w, allData)
}

func MovieHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	path := r.URL.Path
	parts := strings.Split(path, "/")
	idStr := parts[len(parts)-1]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		fmt.Println("Invalid ID")
		http.NotFound(w, r)
		return
	}

	movie, err := services.GetMovie(id)
	if err != nil {
		fmt.Println("ошибка", err)
		return
	} else if movie.Id == 0 {
		http.NotFound(w, r)
		return
	}

	movie.Views, movie.Likes, movie.Dislikes, err = services.GetStatsToDB(id)
	if err != nil {
		fmt.Println(err)
	}

	var allData models.AllData
	allData.GeneralData = models.GeneralData{
		MovieHandler: true,
		SearchAside:  true,
		FilterAside:  true,
	}
	allData.MovieData = append(allData.MovieData, movie)

	services.HandleView(r, movie.Id)

	ParseTemplates(w, allData)
}

func SearchPageHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	textSearch := strings.ToLower(r.Form.Get("search"))

	movies, err := services.GetSearchMovies(textSearch)
	if err != nil {
		fmt.Println("ошибка:", err)
		return
	}

	var allData models.AllData
	allData.GeneralData = models.GeneralData{
		TextSearch:    textSearch,
		SearchHandler: true,
		SearchAside:   true,
		FilterAside:   true,
	}
	allData.MovieData = append(allData.MovieData, movies...)

	ParseTemplates(w, allData)
}

func AdminHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("admin/templates/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func ParseTemplates(w http.ResponseWriter, allData models.AllData) {
	tmpl, err := template.ParseFiles(
		"web/main/templates/index.html",
		"web/main/templates/twitch.html",
		"web/main/templates/searchaside.html",
		"web/main/templates/filteraside.html",
		"web/main/templates/moviecard.html",
		"web/main/templates/filter.html",
		"web/main/templates/movie.html",
		"web/main/templates/bestmovieaside.html",
	)
	if err != nil {
		fmt.Println("ошибка:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, allData)
	if err != nil {
		fmt.Println("ошибка:", err)
		return
	}
}
