package app

import (
	"fmt"
	judioio "judo/internal/io"
	"judo/internal/parse"
	"sync"
	"time"
)

type App struct {
	files      []string
	createJSON bool
}

func NewApp(files []string, isDuplicates, createJSON bool) *App {
	return &App{
		files:      files,
		createJSON: createJSON,
	}
}

func (app *App) Run() error {
	wg := &sync.WaitGroup{}
	start := time.Now()

	table := judioio.InitTable("Сводная таблица")
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
			go judioio.ToJson(wg, data, file)
		}
		wg.Wait()
	}

	fmt.Println("Выполнение заняло ", time.Since(start))
	return nil
}
