package app

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"slices"

	"judo/internal/client/ssh"
	"judo/internal/config"
	dictio "judo/internal/io/excel/dict"
	dupio "judo/internal/io/excel/duplicates"
	parseio "judo/internal/io/excel/parse"
	jsonio "judo/internal/io/json"
	filesutils "judo/internal/lib/utils/files"
	"judo/internal/lib/utils/note/locresolver"
	"judo/internal/lib/utils/note/russifiers"
	"judo/internal/models"
	"judo/internal/repository"
	dbpool "judo/internal/repository/pool"
	"judo/internal/services/dict"
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
	DataTargetCities      string = "cities"
	DataTargetCountries   string = "countries"
	DataTargetSportClubs  string = "sportclubs"

	MigrationTargetServer string = "server"
	MigrationTargetLocal  string = "local"

	DictFileName string = "#SPRAVOCHNIK.xlsx"
)

type Data struct {
	tournaments models.ExcelSheet
	judokas     []models.JudokaDBRow
	cities      []models.CityDBRow
	countries   []models.CountryDBRow
	sportClubs  []models.SportClubDBRow
}

func NewData(logger *slog.Logger, files []string) (Data, error) {
	data := Data{}

	parseService, err := parse.NewParseService(files)
	if err != nil {
		return data, fmt.Errorf("ошибка инициализации ParseService - %w", err)
	}

	tournaments, err := parseService.ParseTournaments()
	if err != nil {
		return data, fmt.Errorf("ошибка парсинга файлов - %w", err)
	}

	filePath, err := filesutils.GetRootFilePath(DictFileName)
	if err != nil {
		return data, err
	}
	reader, err := dictio.NewReader(filePath)
	if err != nil {
		return data, fmt.Errorf("ошибка инициализации Reader - %w", err)
	}

	service := dict.NewDictService(reader, logger)
	judokas, err := service.ParseJudokas()
	if err != nil {
		return data, fmt.Errorf("ошибка парсинга дзюдоистов - %w", err)
	}

	cities, err := service.ParseCities()
	if err != nil {
		return data, fmt.Errorf("ошибка парсинга городов - %w", err)
	}

	countries, err := service.ParseCountries()
	if err != nil {
		return data, fmt.Errorf("ошибка парсинга стран - %w", err)
	}

	sportClubs, err := service.ParseSportClubs()
	if err != nil {
		return data, fmt.Errorf("ошибка парсинга спортивных клубов - %w", err)
	}

	data = Data{
		judokas:     judokas,
		tournaments: tournaments,
		cities:      cities,
		countries:   countries,
		sportClubs:  sportClubs,
	}

	return data, nil
}

type RunOptions struct {
	isDuplicates    bool
	isServerMigrate bool
	isLocalMigrate  bool
	isCreateJSON    bool
	files           []string
	dataTargets     []string
}

func NewRunOptions(isDuplicates, isCreateJSON bool, files, migrationTargets, dataTargets []string) RunOptions {
	return RunOptions{
		isDuplicates:    isDuplicates,
		isServerMigrate: slices.Contains(migrationTargets, MigrationTargetServer),
		isLocalMigrate:  slices.Contains(migrationTargets, MigrationTargetLocal),
		isCreateJSON:    isCreateJSON,
		files:           files,
		dataTargets:     dataTargets,
	}
}

type App struct {
	cfg    config.Config
	opt    RunOptions
	data   Data
	logger *slog.Logger
}

func NewApp(logger *slog.Logger, cfg config.Config, options RunOptions) (*App, error) {
	data, err := NewData(logger, options.files)
	if err != nil {
		return nil, fmt.Errorf("не удалось получить данные. Возникла ошибка %w", err)
	}

	return &App{
		cfg:    cfg,
		opt:    options,
		logger: logger,
		data:   data,
	}, nil
}

func (app *App) Run() error {
	wg := &sync.WaitGroup{}
	start := time.Now()

	// Инициализация русификатора
	judokaNames := models.JudokaRowsToNames(app.data.judokas)
	judokaRussifier := russifiers.NewJudokaRussifier(judokaNames)

	excelWriter := parseio.NewWriter("Сводная таблица", app.logger)
	defer excelWriter.SaveFile()

	columnSplitter, err := locresolver.NewLocationResolver(
		models.GetCountryCodes(app.data.countries),
		models.ToCityInput(app.data.cities),
		models.GetSportClubNames(app.data.sportClubs),
	)
	if err != nil {
		return fmt.Errorf("ошибка инициализации LocationResolver - %w", err)
	}

	wg.Go(func() {
		notes := pivot.ProcessData(app.data.tournaments, judokaRussifier, columnSplitter)
		excelWriter.Write(notes)
	})

	if app.opt.isCreateJSON {
		jsonWriter := jsonio.NewWriter("Data", app.logger)
		defer jsonWriter.SaveFile()

		wg.Go(func() {
			jsonWriter.Write(app.data.tournaments)
		})
	}

	if app.opt.isDuplicates {
		dupWriter := dupio.NewWriter("Дубли", app.logger)
		defer dupWriter.SaveFile()

		wg.Go(func() {
			dupNotes := duplicates.ProcessData(app.data.tournaments, judokaRussifier, columnSplitter)
			dupWriter.Write(dupNotes)
		})
	}

	if app.opt.isLocalMigrate {
		app.logger.Info("Запись в локальную БД")

		var err error
		wg.Go(func() {
			err = app.writeToDB(nil)
		})
		if err != nil {
			return fmt.Errorf("ошибка записи в БД - %w", err)
		}
	}

	if app.opt.isServerMigrate {
		app.logger.Info("Запись в удаленную БД")

		sshClient, err := ssh.NewSSHClient(app.cfg)
		if err != nil {
			return fmt.Errorf("ошибка инициализации SSHClient - %w", err)
		}
		defer sshClient.Close()

		wg.Go(func() {
			err = app.writeToDB(sshClient.ConnectRemoteDB)
		})

		if err != nil {
			return fmt.Errorf("ошибка записи в БД - %w", err)
		}
	}

	wg.Wait()

	app.logger.Info("Выполнение заняло", slog.String("duration", time.Since(start).String()))
	return nil
}

func (app *App) writeToDB(
	dialFunc func(ctx context.Context, network, addr string) (net.Conn, error),
) error {
	dbInitCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	connString := app.cfg.Database.GetConnString()

	dbWriter, err := dbpool.New(dbInitCtx, connString, dialFunc)
	if err != nil {
		return fmt.Errorf("ошибка инициализации DBWriter - %w", err)
	}

	pgRepo := repository.NewCommonRepository(dbWriter, app.logger)
	exportService, err := export.NewExportService(pgRepo)
	if err != nil {
		dbWriter.Close()
		return fmt.Errorf("возникла ошибка при создании сервиса для экспорта данных: %w", err)
	}

	dbWg := sync.WaitGroup{}

	if slices.Contains(app.opt.dataTargets, DataTargetTournaments) {
		dbWg.Go(func() {
			exportService.SaveTournaments(context.Background(), app.data.tournaments)
		})
	}
	if slices.Contains(app.opt.dataTargets, DataTargetJudokas) {
		dbWg.Go(func() {
			exportService.SaveJudokas(context.Background(), app.data.judokas)
		})
	}
	if slices.Contains(app.opt.dataTargets, DataTargetCities) {
		dbWg.Go(func() {
			exportService.SaveCities(context.Background(), app.data.cities)
		})
	}
	if slices.Contains(app.opt.dataTargets, DataTargetCountries) {
		dbWg.Go(func() {
			exportService.SaveCountries(context.Background(), app.data.countries)
		})
	}
	if slices.Contains(app.opt.dataTargets, DataTargetSportClubs) {
		dbWg.Go(func() {
			exportService.SaveSportClubs(context.Background(), app.data.sportClubs)
		})
	}

	dbWg.Wait()
	dbWriter.Close()

	return nil
}
