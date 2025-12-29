package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/example/block-indexer/pkg/models"
)

type StatsRepo struct {
	pool *pgxpool.Pool
}

func NewStatsRepo(pool *pgxpool.Pool) *StatsRepo {
	return &StatsRepo{pool: pool}
}

func (r *StatsRepo) Overview(ctx context.Context) (models.StatsOverview, error) {
	var overview models.StatsOverview
	err := r.pool.QueryRow(ctx, `SELECT COALESCE(SUM(block_count),0), COALESCE(SUM(tx_count),0), COALESCE(AVG(avg_gas_price),0) FROM stats_daily`).
		Scan(&overview.BlockCount, &overview.TransactionCount, &overview.AvgGasPrice)
	if err != nil {
		overview = models.StatsOverview{BlockCount: 0, TransactionCount: 0}
	}
	return overview, nil
}

func (r *StatsRepo) Historical(ctx context.Context, metric string, from, to time.Time) ([]models.HistoricalPoint, error) {
	rows, err := r.pool.Query(ctx, `SELECT day, tx_count FROM stats_daily ORDER BY day DESC LIMIT 30`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	points := []models.HistoricalPoint{}
	for rows.Next() {
		var day time.Time
		var value string
		var txCount int64
		if err := rows.Scan(&day, &txCount); err != nil {
			continue
		}
		value = fmt.Sprintf("\"%d\"", txCount)
		points = append(points, models.HistoricalPoint{Timestamp: day, Value: value, Metric: metric})
	}
	return points, nil
}
