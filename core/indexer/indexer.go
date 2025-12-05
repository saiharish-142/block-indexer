package indexer

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/example/block-indexer/core/config"
	"github.com/example/block-indexer/core/db"
	"github.com/example/block-indexer/core/metrics"
	"github.com/example/block-indexer/core/pb"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

// Indexer coordinates chain ingestion, polling, and gRPC publication.
type Indexer struct {
	logger  *zap.Logger
	cfg     config.Config
	stopCh  chan struct{}
	pool    *pgxpool.Pool
	evmNext uint64
	dagNext uint64
}

// New constructs an Indexer.
func New(logger *zap.Logger, cfg config.Config, pool *pgxpool.Pool) *Indexer {
	return &Indexer{
		logger:  logger,
		cfg:     cfg,
		stopCh:  make(chan struct{}),
		pool:    pool,
		evmNext: cfg.EVMStartBlock,
		dagNext: cfg.DagStartOrder,
	}
}

// Run starts the polling loop and listens for new heads via websocket (placeholder).
func (i *Indexer) Run(ctx context.Context) error {
	if err := i.bootstrapState(ctx); err != nil {
		i.logger.Warn("bootstrap from db failed", zap.Error(err))
	}

	ticker := time.NewTicker(i.cfg.PollInterval)
	defer ticker.Stop()

	go i.streamEthHeads(ctx)

	i.logger.Info("indexer started", zap.Duration("poll_interval", i.cfg.PollInterval),
		zap.Uint64("evm_next", i.evmNext), zap.Uint64("dag_next", i.dagNext))

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

func (i *Indexer) bootstrapState(ctx context.Context) error {
	if i.pool == nil {
		return nil
	}

	if num, err := db.LatestBlockNumber(ctx, i.pool); err == nil {
		i.evmNext = num + 1
	} else if !errors.Is(err, db.ErrNoRows) {
		return fmt.Errorf("latest evm block: %w", err)
	}

	if num, err := db.LatestDagOrder(ctx, i.pool); err == nil {
		i.dagNext = num + 1
	} else if !errors.Is(err, db.ErrNoRows) {
		return fmt.Errorf("latest dag block: %w", err)
	}

	return nil
}

func (i *Indexer) processNextBatch(ctx context.Context) error {
	start := time.Now()

	block, err := i.fetchEthBlockByNumber(ctx, i.evmNext)
	if err != nil {
		return fmt.Errorf("fetch eth block: %w", err)
	}

	dagBlock, err := i.fetchDagBlockByOrder(ctx, i.dagNext, true, true, false)
	if err != nil {
		return fmt.Errorf("fetch dag block: %w", err)
	}

	if i.pool != nil {
		if err := db.CopyBlocks(ctx, i.pool, []pb.BlockSummary{*block}); err != nil {
			return fmt.Errorf("copy blocks: %w", err)
		}
		if err := db.CopyDagBlocks(ctx, i.pool, []pb.BlockSummary{*dagBlock}); err != nil {
			return fmt.Errorf("copy dag blocks: %w", err)
		}
	}

	metrics.BlocksProcessed.Add(1)
	i.logger.Info("processed batch",
		zap.Uint64("block", block.Number),
		zap.String("hash", block.Hash),
		zap.String("miner", block.Miner),
		zap.Uint64("dag_order", dagBlock.Number),
		zap.String("dag_hash", dagBlock.Hash),
		zap.Duration("took", time.Since(start)),
	)
	i.evmNext = block.Number + 1
	i.dagNext = dagBlock.Number + 1
	return nil
}
