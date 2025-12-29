package clients

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/example/block-indexer/pkg/models"
)

type ContractClient struct {
	pool *pgxpool.Pool
}

func (c *ContractClient) Get(_ context.Context, addr string) (models.Contract, error) {
	return models.Contract{Address: addr}, nil
}
