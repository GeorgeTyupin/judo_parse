package dupio

import (
	"fmt"
	"log/slog"

	"judo/internal/models"

	"github.com/xuri/excelize/v2"
)

type ExcelWriter struct {
	Name   string
	File   *excelize.File
	logger *slog.Logger
}

func NewWriter(name string, logger *slog.Logger) *ExcelWriter {
	table := excelize.NewFile()

	dTable := ExcelWriter{
		Name:   name,
		File:   table,
		logger: logger,
	}

	dTable.setHeader()

	return &dTable
}

func (d *ExcelWriter) setHeader() {
	dupHeaders := append([]string{}, models.Headers...)
	dupHeaders = append(dupHeaders, "TYPE", "ORIGINAL")

	for i, header := range dupHeaders {
		cell := fmt.Sprintf("%c1", 'A'+i)
		d.File.SetCellValue("Sheet1", cell, header)
	}
}

func (d *ExcelWriter) SaveFile() {
	if err := d.File.SaveAs(fmt.Sprintf("%s.xlsx", d.Name)); err != nil {
		slog.Error("Ошибка сохранения файла", slog.String("name", d.Name), slog.String("error", err.Error()))
	}
	d.File.Close()
}

func (d *ExcelWriter) Write(dupNotes []models.DuplicateNote) {
	for i, note := range dupNotes {
		note.SaveNote(d.File, i)
	}
}
