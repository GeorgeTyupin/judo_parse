package parseio

import (
	"fmt"
	"strings"
	"sync"

	"judo/internal/lib/replacers"
	"judo/internal/models"

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
	for i, header := range models.Headers {
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

func SafeGet(parts []string, index int) string {
	if index < len(parts) {
		return parts[index]
	}
	return ""
}

func FormatDate(date string) string {
	var result string

	if len(date) < 5 {
		return ""
	}

	if strings.Contains(date, "-") {
		result = strings.TrimSpace(strings.Split(date, "-")[1])
	} else {
		result = date
	}

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
					tourType := strings.TrimSpace(SafeGet(descParts, 0))
					tourPlace := strings.TrimSpace(SafeGet(descParts, 1))

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
						TOUR_CITY_LAST: replacers.NormalizeCityName(tourCity),
						DATE:           tournament.Date,
						YEAR:           year,
						MONTH:          FormatDate(tournament.Date),
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
						COUNTRY_LAST:   replacers.NormalizeCityName(man.Country),
						SO:             man.SO,
					}

					note.SaveNote(t.Table, cnt)
					cnt++
				}
			}
		}
	}
}
