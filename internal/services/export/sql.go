package export

import (
	"context"
	"judo/internal/models"
)

type Repository interface {
	SaveTournament(context.Context, models.TournamentDBRow)
	SaveJudoka(context.Context, models.JudokaDBRow)
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
	for _, sheet := range data {
		for _, tournament := range sheet {
			row := models.NewTournamentDBRow(tournament)

			es.db.SaveTournament(ctx, row)
		}
	}
}

func (es *ExportService) SaveJudokas(ctx context.Context, data []models.JudokaDBRow) {
	for _, judoka := range data {
		es.db.SaveJudoka(ctx, judoka)
	}
}
