package models

import (
	"fmt"
	"judo/internal/lib/replacers"
	noteutils "judo/internal/lib/utils/note"
	"strings"

	"github.com/xuri/excelize/v2"
)

type Note struct {
	TOURNAMENT     string
	TOUR_TYPE      string
	TOUR_PLACE     string
	TOUR_CITY      string
	TOUR_CITY_LAST string
	TOUR_COUNTRY   string
	DATE           string
	YEAR           string
	MONTH          string
	GENDER         string
	WeightCategory string
	WC             string
	RANK           string
	NAME           string
	FIRSTNAME      string
	JUDOKA         string
	NAME_RUS       string
	FIRSTNAME_RUS  string
	JUDOKA_RUS     string
	COUNTRY        string
	COUNTRY_LAST   string
	SO             string
	NAME_COMP      string
}

func NewNote(tournament *Tournament, man *Judoka, categoryName string) *Note {
	descParts := strings.Split(tournament.Description, "—")
	tourType := strings.TrimSpace(noteutils.SafeGet(descParts, 0))
	tourPlace := strings.TrimSpace(noteutils.SafeGet(descParts, 1))

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

	// Определяем значение Name_Comp
	nameComp := "P" // По умолчанию P
	if len(man.FirstName) > 2 && !strings.Contains(man.FirstName, ".") {
		nameComp = "F"
	}

	note := Note{
		TOURNAMENT:     tournament.Name,
		TOUR_TYPE:      tourType,
		TOUR_PLACE:     tourPlace,
		TOUR_CITY:      tourCity,
		TOUR_COUNTRY:   tourCountry,
		TOUR_CITY_LAST: replacers.NormalizeCityName(tourCity),
		DATE:           tournament.Date,
		YEAR:           year,
		MONTH:          noteutils.FormatDate(tournament.Date),
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
		NAME_COMP:      nameComp,
	}

	return &note
}

func (note *Note) SaveNote(table *excelize.File, counter int) {
	rowNum := counter + 2

	table.SetCellValue("Sheet1", fmt.Sprintf("A%d", rowNum), note.TOURNAMENT)
	table.SetCellValue("Sheet1", fmt.Sprintf("B%d", rowNum), note.TOUR_TYPE)
	table.SetCellValue("Sheet1", fmt.Sprintf("C%d", rowNum), note.TOUR_PLACE)
	table.SetCellValue("Sheet1", fmt.Sprintf("D%d", rowNum), note.TOUR_CITY)
	table.SetCellValue("Sheet1", fmt.Sprintf("E%d", rowNum), note.TOUR_COUNTRY)
	table.SetCellValue("Sheet1", fmt.Sprintf("F%d", rowNum), note.TOUR_CITY_LAST)
	table.SetCellValue("Sheet1", fmt.Sprintf("G%d", rowNum), note.DATE)
	table.SetCellValue("Sheet1", fmt.Sprintf("H%d", rowNum), note.YEAR)
	table.SetCellValue("Sheet1", fmt.Sprintf("I%d", rowNum), note.MONTH)
	table.SetCellValue("Sheet1", fmt.Sprintf("J%d", rowNum), note.GENDER)
	table.SetCellValue("Sheet1", fmt.Sprintf("K%d", rowNum), note.WeightCategory)
	table.SetCellValue("Sheet1", fmt.Sprintf("L%d", rowNum), note.WC)
	table.SetCellValue("Sheet1", fmt.Sprintf("M%d", rowNum), note.RANK)
	table.SetCellValue("Sheet1", fmt.Sprintf("N%d", rowNum), note.NAME)
	table.SetCellValue("Sheet1", fmt.Sprintf("O%d", rowNum), note.FIRSTNAME)
	table.SetCellValue("Sheet1", fmt.Sprintf("P%d", rowNum), note.JUDOKA)
	table.SetCellValue("Sheet1", fmt.Sprintf("Q%d", rowNum), note.NAME_RUS)
	table.SetCellValue("Sheet1", fmt.Sprintf("R%d", rowNum), note.FIRSTNAME_RUS)
	table.SetCellValue("Sheet1", fmt.Sprintf("S%d", rowNum), note.JUDOKA_RUS)
	table.SetCellValue("Sheet1", fmt.Sprintf("T%d", rowNum), note.COUNTRY)
	table.SetCellValue("Sheet1", fmt.Sprintf("U%d", rowNum), note.COUNTRY_LAST)
	table.SetCellValue("Sheet1", fmt.Sprintf("V%d", rowNum), note.SO)
	table.SetCellValue("Sheet1", fmt.Sprintf("W%d", rowNum), note.NAME_COMP)
}
