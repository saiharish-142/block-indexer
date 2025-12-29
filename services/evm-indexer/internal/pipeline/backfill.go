package pipeline

import (
	"context"
	"time"

	"go.uber.org/zap"

	"github.com/example/block-indexer/services/evm-indexer/internal"
)

type Backfill struct {
	cfg internal.Config
	log *zap.Logger
}

func NewBackfill(cfg internal.Config, log *zap.Logger) *Backfill {
	return &Backfill{cfg: cfg, log: log}
}

func (b *Backfill) Run(ctx context.Context, out chan<- BlockRef) {
	defer close(out)
	if !b.cfg.Backfill.Enabled {
		return
	}
	interval := time.Duration(b.cfg.Backfill.IntervalSeconds) * time.Second
	if interval == 0 {
		interval = 30 * time.Second
	}
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	var counter uint64 = 1
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			out <- BlockRef{Number: counter}
			counter++
		}
	}
}
