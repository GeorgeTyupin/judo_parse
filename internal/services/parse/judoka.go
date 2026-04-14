package parse

import (
	"fmt"
	judokaio "judo/internal/io/excel/judoka_parse"
	"judo/internal/models"
)

type JudokaService struct {
	reader *judokaio.Reader
}

func NewJudokaService(reader *judokaio.Reader) *JudokaService {
	return &JudokaService{
		reader: reader,
	}
}

func (s *JudokaService) Parse() ([]models.Judoka, error) {
	data, err := s.reader.Read()
	if err != nil {
		return nil, fmt.Errorf("не удалось прочесть данные: %w", err)
	}

	result := data.([]models.Judoka)

	return result, nil
}
