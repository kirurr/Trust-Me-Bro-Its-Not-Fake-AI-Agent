package ws

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type Client struct {
	conn *websocket.Conn
	mu   sync.Mutex
}

func (c *Client) Send(msg []byte) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.conn.WriteMessage(websocket.TextMessage, msg)
}

type Hub struct {
	mu      sync.Mutex
	clients map[*Client]struct{}
}

func NewHub() *Hub {
	return &Hub{clients: make(map[*Client]struct{})}
}
func (h *Hub) Register(c *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.clients[c] = struct{}{}
}

func (h *Hub) Unregister(c *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	delete(h.clients, c)
}

func (h *Hub) Broadcast(msg []byte) {
	h.mu.Lock()
	defer h.mu.Unlock()

	for client := range h.clients {
		if err := client.Send(msg); err != nil {
			log.Printf("send error: %v", err)
		}
	}
}

func GetUpgrader() *websocket.Upgrader {
	return &websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
}

type MessageFromChat struct {
	Text   string `json:"text"`
	UserId string `json:"user_id"`
}

type MessageHandler func(msg []byte)

func WsHandler(hub *Hub, onMessage MessageHandler, w http.ResponseWriter, r *http.Request) {
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
