package parseio

import (
	"fmt"

	"github.com/xuri/excelize/v2"
)

type Reader struct {
	Files     []*excelize.File
	fileNames []string
}

func NewReader(fileNames []string) (*Reader, error) {
	files := make([]*excelize.File, 0)
	for _, name := range fileNames {
		file, err := excelize.OpenFile(fmt.Sprintf("%s.xlsx", name))
		if err != nil {
			return nil, fmt.Errorf("не удалось открыть исходный файл %s, возникла ошибка %w", name, err)
		}
		files = append(files, file)
	}

	reader := &Reader{
		Files:     files,
		fileNames: fileNames,
	}

	return reader, nil
}

func (r *Reader) ReadSheets() (map[string][][]string, error) {
	data := make(map[string][][]string)

	for i, file := range r.Files {
		sheetList := file.GetSheetList()

		for _, curSheet := range sheetList {
			if string(curSheet[0]) == "_" {
				continue
			}

			// Если в разных файлах есть листы с одинаковым названием, они перезапишутся.
			// Поэтому используем ключ вида "filename_sheetname".
			curSheet = fmt.Sprintf("%s_%s", r.fileNames[i], curSheet)

			fmt.Println(curSheet)

			rows, err := file.GetRows(curSheet)
			if err != nil {
				return nil, fmt.Errorf("не удалось прочесть строки в файле %s в листе %s, возникла ошибка %w", r.fileNames[i], curSheet, err)
			}

			rows = rows[4:]

			data[curSheet] = rows
		}

	}

	return data, nil
}
