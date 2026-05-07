package repository

import (
	"context"
	"judo/internal/models"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type cityRepository struct {
	db     *pgxpool.Pool
	logger *slog.Logger
}

func newCityRepository(db *pgxpool.Pool, logger *slog.Logger) *cityRepository {
	return &cityRepository{
		db:     db,
		logger: logger,
	}
}

func (r *cityRepository) SaveAllCities(ctx context.Context, rows []models.CityDBRow) {
	if err := truncateTable(ctx, r.db, "cities"); err != nil {
		r.logger.Error("Ошибка очистки таблицы cities", slog.String("error", err.Error()))
		return
	}

	query := `
		INSERT INTO cities (name, name_rus, name_rus_last, republic_id, description, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7);`

	for _, row := range rows {
		insCtx, insCancel := context.WithTimeout(ctx, 5*time.Second)

		_, err := r.db.Exec(insCtx, query,
			row.Name,
			row.NameRus,
			row.NameRusLast,
			row.RepublicID,
			row.Description,
			row.CreatedAt,
			row.UpdatedAt,
		)
		insCancel()

		if err != nil {
			r.logger.Error("Ошибка сохранения города", slog.String("error", err.Error()))
		}
	}
}
