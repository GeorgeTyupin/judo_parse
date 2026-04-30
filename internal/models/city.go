package models

import (
	"fmt"
	"time"
)

const MinCityRowLen int = 5

// CityDBRow — справочная запись города для сохранения в БД.
type CityDBRow struct {
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

	return CityDBRow{
		Name:            row[0],
		NameRus:         new(row[1]),
		NameRusLast:     new(row[2]),
		RepublicNameRus: row[3],
		RepublicNameEng: row[4],
	}, nil
}

func GetCityNames(cities []CityDBRow) []string {
	names := make([]string, len(cities))
	for i, city := range cities {
		names[i] = city.Name
	}
	return names
}
