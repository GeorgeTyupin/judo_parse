package app

import (
	"fmt"

	jsonio "judo/internal/io/json"
	"judo/internal/parse"
	"judo/internal/services/duplicates"
	"judo/internal/services/pivot"

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

	pivotService := pivot.NewPivotService()
	defer pivotService.File.SaveTable()

	duplicateService := duplicates.NewDuplicateService()
	defer duplicateService.File.SaveTable()

	for _, file := range app.files {
		data, err := parse.RenderExel(file)
		if err != nil {
			return fmt.Errorf("ошибка чтения исходного файла %s - %w", file, err)
		}

		wg.Add(1)
		go func() {
			defer wg.Done()

			notes := pivotService.ProcessData(data)
			pivotService.File.Write(notes)
		}()

		if app.createJSON {
			wg.Add(1)
			go jsonio.ToJson(wg, data, file)
		}
		if app.isDuplicates {
			wg.Add(1)
			go func() {
				dupNotes := duplicateService.ProcessData(data)
				duplicateService.File.Write(dupNotes)
			}()
		}
		wg.Wait()
	}

	fmt.Println("Выполнение заняло ", time.Since(start))
	return nil
}
