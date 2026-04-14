package ws

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

func GetUpgrader() *websocket.Upgrader {
	return &websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
}

type MessageHandler_cb func(msg []byte)

func WsHandler(hub *Hub, onMessage MessageHandler_cb, w http.ResponseWriter, r *http.Request) {
	conn, err := GetUpgrader().Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade:", err)
		return
	}
	client := &Client{conn: conn}
	hub.Register(client)

	defer func() {
		hub.Unregister(client)
		conn.Close()
	}()

	for {
		_, payload, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(
				err,
				websocket.CloseGoingAway,
				websocket.CloseAbnormalClosure,
				websocket.CloseNormalClosure,
			) {
				log.Printf("error: %v", err)
			}
			break
		}

		onMessage(payload)
	}
}
