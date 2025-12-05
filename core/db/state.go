package db

import (
    "context"
    "database/sql"
    "errors"
    "time"

    "github.com/jackc/pgx/v5/pgxpool"
)

var ErrNoRows = errors.New("no rows")

// LatestBlockNumber returns the highest EVM block number stored.
func LatestBlockNumber(ctx context.Context, pool *pgxpool.Pool) (uint64, error) {
    ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
    defer cancel()

    var num sql.NullInt64
    if err := pool.QueryRow(ctx, "SELECT MAX(number) FROM blocks").Scan(&num); err != nil {
        return 0, err
    }
    if !num.Valid {
        return 0, ErrNoRows
    }
    return uint64(num.Int64), nil
}

// LatestDagOrder returns the highest DAG block number stored.
func LatestDagOrder(ctx context.Context, pool *pgxpool.Pool) (uint64, error) {
    ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
    defer cancel()

    var num sql.NullInt64
    if err := pool.QueryRow(ctx, "SELECT MAX(number) FROM dag_blocks").Scan(&num); err != nil {
        return 0, err
    }
    if !num.Valid {
        return 0, ErrNoRows
    }
    return uint64(num.Int64), nil
}

// CountBlocks returns the total EVM blocks stored.
func CountBlocks(ctx context.Context, pool *pgxpool.Pool) (uint64, error) {
    ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
    defer cancel()

    var num sql.NullInt64
    if err := pool.QueryRow(ctx, "SELECT COUNT(*) FROM blocks").Scan(&num); err != nil {
        return 0, err
    }
    if !num.Valid {
        return 0, ErrNoRows
    }
    return uint64(num.Int64), nil
}

// CountDagBlocks returns the total DAG blocks stored.
func CountDagBlocks(ctx context.Context, pool *pgxpool.Pool) (uint64, error) {
    ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
    defer cancel()

    var num sql.NullInt64
    if err := pool.QueryRow(ctx, "SELECT COUNT(*) FROM dag_blocks").Scan(&num); err != nil {
        return 0, err
    }
    if !num.Valid {
        return 0, ErrNoRows
    }
    return uint64(num.Int64), nil
}
