package duplicates

// DuplicateType1 - Тип 1: Name-First Name-Country
//
// Анализируемые поля: NAME, FIRSTNAME, COUNTRY
// Логика: Совпадают фамилии (NAME), страна (COUNTRY) и первая буква имени (FIRSTNAME)
// Критерий: Длина (количество символов) в FIRSTNAME <= 2
// Обозначение в листе "Дубли": Name-First Name-Country
type DuplicateType1 struct {
}

// DuplicateType2 - Тип 2: Name-First Name
//
// Анализируемые поля: NAME, FIRSTNAME
// Логика: Совпадают фамилии (NAME) и первая буква имени (FIRSTNAME)
// Критерий: Длина (количество символов) в FIRSTNAME <= 2
// Обозначение в листе "Дубли": Name-First Name
type DuplicateType2 struct {
}

// DuplicateType3 - Тип 3: Name (90%)
//
// Анализируемые поля: NAME
// Логика: Поиск записей с похожими на 90% фамилиями (NAME)
// Критерий: Сходство фамилий >= 90%
// Обозначение в листе "Дубли": Name (90%)
type DuplicateType3 struct {
}

// DuplicateType4 - Тип 4: Name=Name (First Name (-10%))
//
// Анализируемые поля: NAME, FIRSTNAME
// Логика: Поиск записей с одинаковыми фамилиями (NAME), но отличающимися именами (FIRSTNAME)
// Критерий: NAME совпадают полностью, FIRSTNAME отличаются на ~10%
// Обозначение в листе "Дубли": Name=Name (First Name (-10%))
type DuplicateType4 struct {
}
