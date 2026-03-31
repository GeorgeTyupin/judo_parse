package app

import (
	"context"
	"fmt"

	"judo/internal/client/ssh"
	"judo/internal/config"
	dupio "judo/internal/io/excel/duplicates"
	parseio "judo/internal/io/excel/parse"
	jsonio "judo/internal/io/json"
	"judo/internal/io/sql"
	"judo/internal/repository"
	"judo/internal/services/duplicates"
	"judo/internal/services/export"
	"judo/internal/services/parse"
	"judo/internal/services/pivot"

	"sync"
	"time"
)

type App struct {
	files           []string
	isDuplicates    bool
	isServerMigrate bool
	isLocalMigrate  bool
	cfg             config.Config
}

func NewApp(cfg config.Config, files []string, isDuplicates, isServerMigrate bool) *App {
	isLocalMigrate := false
	if cfg.IsDev {
		isLocalMigrate = true
	}

	return &App{
		cfg:             cfg,
		files:           files,
		isDuplicates:    isDuplicates,
		isServerMigrate: isServerMigrate,
		isLocalMigrate:  isLocalMigrate,
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

	if app.isLocalMigrate {
		dbInitCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		connString := app.cfg.Database.GetConnString()

		dbWriter := sql.NewDBWriter(dbInitCtx, connString)
		defer dbWriter.Close()

		pgRepo := repository.NewTournamentRepository(dbWriter)
		exportService, err := export.NewExportService(pgRepo, data)
		if err != nil {
			return fmt.Errorf("возникла ошибка при создании сервиса для экспорта данных: %w", err)
		}

		dbReqCtx := context.Background()

		wg.Add(1)
		go func() {
			defer wg.Done()
			exportService.ProcessTournament(dbReqCtx)
		}()
	} else if app.isServerMigrate {
		sshClient, err := ssh.NewSSHClient(app.cfg.SSH)
		if err != nil {
			return fmt.Errorf("ошибка инициализации SSHClient - %w", err)
		}
		defer sshClient.Close()

		wg.Add(1)
		go func() {
			defer wg.Done()
			sshClient.MigrateOnServer()
		}()
	}

	wg.Wait()

	fmt.Println("Выполнение заняло ", time.Since(start))
	return nil
}
