package models

type Tournament struct {
	Name             string               `json:"name"`
	Description      string               `json:"description"`
	Date             string               `json:"date"`
	Gender           string               `json:"gender"`
	WeightCategories map[string][]*Judoka `json:"weight_categories"`
}
