package models

import (
	"fmt"

	"github.com/xuri/excelize/v2"
)

type DuplicateNote struct {
	Note           Note
	DuplicateType  string
	OriginalJudoka string
}

func NewDuplicateNote(tournament Tournament, man Judoka, categoryName, dupType, original string) DuplicateNote {
	return DuplicateNote{
		Note:           NewNote(tournament, man, categoryName),
		DuplicateType:  dupType,
		OriginalJudoka: original,
	}
}

func (dupNote DuplicateNote) SaveNote(table *excelize.File, rowIndex int) {
	rowNum := rowIndex + 2
	dupNote.Note.SaveNote(table, rowIndex)
	table.SetCellValue("Sheet1", fmt.Sprintf("W%d", rowNum), dupNote.DuplicateType)
	table.SetCellValue("Sheet1", fmt.Sprintf("X%d", rowNum), dupNote.OriginalJudoka)
}
