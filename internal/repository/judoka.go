package repository

import (
	"context"
	"judo/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type JudokaRepository struct {
	db *pgxpool.Pool
}

func NewJudokaRepository(db *pgxpool.Pool) *JudokaRepository {
	return &JudokaRepository{
		db: db,
	}
}

func (r *JudokaRepository) SaveJudoka(ctx context.Context, judokaRow *models.JudokaDBRow) {
}
