package main

import (
	"fmt"

	"github.com/xuri/excelize/v2"
)

func main() {
	file, _ := excelize.OpenFile("Соревнования.xlsx")
	rows, err := file.GetRows("URS_NC")
	rows = rows[4:]
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, name := range file.GetSheetMap() {
		fmt.Println(name)
	}
	for i, row := range rows {
		fmt.Println(row[1])
		if i > 10 {
			break
		}
	}
}
