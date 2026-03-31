package export

import (
	"context"
	"judo/internal/models"
)

type Repository interface {
	SaveTournament(context.Context, *models.TournamentDBRow)
}

type ExportService struct {
	DB   Repository
	data models.ExcelSheet
}

func NewExportService(db Repository, data models.ExcelSheet) (*ExportService, error) {
	return &ExportService{
		DB:   db,
		data: data,
	}, nil
}

func (es *ExportService) ProcessTournament(ctx context.Context) {
	for _, sheet := range es.data {
		for _, tournament := range sheet {
			row := models.NewTournamentDBRow(tournament)

			es.DB.SaveTournament(ctx, row)
		}
	}

}
