package api

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/example/block-indexer/core/config"
	"github.com/example/block-indexer/core/db"
	"github.com/example/block-indexer/core/metrics"
	"github.com/example/block-indexer/core/pb"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httprate"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.uber.org/zap"
)

// Server holds dependencies for the API service.
type Server struct {
	cfg    config.Config
	logger *zap.Logger
	pool   *pgxpool.Pool
	router chi.Router
}

type blockFetcher func(context.Context, *pgxpool.Pool, int, *uint64) ([]pb.BlockSummary, error)

// NewServer wires the router with middleware and endpoints.
func NewServer(cfg config.Config, logger *zap.Logger, pool *pgxpool.Pool) http.Handler {
	s := &Server{
		cfg:    cfg,
		logger: logger,
		pool:   pool,
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

	blockLimiter := httprate.LimitByIP(60, time.Minute)

	r.Route("/v1", func(r chi.Router) {
		r.With(blockLimiter).Get("/evm/blocks", s.handleListEVMBlocks)
		r.With(blockLimiter).Get("/dag/blocks", s.handleListDagBlocks)
		r.With(blockLimiter).Get("/blocks", s.handleListDagBlocks)
		r.With(blockLimiter).Get("/blocks/{id}", s.handleGetBlock)
		r.Get("/txs/{hash}", s.handleGetTx)
		r.Get("/addresses/{address}", s.handleGetAddress)
		r.Get("/addresses/{address}/txs", s.handleListAddressTxs)
		r.Get("/stats/blocks", s.handleBlockCounts)
	})

	s.router = r
	return r
}

func (s *Server) handleListEVMBlocks(w http.ResponseWriter, r *http.Request) {
	s.handleListBlocks(w, r, db.ListEVMBlocks, "evm")
}

func (s *Server) handleListDagBlocks(w http.ResponseWriter, r *http.Request) {
	s.handleListBlocks(w, r, db.ListDagBlocks, "dag")
}

func (s *Server) handleListBlocks(w http.ResponseWriter, r *http.Request, fetch blockFetcher, chain string) {
	ctx := r.Context()
	if s.pool == nil {
		http.Error(w, "db not configured", http.StatusServiceUnavailable)
		return
	}

	limit := parseLimit(r.URL.Query().Get("limit"), 50)
	cursorParam := r.URL.Query().Get("cursor")

	var before *uint64
	if cursorParam != "" {
		num, err := strconv.ParseUint(cursorParam, 10, 64)
		if err != nil {
			http.Error(w, "invalid cursor", http.StatusBadRequest)
			return
		}
		before = &num
	}

	blocks, err := fetch(ctx, s.pool, limit, before)
	if err != nil {
		s.logger.Error("list blocks failed", zap.String("chain", chain), zap.Error(err))
		http.Error(w, "failed to fetch blocks", http.StatusInternalServerError)
		return
	}

	nextCursor := ""
	if len(blocks) > limit {
		nextCursor = strconv.FormatUint(blocks[limit-1].Number, 10)
		blocks = blocks[:limit]
	}

	writeJSON(ctx, w, http.StatusOK, map[string]any{
		"cursor": nextCursor,
		"items":  blocks,
	})
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

func (s *Server) handleBlockCounts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if s.pool == nil {
		http.Error(w, "db not configured", http.StatusServiceUnavailable)
		return
	}

	evmCount, err := db.CountBlocks(ctx, s.pool)
	if err != nil {
		s.logger.Error("count evm blocks failed", zap.Error(err))
		http.Error(w, "failed to count blocks", http.StatusInternalServerError)
		return
	}

	dagCount, err := db.CountDagBlocks(ctx, s.pool)
	if err != nil {
		s.logger.Error("count dag blocks failed", zap.Error(err))
		http.Error(w, "failed to count dag blocks", http.StatusInternalServerError)
		return
	}

	writeJSON(ctx, w, http.StatusOK, map[string]any{
		"evm_blocks": evmCount,
		"dag_blocks": dagCount,
	})
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
