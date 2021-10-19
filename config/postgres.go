package config

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"goland-hello/pkg/db/postgresql"
)

func OpenPostgresConfig(ctx context.Context, cfg *Config) (*pgx.Conn, error) {
	pgConf, err := pgx.ParseConfig(fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.Postgres.User,
		cfg.Postgres.Password,
		cfg.Postgres.Host,
		cfg.Postgres.Port,
		cfg.Postgres.Database,
		cfg.Postgres.SSLMode))
	if err != nil {
		return nil, err
	}
	return postgresql.OpenPgx(ctx, pgConf)
}

func OpenPostgresPoolConfig(ctx context.Context, cfg *Config) (*pgxpool.Pool, error) {
	pgConf, err := pgxpool.ParseConfig(fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.Postgres.User,
		cfg.Postgres.Password,
		cfg.Postgres.Host,
		cfg.Postgres.Port,
		cfg.Postgres.Database,
		cfg.Postgres.SSLMode))
	if err != nil {
		return nil, err
	}
	return postgresql.OpenPgxPool(ctx, pgConf)
}
