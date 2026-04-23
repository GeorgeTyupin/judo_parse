package repository

import (
	"context"
	"judo/internal/models"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type judokaRepository struct {
	db     *pgxpool.Pool
	logger *slog.Logger
}

func newJudokaRepository(db *pgxpool.Pool, logger *slog.Logger) *judokaRepository {
	return &judokaRepository{
		db:     db,
		logger: logger,
	}
}

func (r *judokaRepository) SaveAllJudokas(ctx context.Context, rows []models.JudokaDBRow) {
	if err := truncateTable(ctx, r.db, "judokas"); err != nil {
		r.logger.Error("Ошибка очистки таблицы judokas", slog.String("error", err.Error()))
		return
	}

	query := `
		INSERT INTO judokas (last_name, first_name, last_name_rus, first_name_rus, weight_category, birth_date, birth_place, gender, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10);`

	for _, row := range rows {
		insCtx, insCancel := context.WithTimeout(ctx, 5*time.Second)

		_, err := r.db.Exec(insCtx, query,
			row.LastName,
			row.FirstName,
			row.LastNameRus,
			row.FirstNameRus,
			row.WeightCategory,
			row.BirthDate,
			row.BirthPlace,
			row.Gender,
			row.CreatedAt,
			row.UpdatedAt,
		)
		insCancel()

		if err != nil {
			r.logger.Error("Ошибка сохранения дзюдоиста", slog.String("error", err.Error()))
		}
	}
}
