package pivot

import (
	"judo/internal/lib/utils/note/colsplit"
	"judo/internal/lib/utils/note/russifiers"
	"judo/internal/models"
)

func ProcessData(
	data models.ExcelSheet,
	judokaRussifier russifiers.JudokaRussifier,
	columnSplitter *colsplit.ColumnSplitter,
) []models.Note {
	notes := make([]models.Note, 0)

	for _, sheet := range data {
		for _, tournament := range sheet {
			for categoryName, category := range tournament.WeightCategories {
				for _, man := range category {
					note := models.NewNote(
						tournament,
						man,
						categoryName,
						judokaRussifier,
						columnSplitter,
					)
					notes = append(notes, note)
				}
			}
		}
	}

	return notes
}
