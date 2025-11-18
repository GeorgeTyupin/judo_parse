package pivot

import (
	"judo/internal/interfaces"
	"judo/internal/models"
)

type PivotService struct {
	Writers map[string]interfaces.Writer
}

func NewPivotService() *PivotService {
	return &PivotService{
		Writers: make(map[string]interfaces.Writer),
	}
}

func (ps *PivotService) RegisterWriter(writerName string, writer interfaces.Writer) {
	ps.Writers[writerName] = writer
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
