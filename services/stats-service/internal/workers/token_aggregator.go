package workers

import (
	"context"
	"time"

	"go.uber.org/zap"
)

type TokenAggregator struct {
	log *zap.Logger
}

func NewTokenAggregator(log *zap.Logger) *TokenAggregator {
	return &TokenAggregator{log: log}
}

func (t *TokenAggregator) Start(ctx context.Context) {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			t.log.Debug("token stats tick")
		}
	}
}
