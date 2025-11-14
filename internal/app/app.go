package app

import (
	"fmt"

	"judo/internal/duplicates"
	parseio "judo/internal/io/parse"
	"judo/internal/parse"

	"sync"
	"time"
)

type App struct {
	files        []string
	isDuplicates bool
	createJSON   bool
}

func NewApp(files []string, isDuplicates, createJSON bool) *App {
	return &App{
		files:        files,
		isDuplicates: isDuplicates,
		createJSON:   createJSON,
	}
}

func (app *App) Run() error {
	wg := &sync.WaitGroup{}
	start := time.Now()

	table := parseio.InitTable("Сводная таблица")
	defer table.SaveTable()

	for _, file := range app.files {
		data, err := parse.RenderExel(file)
		if err != nil {
			return fmt.Errorf("ошибка чтения исходного файла %s - %w", file, err)
		}

		wg.Add(1)
		go table.ToExcel(wg, data)
		if app.createJSON {
			wg.Add(1)
			go parseio.ToJson(wg, data, file)
		}
		if app.isDuplicates {
			wg.Add(1)
			go duplicates.SearchDuplicates(wg, data)
		}
		wg.Wait()
	}

	fmt.Println("Выполнение заняло ", time.Since(start))
	return nil
}
