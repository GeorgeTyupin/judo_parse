package jsonio

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"

	"judo/internal/models"
)

type JsonWriter struct {
	Name   string
	File   *os.File
	logger *slog.Logger
}

func NewWriter(name string, logger *slog.Logger) *JsonWriter {
	newJson, err := os.Create(fmt.Sprintf("%s.json", name))
	if err != nil {
		slog.Error("Ошибка создания файла", slog.String("error", err.Error()))
		os.Exit(1)
	}

	jsonFile := JsonWriter{
		Name:   name,
		File:   newJson,
		logger: logger,
	}

	return &jsonFile
}

func (w *JsonWriter) Write(data models.ExcelSheet) {
	encoder := json.NewEncoder(w.File)
	encoder.SetIndent("", "    ")

	err := encoder.Encode(&data)
	if err != nil {
		slog.Error("Ошибка записи в файл", slog.String("error", err.Error()))
		os.Exit(1)
	}
}

func (w *JsonWriter) SaveFile() {
	w.File.Close()
}
