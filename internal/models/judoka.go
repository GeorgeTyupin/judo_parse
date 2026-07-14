package models

import (
	"fmt"
	"judo/internal/lib/utils/note/russifiers"
	"strconv"
	"strings"
	"time"
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

	return JudokaDBRow{
		ExternalID:   extID,
		LastName:     curRow[1],
		FirstName:    curRow[2],
		LastNameRus:  new(curRow[4]),
		FirstNameRus: new(curRow[5]),
	}, nil
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
