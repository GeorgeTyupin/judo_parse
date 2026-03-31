package models

import (
	noteutils "judo/internal/lib/utils/note"
	"strconv"
	"strings"
)

type Tournament struct {
	Name             string               `json:"name"`
	Description      string               `json:"description"`
	Date             string               `json:"date"`
	Gender           string               `json:"gender"`
	WeightCategories map[string][]*Judoka `json:"weight_categories"`
}

type TournamentDBRow struct {
	Name      string `db:"name"`
	Type      string `db:"type"`
	Place     string `db:"place"`
	Gender    string `db:"gender"`
	Date      string `db:"date"`
	CityID    *int64 `db:"city_id"`
	CountryID *int32 `db:"country_id"`
	Year      int16  `db:"year"`
	Month     int16  `db:"month"`
}

func NewTournamentDBRow(tournament *Tournament) *TournamentDBRow {
	descParts := strings.Split(tournament.Description, "—")
	tourType := strings.TrimSpace(noteutils.SafeGet(descParts, 0))
	tourPlace := strings.TrimSpace(noteutils.SafeGet(descParts, 1))

	year := "-1"
	if len(tournament.Date) >= 4 {
		year = tournament.Date[len(tournament.Date)-4:]
	}
	yearInt, _ := strconv.Atoi(year)

	monthInt, _ := strconv.Atoi(noteutils.FormatDate(tournament.Date))

	return &TournamentDBRow{
		Name:   tournament.Name,
		Type:   tourType,
		Place:  tourPlace,
		Date:   tournament.Date,
		Year:   int16(yearInt),
		Month:  int16(monthInt),
		Gender: tournament.Gender,
	}
}
