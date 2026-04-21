package repository

import (
	"context"
	"judo/internal/models"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type tournamentRepository struct {
	db     *pgxpool.Pool
	logger *slog.Logger
}

func newTournamentRepository(db *pgxpool.Pool, logger *slog.Logger) *tournamentRepository {
	return &tournamentRepository{
		db:     db,
		logger: logger,
	}
}

func (r *tournamentRepository) SaveTournament(ctx context.Context, tourRow models.TournamentDBRow) {
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
		r.logger.Error("Ошибка сохранения турнира", slog.String("error", err.Error()))
	}
}
