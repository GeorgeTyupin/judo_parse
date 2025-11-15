package parseio

import (
	"fmt"

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

func (t *PivotTable) Write(notes []*models.Note) {
	for i, note := range notes {
		note.SaveNote(t.Table, i)
	}
}
