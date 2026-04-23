package form

import (
	"judo/internal/app"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/huh"
)

func Run() (app.RunOptions, error) {
	var (
		isMigrate        bool
		files            []string
		migrationTargets []string
		dataTargets      []string
		isDuplicates     bool
		isCreateJSON     bool
	)

	keyMap := huh.NewDefaultKeyMap()
	keyMap.MultiSelect.Toggle = key.NewBinding(
		key.WithKeys(" ", "x", "ч"),
		key.WithHelp("x/ч", "toggle"),
	)

	err := huh.NewForm(
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
					huh.NewOption("На сервер", app.MigrationTargetServer),
					huh.NewOption("В локальную БД", app.MigrationTargetLocal),
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
	).WithKeyMap(keyMap).Run()

	return app.NewRunOptions(isDuplicates, isCreateJSON, files, migrationTargets, dataTargets), err
}
