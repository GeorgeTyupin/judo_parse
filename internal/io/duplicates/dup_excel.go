package dupio

import (
	"fmt"
	"judo/internal/models"

	"github.com/xuri/excelize/v2"
)

type DuplicateTable struct {
	Name   string
	Table  *excelize.File
	Sheets []string
}

func InitTable(name string) *DuplicateTable {
	table := excelize.NewFile()

	dTable := DuplicateTable{
		Name:  name,
		Table: table,
		Sheets: []string{
			"Тип 1",
			"Тип 2",
			"Тип 3",
			"Тип 4",
		},
	}

	dTable.setHeader()

	return &dTable
}

func (d *DuplicateTable) setHeader() {
	for _, sheet := range d.Sheets {
		for i, header := range models.Headers {
			cell := fmt.Sprintf("%c1", 'A'+i)
			d.Table.SetCellValue(sheet, cell, header)
		}
	}
}
