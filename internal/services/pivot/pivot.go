package pivot

import (
	birthdateResolver "judo/internal/lib/utils/note/birthdate_resolver"
	"judo/internal/lib/utils/note/locresolver"
	ratingResolver "judo/internal/lib/utils/note/rating_resolver"
	"judo/internal/lib/utils/note/russifiers"
	"judo/internal/models"
)

func ProcessData(
	data models.ExcelSheet,
	judokaRussifier russifiers.JudokaRussifier,
	columnSplitter *locresolver.LocationResolver,
	ratingResolver ratingResolver.RatingResolver,
	birthDateResolver birthdateResolver.BirthDateResolver,
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
						ratingResolver,
					)
					note.DOB = birthDateResolver.Resolve(man.FirstName, man.LastName)
					notes = append(notes, note)
				}
			}
		}
	}

	return notes
}
