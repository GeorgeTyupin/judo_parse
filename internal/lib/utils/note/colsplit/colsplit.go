package colsplit


type ColumnSplitter struct {
	countries  map[string]struct{}
	cities     map[string]struct{}
	sportClubs map[string]struct{}
}

func NewColumnSplitter(countries []string, cities []string, sportClubs []string) (*ColumnSplitter, error) {
	countriesSet := make(map[string]struct{})
	citiesSet := make(map[string]struct{})
	sportClubsSet := make(map[string]struct{})

	for _, country := range countries {
		countriesSet[country] = struct{}{}
	}

	for _, city := range cities {
		citiesSet[city] = struct{}{}
	}

	for _, sportClub := range sportClubs {
		sportClubsSet[sportClub] = struct{}{}
	}

	splitter := &ColumnSplitter{
		countries:  countriesSet,
		cities:     citiesSet,
		sportClubs: sportClubsSet,
	}
	return splitter, nil
}

func (c *ColumnSplitter) isSportClub(value string) bool {
	if _, ok := c.sportClubs[value]; ok {
		return true
	}
	return false
}

func (c *ColumnSplitter) isCountry(value string) bool {
	if _, ok := c.countries[value]; ok {
		return true
	}
	return false
}

func (c *ColumnSplitter) isCity(value string) bool {
	if _, ok := c.cities[value]; ok {
		return true
	}
	return false
}

func (c *ColumnSplitter) SplitCountryAndClub(country, so string) (resolvedCountry, resolvedCity, resolvedClub string) {
	if country != "" {
		switch {
		case c.isCountry(country):
			resolvedCountry = country
		case c.isSportClub(country):
			resolvedClub = country
		case c.isCity(country):
			resolvedCity = country
		default:
			resolvedCountry = country
		}
	}
	if so != "" {
		switch {
		case c.isCountry(so):
			resolvedCountry = so
		case c.isSportClub(so):
			resolvedClub = so
		default:
			resolvedClub = so
		}
	}
	return
}
