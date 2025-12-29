package clients

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/example/block-indexer/pkg/models"
)

type TraceClient struct {
	pool *pgxpool.Pool
}

func (c *TraceClient) Trace(_ context.Context, hash string) (models.TxTrace, error) {
	return models.TxTrace{TxHash: hash}, nil
}

func (c *TraceClient) StateDiff(_ context.Context, hash string) ([]models.StateDiff, error) {
	return []models.StateDiff{}, nil
}
