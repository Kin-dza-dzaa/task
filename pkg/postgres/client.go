// Package postgres implements connection to postgresql database.
package postgres

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Client struct {
	Pool    *pgxpool.Pool
	Builder sq.StatementBuilderType
}

func NewClient(ctx context.Context, url string, maxPoolSize int, level pgx.LogLevel, l pgx.Logger) (*Client, error) {
	poolConfig, err := pgxpool.ParseConfig(url)
	if err != nil {
		return nil, fmt.Errorf("Client - NewClient - pgxpool.ParseConfig: %w", err)
	}

	poolConfig.MaxConns = int32(maxPoolSize)
	poolConfig.ConnConfig.Logger = l
	poolConfig.ConnConfig.LogLevel = level

	pool, err := pgxpool.ConnectConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("Client - NewClient - pgxpool.ConnectConfig: %w", err)
	}

	PGclient := new(Client)
	PGclient.Builder = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	PGclient.Pool = pool

	return PGclient, nil
}
