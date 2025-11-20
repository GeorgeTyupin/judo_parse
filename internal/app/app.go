package app

import (
	"fmt"

	dupio "judo/internal/io/excel/duplicates"
	parseio "judo/internal/io/excel/parse"
	jsonio "judo/internal/io/json"
	"judo/internal/services/duplicates"
	"judo/internal/services/parse"
	"judo/internal/services/pivot"

	"sync"
	"time"
)

const (
	ExcelWriterKey = "ExcelWriter"
	DupWriterKey   = "DupWriter"
	JsonWriterKey  = "JsonWriter"
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
	excelWriter := parseio.NewWriter("Сводная таблица")
	pivotService.RegisterWriter(ExcelWriterKey, excelWriter)
	defer pivotService.Writers[ExcelWriterKey].SaveFile()

	wg.Add(1)
	go func() {
		defer wg.Done()

		notes := pivotService.ProcessData(data)
		pivotService.Writers[ExcelWriterKey].Write(notes)
	}()

	if app.createJSON {
		jsonWriter := jsonio.NewWriter("Data")
		pivotService.RegisterWriter(JsonWriterKey, jsonWriter)
		defer pivotService.Writers[JsonWriterKey].SaveFile()

		wg.Add(1)
		go func() {
			defer wg.Done()

			pivotService.Writers[JsonWriterKey].Write(data)
		}()
	}

	if app.isDuplicates {
		duplicateService := duplicates.NewDuplicateService()
		dupWriter := dupio.NewWriter("Дубли")
		duplicateService.RegisterWriter(DupWriterKey, dupWriter)
		defer duplicateService.Writers[DupWriterKey].SaveFile()

		wg.Add(1)
		go func() {
			defer wg.Done()
			dupNotes := duplicateService.ProcessData(data)
			duplicateService.Writers[DupWriterKey].Write(dupNotes)
		}()
	}
	wg.Wait()

	fmt.Println("Выполнение заняло ", time.Since(start))
	return nil
}
