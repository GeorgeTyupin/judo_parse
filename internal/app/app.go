package app

import (
	"fmt"

	jsonio "judo/internal/io/json"
	"judo/internal/services/duplicates"
	"judo/internal/services/parse"
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

	// Создаем ParseService со всеми файлами
	parseService, err := parse.NewParseService(app.files)
	if err != nil {
		return fmt.Errorf("ошибка инициализации ParseService - %w", err)
	}

	// Парсим все файлы сразу
	data, err := parseService.ParseTournaments()
	if err != nil {
		return fmt.Errorf("ошибка парсинга файлов - %w", err)
	}

	pivotService := pivot.NewPivotService()
	defer pivotService.Writer.SaveTable()

	duplicateService := duplicates.NewDuplicateService()
	defer duplicateService.Writer.SaveTable()

	wg.Add(1)
	go func() {
		defer wg.Done()

		notes := pivotService.ProcessData(data)
		pivotService.Writer.Write(notes)
	}()

	if app.createJSON {
		wg.Add(1)
		go jsonio.ToJson(wg, data, app.files[0])
	}
	if app.isDuplicates {
		wg.Add(1)
		go func() {
			defer wg.Done()
			dupNotes := duplicateService.ProcessData(data)
			duplicateService.Writer.Write(dupNotes)
		}()
	}
	wg.Wait()

	fmt.Println("Выполнение заняло ", time.Since(start))
	return nil
}
