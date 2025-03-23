package main

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/xuri/excelize/v2"
)

var re = regexp.MustCompile(`\S+`)
var reNum = regexp.MustCompile(`\d+`)

func findLenTables(row []string) []int {
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

func ExelToJson() {
	file, _ := excelize.OpenFile("Соревнования.xlsx")
	rows, err := file.GetRows("URS_NC")
	rows = rows[3:]
	lenTables := findLenTables(rows[0])
	if err != nil {
		fmt.Println(err)
		return
	}

	left := 1
	cnt := 0
	for _, lenCurRow := range lenTables {
		right := left + lenCurRow
		for _, row := range rows {
			if lenCurRow > len(row) {
				continue
			}

			curRow := row[left:right]
			if !re.MatchString(curRow[0]) || ((reNum.MatchString(curRow[0]) && len(curRow[0]) <= 2) && !re.MatchString(curRow[1])) {
				continue
			}
			fmt.Println(curRow)

		}
		left = right + 1

		if cnt > 0 {
			break
		}
		cnt++
	}

}
