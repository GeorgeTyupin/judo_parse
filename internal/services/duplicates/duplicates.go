package duplicates

import (
	"judo/internal/lib/utils/note/locresolver"
	ratingResolver "judo/internal/lib/utils/note/rating_resolver"
	"judo/internal/lib/utils/note/russifiers"
	"judo/internal/models"
	"judo/internal/services/duplicates/dupfind"
	"sort"
)

func ProcessData(data models.ExcelSheet,
	judokaRussifier russifiers.JudokaRussifier,
	columnSplitter *locresolver.LocationResolver,
	ratingResolver ratingResolver.RatingResolver,
) []models.DuplicateNote {
	dupNotes := make([]models.DuplicateNote, 0)

	finder := dupfind.NewDuplicateFinder()

	for _, sheet := range data {
		for _, tournament := range sheet {
			for categoryName, category := range tournament.WeightCategories {
				for _, man := range category {
					dupType, original := finder.GetDuplicateType(man)
					if dupType == dupfind.NotDuplicate {
						continue
					}
					note := models.NewDuplicateNote(tournament,
						man,
						categoryName,
						dupType,
						original,
						judokaRussifier,
						columnSplitter,
						ratingResolver,
					)
					dupNotes = append(dupNotes, note)
				}
			}
		}
	}

	sort.Slice(dupNotes, func(i, j int) bool {
		return dupNotes[i].Note.JUDOKA < dupNotes[j].Note.JUDOKA
	})

	return dupNotes
}
