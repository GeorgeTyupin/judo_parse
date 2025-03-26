package main

import (
	"fmt"
	"time"
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
	fmt.Scanln(&choise)
	switch choise {
	case "1":
		FILE = "USSR_tours"
	case "2":
		FILE = "INT_tours"
	default:
		fmt.Println("Ошибка ввода, попробуйте еще раз. Программа завершится сама через:")
		for i := range 3 {
			go func() {
				time.Sleep(time.Second * time.Duration(i))
				fmt.Println(3-i, " сек")
			}()
		}
		time.Sleep(time.Millisecond * 3001)
		return
	}

	ExelToJson()
	JsonToExel()
}
