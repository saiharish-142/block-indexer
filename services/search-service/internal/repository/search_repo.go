package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/example/block-indexer/pkg/models"
)

type SearchRepo struct {
	pool *pgxpool.Pool
}

func NewSearchRepo(pool *pgxpool.Pool) *SearchRepo {
	return &SearchRepo{pool: pool}
}

func (r *SearchRepo) Search(ctx context.Context, query string) ([]models.SearchResult, error) {
	rows, err := r.pool.Query(ctx, `SELECT type, key, display_name, metadata FROM search_index WHERE key ILIKE '%' || $1 || '%' LIMIT 20`, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	results := []models.SearchResult{}
	for rows.Next() {
		var res models.SearchResult
		if err := rows.Scan(&res.Type, &res.Key, &res.DisplayName, &res.Metadata); err == nil {
			results = append(results, res)
		}
	}
	return results, nil
}
