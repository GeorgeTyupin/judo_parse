package lsio

import (
	"fmt"

	"github.com/xuri/excelize/v2"
)

type Reader struct {
	File     *excelize.File
	fileName string
}

func NewReader(name string) (*Reader, error) {
	file, err := excelize.OpenFile(fmt.Sprintf("%s.xlsx", name))
	if err != nil {
		return nil, fmt.Errorf("не удалось открыть исходный файл %s, возникла ошибка %w", name, err)
	}

	reader := &Reader{
		File:     file,
		fileName: name,
	}

	return reader, nil
}

func (r *Reader) Read() (any, error) {
	curSheet := r.File.GetSheetList()[0]

	rows, err := r.File.GetRows(curSheet)
	if err != nil {
		return nil, fmt.Errorf("не удалось прочесть строки в файле %s в листе %s, возникла ошибка %w", r.fileName, curSheet, err)
	}

	return rows, nil
}
