package main

import (
	"fmt"
	"os"
	"sync"
	"time"

	judioio "judo/internal/io"
	"judo/internal/parse"
)

var files = make([]string, 0, 2)

type PTable interface {
	SetHeader()
	SaveTable()
}

func main() {
	os.Remove("Сводная таблица.xlsx")
	os.Remove("USSR_tours.json")

	var choise string

	fmt.Println("Выбор исходного файла. Введи:\n1, если исходный USSR_tours\n2, если исходный INT_tours\n3, если оба")
	// fmt.Scanln(&choise)
	choise = "1"

	switch choise {
	case "1":
		files = append(files, "USSR_tours")
	case "2":
		files = append(files, "INT_tours")
	case "3":
		files = append(files, "USSR_tours")
		files = append(files, "INT_tours")
	default:
		fmt.Println("Ошибка ввода, попробуйте еще раз.")
	}

	wg := &sync.WaitGroup{}
	start := time.Now()

	table := judioio.InitTable("Сводная таблица")
	defer table.SaveTable()

	for _, file := range files {
		data, err := parse.RenderExel(file)
		if err != nil {
			panic(fmt.Sprintf("Ошибка чтения исходного файла %s - %v", file, err))
		}

		wg.Add(2)
		go judioio.ToExcel(wg, data, table.Table)
		go judioio.ToJson(wg, data, file)
		wg.Wait()
	}

	fmt.Println("Выполнение заняло ", time.Since(start))
}
