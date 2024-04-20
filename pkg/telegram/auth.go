package telegram

import (
	"fmt"
	"net/http"

	"kinogo/pkg/logger"

	"go.uber.org/zap"
)

func TelegramCallbackHandler(w http.ResponseWriter, r *http.Request) {
	authData := r.FormValue("auth_data")
	botToken := "kinogolang_bot"

	url := fmt.Sprintf("https://oauth.telegram.org/auth?bot_token=%s&auth_data=%s", botToken, authData)
	resp, err := http.Get(url)
	if err != nil {
		logger.Error("Ошибка получения данных при авторизации пользователя", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()
}
