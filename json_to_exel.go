package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/xuri/excelize/v2"
)

func decode(file *os.File) (ExelSheet, error) {
	decode := json.NewDecoder(file)
	var list ExelSheet
	err := decode.Decode(&list)
	if err != nil {
		return nil, errors.New("неудалось прочитать ошибку")
	}

	return list, nil
}

func safeGet(parts []string, index int) string {
	if index < len(parts) {
		return parts[index]
	}
	return ""
}

func saveNote(note Note, f *excelize.File, i int) {
	rowNum := i + 2
	f.SetCellValue("Sheet1", fmt.Sprintf("A%d", rowNum), note.TOURNAMENT)
	f.SetCellValue("Sheet1", fmt.Sprintf("B%d", rowNum), note.TOUR_TYPE)
	f.SetCellValue("Sheet1", fmt.Sprintf("C%d", rowNum), note.TOUR_PLACE)
	f.SetCellValue("Sheet1", fmt.Sprintf("D%d", rowNum), note.TOUR_CITY)
	f.SetCellValue("Sheet1", fmt.Sprintf("E%d", rowNum), note.TOUR_COUNTRY)
	f.SetCellValue("Sheet1", fmt.Sprintf("F%d", rowNum), note.DATE)
	f.SetCellValue("Sheet1", fmt.Sprintf("G%d", rowNum), note.YEAR)
	f.SetCellValue("Sheet1", fmt.Sprintf("H%d", rowNum), note.GENDER)
	f.SetCellValue("Sheet1", fmt.Sprintf("I%d", rowNum), note.WeightCategory)
	f.SetCellValue("Sheet1", fmt.Sprintf("J%d", rowNum), note.WC)
	f.SetCellValue("Sheet1", fmt.Sprintf("K%d", rowNum), note.RANK)
	f.SetCellValue("Sheet1", fmt.Sprintf("L%d", rowNum), note.NAME)
	f.SetCellValue("Sheet1", fmt.Sprintf("M%d", rowNum), note.FIRSTNAME)
	f.SetCellValue("Sheet1", fmt.Sprintf("N%d", rowNum), note.JUDOKA)
	f.SetCellValue("Sheet1", fmt.Sprintf("O%d", rowNum), note.COUNTRY)
}

func JsonToExel() {
	file, errExel := excelize.OpenFile("Сводная таблица.xlsx")
	if errExel != nil {
		file = excelize.NewFile()
	}

	headers := []string{"TOURNAMENT", "TOUR_TYPE", "TOUR_PLACE", "TOUR_CITY", "TOUR_COUNTRY", "DATE", "YEAR", "GENDER", "WEIGHT_CATEGORY", "WC", "RANK", "NAME", "FIRSTNAME", "JUDOKA", "COUNTRY"}

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
					descParts := strings.Split(tournament.Description, "-")
					tourType := strings.TrimSpace(safeGet(descParts, 0))
					tourPlace := strings.TrimSpace(safeGet(descParts, 1))

					tourCity := strings.TrimSpace(safeGet(strings.SplitN(tourPlace, ",", 2), 0))
					tourCountry := strings.Trim(safeGet(strings.SplitN(tourPlace, " ", 4), 3), "() ")

					year := strings.TrimSpace(safeGet(strings.Fields(tournament.Date), 2))
					wc := strings.TrimSpace(safeGet(strings.Fields(categoryName), 1))
					// nameRus, _ := gtranslate.Translate(man.Name, language.English, language.Russian)
					note := Note{
						TOURNAMENT:     tournament.Name,
						TOUR_TYPE:      tourType,
						TOUR_PLACE:     tourPlace,
						TOUR_CITY:      tourCity,
						TOUR_COUNTRY:   tourCountry,
						DATE:           tournament.Date,
						YEAR:           year,
						GENDER:         tournament.Gender,
						WeightCategory: categoryName,
						WC:             wc,
						RANK:           man.Rank,
						NAME:           man.Name,
						FIRSTNAME:      man.FirstName,
						JUDOKA:         man.JUDOKA,
						// NAME_RUS:       nameRus,
						COUNTRY: man.Country,
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
