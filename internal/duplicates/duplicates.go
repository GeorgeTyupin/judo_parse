package duplicates

import (
	"judo/internal/models"
	"sync"

	dupio "judo/internal/io/duplicates"

	"github.com/xuri/excelize/v2"
)

var uniqueJudoka = []*models.Judoka{}

// SearchDuplicates ищет все типы дублей в данных
func SearchDuplicates(wg *sync.WaitGroup, data models.ExelSheet) {
	defer wg.Done()
	dupTable := dupio.InitTable("Дубли")
	searchDuplicatesLoop(data, dupTable.Table)
}

func searchDuplicatesLoop(data models.ExelSheet, table *excelize.File) {
	var cnt = 0
	for _, sheet := range data {
		for _, tournament := range sheet {
			for categoryName, category := range tournament.WeightCategories {
				for _, man := range category {
					note := dupio.NewDuplicateNote(tournament, man, categoryName, "Тип1")

					note.SaveDuplicateNote(table, cnt)
					cnt++
				}
			}
		}
	}
}
