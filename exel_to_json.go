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

func isValidDataRow(curRow []string) bool {
	if len(curRow) == 0 {
		return false
	}
	return re.MatchString(curRow[0]) && !(reNum.MatchString(curRow[0]) && len(curRow[0]) <= 2 && !re.MatchString(curRow[1]))
}

func readTournamentHeader(rows [][]string, i, left, right int) (Tournament, int) {
	var tournament Tournament

	for j := 0; j < 4; j++ {
		if i+j >= len(rows) || left >= len(rows[i+j]) {
			continue
		}

		var tournamentRow []string
		if right > len(rows[i+j]) {
			tournamentRow = rows[i+j][left:]
		} else {
			tournamentRow = rows[i+j][left:right]
		}

		if len(tournamentRow) == 0 {
			continue
		}

		switch j {
		case 0:
			tournament.Name = tournamentRow[0]
		case 1:
			tournament.Description = tournamentRow[0]
		case 2:
			tournament.Date = tournamentRow[0]
		case 3:
			tournament.Gender = tournamentRow[0]
		}
	}

	return tournament, i + 3
}

func createJudoka(curRow []string, lenCurTable int) Judoka {
	athlete := Judoka{
		Rank:      curRow[0],
		Name:      curRow[1],
		FirstName: curRow[2],
		JUDOKA:    curRow[3],
	}

	if lenCurTable > 4 {
		athlete.Country = curRow[4]
	}

	return athlete
}

func createTournament(left, right, lenCurTable int, rows [][]string) *Tournament {
	var tournament Tournament
	var curWeightCategoryName string
	WeightCategories := make(map[string][]Judoka)

	//Проход по турниру
	for i := 0; i < len(rows); i++ {
		row := rows[i]
		if left > len(row) {
			continue
		}

		isNewTournament := false
		curRow := row[left:right]

		//Отсеиваем пустые и другие строки без нужной информации
		if !isValidDataRow(curRow) {
			continue
		}

		if curRow[0] == "_" {
			isNewTournament = true
			i++
		}

		if isNewTournament {
			// Получаем шапку турнира
			tournament, i = readTournamentHeader(rows, i, left, right)
		} else {
			// fmt.Println(curRow[0])
			if reNum.MatchString(curRow[0]) || strings.Contains(curRow[0], "Open") {

				if len(curRow[0]) > 2 {
					curWeightCategoryName = curRow[0]
					WeightCategories[curWeightCategoryName] = make([]Judoka, 0)
				} else {
					athlete := createJudoka(curRow, lenCurTable)
					WeightCategories[curWeightCategoryName] = append(WeightCategories[curWeightCategoryName], athlete)
				}
			} else {
				continue
			}
		}

	}

	tournament.WeightCategories = WeightCategories

	return &tournament
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

		fmt.Println(curSheet)

		if err != nil {
			return make(ExelSheet), err
		}

		left := 1

		//Проход по всей таблице
		for _, lenCurTable := range lenTables {
			right := left + lenCurTable

			tournament := createTournament(left, right, lenCurTable, rows)

			if _, exists := toJson[curSheet]; !exists {
				toJson[curSheet] = make([]Tournament, 0)
			}
			toJson[curSheet] = append(toJson[curSheet], *tournament)

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
