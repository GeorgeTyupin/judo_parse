package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/xuri/excelize/v2"
)

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
	f.SetCellValue("Sheet1", fmt.Sprintf("O%d", rowNum), note.NAME_RUS)
	f.SetCellValue("Sheet1", fmt.Sprintf("P%d", rowNum), note.FIRSTNAME_RUS)
	f.SetCellValue("Sheet1", fmt.Sprintf("Q%d", rowNum), note.JUDOKA_RUS)
	f.SetCellValue("Sheet1", fmt.Sprintf("R%d", rowNum), note.COUNTRY)
	f.SetCellValue("Sheet1", fmt.Sprintf("S%d", rowNum), note.SO)
}

func formatDate(date string) string {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Поймана паника, ошибка с датой ", date, strings.Contains(date, "-"), r)
		}
	}()
	var result string

	if len(date) < 6 {
		fmt.Println(date)
		return date
	}

	if strings.Contains(date, "-") {
		result = strings.TrimSpace(strings.Split(date, "-")[1])
	} else {
		result = date
	}

	// fmt.Println(result)

	for month, num := range monthMap {
		result = strings.Replace(result, month, num, -1)
	}
	return result
}

func JsonToExel() {
	file, errExel := excelize.OpenFile("Сводная таблица.xlsx")
	if errExel != nil {
		file = excelize.NewFile()
	}

	headers := []string{"TOURNAMENT", "TOUR_TYPE", "TOUR_PLACE", "TOUR_CITY", "TOUR_COUNTRY", "DATE", "YEAR", "GENDER", "WEIGHT_CATEGORY", "WC", "RANK", "NAME", "FIRSTNAME", "JUDOKA", "NAME_RUS", "FIRSTNAME_RUS", "JUDOKA_RUS", "COUNTRY", "SO"}

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
					descParts := strings.Split(tournament.Description, "—")
					tourType := strings.TrimSpace(safeGet(descParts, 0))
					tourPlace := strings.TrimSpace(safeGet(descParts, 1))

					tourCity := ""
					if placeParts := strings.SplitN(tourPlace, ",", 2); len(placeParts) > 0 {
						tourCity = strings.TrimSpace(placeParts[0])
					}

					tourCountry := ""
					if startIdx := strings.Index(tourPlace, "("); startIdx != -1 {
						if endIdx := strings.Index(tourPlace[startIdx:], ")"); endIdx != -1 {
							countryPart := tourPlace[startIdx+1 : startIdx+endIdx]
							if commaParts := strings.Split(countryPart, ","); len(commaParts) > 0 {
								tourCountry = strings.TrimSpace(commaParts[0])
							}
						}
					}

					year := ""
					if len(tournament.Date) >= 4 {
						year = tournament.Date[len(tournament.Date)-4:]
					}

					wc := ""
					if catParts := strings.Fields(categoryName); len(catParts) > 1 {
						wc = strings.Join(catParts[1:], " ")
					}

					nameRus := transliterate(man.Name)
					firstName := transliterate(man.FirstName)
					judokaRus := transliterate(man.JUDOKA)

					// fmt.Println(formatDate(tournament.Date))
					formatDate(tournament.Date)

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
						NAME_RUS:       nameRus,
						FIRSTNAME_RUS:  firstName,
						JUDOKA_RUS:     judokaRus,
						COUNTRY:        man.Country,
						SO:             man.SO,
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
