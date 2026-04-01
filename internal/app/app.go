package app

import (
	"context"
	"fmt"
	"net"

	"judo/internal/client/ssh"
	"judo/internal/config"
	dupio "judo/internal/io/excel/duplicates"
	parseio "judo/internal/io/excel/parse"
	jsonio "judo/internal/io/json"
	"judo/internal/io/sql"
	"judo/internal/models"
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

	wg.Go(func() {
		notes := pivot.ProcessData(data)
		excelWriter.Write(notes)
	})

	if app.cfg.CreateJSON {
		jsonWriter := jsonio.NewWriter("Data")
		defer jsonWriter.SaveFile()

		wg.Go(func() {
			jsonWriter.Write(data)
		})
	}

	if app.isDuplicates {
		dupWriter := dupio.NewWriter("Дубли")
		defer dupWriter.SaveFile()

		wg.Go(func() {
			dupNotes := duplicates.ProcessData(data)
			dupWriter.Write(dupNotes)
		})
	}

	if app.isLocalMigrate {
		fmt.Println("Запись в локальную БД")

		if err := writeToDB(wg, data, nil, app.cfg); err != nil {
			return fmt.Errorf("ошибка записи в БД - %w", err)
		}
	} else if app.isServerMigrate {
		fmt.Println("Запись в удаленную БД")

		sshClient, err := ssh.NewSSHClient(app.cfg)
		if err != nil {
			return fmt.Errorf("ошибка инициализации SSHClient - %w", err)
		}
		defer sshClient.Close()

		if err := writeToDB(wg, data, sshClient.ConnectRemoteDB, app.cfg); err != nil {
			return fmt.Errorf("ошибка записи в БД - %w", err)
		}
	}

	wg.Wait()

	fmt.Println("Выполнение заняло ", time.Since(start))
	return nil
}

func writeToDB(
	wg *sync.WaitGroup,
	data models.ExcelSheet,
	dialFunc func(ctx context.Context, network, addr string) (net.Conn, error),
	cfg config.Config,
) error {
	dbInitCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	connString := cfg.Database.GetConnString()

	dbWriter, err := sql.NewDBWriter(dbInitCtx, connString, dialFunc)
	if err != nil {
		return fmt.Errorf("ошибка инициализации DBWriter - %w", err)
	}

	pgRepo := repository.NewTournamentRepository(dbWriter)
	exportService, err := export.NewExportService(pgRepo, data)
	if err != nil {
		dbWriter.Close()
		return fmt.Errorf("возникла ошибка при создании сервиса для экспорта данных: %w", err)
	}

	wg.Go(func() {
		defer dbWriter.Close()
		exportService.ProcessTournament(context.Background())
	})

	return nil
}
