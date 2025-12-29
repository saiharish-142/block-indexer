package pipeline

import (
	"context"
	"sync/atomic"

	"go.uber.org/zap"

	"github.com/example/block-indexer/services/dag-indexer/internal/repository"
)

type Persistence struct {
	repos *repository.Repositories
	log   *zap.Logger
}

func NewPersistence(repos *repository.Repositories, log *zap.Logger) *Persistence {
	return &Persistence{repos: repos, log: log}
}

func (p *Persistence) Run(ctx context.Context, in <-chan DagBlockPayload, head *atomic.Uint64) {
	for payload := range in {
		if err := p.repos.InsertBlock(ctx, payload.Block, payload.Parents); err != nil {
			p.log.Warn("failed to persist dag block", zap.String("hash", payload.Block.Hash), zap.Error(err))
			continue
		}
		// Update head if order is higher?
		// atomic.Max or similar?
		// for now simple store
		current := head.Load()
		if uint64(payload.Block.Order) > current {
			head.Store(uint64(payload.Block.Order))
		}
	}
}
