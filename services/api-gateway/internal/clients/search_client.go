package clients

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/example/block-indexer/pkg/models"
)

type SearchClient struct {
	pool *pgxpool.Pool
}

func (c *SearchClient) Search(_ context.Context, query string) ([]models.SearchResult, error) {
	return []models.SearchResult{{Type: "placeholder", Key: query}}, nil
}
