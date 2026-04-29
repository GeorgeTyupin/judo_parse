package dictio

import (
	"fmt"

	"github.com/xuri/excelize/v2"
)

type Reader struct {
	File     *excelize.File
	fileName string
}

func NewReader(fileName string) (*Reader, error) {
	file, err := excelize.OpenFile(fileName)
	if err != nil {
		return nil, fmt.Errorf("не удалось открыть исходный файл %s, возникла ошибка %w", fileName, err)
	}

	reader := &Reader{
		File:     file,
		fileName: fileName,
	}

	return reader, nil
}

func (r *Reader) Read() (map[string][][]string, error) {
	data := make(map[string][][]string)

	sheetList := r.File.GetSheetList()

	for _, curSheet := range sheetList {
		rows, err := r.File.GetRows(curSheet)
		if err != nil {
			return nil, fmt.Errorf("не удалось прочесть строки в файле %s в листе %s, возникла ошибка %w", r.fileName, curSheet, err)
		}

		data[curSheet] = rows[1:]
	}

	return data, nil
}
