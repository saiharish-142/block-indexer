package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/example/block-indexer/core/config"
	"github.com/example/block-indexer/core/logging"
	"github.com/example/block-indexer/core/metrics"
	"github.com/example/block-indexer/core/telemetry"
	"github.com/example/block-indexer/core/ws"
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

	handler := ws.NewServer(cfg, logger)
	srv := &http.Server{
		Addr:         cfg.WSAddr,
		Handler:      handler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	go func() {
		logger.Info("ws server starting", zap.String("addr", cfg.WSAddr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("ws server failed", zap.Error(err))
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
