package models

type Judoka struct {
	Rank      string `json:"RANK"`
	Name      string `json:"NAME"`
	FirstName string `json:"FIRSTNAME"`
	JUDOKA    string `json:"JUDOKA"`
	Country   string `json:"COUNTRY"`
	SO        string `json:"SO"`
}

type Tournament struct {
	Name             string               `json:"name"`
	Description      string               `json:"description"`
	Date             string               `json:"date"`
	Gender           string               `json:"gender"`
	WeightCategories map[string][]*Judoka `json:"weight_categories"`
}

type ExelSheet map[string][]*Tournament

type Note struct {
	TOURNAMENT     string
	TOUR_TYPE      string
	TOUR_PLACE     string
	TOUR_CITY      string
	TOUR_COUNTRY   string
	DATE           string
	YEAR           string
	GENDER         string
	WeightCategory string
	WC             string
	RANK           string
	NAME           string
	FIRSTNAME      string
	JUDOKA         string
	NAME_RUS       string
	FIRSTNAME_RUS  string
	JUDOKA_RUS     string
	COUNTRY        string
	SO             string
}
