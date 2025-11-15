package parseio

import (
	"fmt"
	"sync"

	"judo/internal/models"

	"github.com/xuri/excelize/v2"
)

type PivotTable struct {
	Name  string
	Table *excelize.File
}

func InitTable(name string) *PivotTable {
	table := excelize.NewFile()

	pTable := PivotTable{
		Name:  name,
		Table: table,
	}

	pTable.setHeader()

	return &pTable
}

func (t *PivotTable) setHeader() {
	for i, header := range models.Headers {
		cell := fmt.Sprintf("%c1", 'A'+i)
		t.Table.SetCellValue("Sheet1", cell, header)
	}
}

func (t *PivotTable) SaveTable() {
	if err := t.Table.SaveAs(fmt.Sprintf("%s.xlsx", t.Name)); err != nil {
		fmt.Printf("Ошибка сохранения файла %s: %v\n", t.Name, err)
	}
	t.Table.Close()
}

func (t *PivotTable) ToExcel(wg *sync.WaitGroup, data models.ExelSheet) {
	defer wg.Done()

	var cnt = 0
	for _, sheet := range data {
		for _, tournament := range sheet {
			for categoryName, category := range tournament.WeightCategories {
				for _, man := range category {
					note := models.NewNote(tournament, man, categoryName)

					note.SaveNote(t.Table, cnt)
					cnt++
				}
			}
		}
	}
}
