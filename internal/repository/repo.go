package repository

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func truncateTable(ctx context.Context, db *pgxpool.Pool, table string) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_, err := db.Exec(ctx, fmt.Sprintf(`TRUNCATE TABLE %s RESTART IDENTITY CASCADE;`, table))
	return err
}

type CommonRepository struct {
	*tournamentRepository
	*judokaRepository
	*cityRepository
	*countryRepository
	*sportClubRepository
}

func NewCommonRepository(pool *pgxpool.Pool, logger *slog.Logger) *CommonRepository {
	return &CommonRepository{
		tournamentRepository: newTournamentRepository(pool, logger),
		judokaRepository:     newJudokaRepository(pool, logger),
		cityRepository:       newCityRepository(pool, logger),
		countryRepository:    newCountryRepository(pool, logger),
		sportClubRepository:  newSportClubRepository(pool, logger),
	}
}
