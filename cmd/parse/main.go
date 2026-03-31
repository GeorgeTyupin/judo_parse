package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	app "judo/internal/app"
	"judo/internal/config"
)

func yesNoChoice(choice string) (bool, error) {
	switch strings.ToLower(choice) {
	case "y":
		return true, nil
	case "n":
		return false, nil
	default:
		return false, fmt.Errorf("неверный выбор %s", choice)
	}
}

func filesChoice(choice string) ([]string, error) {
	files := make([]string, 0, 2)

	switch choice {
	case "1":
		files = append(files, "USSR_tours")
	case "2":
		files = append(files, "INT_tours")
	case "3":
		files = append(files, "USSR_tours")
		files = append(files, "INT_tours")
	default:
		return nil, errors.New("ошибка ввода файла")
	}

	return files, nil
}

func main() {
	os.Remove("Сводная таблица.xlsx")
	os.Remove("Дубли.xlsx")
	os.Remove("USSR_tours.json")

	cfg := config.MustLoad()
	var choiceFile, choiceDuplicates, choiceMigrate string

	if cfg.IsDev {
		choiceFile = "1"
		choiceDuplicates = "n"
		choiceMigrate = "n"
	} else {
		fmt.Println("Выбор исходного файла. Введи:\n1, если исходный USSR_tours\n2, если исходный INT_tours\n3, если оба")
		fmt.Scanln(&choiceFile)
		fmt.Println("Проверять на дубли. Y/n")
		fmt.Scanln(&choiceDuplicates)
		fmt.Println("Мигрировать данные на сервер? Y/n")
		fmt.Scanln(&choiceMigrate)
	}

	files, err := filesChoice(choiceFile)
	if err != nil {
		log.Fatalf("Ошибка: %v, попробуйте еще раз", err)
	}

	isDuplicates, err := yesNoChoice(choiceDuplicates)
	if err != nil {
		log.Fatalf("Ошибка: %v, попробуйте еще раз", err)
	}

	isServerMigrate, err := yesNoChoice(choiceMigrate)
	application := app.NewApp(cfg, files, isDuplicates, isServerMigrate)
	if err = application.Run(); err != nil {
		log.Fatal(err)
	}
}
