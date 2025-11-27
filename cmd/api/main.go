package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/example/block-indexer/internal/api"
	"github.com/example/block-indexer/internal/config"
	"github.com/example/block-indexer/internal/logging"
	"github.com/example/block-indexer/internal/metrics"
	"github.com/example/block-indexer/internal/telemetry"
	"go.uber.org/zap"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg := config.Load()
	logger := logging.New(cfg.Env)
	defer logger.Sync() //nolint:errcheck // best-effort

	metricsSrv := metrics.StartServer(cfg.MetricsAddr, logger)
	defer metricsSrv.Shutdown(ctx) //nolint:errcheck

	tp, shutdownTrace := telemetry.InitProvider(ctx, cfg)
	defer shutdownTrace(context.Background()) //nolint:errcheck
	_ = tp

	router := api.NewServer(cfg, logger)
	srv := &http.Server{
		Addr:         cfg.APIAddr,
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		logger.Info("api server starting", zap.String("addr", cfg.APIAddr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("api server failed", zap.Error(err))
		}
	}()

	waitForSignal(logger)
	_ = srv.Shutdown(ctx)
}

func waitForSignal(logger *zap.Logger) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	s := <-c
	logger.Info("received signal, shutting down", zap.String("signal", s.String()))
}
