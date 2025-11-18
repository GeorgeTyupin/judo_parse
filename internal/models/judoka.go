package models

import "strings"

type Judoka struct {
	Rank      string `json:"RANK"`
	Name      string `json:"NAME"`
	FirstName string `json:"FIRSTNAME"`
	JUDOKA    string `json:"JUDOKA"`
	Country   string `json:"COUNTRY"`
	SO        string `json:"SO"`
}

func NewJudoka(curRow []string, lenCurTable int) *Judoka {
	athlete := Judoka{
		Rank:      curRow[0],
		Name:      curRow[1],
		FirstName: curRow[2],
		JUDOKA:    curRow[3],
	}

	if lenCurTable > 5 {
		athlete.SO = strings.Trim(curRow[4], `'`)
		athlete.Country = curRow[5]
	} else if lenCurTable > 4 {
		athlete.Country = curRow[4]
	}

	return &athlete
}
