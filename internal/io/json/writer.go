package jsonio

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"judo/internal/models"
)

type JsonWriter struct {
	Name string
	File *os.File
}

func NewWriter(name string) *JsonWriter {
	newJson, err := os.Create(fmt.Sprintf("%s.json", name))
	if err != nil {
		log.Fatalf("Ошибка создания файла: %v", err)
	}

	jsonFile := JsonWriter{
		Name: name,
		File: newJson,
	}

	return &jsonFile
}

func (w *JsonWriter) Write(data models.ExcelSheet) {
	encoder := json.NewEncoder(w.File)
	encoder.SetIndent("", "    ")

	err := encoder.Encode(&data)
	if err != nil {
		log.Fatalf("Ошибка записи в файл: %v", err)
	}
}

func (w *JsonWriter) SaveFile() {
	w.File.Close()
}
