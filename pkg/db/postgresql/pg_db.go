package postgresql

import (
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)


func OpenPgx(ctx context.Context, conf *pgx.ConnConfig) (*pgx.Conn, error) {
	conn, err := pgx.ConnectConfig(ctx, conf)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func OpenPgxPool(ctx context.Context, conf *pgxpool.Config) (*pgxpool.Pool, error) {
	conn, err := pgxpool.ConnectConfig(ctx, conf)
	if err != nil {
		return nil, err
	}
	return conn, nil
}