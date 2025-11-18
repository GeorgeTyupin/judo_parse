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
