package export

import (
	"context"
	"judo/internal/models"
)

type Repository interface {
	SaveAllTournaments(context.Context, []models.TournamentDBRow)
	SaveAllJudokas(context.Context, []models.JudokaDBRow)
	SaveAllCities(context.Context, []models.CityDBRow)
	SaveAllCountries(context.Context, []models.CountryDBRow)
	SaveAllSportClubs(context.Context, []models.SportClubDBRow)
}

type ExportService struct {
	db Repository
}

func NewExportService(repo Repository) (*ExportService, error) {
	return &ExportService{
		db: repo,
	}, nil
}

func (es *ExportService) SaveTournaments(ctx context.Context, data models.ExcelSheet) {
	var rows []models.TournamentDBRow

	for _, sheet := range data {
		for _, tournament := range sheet {
			rows = append(rows, models.NewTournamentDBRow(tournament))
		}
	}

	es.db.SaveAllTournaments(ctx, rows)
}

func (es *ExportService) SaveJudokas(ctx context.Context, data []models.JudokaDBRow) {
	es.db.SaveAllJudokas(ctx, data)
}

func (es *ExportService) SaveCities(ctx context.Context, data []models.CityDBRow) {
	es.db.SaveAllCities(ctx, data)
}

func (es *ExportService) SaveCountries(ctx context.Context, data []models.CountryDBRow) {
	es.db.SaveAllCountries(ctx, data)
}

func (es *ExportService) SaveSportClubs(ctx context.Context, data []models.SportClubDBRow) {
	es.db.SaveAllSportClubs(ctx, data)
}
