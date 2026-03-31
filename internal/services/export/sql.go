package export

import (
	"context"
	"judo/internal/models"
)

type Repository interface {
	SaveTournament(context.Context, *models.TournamentDBRow)
}

type ExportService struct {
	db   Repository
	data models.ExcelSheet
}

func NewExportService(repo Repository, data models.ExcelSheet) (*ExportService, error) {
	return &ExportService{
		db:   repo,
		data: data,
	}, nil
}

func (es *ExportService) ProcessTournament(ctx context.Context) {
	for _, sheet := range es.data {
		for _, tournament := range sheet {
			row := models.NewTournamentDBRow(tournament)

			es.db.SaveTournament(ctx, row)
		}
	}

}
