package repository

import (
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
)

type CommonRepository struct {
	*tournamentRepository
	*judokaRepository
}

func NewCommonRepository(pool *pgxpool.Pool, logger *slog.Logger) *CommonRepository {
	tournamentRepository := newTournamentRepository(pool, logger)
	judokaRepository := newJudokaRepository(pool, logger)

	return &CommonRepository{
		tournamentRepository: tournamentRepository,
		judokaRepository:     judokaRepository,
	}
}
