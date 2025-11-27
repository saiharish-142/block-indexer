package indexer

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/example/block-indexer/core/config"
	"github.com/example/block-indexer/core/metrics"
	"go.uber.org/zap"
)

// Indexer coordinates chain ingestion, polling, and gRPC publication.
type Indexer struct {
	logger *zap.Logger
	cfg    config.Config
	stopCh chan struct{}
	evmNext uint64
	dagNext uint64
}

// New constructs an Indexer.
func New(logger *zap.Logger, cfg config.Config) *Indexer {
	return &Indexer{
		logger: logger,
		cfg:    cfg,
		stopCh: make(chan struct{}),
		evmNext: cfg.EVMStartBlock,
		dagNext: cfg.DagStartOrder,
	}
}

// Run starts the polling loop and listens for new heads via websocket (placeholder).
func (i *Indexer) Run(ctx context.Context) error {
	ticker := time.NewTicker(i.cfg.PollInterval)
	defer ticker.Stop()

	go i.streamEthHeads(ctx)

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
	// TODO: connect to RPC, backfill gaps, apply confirmations.
	start := time.Now()

	block, err := i.fetchEthBlockByNumber(ctx, i.evmNext)
	if err != nil {
		return fmt.Errorf("fetch eth block: %w", err)
	}

	dagBlock, err := i.fetchDagBlockByOrder(ctx, i.dagNext, true, true, false)
	if err != nil {
		return fmt.Errorf("fetch dag block: %w", err)
	}

	metrics.BlocksProcessed.Add(1)
	i.logger.Info("processed batch",
		zap.Uint64("block", block.Number),
		zap.String("hash", block.Hash),
		zap.Uint64("dag_order", dagBlock.Number),
		zap.String("dag_hash", dagBlock.Hash),
		zap.Duration("took", time.Since(start)),
	)
	i.evmNext++
	i.dagNext++
	return nil
}
