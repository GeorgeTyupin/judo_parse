package duplicates

import (
	"judo/internal/interfaces"
	"judo/internal/models"

	dupio "judo/internal/io/excel/duplicates"
)

type DuplicateService struct {
	uniqueJudoka []*models.Judoka
	Writers      map[string]interfaces.Writer
}

func NewDuplicateService() *DuplicateService {
	return &DuplicateService{
		uniqueJudoka: make([]*models.Judoka, 0),
		Writers:      make(map[string]interfaces.Writer),
	}
}

func (ds *DuplicateService) RegisterWriter(writerName string, writer interfaces.Writer) {
	ds.Writers[writerName] = writer
}

func (ds *DuplicateService) ProcessData(data models.ExсelSheet) []*dupio.DuplicateNote {
	dupNotes := make([]*dupio.DuplicateNote, 0)

	for _, sheet := range data {
		for _, tournament := range sheet {
			for categoryName, category := range tournament.WeightCategories {
				for _, man := range category {
					note := dupio.NewDuplicateNote(tournament, man, categoryName, "Тип1")

					dupNotes = append(dupNotes, note)
				}
			}
		}
	}

	return dupNotes
}
