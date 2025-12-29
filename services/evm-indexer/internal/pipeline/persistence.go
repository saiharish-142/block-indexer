package pipeline

import (
	"context"

	"go.uber.org/zap"

	"github.com/example/block-indexer/pkg/models"
	"github.com/example/block-indexer/services/evm-indexer/internal/repository"
)

type Persistence struct {
	repos *repository.Repositories
	log   *zap.Logger
}

func NewPersistence(repos *repository.Repositories, log *zap.Logger) *Persistence {
	return &Persistence{repos: repos, log: log}
}

func (p *Persistence) Run(ctx context.Context, in <-chan models.BlockEnvelope) {
	for env := range in {
		if err := p.repos.InsertEnvelope(ctx, env); err != nil {
			p.log.Warn("failed to persist block", zap.Uint64("number", env.Block.Number), zap.Error(err))
		}
	}
}
