package dupfind

import "judo/internal/models"

const (
	NotDuplicate = "Не дубликат"
	Type1        = "Name-First Name-Country"
	Type2        = "Name-First Name"
	Type3        = "Name (90%)"
	Type4        = "Name=Name (First Name (-10%))"
)

type DuplicateFinder struct {
	uniqueJudoka []*models.Judoka
}

func NewDuplicateFinder() *DuplicateFinder {
	return &DuplicateFinder{
		make([]*models.Judoka, 0),
	}
}

func (df *DuplicateFinder) GetDuplicateType(judoka *models.Judoka) string {
	for _, uJudoka := range df.uniqueJudoka {
		switch {
		case CheckType1(judoka, uJudoka):
			return Type1
		case CheckType2(judoka, uJudoka):
			return Type2
		case CheckType3(judoka, uJudoka):
			return Type3
		case CheckType4(judoka, uJudoka):
			return Type4
		}
	}

	df.uniqueJudoka = append(df.uniqueJudoka, judoka)

	return NotDuplicate
}
