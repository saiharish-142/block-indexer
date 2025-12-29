package clients

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/example/block-indexer/pkg/models"
)

type EVMIndexerClient struct {
	pool *pgxpool.Pool
}

func (c *EVMIndexerClient) ListBlocks(_ context.Context, limit, _ int) ([]models.Block, error) {
	if limit == 0 {
		limit = 10
	}
	return make([]models.Block, 0, limit), nil
}

func (c *EVMIndexerClient) GetBlock(_ context.Context, id string) (models.Block, error) {
	return models.Block{Hash: id}, nil
}

func (c *EVMIndexerClient) GetTransaction(_ context.Context, hash string) (models.Transaction, error) {
	return models.Transaction{Hash: hash}, nil
}

func (c *EVMIndexerClient) GetAddressTxs(_ context.Context, address string) ([]models.Transaction, error) {
	return []models.Transaction{{From: address}}, nil
}
