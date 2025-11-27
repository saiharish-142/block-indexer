package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/example/block-indexer/core/config"
	"github.com/example/block-indexer/core/indexer"
	"github.com/example/block-indexer/core/logging"
	"github.com/example/block-indexer/core/metrics"
	"github.com/example/block-indexer/core/telemetry"
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

	idx := indexer.New(logger, cfg)

	go func() {
		if err := idx.Run(ctx); err != nil {
			logger.Fatal("indexer failed", zap.Error(err))
		}
	}()

	waitForSignal(logger)
	cancel()
	idx.Stop()
}

func waitForSignal(logger *zap.Logger) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	s := <-c
	logger.Info("received signal, shutting down", zap.String("signal", s.String()))
}
