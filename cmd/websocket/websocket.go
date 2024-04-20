package websocket

import (
	"kinogo/pkg/logger"
	"net/http"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

var Conn *websocket.Conn

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func Start() {
	ws := http.NewServeMux()
	ws.HandleFunc("/progress", ProgressHandler)
	err := http.ListenAndServe(":8080", ws)
	if err != nil {
		logger.Error("Ошибка при запуске WebSocket-сервера", zap.Error(err))
	}
	logger.Debug("Запуск WebSocket-сервера на порту 8080")
}

func ProgressHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	Conn, err = upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Error("Ошибка при обновлении прогресса", zap.Error(err))
	}
}
