package main

import (
	"log/slog"
	"os"

	app "judo/internal/app"
	"judo/internal/config"
	"judo/internal/ui/form"
)

func main() {
	logger := slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	)

	logger.Info("Удаление старых файлов")
	os.Remove("Сводная таблица.xlsx")
	os.Remove("Дубли.xlsx")
	os.Remove("USSR_tours.json")

	cfg := config.MustLoad()

	options, err := form.Run()
	if err != nil {
		logger.Error("Ошибка создания формы, попробуйте еще раз", slog.String("error", err.Error()))
		os.Exit(1)
	}

	application, err := app.NewApp(logger, cfg, options)
	if err != nil {
		logger.Error("Ошибка при создании приложения", slog.String("error", err.Error()))
		os.Exit(1)
	}

	if err := application.Run(); err != nil {
		logger.Error("Ошибка при запуске приложения", slog.String("error", err.Error()))
		os.Exit(1)
	}
}
