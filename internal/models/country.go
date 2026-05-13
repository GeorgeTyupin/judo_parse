package models

import (
	"fmt"
	"strconv"
	"time"
)

const MinCountryRowLen int = 4

// CountryDBRow — справочная запись страны для сохранения в БД.
type CountryDBRow struct {
	ExternalID  int64      `db:"external_id"`
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

	extID, err := strconv.ParseInt(row[0], 10, 64)
	if err != nil {
		return CountryDBRow{}, fmt.Errorf("недопустимый external_id %q: %w", row[0], err)
	}

	return CountryDBRow{
		ExternalID: extID,
		ISOCode:    new(row[1]),
		NameRus:    new(row[2]),
		Name:       row[3],
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
