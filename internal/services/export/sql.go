package export

import (
	"context"
	"judo/internal/models"
)

type Repository interface {
	SaveAllTournaments(context.Context, []models.TournamentDBRow)
	SaveAllJudokas(context.Context, []models.JudokaDBRow)
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
