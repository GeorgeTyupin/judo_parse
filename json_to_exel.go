package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/xuri/excelize/v2"
)

type Note struct {
	TOURNAMENT     string
	TOUR_TYPE      string
	DATE           string
	GENDER         string
	WeightCategory string
	RANK           string
	NAME           string
	FIRSTNAME      string
	JUDOKA         string
	COUNTRY        string
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

func saveNote(note Note, f *excelize.File, i int) {
	rowNum := i + 2
	f.SetCellValue("Sheet1", fmt.Sprintf("A%d", rowNum), note.TOURNAMENT)
	f.SetCellValue("Sheet1", fmt.Sprintf("B%d", rowNum), note.TOUR_TYPE)
	f.SetCellValue("Sheet1", fmt.Sprintf("C%d", rowNum), note.DATE)
	f.SetCellValue("Sheet1", fmt.Sprintf("D%d", rowNum), note.GENDER)
	f.SetCellValue("Sheet1", fmt.Sprintf("E%d", rowNum), note.WeightCategory)
	f.SetCellValue("Sheet1", fmt.Sprintf("F%d", rowNum), note.RANK)
	f.SetCellValue("Sheet1", fmt.Sprintf("G%d", rowNum), note.NAME)
	f.SetCellValue("Sheet1", fmt.Sprintf("H%d", rowNum), note.FIRSTNAME)
	f.SetCellValue("Sheet1", fmt.Sprintf("I%d", rowNum), note.JUDOKA)
	f.SetCellValue("Sheet1", fmt.Sprintf("J%d", rowNum), note.COUNTRY)
}

func JsonToExel() {
	file, errExel := excelize.OpenFile("Сводная таблица.xlsx")
	if errExel != nil {
		file = excelize.NewFile()
	}

	headers := []string{"TOURNAMENT", "TOUR TYPE", "DATE", "GENDER", "WEIGHT_CATEGORY", "RANK", "NAME", "FIRSTNAME", "JUDOKA", "COUNTRY"}

	// Записываем заголовки в первую строку
	for i, header := range headers {
		cell := fmt.Sprintf("%c1", 'A'+i)
		file.SetCellValue("Sheet1", cell, header)
	}

	jsonFile, errJson := os.Open(fmt.Sprintf("%s.json", FILE))
	if errJson != nil {
		panic("Файл json не удалось прочесть")
	}
	defer jsonFile.Close()

	data, err := decode(jsonFile)
	if err != nil {
		panic(err)
	}

	// fmt.Println(data["URS_NC"])
	var cnt = 0
	for _, sheet := range data {
		for _, tournament := range sheet {
			for categoryName, category := range tournament.WeightCategories {
				for _, man := range category {
					note := Note{
						TOURNAMENT:     tournament.Name,
						TOUR_TYPE:      tournament.Description,
						DATE:           tournament.Date,
						GENDER:         tournament.Gender,
						WeightCategory: categoryName,
						RANK:           man.Rank,
						NAME:           man.Name,
						FIRSTNAME:      man.FirstName,
						JUDOKA:         man.JUDOKA,
						COUNTRY:        man.Country,
					}
					saveNote(note, file, cnt)
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
