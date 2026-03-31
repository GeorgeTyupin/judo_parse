package repository

import (
	"context"
	"judo/internal/models"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type TournamentRepository struct {
	db *pgxpool.Pool
}

func NewTournamentRepository(db *pgxpool.Pool) *TournamentRepository {
	return &TournamentRepository{
		db: db,
	}
}

func (r *TournamentRepository) SaveTournament(ctx context.Context, tourRow *models.TournamentDBRow) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `
		INSERT INTO tournaments (name, type, place, date, year, month, gender)
		VALUES ($1, $2, $3, $4, $5, $6, $7);`

	_, err := r.db.Exec(ctx, query,
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
