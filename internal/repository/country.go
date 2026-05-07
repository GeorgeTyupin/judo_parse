package repository

import (
	"context"
	"judo/internal/models"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type countryRepository struct {
	db     *pgxpool.Pool
	logger *slog.Logger
}

func newCountryRepository(db *pgxpool.Pool, logger *slog.Logger) *countryRepository {
	return &countryRepository{
		db:     db,
		logger: logger,
	}
}

func (r *countryRepository) SaveAllCountries(ctx context.Context, rows []models.CountryDBRow) {
	if err := truncateTable(ctx, r.db, "countries"); err != nil {
		r.logger.Error("Ошибка очистки таблицы countries", slog.String("error", err.Error()))
		return
	}

	query := `
		INSERT INTO countries (name, iso_code, name_rus, description, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6);`

	for _, row := range rows {
		insCtx, insCancel := context.WithTimeout(ctx, 5*time.Second)

		_, err := r.db.Exec(insCtx, query,
			row.Name,
			row.ISOCode,
			row.NameRus,
			row.Description,
			row.CreatedAt,
			row.UpdatedAt,
		)
		insCancel()

		if err != nil {
			r.logger.Error("Ошибка сохранения страны", slog.String("error", err.Error()))
		}
	}
}
