package ws

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/coder/websocket"
	"github.com/example/block-indexer/core/config"
	"github.com/example/block-indexer/core/metrics"
	"github.com/example/block-indexer/core/pb"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.uber.org/zap"
)

// Server exposes websocket endpoints for real-time updates.
type Server struct {
	cfg    config.Config
	logger *zap.Logger
}

// NewServer returns an http.Handler with websocket routes.
func NewServer(cfg config.Config, logger *zap.Logger) http.Handler {
	s := &Server{cfg: cfg, logger: logger}
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(otelhttp.NewMiddleware("ws"))

	r.Get("/ws/heads", s.handleHeads)

	return r
}

func (s *Server) handleHeads(w http.ResponseWriter, r *http.Request) {
	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{InsecureSkipVerify: true})
	if err != nil {
		s.logger.Error("accept ws", zap.Error(err))
		return
	}
	defer conn.Close(websocket.StatusNormalClosure, "bye")
	metrics.WSConnections.Inc()
	defer metrics.WSConnections.Dec()

	ctx := r.Context()
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			msg := pb.BlockSummary{
				Number:     uint64(time.Now().Unix()),
				Hash:       "0xhead",
				ParentHash: "0xparent",
				Timestamp:  time.Now().Unix(),
			}
			if err := writeJSON(ctx, conn, msg); err != nil {
				s.logger.Error("write ws", zap.Error(err))
				return
			}
		}
	}
}

func writeJSON(ctx context.Context, c *websocket.Conn, v interface{}) error {
	writer, err := c.Writer(ctx, websocket.MessageText)
	if err != nil {
		return err
	}
	defer writer.Close()
	return json.NewEncoder(writer).Encode(v)
}
