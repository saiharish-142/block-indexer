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
	"github.com/example/block-indexer/services/contract-service/internal"
	"github.com/example/block-indexer/services/contract-service/internal/api"
	"github.com/example/block-indexer/services/contract-service/internal/repository"
	"github.com/example/block-indexer/services/contract-service/internal/verifier"
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
	repo := repository.NewContractsRepo(pool)
	compiler := &verifier.SolidityVerifier{}

	router := internal.NewRouter()
	router.Get("/health", api.Health)
	router.Get("/ready", api.Ready)
	router.Handle("/metrics", metrics.Handler())
	router.Get("/api/v1/contracts/{addr}", api.GetContract(repo))
	router.Post("/api/v1/contracts/{addr}/verify", api.VerifyContract(repo, compiler))

	srv := internal.NewHTTPServer(cfg, router)
	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		_ = srv.Shutdown(shutdownCtx)
	}()

	log.Info("contract-service starting", zap.Int("port", cfg.HTTPPort))
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal("server error", zap.Error(err))
	}
}
