package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/xuri/excelize/v2"
)

type Note struct {
	TOURNAMENT      string
	GENDER          string
	WEIGHT_CATEGORY string
	RANK            int
	NAME            string
	FIRSTNAME       string
	JUDOKA          string
	COUNTRY         string
}

func decode(file *os.File) (ExelSheet, error) {
	decode := json.NewDecoder(file)
	var list ExelSheet
	err := decode.Decode(&list)
	if err != nil {
		return nil, errors.New("неудалось прочитать ошибку")
	}

	return list, nil
}

func save_note(note Note, f *excelize.File, i int) {
	rowNum := i + 2
	f.SetCellValue("Sheet1", fmt.Sprintf("A%d", rowNum), note.TOURNAMENT)
	f.SetCellValue("Sheet1", fmt.Sprintf("B%d", rowNum), note.GENDER)
	f.SetCellValue("Sheet1", fmt.Sprintf("C%d", rowNum), note.WEIGHT_CATEGORY)
	f.SetCellValue("Sheet1", fmt.Sprintf("D%d", rowNum), note.RANK)
	f.SetCellValue("Sheet1", fmt.Sprintf("E%d", rowNum), note.NAME)
	f.SetCellValue("Sheet1", fmt.Sprintf("F%d", rowNum), note.FIRSTNAME)
	f.SetCellValue("Sheet1", fmt.Sprintf("G%d", rowNum), note.JUDOKA)
	f.SetCellValue("Sheet1", fmt.Sprintf("H%d", rowNum), note.COUNTRY)
}

func JsonToExel() {
	// exel_to_json()

	file, err_exel := excelize.OpenFile("Сводная таблица.xlsx")
	if err_exel != nil {
		file = excelize.NewFile()
	}

	headers := []string{"TOURNAMENT", "GENDER", "WEIGHT_CATEGORY", "RANK", "NAME", "FIRSTNAME", "JUDOKA", "COUNTRY"}

	// Записываем заголовки в первую строку
	for i, header := range headers {
		cell := fmt.Sprintf("%c1", 'A'+i)
		file.SetCellValue("Sheet1", cell, header)
	}

	json_file, err_json := os.Open("Соревнования.json")
	if err_json != nil {
		panic("Файл json не удалось прочесть")
	}
	defer json_file.Close()

	data, err := decode(json_file)
	if err != nil {
		panic(err)
	}

	// fmt.Println(data["URS_NC"])
	var cnt int = 0
	for _, sheet := range data {
		for _, tournament := range sheet {
			for category_name, category := range tournament.WeightCategories {
				for _, man := range category {
					var note Note
					note.TOURNAMENT = tournament.Name
					note.GENDER = tournament.Gender
					note.WEIGHT_CATEGORY = category_name
					note.RANK = man.Rank
					note.NAME = man.Name
					note.FIRSTNAME = man.FirstName
					note.JUDOKA = man.JUDOKA
					note.COUNTRY = man.Country
					save_note(note, file, cnt)
					cnt++
				}
			}
		}
	}

	if err := file.SaveAs("Сводная таблица.xlsx"); err != nil {
		panic("Изменения не записаны")
	} else {
		fmt.Println("Изменения записаны")
	}
}
