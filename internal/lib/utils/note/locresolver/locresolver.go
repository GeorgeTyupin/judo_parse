package locresolver

type CityInput struct {
	City     string
	Republic string
}

type LocationResolver struct {
	countries  map[string]struct{}
	cities     map[string]string
	sportClubs map[string]struct{}
}

func NewLocationResolver(countries []string, cities []CityInput, sportClubs []string) (*LocationResolver, error) {
	countriesSet := make(map[string]struct{})
	citiesSet := make(map[string]string)
	sportClubsSet := make(map[string]struct{})

	for _, country := range countries {
		countriesSet[country] = struct{}{}
	}

	for _, city := range cities {
		citiesSet[city.City] = city.Republic
	}

	for _, sportClub := range sportClubs {
		sportClubsSet[sportClub] = struct{}{}
	}

	resolver := &LocationResolver{
		countries:  countriesSet,
		cities:     citiesSet,
		sportClubs: sportClubsSet,
	}
	return resolver, nil
}

func (r *LocationResolver) isSportClub(value string) bool {
	if _, ok := r.sportClubs[value]; ok {
		return true
	}
	return false
}

func (r *LocationResolver) isCountry(value string) bool {
	if _, ok := r.countries[value]; ok {
		return true
	}
	return false
}

func (r *LocationResolver) isCity(value string) bool {
	if _, ok := r.cities[value]; ok {
		return true
	}
	return false
}

func (r *LocationResolver) SplitCountryAndClub(country, so string) (resolvedCountry, resolvedCity, resolvedClub string) {
	if country != "" {
		switch {
		case r.isCountry(country):
			resolvedCountry = country
		case r.isSportClub(country):
			resolvedClub = country
		case r.isCity(country):
			resolvedCity = country
		default:
			resolvedCountry = country
		}
	}
	if so != "" {
		switch {
		case r.isCountry(so):
			resolvedCountry = so
		case r.isSportClub(so):
			resolvedClub = so
		default:
			resolvedClub = so
		}
	}
	return
}

func (r *LocationResolver) GetRepublicByCity(city string) string {
	if republic, ok := r.cities[city]; ok {
		return republic
	}
	return ""
}
