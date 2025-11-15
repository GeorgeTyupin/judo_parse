package models

import (
	"fmt"

	"github.com/xuri/excelize/v2"
)

type Judoka struct {
	Rank      string `json:"RANK"`
	Name      string `json:"NAME"`
	FirstName string `json:"FIRSTNAME"`
	JUDOKA    string `json:"JUDOKA"`
	Country   string `json:"COUNTRY"`
	SO        string `json:"SO"`
}

type Tournament struct {
	Name             string               `json:"name"`
	Description      string               `json:"description"`
	Date             string               `json:"date"`
	Gender           string               `json:"gender"`
	WeightCategories map[string][]*Judoka `json:"weight_categories"`
}

type ExelSheet map[string][]*Tournament

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
}

// Headers defines the Excel column headers matching Note struct fields
var Headers = []string{
	"TOURNAMENT", "TOUR_TYPE", "TOUR_PLACE", "TOUR_CITY", "TOUR_COUNTRY",
	"TOUR_CITY_LAST", "DATE", "YEAR", "MONTH", "GENDER", "WEIGHT_CATEGORY",
	"WC", "RANK", "NAME", "FIRSTNAME", "JUDOKA", "NAME_RUS", "FIRSTNAME_RUS",
	"JUDOKA_RUS", "COUNTRY", "COUNTRY_LAST", "SO",
}
