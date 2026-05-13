package models

import (
	"fmt"
	"judo/internal/lib/utils/note/locresolver"
	"strconv"
	"time"
)

const MinCityRowLen int = 6

// CityDBRow — справочная запись города для сохранения в БД.
type CityDBRow struct {
	ExternalID  int64      `db:"external_id"`
	Name        string     `db:"name"`
	NameRus     *string    `db:"name_rus"`
	NameRusLast *string    `db:"name_rus_last"`
	Description *string    `db:"description"`
	RepublicID  *int64     `db:"republic_id"`
	CreatedAt   *time.Time `db:"created_at"`
	UpdatedAt   *time.Time `db:"updated_at"`

	// Поля для сводной таблицы
	RepublicNameEng string
	RepublicNameRus string
}

func NewCityDBRow(row []string) (CityDBRow, error) {
	if len(row) < MinCityRowLen {
		return CityDBRow{}, fmt.Errorf("недопустимый формат строки: %+v", row)
	}

	extID, err := strconv.ParseInt(row[0], 10, 64)
	if err != nil {
		return CityDBRow{}, fmt.Errorf("недопустимый external_id %q: %w", row[0], err)
	}

	return CityDBRow{
		ExternalID:      extID,
		Name:            row[1],
		NameRus:         new(row[2]),
		NameRusLast:     new(row[3]),
		RepublicNameRus: row[4],
		RepublicNameEng: row[5],
	}, nil
}

func ToCityInput(cities []CityDBRow) []locresolver.CityInput {
	names := make([]locresolver.CityInput, len(cities))
	for i, city := range cities {
		names[i] = locresolver.CityInput{
			City:     city.Name,
			Republic: city.RepublicNameEng,
		}
	}
	return names
}
