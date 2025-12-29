package main

import (
	"context"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"go.uber.org/zap"

	"github.com/example/block-indexer/pkg/db"
	"github.com/example/block-indexer/pkg/events"
	"github.com/example/block-indexer/pkg/logger"
	"github.com/example/block-indexer/pkg/metrics"
	"github.com/example/block-indexer/pkg/rpc"
	"github.com/example/block-indexer/services/evm-indexer/internal"
	"github.com/example/block-indexer/services/evm-indexer/internal/api"
	"github.com/example/block-indexer/services/evm-indexer/internal/pipeline"
	"github.com/example/block-indexer/services/evm-indexer/internal/repository"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	cfg, err := internal.LoadConfig()
	if err != nil {
		panic(err)
	}
	log, err := logger.NewLogger(cfg.ServiceName, cfg.LogLevel)
	if err != nil {
		panic(err)
	}
	defer log.Sync()

	pool, err := db.NewPostgresPool(cfg.DB)
	if err != nil {
		log.Fatal("failed to init db", zap.Error(err))
	}
	natsClient, err := events.NewNATSClient(cfg.NATS.URL, cfg.NATS.SubjectPrefix)
	if err != nil {
		log.Fatal("failed to init nats", zap.Error(err))
	}
	evmClient := rpc.NewEVMClient(cfg.RPC.EVMEndpoint)

	repos := repository.NewRepositories(pool)
	p := pipeline.NewManager(cfg, log, repos, evmClient, natsClient)
	p.Start(ctx)

	router := newRouter(log, p)
	srv := internal.NewHTTPServer(cfg, router)
	go internal.GracefulShutdown(ctx, srv)

	log.Info("evm-indexer starting", zap.Int("port", cfg.HTTPPort))
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal("server error", zap.Error(err))
	}
	time.Sleep(100 * time.Millisecond)
}

func newRouter(log *zap.Logger, p *pipeline.Manager) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Logger)
	r.Use(cors.AllowAll().Handler)

	r.Get("/health", api.Health)
	r.Get("/ready", api.Ready)
	r.Handle("/metrics", metrics.Handler())

	r.Route("/api/v1/internal", func(r chi.Router) {
		r.Get("/tip", api.TipHandler(p))
	})
	return r
}
