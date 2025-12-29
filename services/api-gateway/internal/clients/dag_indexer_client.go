package clients

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/example/block-indexer/pkg/models"
)

type DAGIndexerClient struct {
	pool *pgxpool.Pool
}

func (c *DAGIndexerClient) GetBlock(_ context.Context, hash string) (models.DagBlock, error) {
	return models.DagBlock{Hash: hash}, nil
}

func (c *DAGIndexerClient) GetRange(_ context.Context, from, to int64) ([]models.DagBlock, error) {
	return []models.DagBlock{{Order: uint64(from)}, {Order: uint64(to)}}, nil
}
