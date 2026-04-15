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
	"judo/internal/models"
	"judo/internal/repository"
	dbpool "judo/internal/repository/pool"
	"judo/internal/services/duplicates"
	"judo/internal/services/export"
	"judo/internal/services/parse"
	"judo/internal/services/pivot"

	"sync"
	"time"
)

type RunOptions struct {
	isDuplicates    bool
	isServerMigrate bool
	isLocalMigrate  bool
	isCreateJSON    bool
}

func NewRunOptions(isDuplicates, isServerMigrate, isLocalMigrate, isCreateJSON bool) RunOptions {
	return RunOptions{
		isDuplicates:    isDuplicates,
		isServerMigrate: isServerMigrate,
		isLocalMigrate:  isLocalMigrate,
		isCreateJSON:    isCreateJSON,
	}
}

type App struct {
	files []string
	cfg   config.Config
	opt   RunOptions
}

func NewApp(cfg config.Config, options RunOptions, files []string) *App {
	return &App{
		cfg:   cfg,
		files: files,
		opt:   options,
	}
}

func (app *App) Run() error {
	wg := &sync.WaitGroup{}
	start := time.Now()

	parseService, err := parse.NewParseService(app.files)
	if err != nil {
		return fmt.Errorf("ошибка инициализации ParseService - %w", err)
	}

	tournaments, err := parseService.ParseTournaments()
	if err != nil {
		return fmt.Errorf("ошибка парсинга файлов - %w", err)
	}

	excelWriter := parseio.NewWriter("Сводная таблица")
	defer excelWriter.SaveFile()

	wg.Go(func() {
		notes := pivot.ProcessData(tournaments)
		excelWriter.Write(notes)
	})

	if app.opt.isCreateJSON {
		jsonWriter := jsonio.NewWriter("Data")
		defer jsonWriter.SaveFile()

		wg.Go(func() {
			jsonWriter.Write(tournaments)
		})
	}

	if app.opt.isDuplicates {
		dupWriter := dupio.NewWriter("Дубли")
		defer dupWriter.SaveFile()

		wg.Go(func() {
			dupNotes := duplicates.ProcessData(tournaments)
			dupWriter.Write(dupNotes)
		})
	}

	if app.opt.isLocalMigrate {
		fmt.Println("Запись в локальную БД")

		if err := writeToDB(wg, tournaments, nil, app.cfg); err != nil {
			return fmt.Errorf("ошибка записи в БД - %w", err)
		}
	} else if app.opt.isServerMigrate {
		fmt.Println("Запись в удаленную БД")

		sshClient, err := ssh.NewSSHClient(app.cfg)
		if err != nil {
			return fmt.Errorf("ошибка инициализации SSHClient - %w", err)
		}
		defer sshClient.Close()

		if err := writeToDB(wg, tournaments, sshClient.ConnectRemoteDB, app.cfg); err != nil {
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

	dbWriter, err := dbpool.New(dbInitCtx, connString, dialFunc)
	if err != nil {
		return fmt.Errorf("ошибка инициализации DBWriter - %w", err)
	}

	pgRepo := repository.NewTournamentRepository(dbWriter)
	exportService, err := export.NewExportService(pgRepo)
	if err != nil {
		dbWriter.Close()
		return fmt.Errorf("возникла ошибка при создании сервиса для экспорта данных: %w", err)
	}

	wg.Go(func() {
		defer dbWriter.Close()
		exportService.SaveTournaments(context.Background(), data)
	})

	return nil
}
