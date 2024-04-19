package server

import (
	"fmt"
	"kinogo/cmd/metrics"
	"kinogo/internal/handlers"
	"kinogo/internal/services"
	"kinogo/pkg/logger"
	"net/http"
	"strconv"
)

func Start() {
	mux := http.NewServeMux()

	// Добавление видео
	mux.HandleFunc("/resultmovie", services.ResultMovieHandler)
	mux.HandleFunc("/addmovie", services.AddMovieHandler)

	mux.HandleFunc("/like", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		id, err := strconv.Atoi(r.Form.Get("like"))
		if err != nil {
			fmt.Println("Invalid ID")
			return
		}
		services.HandleLike(r, int64(id))
	})
	mux.HandleFunc("/dislike", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		id, err := strconv.Atoi(r.Form.Get("dislike"))
		if err != nil {
			fmt.Println("Invalid ID")
			return
		}
		services.HandleDislike(r, int64(id))
	})

	// Фильтры
	mux.HandleFunc("/filter", handlers.FilterIndexHandler)

	// Поиск
	mux.HandleFunc("/searchpage", services.SearchHandler)
	mux.HandleFunc("/search", handlers.SearchPageHandler)

	// Админка
	fsAdmin := http.FileServer(http.Dir("./web/admin/assets/"))
	mux.Handle("/web/admin/assets/", http.StripPrefix("/web/admin/assets/", fsAdmin))
	mux.HandleFunc("/admin", handlers.AdminHandler)

	// Фильм
	mux.HandleFunc("/id/", handlers.MovieHandler)

	// Общее
	fsAssets := http.FileServer(http.Dir("./web/main/assets/"))
	mux.Handle("/web/main/assets/", http.StripPrefix("/web/main/assets/", fsAssets))
	fsMedia := http.FileServer(http.Dir("media"))
	mux.Handle("/media/", http.StripPrefix("/media/", fsMedia))

	// JSON логи
	fsLogs := http.FileServer(http.Dir("logs"))
	mux.Handle("/logs/", http.StripPrefix("/logs/", fsLogs))

	mux.HandleFunc("/films", handlers.IndexHandler)
	mux.HandleFunc("/series", handlers.IndexHandler)
	mux.HandleFunc("/telecasts", handlers.IndexHandler)
	mux.HandleFunc("/", handlers.IndexHandler)

	metrics.Init(mux)

	err := http.ListenAndServe(":4000", mux)
	if err != nil {
		logger.Error("Ошибка при запуске HTTP-сервера", err)
	}
}
