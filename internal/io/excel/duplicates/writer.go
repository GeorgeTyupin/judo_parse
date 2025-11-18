package dupio

import (
	"fmt"

	"judo/internal/models"

	"github.com/xuri/excelize/v2"
)

type ExcelWriter struct {
	Name string
	File *excelize.File
}

func NewWriter(name string) *ExcelWriter {
	table := excelize.NewFile()

	dTable := ExcelWriter{
		Name: name,
		File: table,
	}

	dTable.setHeader()

	return &dTable
}

func (d *ExcelWriter) setHeader() {
	for i, header := range models.Headers {
		cell := fmt.Sprintf("%c1", 'A'+i)
		d.File.SetCellValue("Sheet1", cell, header)
	}
}

func (d *ExcelWriter) SaveFile() {
	if err := d.File.SaveAs(fmt.Sprintf("%s.xlsx", d.Name)); err != nil {
		fmt.Printf("Ошибка сохранения файла %s: %v\n", d.Name, err)
	}
	d.File.Close()
}

func (d *ExcelWriter) Write(data any) {
	dupNotes, ok := data.([]*DuplicateNote)
	if !ok {
		fmt.Printf("Ошибка: ожидался тип []*DuplicateNote, получен %T\n", data)
		return
	}

	for i, note := range dupNotes {
		note.SaveNote(d.File, i)
	}
}

type DuplicateNote struct {
	Note          *models.Note
	DuplicateType string
}

func NewDuplicateNote(tournament *models.Tournament, man *models.Judoka, categoryName, dupType string) *DuplicateNote {
	return &DuplicateNote{
		Note:          models.NewNote(tournament, man, categoryName),
		DuplicateType: dupType,
	}
}

func (dupNote *DuplicateNote) SaveNote(table *excelize.File, rowIndex int) {
	rowNum := rowIndex + 2

	dupNote.Note.SaveNote(table, rowIndex)
	table.SetCellValue("Sheet1", fmt.Sprintf("W%d", rowNum), dupNote.DuplicateType)
}
