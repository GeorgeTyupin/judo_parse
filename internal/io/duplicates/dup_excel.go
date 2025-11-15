package dupio

import (
	"fmt"

	"judo/internal/models"

	"github.com/xuri/excelize/v2"
)

type DuplicateTable struct {
	Name  string
	Table *excelize.File
}

func InitTable(name string) *DuplicateTable {
	table := excelize.NewFile()

	dTable := DuplicateTable{
		Name:  name,
		Table: table,
	}

	dTable.setHeader()

	return &dTable
}

func (d *DuplicateTable) setHeader() {
	for i, header := range models.Headers {
		cell := fmt.Sprintf("%c1", 'A'+i)
		d.Table.SetCellValue("Sheet1", cell, header)
	}
}

type DuplicateNote struct {
	Note          models.Note
	DuplicateType string
}

func (dupNote *DuplicateNote) SaveDuplicateNote(table *excelize.File, rowIndex int) {
	rowNum := rowIndex + 2

	dupNote.Note.SaveNote(table, rowIndex)
	table.SetCellValue("Sheet1", fmt.Sprintf("W%d", rowNum), dupNote.DuplicateType)
}

func (d *DuplicateTable) SaveTable() {
	if err := d.Table.SaveAs(fmt.Sprintf("%s.xlsx", d.Name)); err != nil {
		fmt.Printf("Ошибка сохранения файла %s: %v\n", d.Name, err)
	}
	d.Table.Close()
}
