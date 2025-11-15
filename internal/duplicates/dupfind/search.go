package dupfind

import (
	"judo/internal/models"
	"sync"
)

// SearchDuplicates ищет все типы дублей в данных
func SearchDuplicates(wg *sync.WaitGroup, data models.ExelSheet) {
	defer wg.Done()

	// TODO: реализовать поиск дублей
}

// FindType1 - Тип 1: Name-First Name-Country
//
// Анализируемые поля: NAME, FIRSTNAME, COUNTRY
// Логика: Совпадают фамилии (NAME), страна (COUNTRY) и первая буква имени (FIRSTNAME)
// Критерий: Длина (количество символов) в FIRSTNAME <= 2
// Обозначение в листе "Дубли": Name-First Name-Country
func FindType1(data models.ExelSheet) {
	// TODO: реализовать
}

// FindType2 - Тип 2: Name-First Name
//
// Анализируемые поля: NAME, FIRSTNAME
// Логика: Совпадают фамилии (NAME) и первая буква имени (FIRSTNAME)
// Критерий: Длина (количество символов) в FIRSTNAME <= 2
// Обозначение в листе "Дубли": Name-First Name
func FindType2(data models.ExelSheet) {
	// TODO: реализовать
}

// FindType3 - Тип 3: Name (90%)
//
// Анализируемые поля: NAME
// Логика: Поиск записей с похожими на 90% фамилиями (NAME)
// Критерий: Сходство фамилий >= 90%
// Обозначение в листе "Дубли": Name (90%)
func FindType3(data models.ExelSheet) {
	// TODO: реализовать
}

// FindType4 - Тип 4: Name=Name (First Name (-10%))
//
// Анализируемые поля: NAME, FIRSTNAME
// Логика: Поиск записей с одинаковыми фамилиями (NAME), но отличающимися именами (FIRSTNAME)
// Критерий: NAME совпадают полностью, FIRSTNAME отличаются на ~10%
// Обозначение в листе "Дубли": Name=Name (First Name (-10%))
func FindType4(data models.ExelSheet) {
	// TODO: реализовать
}
