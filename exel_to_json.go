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

func ExelToJson() {
	file, _ := excelize.OpenFile("Соревнования.xlsx")
	sheetList := file.GetSheetList()
	toJson := make(ExelSheet)

	for _, curSheet := range sheetList {
		rows, err := file.GetRows(curSheet)
		rows = rows[3:]
		lenTables := findLenTables(rows[0])
		if err != nil {
			fmt.Println(err)
			return
		}

		left := 1
		cnt := 0

		for _, lenCurRow := range lenTables {
			var tournament Tournament
			curWeightCategory := make(map[string][]Judoka)
			right := left + lenCurRow

			for i, row := range rows {
				if lenCurRow > len(row) {
					continue
				}

				curRow := row[left:right]
				if !re.MatchString(curRow[0]) || ((reNum.MatchString(curRow[0]) && len(curRow[0]) <= 2) && !re.MatchString(curRow[1])) {
					continue
				}

				switch i {
				case 0:
					tournament.Name = curRow[0]
				case 1:
					tournament.Description = curRow[0]
				case 2:
					tournament.Date = curRow[0]
				case 3:
					tournament.Gender = curRow[0]
				default:
					var curWeightCategoryName string
					if reNum.MatchString(curRow[0]) {
						if len(curRow[0]) > 2 {
							curWeightCategoryName = curRow[0]
							curWeightCategory[curWeightCategoryName] = make([]Judoka, 0)
						} else {
							athlete := Judoka{
								Rank:      curRow[0],
								Name:      curRow[1],
								FirstName: curRow[2],
								JUDOKA:    curRow[3],
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

			if cnt > 0 {
				break
			}
			cnt++
		}
	}

	jsonData, err := json.MarshalIndent(toJson, "", "    ")
	if err != nil {
		log.Fatalf("Ошибка маршалинга: %v", err)
	}

	newJson, err := os.Create("Соревнования.json")
	if err != nil {
		log.Fatalf("Ошибка создания файла: %v", err)
	}
	defer newJson.Close()

	_, err = newJson.Write(jsonData)
	if err != nil {
		log.Fatalf("Ошибка записи в файл: %v", err)
	}

}
