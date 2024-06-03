package server

import (
	"go.uber.org/zap"
	"kinogo/pkg/logger"
	"net/http"
)

func Start() {
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
	})

	// Фильтры
	mux.HandleFunc("/filter", handlers.FilterIndexHandler)

	// Авторизация в TG
	mux.HandleFunc("/auth/telegram/callback", auth.TelegramCallbackHandler)
	mux.HandleFunc("/auth/telegram/logout", auth.TelegramLogoutHandler)

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

	metrics.Init(mux)*/

	err := http.ListenAndServe(":4000", mux)
	if err != nil {
		logger.Fatal("Ошибка при запуске HTTP-сервера", zap.Error(err))
	}
	logger.Debug("HTTP-сервер запущен")
}
