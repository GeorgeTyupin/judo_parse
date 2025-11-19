package dupfind

import (
	"judo/internal/models"

	"github.com/hbollon/go-edlib"
)

const (
	MaxShortNameLength     = 2   // Максимальная длина сокращенного имени
	MinSimilarityThreshold = 0.9 // Минимальная необходимая схожесть строк
	MaxSimilarityThreshold = 1   // Максимальная допустимая
)

// CheckType1 - Тип 1: Name-First Name-Country
//
// Анализируемые поля: NAME, FIRSTNAME, COUNTRY
// Логика: Совпадают фамилии (NAME), страна (COUNTRY) и первая буква имени (FIRSTNAME)
// Критерий: Длина (количество символов) в FIRSTNAME <= 2
// Обозначение в листе "Дубли": Name-First Name-Country
func CheckType1(judoka, uJudoka *models.Judoka) bool {
	if len(judoka.FirstName) > MaxShortNameLength {
		return false
	}

	return judoka.Name == uJudoka.Name &&
		judoka.Country == uJudoka.Country &&
		judoka.FirstName[0] == uJudoka.FirstName[0]
}

// CheckType2 - Тип 2: Name-First Name
//
// Анализируемые поля: NAME, FIRSTNAME
// Логика: Совпадают фамилии (NAME) и первая буква имени (FIRSTNAME)
// Критерий: Длина (количество символов) в FIRSTNAME <= 2
// Обозначение в листе "Дубли": Name-First Name
func CheckType2(judoka, uJudoka *models.Judoka) bool {
	if len(judoka.FirstName) > MaxShortNameLength {
		return false
	}

	return judoka.Name == uJudoka.Name &&
		judoka.FirstName[0] == uJudoka.FirstName[0]
}

// CheckType3 - Тип 3: Name=Name (First Name (-10%))
//
// Анализируемые поля: NAME, FIRSTNAME
// Логика: Поиск записей с одинаковыми фамилиями (NAME), но отличающимися именами (FIRSTNAME)
// Критерий: NAME совпадают полностью, FIRSTNAME отличаются на ~10%
// Обозначение в листе "Дубли": Name=Name (First Name (-10%))
func CheckType3(judoka, uJudoka *models.Judoka) bool {
	similarity, err := edlib.StringsSimilarity(judoka.FirstName, uJudoka.FirstName, edlib.Levenshtein)
	return err == nil &&
		judoka.Name == uJudoka.Name &&
		MinSimilarityThreshold <= similarity && similarity < MaxSimilarityThreshold
}

// CheckType4 - Тип 4: Name (90%)
//
// Анализируемые поля: NAME
// Логика: Поиск записей с похожими на 90% фамилиями (NAME)
// Критерий: Сходство фамилий >= 90%
// Обозначение в листе "Дубли": Name (90%)
func CheckType4(judoka, uJudoka *models.Judoka) bool {
	similarity, err := edlib.StringsSimilarity(judoka.Name, uJudoka.Name, edlib.Levenshtein)
	return err == nil &&
		MinSimilarityThreshold <= similarity && similarity < MaxSimilarityThreshold

}
