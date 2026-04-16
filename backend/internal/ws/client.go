package ws

import (
	"log"
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

type Hub interface {
	Register(c *Client)
	Unregister(c *Client)
	Broadcast(msg []byte)
}

type WsHub struct {
	mu      sync.Mutex
	clients map[*Client]struct{}
}

func NewWsHub() Hub {
	return &WsHub{clients: make(map[*Client]struct{})}
}

func (h *WsHub) Register(c *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.clients[c] = struct{}{}
}

func (h *WsHub) Unregister(c *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	delete(h.clients, c)
}

func (h *WsHub) Broadcast(msg []byte) {
	h.mu.Lock()
	defer h.mu.Unlock()

	for client := range h.clients {
		if err := client.Send(msg); err != nil {
			log.Printf("send error: %v", err)
		}
	}
}
