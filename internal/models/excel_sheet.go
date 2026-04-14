package models

type ExcelSheet map[string][]Tournament

var Headers = []string{
	"TOURNAMENT", "TOUR_TYPE", "TOUR_PLACE", "TOUR_CITY", "TOUR_COUNTRY",
	"TOUR_CITY_LAST", "DATE", "YEAR", "MONTH", "GENDER", "GENDER_FULL", "AGE", "B_RANGE", "WEIGHT_CATEGORY",
	"WC", "RANK", "NAME", "FIRSTNAME", "JUDOKA", "NAME_RUS", "FIRSTNAME_RUS",
	"JUDOKA_RUS", "COUNTRY", "COUNTRY_LAST", "SO", "NAME_COMP",
}
