package main

import (
	"fmt"
	"os"
	"time"
)

var FILE string

func main() {
	os.Remove("Сводная таблица.xlsx")
	os.Remove("USSR_tours.json")

	var choise string

	fmt.Println("Выбор исходного файла. Введи:\n1, если исходный USSR_tours\n2, если исходный INT_tours")
	// fmt.Scanln(&choise)
	choise = "1"

	switch choise {
	case "1":
		FILE = "USSR_tours"
	case "2":
		FILE = "INT_tours"
	default:
		fmt.Println("Ошибка ввода, попробуйте еще раз.")
	}

	start := time.Now()

	ExelToJson()
	JsonToExel()

	fmt.Println("Выполнение заняло ", time.Since(start))
}
