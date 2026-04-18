package parse

import (
	"fmt"
	judokaio "judo/internal/io/excel/judoka_parse"
	"judo/internal/models"
	"log/slog"
)

type JudokaService struct {
	reader *judokaio.Reader
	logger *slog.Logger
}

func NewJudokaService(reader *judokaio.Reader, logger *slog.Logger) *JudokaService {
	return &JudokaService{
		reader: reader,
		logger: logger,
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
			slog.Debug("Пропуск строки", slog.Int("длина", len(row)))
			continue
		}

		judoka, err := models.NewJudokaDBRow(row)
		if err != nil {
			slog.Debug("Ошибка создания дзюдоиста", slog.String("error", err.Error()))
			continue
		}

		result = append(result, judoka)
	}

	return result, nil
}
