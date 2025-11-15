package duplicates

import (
	"judo/internal/lib/replacers"
	"judo/internal/models"
	"strings"
	"sync"

	dupio "judo/internal/io/duplicates"
	parseio "judo/internal/io/parse"

	"github.com/xuri/excelize/v2"
)

var uniqueJudoka = []*models.Judoka{}

// SearchDuplicates ищет все типы дублей в данных
func SearchDuplicates(wg *sync.WaitGroup, data models.ExelSheet) {
	defer wg.Done()
	dupTable := dupio.InitTable("Дубли")
	searchDuplicatesLoop(data, dupTable.Table)
}

func searchDuplicatesLoop(data models.ExelSheet, table *excelize.File) {
	var cnt = 0
	for _, sheet := range data {
		for _, tournament := range sheet {
			for categoryName, category := range tournament.WeightCategories {
				for _, man := range category {
					descParts := strings.Split(tournament.Description, "—")
					tourType := strings.TrimSpace(parseio.SafeGet(descParts, 0))
					tourPlace := strings.TrimSpace(parseio.SafeGet(descParts, 1))

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

					note := dupio.DuplicateNote{
						Note: models.Note{
							TOURNAMENT:     tournament.Name,
							TOUR_TYPE:      tourType,
							TOUR_PLACE:     tourPlace,
							TOUR_CITY:      tourCity,
							TOUR_COUNTRY:   tourCountry,
							TOUR_CITY_LAST: replacers.NormalizeCityName(tourCity),
							DATE:           tournament.Date,
							YEAR:           year,
							MONTH:          parseio.FormatDate(tournament.Date),
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
						},
						DuplicateType: "Тип1",
					}

					note.SaveDuplicateNote(table, cnt)
					cnt++
				}
			}
		}
	}
}
