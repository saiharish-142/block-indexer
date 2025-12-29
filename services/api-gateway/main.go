package main

import (
	"context"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"

	"github.com/example/block-indexer/pkg/db"
	"github.com/example/block-indexer/pkg/logger"
	"github.com/example/block-indexer/pkg/metrics"
	"github.com/example/block-indexer/services/api-gateway/internal"
	"github.com/example/block-indexer/services/api-gateway/internal/clients"
	"github.com/example/block-indexer/services/api-gateway/internal/handlers"
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
		log.Fatal("db init failed", zap.Error(err))
	}
	deps := clients.NewClientSet(pool)
	router := internal.NewRouter(log, cfg.CORS)
	internal.RegisterRoutes(router, deps)

	router.Get("/health", handlers.Health)
	router.Get("/ready", handlers.Ready)
	router.Handle("/metrics", metrics.Handler())

	srv := internal.NewHTTPServer(cfg, router)
	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		_ = srv.Shutdown(shutdownCtx)
	}()

	log.Info("api-gateway starting", zap.Int("port", cfg.HTTPPort))
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal("server error", zap.Error(err))
	}
}
