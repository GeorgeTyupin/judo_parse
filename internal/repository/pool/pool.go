package dbpool

import (
	"context"
	"net"

	"github.com/jackc/pgx/v5/pgxpool"
)

func New(
	ctx context.Context,
	connString string,
	dialFunc func(ctx context.Context, network, addr string) (net.Conn, error),
) (*pgxpool.Pool, error) {
	pgCfg, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, err
	}
	if dialFunc != nil {
		pgCfg.ConnConfig.DialFunc = dialFunc
	}

	dbPool, err := pgxpool.NewWithConfig(ctx, pgCfg)
	if err != nil {
		return nil, err
	}

	return dbPool, nil
}
