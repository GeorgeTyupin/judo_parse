package parseutils

import (
	"regexp"
	"strings"
)

var (
	re    = regexp.MustCompile(`\S+`)
	ReNum = regexp.MustCompile(`\d+`)
)

func FindLenTables(row []string) []int {
	var lenTables []int
	var arr []int

	for _, elem := range row {
		if elem == "|" {
			arr = append(arr, 1)
		} else if strings.ToLower(elem) == "end" {
			arr = append(arr, 2)
		} else {
			arr = append(arr, 0)
		}
	}

	lenArr := len(arr)
	for i := 1; i < lenArr; i++ {
		if arr[i] == 0 {
			start := i
			for i < lenArr && arr[i] == 0 {
				i++
			}
			lenTables = append(lenTables, i-start)
			if i < lenArr && arr[i] == 2 {
				break
			}
		} else {
			i++
		}
	}

	return lenTables
}

func IsValidDataRow(curRow []string) bool {
	if len(curRow) == 0 {
		return false
	}
	return re.MatchString(curRow[0]) && !(ReNum.MatchString(curRow[0]) && len(curRow[0]) <= 2 && !re.MatchString(curRow[1]))
}
