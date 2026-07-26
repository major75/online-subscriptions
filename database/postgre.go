package database

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/major75/online-subscriptions/pkg/logger"
)

type DBCfg struct {
	Dsn                string
	MaxConnections     int32
	ConnectionLifetime time.Duration
}

type DB struct {
	pool *pgxpool.Pool
}

func NewPostgreDB(ctx context.Context, cfg DBCfg, log logger.Logger) (*DB, error) {
	if cfg.Dsn == "" {
		return nil, fmt.Errorf("database DSN is not configured")
	}

	pgConfig, err := pgxpool.ParseConfig(cfg.Dsn)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", "database.ParseConfig", err)
	}

	pgConfig.MaxConns = cfg.MaxConnections
	pgConfig.MaxConnLifetime = cfg.ConnectionLifetime

	pool, err := pgxpool.NewWithConfig(ctx, pgConfig)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", "database.Connect", err)
	}

	if err = pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("%s: %w", "database.Ping", err)
	}

	log.Info("Successfully connected to database")

	return &DB{
		pool: pool,
	}, nil
}

func (db *DB) Close() {
	db.pool.Close()
}

type Pool interface {
	Exec(ctx context.Context, sql string, arguments ...any) (commandTag pgconn.CommandTag, err error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

func (db *DB) Pool(ctx context.Context) Pool {
	if tx, exist := ctx.Value("tx").(pgx.Tx); exist {
		return tx
	}

	return db.pool
}
