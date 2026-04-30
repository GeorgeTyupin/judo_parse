package models

import (
	"fmt"
	"time"
)

const MinCountryRowLen int = 3

// CountryDBRow — справочная запись страны для сохранения в БД.
type CountryDBRow struct {
	Name        string     `db:"name"`
	ISOCode     *string    `db:"iso_code"`
	NameRus     *string    `db:"name_rus"`
	Description *string    `db:"description"`
	CreatedAt   *time.Time `db:"created_at"`
	UpdatedAt   *time.Time `db:"updated_at"`
}

func NewCountryDBRow(row []string) (CountryDBRow, error) {
	if len(row) < MinCountryRowLen {
		return CountryDBRow{}, fmt.Errorf("недопустимый формат строки: %+v", row)
	}

	return CountryDBRow{
		ISOCode: new(row[0]),
		NameRus: new(row[1]),
		Name:    row[2],
	}, nil
}

func GetCountryCodes(countries []CountryDBRow) []string {
	countryCodes := make([]string, 0, len(countries))
	for _, c := range countries {
		if c.ISOCode != nil {
			countryCodes = append(countryCodes, *c.ISOCode)
		}
	}
	return countryCodes
}
