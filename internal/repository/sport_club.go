package repository

import (
	"context"
	"judo/internal/models"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type sportClubRepository struct {
	db     *pgxpool.Pool
	logger *slog.Logger
}

func newSportClubRepository(db *pgxpool.Pool, logger *slog.Logger) *sportClubRepository {
	return &sportClubRepository{
		db:     db,
		logger: logger,
	}
}

func (r *sportClubRepository) SaveAllSportClubs(ctx context.Context, rows []models.SportClubDBRow) {
	if err := truncateTable(ctx, r.db, "sport_clubs"); err != nil {
		r.logger.Error("Ошибка очистки таблицы sport_clubs", slog.String("error", err.Error()))
		return
	}

	query := `
		INSERT INTO sport_clubs (name, name_rus, full_name, founded, city_id, region, head_coach, description, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10);`

	for _, row := range rows {
		insCtx, insCancel := context.WithTimeout(ctx, 5*time.Second)

		_, err := r.db.Exec(insCtx, query,
			row.Name,
			row.NameRus,
			row.FullName,
			row.Founded,
			row.CityID,
			row.Region,
			row.HeadCoach,
			row.Description,
			row.CreatedAt,
			row.UpdatedAt,
		)
		insCancel()

		if err != nil {
			r.logger.Error("Ошибка сохранения спортивного общества", slog.String("error", err.Error()))
		}
	}
}
