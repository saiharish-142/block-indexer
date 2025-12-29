package main

import (
	"context"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"

	"github.com/example/block-indexer/pkg/events"
	"github.com/example/block-indexer/pkg/logger"
	"github.com/example/block-indexer/pkg/metrics"
	"github.com/example/block-indexer/services/ws-service/internal"
	"github.com/example/block-indexer/services/ws-service/internal/api"
	"github.com/example/block-indexer/services/ws-service/internal/broadcaster"
	"github.com/example/block-indexer/services/ws-service/internal/hub"
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

	natsClient, err := events.NewNATSClient(cfg.NATS.URL, cfg.NATS.SubjectPrefix)
	if err != nil {
		log.Warn("nats unavailable", zap.Error(err))
	}
	h := hub.New(log)
	b := broadcaster.New(log, h, natsClient)
	go h.Run()
	b.Start(ctx)

	router := internal.NewRouter()
	router.Get("/health", api.Health)
	router.Get("/ready", api.Ready)
	router.Handle("/metrics", metrics.Handler())
	router.Get("/ws", api.WSHandler(h))

	srv := internal.NewHTTPServer(cfg, router)
	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		_ = srv.Shutdown(shutdownCtx)
	}()

	log.Info("ws-service starting", zap.Int("port", cfg.HTTPPort))
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal("server error", zap.Error(err))
	}
}
