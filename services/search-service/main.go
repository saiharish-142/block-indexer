package main

import (
	"context"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"

	"github.com/example/block-indexer/pkg/db"
	"github.com/example/block-indexer/pkg/events"
	"github.com/example/block-indexer/pkg/logger"
	"github.com/example/block-indexer/pkg/metrics"
	"github.com/example/block-indexer/services/search-service/internal"
	"github.com/example/block-indexer/services/search-service/internal/api"
	"github.com/example/block-indexer/services/search-service/internal/indexer"
	"github.com/example/block-indexer/services/search-service/internal/repository"
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
	repo := repository.NewSearchRepo(pool)
	natsClient, err := events.NewNATSClient(cfg.NATS.URL, cfg.NATS.SubjectPrefix)
	if err != nil {
		log.Warn("nats unavailable", zap.Error(err))
	}
	consumer := indexer.NewIndexConsumer(log, repo, natsClient)
	consumer.Start(ctx)

	router := internal.NewRouter()
	router.Get("/health", api.Health)
	router.Get("/ready", api.Ready)
	router.Handle("/metrics", metrics.Handler())
	router.Post("/api/v1/search", api.Search(repo))
	router.Get("/api/v1/search", api.Search(repo))

	srv := internal.NewHTTPServer(cfg, router)
	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		_ = srv.Shutdown(shutdownCtx)
	}()

	log.Info("search-service starting", zap.Int("port", cfg.HTTPPort))
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal("server error", zap.Error(err))
	}
}
