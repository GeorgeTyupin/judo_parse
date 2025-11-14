package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	app "judo/internal/app"

	"github.com/joho/godotenv"
)

var createJSON = true

func duplicatesCheckChoice(choice string) (bool, error) {
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
	os.Remove("USSR_tours.json")

	var isDev bool
	err := godotenv.Load("configs/.env")
	if err != nil {
		isDev = true
	}
	isDev, _ = strconv.ParseBool(os.Getenv("IS_DEV"))

	var choiceFile, choiceDuplicates string

	if isDev {
		choiceFile = "1"
		choiceDuplicates = "нет"
	} else {
		fmt.Println("Выбор исходного файла. Введи:\n1, если исходный USSR_tours\n2, если исходный INT_tours\n3, если оба")
		fmt.Scanln(&choiceFile)
		fmt.Println("Проверять на дубли. Y/n")
		fmt.Scanln(&choiceDuplicates)
	}

	files, err := filesChoice(choiceFile)
	if err != nil {
		panic(fmt.Sprintf("Ошибка: %v, попробуйте еще раз", err))
	}

	isDuplicates, err := duplicatesCheckChoice(choiceDuplicates)
	if err != nil {
		panic(fmt.Sprintf("Ошибка: %v, попробуйте еще раз", err))
	}

	application := app.NewApp(files, isDuplicates, createJSON)
	if err = application.Run(); err != nil {
		panic(err)
	}
}
