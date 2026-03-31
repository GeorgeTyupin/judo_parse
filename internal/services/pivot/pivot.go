package pivot

import (
	"judo/internal/models"
)

func ProcessData(data models.ExcelSheet) []*models.Note {
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
