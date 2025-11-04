package judioio

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"

	"judo/internal/models"
)

func ToJson(wg *sync.WaitGroup, data models.ExelSheet, file string) {
	defer wg.Done()
	fmt.Println(1)

	newJson, err := os.Create(fmt.Sprintf("%s.json", file))
	if err != nil {
		log.Fatalf("Ошибка создания файла: %v", err)
	}
	defer newJson.Close()

	encoder := json.NewEncoder(newJson)
	encoder.SetIndent("", "    ")

	err = encoder.Encode(&data)

	if err != nil {
		log.Fatalf("Ошибка записи в файл: %v", err)
	}
}
