package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
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

func renderExel() (ExelSheet, error) {
	file, _ := excelize.OpenFile(fmt.Sprintf("%s.xlsx", FILE))
	sheetList := file.GetSheetList()
	toJson := make(ExelSheet)

	//Проход по всем листам
	for _, curSheet := range sheetList {
		if string(curSheet[0]) == "_" {
			continue
		}
		rows, err := file.GetRows(curSheet)
		rows = rows[4:]
		lenTables := findLenTables(rows[1])

		if lenTables[0] > 5 {
			fmt.Println(curSheet)
		}

		if err != nil {
			return make(ExelSheet), err
		}

		left := 1

		//Проход по всей таблице
		for _, lenCurRow := range lenTables {
			var tournament Tournament
			var curWeightCategoryName string
			curWeightCategory := make(map[string][]Judoka)
			right := left + lenCurRow

			//Проход по турниру
			// for i, row := range rows {
			for i := 0; i < len(rows); i++ {
				row := rows[i]
				if lenCurRow > len(row) || left > len(row) {
					continue
				}

				isNewTournament := false
				curRow := row[left:right]

				if !re.MatchString(curRow[0]) || ((reNum.MatchString(curRow[0]) && len(curRow[0]) <= 2) && !re.MatchString(curRow[1])) {
					continue
				}

				if curRow[0] == "_" {
					isNewTournament = true
					// fmt.Println("начало турнира", curSheet)
					i++
				}

				if isNewTournament {
					//мнимый цикл для прохода по шапке турнира
					for j := range 4 {
						tournamentRow := rows[i+j][left:right]
						switch j {
						case 0:
							tournament.Name = tournamentRow[0]
						case 1:
							tournament.Description = tournamentRow[0]
						case 2:
							tournament.Date = tournamentRow[0]
						case 3:
							tournament.Gender = tournamentRow[0]
							isNewTournament = false
						}
					}
					i += 3
				} else {
					if reNum.MatchString(curRow[0]) {

						if len(curRow[0]) > 2 {
							curWeightCategoryName = curRow[0]
							if curWeightCategoryName == "National Tournament Tbilisi - 1993" {
								fmt.Println(rows[i-1])
							}
							curWeightCategory[curWeightCategoryName] = make([]Judoka, 0)
						} else {
							athlete := Judoka{
								Rank:      curRow[0],
								Name:      curRow[1],
								FirstName: curRow[2],
								JUDOKA:    curRow[3],
							}

							if lenCurRow > 4 {
								athlete.Country = curRow[4]
							}

							curWeightCategory[curWeightCategoryName] = append(curWeightCategory[curWeightCategoryName], athlete)
						}
					} else {
						continue
					}
				}

			}
			tournament.WeightCategories = curWeightCategory
			if _, exists := toJson[curSheet]; !exists {
				toJson[curSheet] = make([]Tournament, 0)
			}
			toJson[curSheet] = append(toJson[curSheet], tournament)

			left = right + 1
		}
	}

	return toJson, nil
}

func ExelToJson() {
	file, err := renderExel()
	if err != nil {
		fmt.Println(err)
	}

	jsonData, err := json.MarshalIndent(file, "", "    ")
	if err != nil {
		log.Fatalf("Ошибка маршалинга: %v", err)
	}

	newJson, err := os.Create(fmt.Sprintf("%s.json", FILE))
	if err != nil {
		log.Fatalf("Ошибка создания файла: %v", err)
	}
	defer newJson.Close()

	_, err = newJson.Write(jsonData)
	if err != nil {
		log.Fatalf("Ошибка записи в файл: %v", err)
	}
}
