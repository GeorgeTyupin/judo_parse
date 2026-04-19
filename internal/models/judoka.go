package models

import (
	"fmt"
	"judo/internal/lib/utils/note/russifiers"
	"strings"
	"time"
)

type Judoka struct {
	Rank      string `json:"RANK"`
	LastName  string `json:"NAME"`
	FirstName string `json:"FIRSTNAME"`
	JUDOKA    string `json:"JUDOKA"`
	Country   string `json:"COUNTRY"`
	SO        string `json:"SO"`

	FirstNameRus *string
	LastNameRus  *string
}

func NewJudoka(curRow []string, lenCurTable int) Judoka {
	athlete := Judoka{
		Rank:      curRow[0],
		LastName:  curRow[1],
		FirstName: curRow[2],
		JUDOKA:    curRow[3],
	}

	if lenCurTable > 5 {
		athlete.SO = strings.Trim(curRow[4], `'`)
		athlete.Country = curRow[5]
	} else if lenCurTable > 4 {
		athlete.Country = curRow[4]
	}

	return athlete
}

type JudokaDBRow struct {
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
	if len(curRow) < 4 {
		return JudokaDBRow{}, fmt.Errorf("недопустимый формат строки: %+v", curRow)
	}

	return JudokaDBRow{
		LastName:       curRow[0],
		FirstName:      curRow[1],
		LastNameRus:    &curRow[2],
		FirstNameRus:   &curRow[3],
		WeightCategory: nil,
		BirthDate:      nil,
		BirthPlace:     nil,
		Gender:         nil,
		CreatedAt:      nil,
		UpdatedAt:      nil,
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
