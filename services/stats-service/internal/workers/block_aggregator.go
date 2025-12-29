package workers

import (
	"context"
	"time"

	"go.uber.org/zap"

	"github.com/example/block-indexer/services/stats-service/internal"
	"github.com/example/block-indexer/services/stats-service/internal/repository"
)

type Aggregator struct {
	cfg    internal.Config
	log    *zap.Logger
	repo   *repository.StatsRepo
	tokens *TokenAggregator
}

func NewAggregator(cfg internal.Config, log *zap.Logger, repo *repository.StatsRepo) *Aggregator {
	return &Aggregator{
		cfg:    cfg,
		log:    log,
		repo:   repo,
		tokens: NewTokenAggregator(log),
	}
}

func (a *Aggregator) Start(ctx context.Context) {
	go a.tokens.Start(ctx)
	interval := time.Duration(a.cfg.Aggregation.IntervalSeconds) * time.Second
	if interval == 0 {
		interval = time.Minute
	}
	ticker := time.NewTicker(interval)
	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				a.refresh(ctx)
			}
		}
	}()
}

func (a *Aggregator) refresh(ctx context.Context) {
	_, _ = a.repo.Overview(ctx)
	a.log.Debug("refreshed stats snapshot")
}
