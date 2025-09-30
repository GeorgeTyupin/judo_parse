package main

type Judoka struct {
	Rank      string `json:"RANK"`
	Name      string `json:"NAME"`
	FirstName string `json:"FIRSTNAME"`
	JUDOKA    string `json:"JUDOKA"`
	Country   string `json:"COUNTRY"`
}

type Tournament struct {
	Name             string              `json:"name"`
	Description      string              `json:"description"`
	Date             string              `json:"date"`
	Gender           string              `json:"gender"`
	WeightCategories map[string][]Judoka `json:"weight_categories"`
}

type ExelSheet map[string][]Tournament
