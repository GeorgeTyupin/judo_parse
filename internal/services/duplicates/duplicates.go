package duplicates

import (
	"judo/internal/models"
	"judo/internal/services/duplicates/dupfind"
	"sort"

	dupio "judo/internal/io/excel/duplicates"
)

type Writer interface {
	Write(data any)
	SaveFile()
}

type DuplicateService struct {
	Writers map[string]Writer
}

func NewDuplicateService() *DuplicateService {
	return &DuplicateService{
		Writers: make(map[string]Writer),
	}
}

func (ds *DuplicateService) RegisterWriter(writerName string, writer Writer) {
	ds.Writers[writerName] = writer
}

func (ds *DuplicateService) ProcessData(data models.ExсelSheet) []*dupio.DuplicateNote {
	dupNotes := make([]*dupio.DuplicateNote, 0)

	finder := dupfind.NewDuplicateFinder()

	for _, sheet := range data {
		for _, tournament := range sheet {
			for categoryName, category := range tournament.WeightCategories {
				for _, man := range category {
					dupType, original := finder.GetDuplicateType(man)
					if dupType == dupfind.NotDuplicate {
						continue
					}
					note := dupio.NewDuplicateNote(tournament, man, categoryName, dupType, original)

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
