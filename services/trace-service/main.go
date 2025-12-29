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
	"github.com/example/block-indexer/pkg/rpc"
	"github.com/example/block-indexer/services/trace-service/internal"
	"github.com/example/block-indexer/services/trace-service/internal/api"
	"github.com/example/block-indexer/services/trace-service/internal/repository"
	"github.com/example/block-indexer/services/trace-service/internal/workers"
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
	traceRepo := repository.NewTraceRepo(pool)
	diffRepo := repository.NewStateDiffRepo(pool)
	client := rpc.NewEVMClient(cfg.RPC.EVMEndpoint)
	worker := workers.NewTraceWorker(cfg, log, traceRepo, diffRepo, client)
	worker.Start(ctx)

	router := internal.NewRouter()
	router.Get("/health", api.Health)
	router.Get("/ready", api.Ready)
	router.Handle("/metrics", metrics.Handler())
	router.Get("/api/v1/tx/{hash}/trace", api.GetTrace(traceRepo))
	router.Get("/api/v1/tx/{hash}/state-diff", api.GetStateDiff(diffRepo))

	srv := internal.NewHTTPServer(cfg, router)
	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		_ = srv.Shutdown(shutdownCtx)
	}()

	log.Info("trace-service starting", zap.Int("port", cfg.HTTPPort))
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal("server error", zap.Error(err))
	}
}
