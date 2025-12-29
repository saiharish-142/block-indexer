package pipeline

import (
	"context"
	"sync/atomic"

	"go.uber.org/zap"

	"github.com/example/block-indexer/pkg/models"
)

type Canonical struct {
	log *zap.Logger
}

func NewCanonical(log *zap.Logger) *Canonical {
	return &Canonical{log: log}
}

func (c *Canonical) Run(ctx context.Context, in <-chan models.BlockEnvelope, head *atomic.Uint64) <-chan models.BlockEnvelope {
	out := make(chan models.BlockEnvelope, 8)
	go func() {
		defer close(out)
		for env := range in {
			head.Store(env.Block.Number)
			select {
			case <-ctx.Done():
				return
			case out <- env:
			}
		}
	}()
	return out
}
