package hub

import (
	"encoding/json"
	"sync"

	"github.com/gorilla/websocket"
)

type Client struct {
	conn *websocket.Conn
	mu   sync.Mutex
}

func NewClient(conn *websocket.Conn) *Client {
	return &Client{conn: conn}
}

func (c *Client) Send(msg Message) {
	c.mu.Lock()
	defer c.mu.Unlock()
	_ = c.conn.WriteJSON(msg)
}

func (c *Client) Read() (Message, error) {
	var msg Message
	err := c.conn.ReadJSON(&msg)
	return msg, err
}

func (c *Client) SendRaw(data any) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.conn.WriteJSON(data)
}

func (c *Client) Close() {
	_ = c.conn.Close()
}

func (c *Client) SendBytes(payload []byte) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.conn.WriteMessage(websocket.TextMessage, payload)
}

func (c *Client) SendJSON(payload any) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.conn.WriteMessage(websocket.TextMessage, mustJSON(payload))
}

func mustJSON(payload any) []byte {
	b, _ := json.Marshal(payload)
	return b
}
