package hub

import (
	"sync"

	"go.uber.org/zap"
)

type Message struct {
	Type string `json:"type"`
	Data any    `json:"data"`
}

type Hub struct {
	log       *zap.Logger
	clients   map[*Client]struct{}
	register  chan *Client
	remove    chan *Client
	broadcast chan Message
	mu        sync.RWMutex
}

func New(log *zap.Logger) *Hub {
	return &Hub{
		log:       log,
		clients:   map[*Client]struct{}{},
		register:  make(chan *Client),
		remove:    make(chan *Client),
		broadcast: make(chan Message, 16),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case c := <-h.register:
			h.mu.Lock()
			h.clients[c] = struct{}{}
			h.mu.Unlock()
		case c := <-h.remove:
			h.mu.Lock()
			delete(h.clients, c)
			h.mu.Unlock()
		case msg := <-h.broadcast:
			h.mu.RLock()
			for c := range h.clients {
				c.Send(msg)
			}
			h.mu.RUnlock()
		}
	}
}

func (h *Hub) Broadcast(msg Message) {
	h.broadcast <- msg
}

func (h *Hub) Register(c *Client) {
	h.register <- c
}

func (h *Hub) Remove(c *Client) {
	h.remove <- c
}
