package models

import (
	"fmt"
	"judo/internal/lib/utils/note/russifiers"
	"strconv"
	"strings"
	"time"

	"github.com/xuri/excelize/v2"
)

// Judoka — структура для парсинга данных из турнирных таблиц.
type Judoka struct {
	Rank      string `json:"RANK"`
	LastName  string `json:"NAME"`
	FirstName string `json:"FIRSTNAME"`
	JUDOKA    string `json:"JUDOKA"`
	Country   string `json:"COUNTRY"`
	SO        string `json:"SO"`
}

func NewJudoka(curRow []string, lenCurTable int) Judoka {
	athlete := Judoka{
		Rank:      curRow[0],
		LastName:  curRow[1],
		FirstName: curRow[2],
		JUDOKA:    curRow[3],
	}

	if lenCurTable > 5 {
		athlete.SO = strings.Trim(curRow[4], `'"`)
		athlete.Country = strings.Trim(curRow[5], `'"`)
	} else if lenCurTable > 4 {
		athlete.Country = strings.Trim(curRow[4], `'"`)
	}

	return athlete
}

const MinJudokaRowLen = 6

// judokaIDPrefix — префикс идентификатора дзюдоиста в справочнике (JUD00001).
const judokaIDPrefix = "JUD"

// JudokaDBRow — справочная запись дзюдоиста для сохранения в БД.
type JudokaDBRow struct {
	ExternalID     int64      `db:"external_id"`
	LastName       string     `db:"last_name"`
	FirstName      string     `db:"first_name"`
	LastNameRus    *string    `db:"last_name_rus"`
	FirstNameRus   *string    `db:"first_name_rus"`
	WeightCategory *string    `db:"weight_category"`
	BirthDate      *string    `db:"birth_date"`
	BirthPlace     *string    `db:"birth_place"`
	Gender         *string    `db:"gender"`
	CreatedAt      *time.Time `db:"created_at"`
	UpdatedAt      *time.Time `db:"updated_at"`
}

func NewJudokaDBRow(curRow []string) (JudokaDBRow, error) {
	if len(curRow) < MinJudokaRowLen {
		return JudokaDBRow{}, fmt.Errorf("недопустимый формат строки: %+v", curRow)
	}

	extID, err := strconv.ParseInt(strings.TrimPrefix(curRow[0], judokaIDPrefix), 10, 64)
	if err != nil {
		return JudokaDBRow{}, fmt.Errorf("недопустимый external_id %q: %w", curRow[0], err)
	}

	var birthDate *string
	if len(curRow) > 6 {
		birthDate = parseJudokaBirthDate(curRow[6])
	}

	return JudokaDBRow{
		ExternalID:   extID,
		LastName:     curRow[1],
		FirstName:    curRow[2],
		LastNameRus:  new(curRow[4]),
		FirstNameRus: new(curRow[5]),
		BirthDate:    birthDate,
	}, nil
}

// parseJudokaBirthDate возвращает полную дату рождения в формате YYYY-MM-DD.
func parseJudokaBirthDate(value string) *string {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil
	}

	var date time.Time
	var err error
	if serial, parseErr := strconv.ParseFloat(value, 64); parseErr == nil {
		if serial <= 0 {
			return nil
		}
		date, err = excelize.ExcelDateToTime(serial, false)
	} else {
		// В справочнике встречаются периоды жизни: DOB — дата смерти.
		birthDate := strings.FieldsFunc(value, func(r rune) bool {
			return r == '-' || r == '–' || r == '—'
		})
		if len(birthDate) == 0 {
			return nil
		}
		date, err = time.Parse("02.01.2006", strings.ReplaceAll(birthDate[0], " ", ""))
	}
	if err != nil {
		return nil
	}

	return new(date.Format(time.DateOnly))
}

func JudokaRowsToNames(judokaDBRows []JudokaDBRow) []russifiers.JudokaName {
	judokaNames := make([]russifiers.JudokaName, 0)
	for _, judokaDBRow := range judokaDBRows {
		name := russifiers.NewJudokaName(
			judokaDBRow.FirstName,
			judokaDBRow.LastName,
			judokaDBRow.FirstNameRus,
			judokaDBRow.LastNameRus,
		)
		judokaNames = append(judokaNames, name)
	}
	return judokaNames
}
