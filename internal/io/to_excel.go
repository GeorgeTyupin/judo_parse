package judioio

import (
	"fmt"
	"strings"
	"sync"

	"judo/internal/models"
	"judo/pkg/replacers"

	"github.com/xuri/excelize/v2"
)

type PivotTable struct {
	Name  string
	Table *excelize.File
}

func InitTable(name string) *PivotTable {
	table := excelize.NewFile()

	pTable := PivotTable{
		Name:  name,
		Table: table,
	}

	pTable.setHeader()

	return &pTable
}

func (t *PivotTable) setHeader() {
	headers := []string{"TOURNAMENT", "TOUR_TYPE", "TOUR_PLACE", "TOUR_CITY", "TOUR_COUNTRY", "COUNTRY_LAST", "DATE", "YEAR", "MONTH", "GENDER", "WEIGHT_CATEGORY", "WC", "RANK", "NAME", "FIRSTNAME", "JUDOKA", "NAME_RUS", "FIRSTNAME_RUS", "JUDOKA_RUS", "COUNTRY", "SO"}

	for i, header := range headers {
		cell := fmt.Sprintf("%c1", 'A'+i)
		t.Table.SetCellValue("Sheet1", cell, header)
	}
}

func (t *PivotTable) SaveTable() {
	if err := t.Table.SaveAs(fmt.Sprintf("%s.xlsx", t.Name)); err != nil {
		fmt.Printf("Ошибка сохранения файла %s: %v\n", t.Name, err)
	}
	t.Table.Close()
}

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

func safeGet(parts []string, index int) string {
	if index < len(parts) {
		return parts[index]
	}
	return ""
}

func saveNote(note models.Note, f *excelize.File, i int) {
	rowNum := i + 2
	f.SetCellValue("Sheet1", fmt.Sprintf("A%d", rowNum), note.TOURNAMENT)
	f.SetCellValue("Sheet1", fmt.Sprintf("B%d", rowNum), note.TOUR_TYPE)
	f.SetCellValue("Sheet1", fmt.Sprintf("C%d", rowNum), note.TOUR_PLACE)
	f.SetCellValue("Sheet1", fmt.Sprintf("D%d", rowNum), note.TOUR_CITY)
	f.SetCellValue("Sheet1", fmt.Sprintf("E%d", rowNum), note.TOUR_COUNTRY)
	f.SetCellValue("Sheet1", fmt.Sprintf("F%d", rowNum), note.COUNTRY_LAST)
	f.SetCellValue("Sheet1", fmt.Sprintf("G%d", rowNum), note.DATE)
	f.SetCellValue("Sheet1", fmt.Sprintf("H%d", rowNum), note.YEAR)
	f.SetCellValue("Sheet1", fmt.Sprintf("I%d", rowNum), note.MONTH)
	f.SetCellValue("Sheet1", fmt.Sprintf("J%d", rowNum), note.GENDER)
	f.SetCellValue("Sheet1", fmt.Sprintf("K%d", rowNum), note.WeightCategory)
	f.SetCellValue("Sheet1", fmt.Sprintf("L%d", rowNum), note.WC)
	f.SetCellValue("Sheet1", fmt.Sprintf("M%d", rowNum), note.RANK)
	f.SetCellValue("Sheet1", fmt.Sprintf("N%d", rowNum), note.NAME)
	f.SetCellValue("Sheet1", fmt.Sprintf("O%d", rowNum), note.FIRSTNAME)
	f.SetCellValue("Sheet1", fmt.Sprintf("P%d", rowNum), note.JUDOKA)
	f.SetCellValue("Sheet1", fmt.Sprintf("Q%d", rowNum), note.NAME_RUS)
	f.SetCellValue("Sheet1", fmt.Sprintf("R%d", rowNum), note.FIRSTNAME_RUS)
	f.SetCellValue("Sheet1", fmt.Sprintf("S%d", rowNum), note.JUDOKA_RUS)
	f.SetCellValue("Sheet1", fmt.Sprintf("T%d", rowNum), note.COUNTRY)
	f.SetCellValue("Sheet1", fmt.Sprintf("U%d", rowNum), note.SO)
}

func formatDate(date string) string {
	var result string

	if len(date) < 5 {
		return date
	}

	if strings.Contains(date, "-") {
		result = strings.TrimSpace(strings.Split(date, "-")[1])
	} else {
		result = date
	}

	// result = strings.Join(strings.Fields(result), ".")

	result = strings.TrimFunc(result, func(r rune) bool {
		return !((r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z'))
	})

	for month, num := range monthMap {
		result = strings.Replace(result, month, num, -1)
	}
	return result
}

func (t *PivotTable) ToExcel(wg *sync.WaitGroup, data models.ExelSheet) {
	defer wg.Done()

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

					nameRus := replacers.Transliterate(man.Name)
					firstName := replacers.Transliterate(man.FirstName)
					judokaRus := replacers.Transliterate(man.JUDOKA)

					note := models.Note{
						TOURNAMENT:     tournament.Name,
						TOUR_TYPE:      tourType,
						TOUR_PLACE:     tourPlace,
						TOUR_CITY:      tourCity,
						TOUR_COUNTRY:   tourCountry,
						COUNTRY_LAST:   replacers.NormalizeCityName(tourCity),
						DATE:           tournament.Date,
						YEAR:           year,
						MONTH:          formatDate(tournament.Date),
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
					saveNote(note, t.Table, cnt)
					cnt++
				}
			}
		}
	}
}
