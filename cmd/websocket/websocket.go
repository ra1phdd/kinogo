package websocket

import (
	"fmt"
	"kinogo/pkg/logger"
	"net/http"

	"github.com/gorilla/websocket"
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
		logger.Error("Ошибка при запуске WebSocket-сервера", err)
	}
}

func ProgressHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	Conn, err = upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Ошибка при обновлении:", err)
		return
	}
}
