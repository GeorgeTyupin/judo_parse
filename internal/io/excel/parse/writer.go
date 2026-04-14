package parseio

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

	pTable := ExcelWriter{
		Name: name,
		File: table,
	}

	pTable.setHeader()

	return &pTable
}

func (t *ExcelWriter) setHeader() {
	for i, header := range models.Headers {
		cell := fmt.Sprintf("%c1", 'A'+i)
		t.File.SetCellValue("Sheet1", cell, header)
	}
}

func (t *ExcelWriter) SaveFile() {
	if err := t.File.SaveAs(fmt.Sprintf("%s.xlsx", t.Name)); err != nil {
		fmt.Printf("Ошибка сохранения файла %s: %v\n", t.Name, err)
	}
	t.File.Close()
}

func (t *ExcelWriter) Write(notes []models.Note) {
	for i, note := range notes {
		note.SaveNote(t.File, i)
	}
}
