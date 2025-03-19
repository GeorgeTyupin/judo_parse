package main

import (
	"fmt"
	"strings"

	"github.com/xuri/excelize/v2"
)

// var re = regexp.MustCompile(`\S+`)

// func deleteSpaces(row []string) []string {
// 	var newRow []string
// 	for _, elem := range row {
// 		if !re.MatchString(elem) {
// 			continue
// 		}
// 		newRow = append(newRow, elem)
// 	}
// 	return newRow
// }

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
	for i := 0; i < lenArr; i++ {
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
	if err != nil {
		fmt.Println(err)
		return
	}
	// for _, name := range file.GetSheetMap() {
	// 	fmt.Println(name)
	// }
	for i := 0; i < len(rows); i++ {
		curRow := rows[i][1:6]

		// for _, col := range curRow {
		// 	fmt.Println(col)
		// }
		fmt.Println(curRow)
		// fmt.Printf("%T", rows[i])
		if i > 2 {
			break
		}
	}

}
