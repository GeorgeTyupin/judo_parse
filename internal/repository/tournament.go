package repository

import (
	"context"
	"judo/internal/models"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type TournamentRepository struct {
	db     *pgxpool.Pool
	logger *slog.Logger
}

func NewTournamentRepository(db *pgxpool.Pool, logger *slog.Logger) *TournamentRepository {
	return &TournamentRepository{
		db:     db,
		logger: logger,
	}
}

func (r *TournamentRepository) SaveTournament(ctx context.Context, tourRow models.TournamentDBRow) {
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

func (r *TournamentRepository) SaveJudoka(ctx context.Context, judokaRow models.JudokaDBRow) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `
		INSERT INTO judokas (last_name, first_name, last_name_rus, first_name_rus, weight_category, birth_date, birth_place, gender, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10);`

	_, err := r.db.Exec(ctx, query,
		judokaRow.LastName,
		judokaRow.FirstName,
		judokaRow.LastNameRus,
		judokaRow.FirstNameRus,
		judokaRow.WeightCategory,
		judokaRow.BirthDate,
		judokaRow.BirthPlace,
		judokaRow.Gender,
		judokaRow.CreatedAt,
		judokaRow.UpdatedAt,
	)
	if err != nil {
		r.logger.Error("Ошибка сохранения дзюдоиста", slog.String("error", err.Error()))
	}
}
