package main

import (
	"fmt"
	"os"
	"time"

	"judo/internal/parse"
)

var FILES = make([]string, 0, 2)

func main() {
	os.Remove("Сводная таблица.xlsx")
	os.Remove("USSR_tours.json")

	var choise string

	fmt.Println("Выбор исходного файла. Введи:\n1, если исходный USSR_tours\n2, если исходный INT_tours\n3, если оба")
	// fmt.Scanln(&choise)
	choise = "1"

	switch choise {
	case "1":
		FILES = append(FILES, "USSR_tours")
	case "2":
		FILES = append(FILES, "INT_tours")
	case "3":
		FILES = append(FILES, "USSR_tours")
		FILES = append(FILES, "INT_tours")
	default:
		fmt.Println("Ошибка ввода, попробуйте еще раз.")
	}

	start := time.Now()

	parse.ExelToJson(FILES)
	parse.JsonToExel(FILES)

	fmt.Println("Выполнение заняло ", time.Since(start))
}
