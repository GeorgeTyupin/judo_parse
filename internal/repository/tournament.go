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

func (r *tournamentRepository) SaveAllTournaments(ctx context.Context, rows []models.TournamentDBRow) {
	if err := truncateTable(ctx, r.db, "tournaments"); err != nil {
		r.logger.Error("Ошибка очистки таблицы tournaments", slog.String("error", err.Error()))
		return
	}

	query := `
		INSERT INTO tournaments (name, type, place, date, year, month, gender)
		VALUES ($1, $2, $3, $4, $5, $6, $7);`

	for _, row := range rows {
		insCtx, insCancel := context.WithTimeout(ctx, 5*time.Second)

		_, err := r.db.Exec(insCtx, query,
			row.Name,
			row.Type,
			row.Place,
			row.Date,
			row.Year,
			row.Month,
			row.Gender,
		)
		insCancel()

		if err != nil {
			r.logger.Error("Ошибка сохранения турнира", slog.String("error", err.Error()))
		}
	}
}
