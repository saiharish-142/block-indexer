package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/example/block-indexer/pkg/models"
)

type Repositories struct {
	Blocks *DagBlocksRepo
	Edges  *DagEdgesRepo
	Tips   *DagTipsRepo
}

func NewRepositories(pool *pgxpool.Pool) *Repositories {
	return &Repositories{
		Blocks: &DagBlocksRepo{pool: pool},
		Edges:  &DagEdgesRepo{pool: pool},
		Tips:   &DagTipsRepo{pool: pool},
	}
}

func (r *Repositories) InsertBlock(ctx context.Context, block models.DagBlock, parents []string) error {
	// Transactional insert would be better but keeping it simple as per previous pattern
	if err := r.Blocks.Insert(ctx, block); err != nil {
		return err
	}
	for _, parent := range parents {
		if err := r.Edges.Insert(ctx, block.Hash, parent); err != nil {
			return err
		}
	}
	return nil
}
