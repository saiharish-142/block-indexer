package indexer

import (
	"context"
	"errors"
	"time"

	"github.com/example/block-indexer/core/config"
	"github.com/example/block-indexer/core/metrics"
	"github.com/example/block-indexer/core/pb"
	"go.uber.org/zap"
)

// Indexer coordinates chain ingestion, polling, and gRPC publication.
type Indexer struct {
	logger *zap.Logger
	cfg    config.Config
	stopCh chan struct{}
}

// New constructs an Indexer.
func New(logger *zap.Logger, cfg config.Config) *Indexer {
	return &Indexer{
		logger: logger,
		cfg:    cfg,
		stopCh: make(chan struct{}),
	}
}

// Run starts the polling loop and listens for new heads via websocket (placeholder).
func (i *Indexer) Run(ctx context.Context) error {
	ticker := time.NewTicker(i.cfg.PollInterval)
	defer ticker.Stop()

	i.logger.Info("indexer started", zap.Duration("poll_interval", i.cfg.PollInterval))

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-i.stopCh:
			return errors.New("stopped")
		case <-ticker.C:
			if err := i.processNextBatch(ctx); err != nil {
				i.logger.Error("process batch failed", zap.Error(err))
			}
		}
	}
}

// Stop signals the indexer to exit.
func (i *Indexer) Stop() {
	close(i.stopCh)
}

func (i *Indexer) processNextBatch(ctx context.Context) error {
	// TODO: connect to RPC, fetch latest head, backfill gaps, apply confirmations.
	start := time.Now()
	_ = pb.BlockSummary{
		Number:     0,
		Hash:       "0x0",
		ParentHash: "0x0",
		Timestamp:  time.Now().Unix(),
	}
	metrics.BlocksProcessed.Add(1)
	i.logger.Info("processed batch", zap.Duration("took", time.Since(start)))
	return nil
}
