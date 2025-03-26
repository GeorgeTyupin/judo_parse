package main

import (
	"fmt"
	"sync"
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

func timer(lastSeconds int, wg *sync.WaitGroup) {

	for i := range lastSeconds {
		wg.Add(1)
		go func() {
			defer wg.Done()
			time.Sleep(time.Second * time.Duration(i))
			fmt.Println(lastSeconds-i, " сек")
		}()
	}

}

func main() {
	var choise string

	fmt.Println("Выбор исходного файла. Введи:\n1, если исходный USSR_tours\n2, если исходный INT_tours")
	fmt.Scanln(&choise)
	lastSeconds := 2
	wg := &sync.WaitGroup{}

	switch choise {
	case "1":
		FILE = "USSR_tours"
		fmt.Println("Программа выполнилась. И автоматически завершится через:")
		timer(lastSeconds, wg)
	case "2":
		FILE = "INT_tours"
		fmt.Println("Программа выполнилась. И автоматически завершится через:")
		timer(lastSeconds, wg)
	default:
		lastSeconds = 3
		fmt.Println("Ошибка ввода, попробуйте еще раз. Программа завершится сама через:")
		timer(lastSeconds, wg)
		return
	}
	wg.Wait()
	ExelToJson()
	JsonToExel()
}
