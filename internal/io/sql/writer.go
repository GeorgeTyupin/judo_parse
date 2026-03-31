package sql

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewDBWriter(ctx context.Context, connString string) *pgxpool.Pool {
	dbPool, err := pgxpool.New(ctx, connString)
	if err != nil {
		dbPool.Close()
		log.Fatalf("Ошибка подключения к бд: %v", err)
	}

	return dbPool
}
