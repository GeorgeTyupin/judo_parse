package sql

import (
	"context"
	"judo/internal/models"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DBWriter struct {
	Pool *pgxpool.Pool
}

func NewDBWriter(ctx context.Context, connString string) *DBWriter {
	dbPool, err := pgxpool.New(ctx, connString)
	if err != nil {
		log.Fatalf("Ошибка подключения к бд: %v", err)
	}

	return &DBWriter{
		Pool: dbPool,
	}
}

func (db *DBWriter) SaveTournament(ctx context.Context, tourRow *models.TournamentDBRow) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `
		INSERT INTO tournaments (name, type, place, date, year, month, gender)
		VALUES ($1, $2, $3, $4, $5, $6, $7);`

	_, err := db.Pool.Exec(ctx, query,
		tourRow.Name,
		tourRow.Type,
		tourRow.Place,
		tourRow.Date,
		tourRow.Year,
		tourRow.Month,
		tourRow.Gender,
	)
	if err != nil {
		log.Printf("Ошибка сохранения турнира: %v", err)
	}
}
