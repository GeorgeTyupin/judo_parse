package birthdate_resolver

import (
	"judo/internal/models"
	"strings"
	"time"
)

type BirthDateResolver map[string]time.Time

func normalizeName(firstName, lastName string) string {
	return strings.ToLower(strings.TrimSpace(lastName) + " " + strings.TrimSpace(firstName))
}

func NewBirthDateResolver(judokas []models.JudokaDBRow) BirthDateResolver {
	resolver := make(BirthDateResolver)
	for _, judoka := range judokas {
		if judoka.BirthDate == nil {
			continue
		}
		birthDate, err := time.Parse(time.DateOnly, *judoka.BirthDate)
		if err != nil {
			continue
		}
		resolver[normalizeName(judoka.FirstName, judoka.LastName)] = birthDate
	}
	return resolver
}

func (r BirthDateResolver) Resolve(firstName, lastName string) *time.Time {
	birthDate, found := r[normalizeName(firstName, lastName)]
	if !found {
		return nil
	}
	return new(birthDate)
}
