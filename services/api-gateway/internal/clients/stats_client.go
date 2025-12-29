package clients

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/example/block-indexer/pkg/models"
)

type StatsClient struct {
	pool *pgxpool.Pool
}

func (c *StatsClient) Overview(_ context.Context) (models.StatsOverview, error) {
	return models.StatsOverview{BlockCount: 0, TransactionCount: 0}, nil
}
