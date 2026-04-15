package main

import (
	"log"
	"os"

	app "judo/internal/app"
	"judo/internal/config"

	"github.com/charmbracelet/huh"
)

const (
	migrationServer = "server"
	migrationLocal  = "local"
)

var (
	isDuplicates     bool
	isMigrate        bool
	isServerMigrate  bool
	isLocalMigrate   bool
	isCreateJSON     bool
	files            []string
	migrationTargets []string
	dataTargets      []string
)

func main() {
	os.Remove("Сводная таблица.xlsx")
	os.Remove("Дубли.xlsx")
	os.Remove("USSR_tours.json")

	cfg := config.MustLoad()

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewMultiSelect[string]().
				Title("Выбор исходного файла").
				Options(
					huh.NewOption("USSR_tours", "USSR_tours").Selected(true),
					huh.NewOption("INT_tours", "INT_tours").Selected(true),
				).Value(&files),
		),

		huh.NewGroup(
			huh.NewConfirm().Title("Мигрировать данные?").Value(&isMigrate),
		),

		huh.NewGroup(
			huh.NewMultiSelect[string]().
				Title("Как мигрировать данные?").
				Options(
					huh.NewOption("На сервер", migrationServer),
					huh.NewOption("В локальную БД", migrationLocal),
				).Value(&migrationTargets),

			huh.NewMultiSelect[string]().
				Title("Какие данные мигрировать?").
				Options(
					huh.NewOption("Турниры", app.DataTargetTournaments),
					huh.NewOption("Дзюдоистов", app.DataTargetJudokas),
				).Value(&dataTargets),
		).WithHideFunc(func() bool {
			return !isMigrate
		}),

		huh.NewGroup(
			huh.NewConfirm().Title("Проверять на дубли?").Value(&isDuplicates),
		),

		huh.NewGroup(
			huh.NewConfirm().Title("Создать JSON?").Value(&isCreateJSON),
		),
	)

	if err := form.Run(); err != nil {
		log.Fatalf("Ошибка: %v, попробуйте еще раз", err)
	}

	for _, t := range migrationTargets {
		switch t {
		case migrationServer:
			isServerMigrate = true
		case migrationLocal:
			isLocalMigrate = true
		}
	}

	options := app.NewRunOptions(isDuplicates, isServerMigrate, isLocalMigrate, isCreateJSON, files, dataTargets)

	application := app.NewApp(cfg, options)
	if err := application.Run(); err != nil {
		log.Fatal(err)
	}
}
