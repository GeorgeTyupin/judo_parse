package app

import (
	"context"
	"fmt"

	"judo/internal/config"
	dupio "judo/internal/io/excel/duplicates"
	parseio "judo/internal/io/excel/parse"
	jsonio "judo/internal/io/json"
	"judo/internal/io/sql"
	"judo/internal/services/duplicates"
	"judo/internal/services/export"
	"judo/internal/services/parse"
	"judo/internal/services/pivot"

	"sync"
	"time"
)

type App struct {
	files        []string
	isDuplicates bool
	cfg          *config.Config
}

func NewApp(cfg *config.Config, files []string, isDuplicates bool) *App {
	return &App{
		cfg:          cfg,
		files:        files,
		isDuplicates: isDuplicates,
	}
}

func (app *App) Run() error {
	wg := &sync.WaitGroup{}
	start := time.Now()

	parseService, err := parse.NewParseService(app.files)
	if err != nil {
		return fmt.Errorf("ошибка инициализации ParseService - %w", err)
	}

	data, err := parseService.ParseTournaments()
	if err != nil {
		return fmt.Errorf("ошибка парсинга файлов - %w", err)
	}

	excelWriter := parseio.NewWriter("Сводная таблица")
	defer excelWriter.SaveFile()

	wg.Add(1)
	go func() {
		defer wg.Done()
		notes := pivot.ProcessData(data)
		excelWriter.Write(notes)
	}()

	if app.cfg.CreateJSON {
		jsonWriter := jsonio.NewWriter("Data")
		defer jsonWriter.SaveFile()

		wg.Add(1)
		go func() {
			defer wg.Done()
			jsonWriter.Write(data)
		}()
	}

	if app.isDuplicates {
		dupWriter := dupio.NewWriter("Дубли")
		defer dupWriter.SaveFile()

		wg.Add(1)
		go func() {
			defer wg.Done()
			dupNotes := duplicates.ProcessData(data)
			dupWriter.Write(dupNotes)
		}()
	}

	if app.cfg.IsDev {
		repoInitCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		connString := app.cfg.Database.GetConnString()

		pgWriter := sql.NewDBWriter(repoInitCtx, connString)
		exportService, err := export.NewExportService(pgWriter, data)
		if err != nil {
			return fmt.Errorf("возникла ошибка при создании сервиса для экспорта данных: %w", err)
		}

		dbReqCtx := context.Background()
		exportService.ProcessTournament(dbReqCtx)
	}

	wg.Wait()

	fmt.Println("Выполнение заняло ", time.Since(start))
	return nil
}
