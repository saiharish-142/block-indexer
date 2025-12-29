package db

import (
	"context"
	// "database/sql"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"

	"github.com/example/block-indexer/pkg/config"
)

// NewPostgresPool builds a pgx connection pool tuned for high-throughput writes.
func NewPostgresPool(cfg config.DBConfig) (*pgxpool.Pool, error) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Name, cfg.SSLMode)
	pgxCfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("parse pgx config: %w", err)
	}
	if cfg.MaxConns > 0 {
		pgxCfg.MaxConns = cfg.MaxConns
	}
	if cfg.MinConns > 0 {
		pgxCfg.MinConns = cfg.MinConns
	}
	pgxCfg.MaxConnLifetime = time.Hour
	pgxCfg.MaxConnIdleTime = 15 * time.Minute

	pool, err := pgxpool.NewWithConfig(context.Background(), pgxCfg)
	if err != nil {
		return nil, fmt.Errorf("create pool: %w", err)
	}
	return pool, nil
}

// RunMigrations applies goose migrations located at dir.
func RunMigrations(ctx context.Context, cfg config.DBConfig, dir string) error {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Name, cfg.SSLMode)
	pgxCfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return fmt.Errorf("parse pgx config: %w", err)
	}
	sqlDB := stdlib.OpenDB(*pgxCfg.ConnConfig)
	defer sqlDB.Close()
	sqlDB.SetConnMaxLifetime(time.Hour)
	if cfg.MaxConns > 0 {
		sqlDB.SetMaxOpenConns(int(cfg.MaxConns))
	}

	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("set goose dialect: %w", err)
	}
	if err := goose.UpContext(ctx, sqlDB, dir); err != nil {
		return fmt.Errorf("apply migrations: %w", err)
	}
	return nil
}
