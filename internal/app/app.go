package app

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"slices"

	"judo/internal/client/ssh"
	"judo/internal/config"
	dupio "judo/internal/io/excel/duplicates"
	judokaio "judo/internal/io/excel/judoka_parse"
	parseio "judo/internal/io/excel/parse"
	jsonio "judo/internal/io/json"
	filesutils "judo/internal/lib/utils/files"
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

const (
	DataTargetTournaments string = "tournaments"
	DataTargetJudokas     string = "judokas"

	JudokasFileName string = "#JUDOKA.xlsx"
)

type RunOptions struct {
	isDuplicates    bool
	isServerMigrate bool
	isLocalMigrate  bool
	isCreateJSON    bool
	files           []string
	dataTargets     []string
}

func NewRunOptions(isDuplicates, isServerMigrate, isLocalMigrate, isCreateJSON bool,
	files []string,
	dataTargets []string) RunOptions {
	return RunOptions{
		isDuplicates:    isDuplicates,
		isServerMigrate: isServerMigrate,
		isLocalMigrate:  isLocalMigrate,
		isCreateJSON:    isCreateJSON,
		files:           files,
		dataTargets:     dataTargets,
	}
}

type App struct {
	files  []string
	cfg    config.Config
	opt    RunOptions
	logger *slog.Logger
}

func NewApp(logger *slog.Logger, cfg config.Config, options RunOptions) *App {
	return &App{
		cfg:    cfg,
		opt:    options,
		logger: logger,
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

	excelWriter := parseio.NewWriter("Сводная таблица", app.logger)
	defer excelWriter.SaveFile()

	wg.Go(func() {
		notes := pivot.ProcessData(tournaments)
		excelWriter.Write(notes)
	})

	if app.opt.isCreateJSON {
		jsonWriter := jsonio.NewWriter("Data", app.logger)
		defer jsonWriter.SaveFile()

		wg.Go(func() {
			jsonWriter.Write(tournaments)
		})
	}

	if app.opt.isDuplicates {
		dupWriter := dupio.NewWriter("Дубли", app.logger)
		defer dupWriter.SaveFile()

		wg.Go(func() {
			dupNotes := duplicates.ProcessData(tournaments)
			dupWriter.Write(dupNotes)
		})
	}

	if app.opt.isLocalMigrate {
		app.logger.Info("Запись в локальную БД")

		if err := app.writeToDB(wg, tournaments, nil); err != nil {
			return fmt.Errorf("ошибка записи в БД - %w", err)
		}
	} else if app.opt.isServerMigrate {
		app.logger.Info("Запись в удаленную БД")

		sshClient, err := ssh.NewSSHClient(app.cfg)
		if err != nil {
			return fmt.Errorf("ошибка инициализации SSHClient - %w", err)
		}
		defer sshClient.Close()

		if err := app.writeToDB(wg, tournaments, sshClient.ConnectRemoteDB); err != nil {
			return fmt.Errorf("ошибка записи в БД - %w", err)
		}
	}

	wg.Wait()

	app.logger.Info("Выполнение занято", slog.String("duration", time.Since(start).String()))
	return nil
}

func (app *App) writeToDB(
	wg *sync.WaitGroup,
	tournaments models.ExcelSheet,
	dialFunc func(ctx context.Context, network, addr string) (net.Conn, error),
) error {
	dbInitCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	connString := app.cfg.Database.GetConnString()

	dbWriter, err := dbpool.New(dbInitCtx, connString, dialFunc)
	if err != nil {
		return fmt.Errorf("ошибка инициализации DBWriter - %w", err)
	}

	pgRepo := repository.NewTournamentRepository(dbWriter, app.logger)
	exportService, err := export.NewExportService(pgRepo)
	if err != nil {
		dbWriter.Close()
		return fmt.Errorf("возникла ошибка при создании сервиса для экспорта данных: %w", err)
	}

	switch {
	case slices.Contains(app.opt.dataTargets, DataTargetTournaments):
		wg.Go(func() {
			defer dbWriter.Close()
			exportService.SaveTournaments(context.Background(), tournaments)
		})
	case slices.Contains(app.opt.dataTargets, DataTargetJudokas):
		filePath, err := filesutils.GetRootFilePath(JudokasFileName)
		if err != nil {
			return err
		}
		reader, err := judokaio.NewReader(filePath)
		if err != nil {
			return fmt.Errorf("ошибка инициализации Reader - %w", err)
		}

		service := parse.NewJudokaService(reader, app.logger)
		judokas, err := service.Parse()
		if err != nil {
			return fmt.Errorf("ошибка парсинга файлов - %w", err)
		}

		wg.Go(func() {
			defer dbWriter.Close()
			exportService.SaveJudokas(context.Background(), judokas)
		})
	}

	return nil
}
