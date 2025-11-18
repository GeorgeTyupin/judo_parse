package pivot

import (
	parseio "judo/internal/io/excel/parse"
	"judo/internal/models"
)

type PivotService struct {
	Writer *parseio.Writer
}

func NewPivotService() *PivotService {
	return &PivotService{
		Writer: parseio.NewWriter("Сводная таблица"),
	}
}

func (ps *PivotService) ProcessData(data models.ExсelSheet) []*models.Note {
	notes := make([]*models.Note, 0)

	for _, sheet := range data {
		for _, tournament := range sheet {
			for categoryName, category := range tournament.WeightCategories {
				for _, man := range category {
					note := models.NewNote(tournament, man, categoryName)

					notes = append(notes, note)
				}
			}
		}
	}

	return notes
}
