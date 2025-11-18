package parseio

import (
	"fmt"

	"judo/internal/models"

	"github.com/xuri/excelize/v2"
)

type Writer struct {
	Name string
	File *excelize.File
}

func NewWriter(name string) *Writer {
	table := excelize.NewFile()

	pTable := Writer{
		Name: name,
		File: table,
	}

	pTable.setHeader()

	return &pTable
}

func (t *Writer) setHeader() {
	for i, header := range models.Headers {
		cell := fmt.Sprintf("%c1", 'A'+i)
		t.File.SetCellValue("Sheet1", cell, header)
	}
}

func (t *Writer) SaveTable() {
	if err := t.File.SaveAs(fmt.Sprintf("%s.xlsx", t.Name)); err != nil {
		fmt.Printf("Ошибка сохранения файла %s: %v\n", t.Name, err)
	}
	t.File.Close()
}

func (t *Writer) Write(notes []*models.Note) {
	for i, note := range notes {
		note.SaveNote(t.File, i)
	}
}
