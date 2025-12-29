package api

import (
	"net/http"

	"github.com/gorilla/websocket"

	"github.com/example/block-indexer/services/ws-service/internal/hub"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func WSHandler(h *hub.Hub) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			http.Error(w, "upgrade failed", http.StatusBadRequest)
			return
		}
		client := hub.NewClient(conn)
		h.Register(client)
		go func() {
			defer func() {
				h.Remove(client)
				client.Close()
			}()
			for {
				if _, err := client.Read(); err != nil {
					return
				}
			}
		}()
	}
}
