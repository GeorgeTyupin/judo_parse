package noteutils

import "strings"

// SafeGet безопасно извлекает элемент из среза строк по указанному индексу.
// Возвращает пустую строку, если индекс выходит за границы среза.
func SafeGet(parts []string, index int) string {
	if index < len(parts) {
		return parts[index]
	}
	return ""
}

var monthMap = map[string]string{
	"January":   "01",
	"February":  "02",
	"March":     "03",
	"April":     "04",
	"May":       "05",
	"June":      "06",
	"July":      "07",
	"August":    "08",
	"September": "09",
	"October":   "10",
	"November":  "11",
	"December":  "12",
}

// FormatDate форматирует строку с датой, преобразуя названия месяцев в числа.
// Обрабатывает даты с разделителем '-' и без него.
func FormatDate(date string) string {
	var result string

	if len(date) < 5 {
		return ""
	}

	if strings.Contains(date, "-") {
		result = strings.TrimSpace(strings.Split(date, "-")[1])
	} else {
		result = date
	}

	result = strings.TrimFunc(result, func(r rune) bool {
		return !((r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z'))
	})

	for month, num := range monthMap {
		result = strings.Replace(result, month, num, -1)
	}
	return result
}

// GenderInfo содержит распарсенную информацию о поле и возрасте
type GenderInfo struct {
	Gender     string // Men или Women
	Age        string // Senior, Junior, Cadets, Women
	BRange     string // Диапазон годов или возрастная категория
	GenderFull string // Полная строка
}

// ParseGenderInfo парсит строку формата "Senior Men", "Juniors Men (1961-1962)" и т.д.
func ParseGenderInfo(genderStr string) GenderInfo {
	info := GenderInfo{
		GenderFull: genderStr,
	}

	// Извлекаем диапазон годов в скобках, если есть
	if startIdx := strings.Index(genderStr, "("); startIdx != -1 {
		if endIdx := strings.Index(genderStr[startIdx:], ")"); endIdx != -1 {
			info.BRange = genderStr[startIdx : startIdx+endIdx+1]
			genderStr = strings.TrimSpace(genderStr[:startIdx])
		}
	}

	// Разделяем на части
	parts := strings.Fields(genderStr)

	if len(parts) == 0 {
		return info
	}

	// Определяем Gender (Men или Women)
	lastPart := parts[len(parts)-1]
	if lastPart == "Men" || lastPart == "Women" {
		info.Gender = lastPart
	}

	// Определяем Age
	if len(parts) > 0 {
		firstPart := parts[0]
		if firstPart == "Senior" || firstPart == "Seniors" {
			info.Age = "Senior"
			if info.BRange == "" {
				info.BRange = "Senior"
			}
		} else if firstPart == "Junior" || firstPart == "Juniors" {
			info.Age = "Junior"
		} else if firstPart == "Cadet" || firstPart == "Cadets" {
			info.Age = "Cadets"
		} else if firstPart == "Women" {
			info.Age = "Women"
		}
	}

	return info
}
