package parse

import (
	"fmt"
	judokaio "judo/internal/io/excel/judoka_parse"
	"judo/internal/models"
	"log"
)

type JudokaService struct {
	reader *judokaio.Reader
}

func NewJudokaService(reader *judokaio.Reader) *JudokaService {
	return &JudokaService{
		reader: reader,
	}
}

func (s *JudokaService) Parse() ([]models.JudokaDBRow, error) {
	data, err := s.reader.Read()
	if err != nil {
		return nil, fmt.Errorf("не удалось прочесть данные: %w", err)
	}

	result := make([]models.JudokaDBRow, 0)

	for _, row := range data {
		if len(row) < 4 {
			log.Printf("Пропуск строки %v, длина %d\n", row, len(row))
			continue
		}

		judoka, err := models.NewJudokaDBRow(row)
		if err != nil {
			log.Printf("Ошибка создания дзюдоиста: %v\n", err)
			continue
		}

		result = append(result, judoka)
	}

	return result, nil
}
