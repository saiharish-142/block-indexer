package api

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/example/block-indexer/internal/config"
	"github.com/example/block-indexer/internal/metrics"
	"github.com/example/block-indexer/internal/pb"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httprate"
	"github.com/google/uuid"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.uber.org/zap"
)

// Server holds dependencies for the API service.
type Server struct {
	cfg    config.Config
	logger *zap.Logger
	router chi.Router
	// TODO: inject db pools/cache clients here.
}

// NewServer wires the router with middleware and endpoints.
func NewServer(cfg config.Config, logger *zap.Logger) http.Handler {
	s := &Server{
		cfg:    cfg,
		logger: logger,
	}

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(30 * time.Second))
	r.Use(httprate.LimitByIP(200, time.Minute))
	r.Use(middleware.Compress(5))
	r.Use(recordLatency)
	r.Use(otelhttp.NewMiddleware("api"))

	r.Route("/v1", func(r chi.Router) {
		r.Get("/blocks", s.handleListBlocks)
		r.Get("/blocks/{id}", s.handleGetBlock)
		r.Get("/txs/{hash}", s.handleGetTx)
		r.Get("/addresses/{address}", s.handleGetAddress)
		r.Get("/addresses/{address}/txs", s.handleListAddressTxs)
	})

	s.router = r
	return r
}

func (s *Server) handleListBlocks(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	cursor := r.URL.Query().Get("cursor")
	// limit := parseLimit(r.URL.Query().Get("limit"), 50)

	// Placeholder: this would query Postgres or Redis cache.
	resp := struct {
		Cursor string            `json:"cursor"`
		Items  []pb.BlockSummary `json:"items"`
	}{
		Cursor: cursor,
		Items: []pb.BlockSummary{
			{Number: 1, Hash: uuid.NewString(), ParentHash: "0x0", Timestamp: time.Now().Unix()},
		},
	}
	writeJSON(ctx, w, http.StatusOK, resp)
}

func (s *Server) handleGetBlock(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")
	_ = id // TODO: lookup by hash or number.
	block := pb.BlockSummary{Number: 1, Hash: uuid.NewString(), ParentHash: "0x0", Timestamp: time.Now().Unix()}
	writeJSON(ctx, w, http.StatusOK, block)
}

func (s *Server) handleGetTx(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	hash := chi.URLParam(r, "hash")
	tx := pb.TxSummary{Hash: hash, From: "0xabc", To: "0xdef", BlockNumber: 1, Status: "success"}
	writeJSON(ctx, w, http.StatusOK, tx)
}

func (s *Server) handleGetAddress(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	address := chi.URLParam(r, "address")
	resp := map[string]interface{}{
		"address": address,
		"balance": "0",
	}
	writeJSON(ctx, w, http.StatusOK, resp)
}

func (s *Server) handleListAddressTxs(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	address := chi.URLParam(r, "address")
	limit := parseLimit(r.URL.Query().Get("limit"), 50)
	_ = address
	_ = limit
	resp := []pb.TxSummary{
		{Hash: uuid.NewString(), From: address, To: "0xdef", BlockNumber: 1, Status: "success"},
	}
	writeJSON(ctx, w, http.StatusOK, resp)
}

func writeJSON(ctx context.Context, w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	enc := json.NewEncoder(w)
	_ = enc.Encode(v)
}

func parseLimit(raw string, def int) int {
	if raw == "" {
		return def
	}
	val, err := strconv.Atoi(raw)
	if err != nil || val <= 0 {
		return def
	}
	if val > 200 {
		return 200
	}
	return val
}

// recordLatency instruments HTTP handlers with a simple histogram.
func recordLatency(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		metrics.APILatency.Observe(time.Since(start).Seconds())
	})
}
