package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"kinogo/internal/app/models"
	"net/http"
	"strconv"
	"strings"

	"kinogo/pkg/logger"

	"go.uber.org/zap"
)

var BotToken = "6438366128:AAHskqGtvfsn98DonkLy-Hdjhcne6VzaFXM"
var Auth models.User

func TelegramCallbackHandler(w http.ResponseWriter, r *http.Request) {
	Auth = ParseParams(r)

	logger.Debug("Данные о Telegram-пользователе", zap.Int64("id", Auth.ID), zap.String("first_name", Auth.FirstName), zap.String("last_name", Auth.LastName), zap.String("username", Auth.Username), zap.String("photo_url", Auth.PhotoURL), zap.Int64("auth_date", Auth.AuthDate), zap.String("hash", Auth.Hash))

	verify := VerifyTelegramData(&Auth, BotToken)
	if !verify {
		logger.Warn("Зафиксирована попытка поддельной авторизации")
		w.WriteHeader(http.StatusUnauthorized)
		return
	} else {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func TelegramLogoutHandler(w http.ResponseWriter, r *http.Request) {
	var logout = ParseParams(r)

	logger.Debug("Данные о Telegram-пользователе", zap.Int64("id", logout.ID), zap.String("first_name", logout.FirstName), zap.String("last_name", logout.LastName), zap.String("username", logout.Username), zap.String("photo_url", logout.PhotoURL), zap.Int64("auth_date", logout.AuthDate), zap.String("hash", logout.Hash))

	verify := VerifyTelegramData(&logout, BotToken)
	if !verify {
		logger.Warn("Зафиксирована попытка поддельной деаутентификации")
		http.Redirect(w, r, "/", http.StatusSeeOther)
	} else {
		Auth = models.User{}
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func ParseParams(r *http.Request) models.User {
	// Получаем параметры из GET-запроса
	params := r.URL.Query()

	// Извлекаем значения параметров
	idStr := params.Get("id")
	logger.Debug("Получение id из формы", zap.String("idStr", idStr))
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		logger.Error("Ошибка преобразования id из string в int64", zap.Error(err))
		return models.User{}
	}

	authDateStr := params.Get("auth_date")
	logger.Debug("Получение AuthDate из формы", zap.String("AuthDateStr", authDateStr))
	authDate, err := strconv.ParseInt(authDateStr, 10, 64)
	if err != nil {
		logger.Error("Ошибка преобразования AuthDate из string в int64", zap.Error(err))
		return models.User{}
	}

	var AuthResp = models.User{
		ID:        id,
		FirstName: params.Get("first_name"),
		LastName:  params.Get("last_name"),
		Username:  params.Get("username"),
		PhotoURL:  params.Get("photo_url"),
		AuthDate:  authDate,
		Hash:      params.Get("hash"),
	}

	return AuthResp
}

func VerifyTelegramData(auth *models.User, botToken string) bool {
	// Создаем строку для проверки
	dataCheckString := buildDataCheckString(auth)
	logger.Debug("Строка для проверки", zap.String("dataCheckString", dataCheckString))

	// Вычисляем секретный ключ
	secretKey := sha256.Sum256([]byte(botToken))
	logger.Debug("secret key", zap.Any("secretKey", secretKey))

	// Вычисляем HMAC-SHA256 подпись
	mac := hmac.New(sha256.New, secretKey[:])
	mac.Write([]byte(dataCheckString))
	expectedHash := hex.EncodeToString(mac.Sum(nil))
	logger.Debug("HMAC-SHA256 подпись", zap.String("expectedHash", expectedHash))

	// Сравниваем полученный хэш с ожидаемым
	return auth.Hash == expectedHash
}

func buildDataCheckString(auth *models.User) string {
	var fields []string
	fields = append(fields, fmt.Sprintf("auth_date=%d", auth.AuthDate))
	fields = append(fields, fmt.Sprintf("first_name=%s", auth.FirstName))
	fields = append(fields, fmt.Sprintf("id=%d", auth.ID))
	if auth.LastName != "" {
		fields = append(fields, fmt.Sprintf("last_name=%s", auth.LastName))
	}
	if auth.PhotoURL != "" {
		fields = append(fields, fmt.Sprintf("photo_url=%s", auth.PhotoURL))
	}
	fields = append(fields, fmt.Sprintf("username=%s", auth.Username))

	return strings.Join(fields, "\n")
}
