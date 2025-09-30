package main

import (
	"fmt"
)

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

var FILE string

func main() {
	var choise string

	fmt.Println("Выбор исходного файла. Введи:\n1, если исходный USSR_tours\n2, если исходный INT_tours")
	// fmt.Scanln(&choise)
	choise = "1"

	switch choise {
	case "1":
		FILE = "USSR_tours"
		fmt.Println("Программа выполнилась.")
	case "2":
		FILE = "INT_tours"
		fmt.Println("Программа выполнилась.")
	default:
		fmt.Println("Ошибка ввода, попробуйте еще раз.")
	}
	ExelToJson()
	JsonToExel()
}
