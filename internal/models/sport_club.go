package models

import (
	"fmt"
	"time"
)

const MinSportClubRowLen int = 2

// SportClubDBRow — справочная запись спортивного клуба для сохранения в БД.
type SportClubDBRow struct {
	Name        string     `db:"name"`
	NameRus     *string    `db:"name_rus"`
	FullName    *string    `db:"full_name"`
	Founded     *string    `db:"founded"`
	CityID      *int32     `db:"city_id"`
	Region      *string    `db:"region"`
	HeadCoach   *string    `db:"head_coach"`
	Description *string    `db:"description"`
	CreatedAt   *time.Time `db:"created_at"`
	UpdatedAt   *time.Time `db:"updated_at"`
}

func NewSportClubDBRow(row []string) (SportClubDBRow, error) {
	if len(row) < MinSportClubRowLen {
		return SportClubDBRow{}, fmt.Errorf("недопустимый формат строки: %+v", row)
	}

	return SportClubDBRow{
		Name:    row[0],
		NameRus: new(row[1]),
	}, nil
}

func GetSportClubNames(sportClubs []SportClubDBRow) []string {
	sportClubNames := make([]string, 0, len(sportClubs))
	for _, sc := range sportClubs {
		sportClubNames = append(sportClubNames, sc.Name)
	}
	return sportClubNames
}
