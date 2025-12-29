package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/example/block-indexer/pkg/models"
)

type StateDiffRepo struct {
	pool *pgxpool.Pool
}

func NewStateDiffRepo(pool *pgxpool.Pool) *StateDiffRepo {
	return &StateDiffRepo{pool: pool}
}

func (r *StateDiffRepo) Store(ctx context.Context, txHash string, diffs []models.StateDiff) error {
	for _, diff := range diffs {
		if _, err := r.pool.Exec(ctx, `INSERT INTO state_diffs (tx_hash, address, slot, before_value, after_value) VALUES ($1,$2,$3,$4,$5)`,
			txHash, diff.Address, diff.Slot, diff.Before, diff.After); err != nil {
			return err
		}
	}
	return nil
}

func (r *StateDiffRepo) GetByTx(ctx context.Context, hash string) ([]models.StateDiff, error) {
	rows, err := r.pool.Query(ctx, `SELECT address, slot, before_value, after_value FROM state_diffs WHERE tx_hash = $1`, hash)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var diffs []models.StateDiff
	for rows.Next() {
		var diff models.StateDiff
		if err := rows.Scan(&diff.Address, &diff.Slot, &diff.Before, &diff.After); err == nil {
			diffs = append(diffs, diff)
		}
	}
	return diffs, nil
}
