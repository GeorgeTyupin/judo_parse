package dict

import (
	"fmt"
	dictio "judo/internal/io/excel/dict"
	"judo/internal/models"
	"log/slog"
)

const (
	judokaSheetKey    string = "JUDOKA"
	citySheetKey      string = "CITY_COMB"
	countrySheetKey   string = "COUNTRY"
	sportClubSheetKey string = "SO"
)

type DictService struct {
	reader *dictio.Reader
	logger *slog.Logger
	data   map[string][][]string
}

func NewDictService(reader *dictio.Reader, logger *slog.Logger) *DictService {
	service := &DictService{
		reader: reader,
		logger: logger,
	}

	if err := service.readData(); err != nil {
		logger.Error("Ошибка чтения данных", slog.String("error", err.Error()))
		return nil
	}

	return service
}

func (s *DictService) readData() error {
	data, err := s.reader.Read()
	if err != nil {
		return fmt.Errorf("не удалось прочесть данные: %w", err)
	}

	s.data = data

	return nil
}

func (s *DictService) ParseJudokas() ([]models.JudokaDBRow, error) {
	result := make([]models.JudokaDBRow, 0)

	for _, row := range s.data[judokaSheetKey] {
		if len(row) < models.MinJudokaRowLen {
			slog.Warn("Пропуск строки", slog.Int("длина", len(row)))
			continue
		}

		judoka, err := models.NewJudokaDBRow(row)
		if err != nil {
			slog.Warn("Ошибка создания дзюдоиста", slog.String("error", err.Error()))
			continue
		}

		result = append(result, judoka)
	}

	return result, nil
}

func (s *DictService) ParseCities() ([]models.CityDBRow, error) {
	result := make([]models.CityDBRow, 0)

	for _, row := range s.data[citySheetKey] {
		if len(row) < models.MinCityRowLen {
			slog.Warn("Пропуск строки", slog.Int("длина", len(row)))
			continue
		}

		city, err := models.NewCityDBRow(row)
		if err != nil {
			slog.Warn("Ошибка создания города", slog.String("error", err.Error()))
			continue
		}

		result = append(result, city)
	}

	return result, nil
}

func (s *DictService) ParseCountries() ([]models.CountryDBRow, error) {
	result := make([]models.CountryDBRow, 0)

	for _, row := range s.data[countrySheetKey] {
		if len(row) < models.MinCountryRowLen {
			slog.Warn("Пропуск строки", slog.Int("длина", len(row)))
			continue
		}

		country, err := models.NewCountryDBRow(row)
		if err != nil {
			slog.Warn("Ошибка создания страны", slog.String("error", err.Error()))
			continue
		}

		result = append(result, country)
	}

	return result, nil
}

func (s *DictService) ParseSportClubs() ([]models.SportClubDBRow, error) {
	result := make([]models.SportClubDBRow, 0)

	for _, row := range s.data[sportClubSheetKey] {
		if len(row) < models.MinSportClubRowLen {
			slog.Warn("Пропуск строки", slog.Int("длина", len(row)))
			continue
		}

		sportClub, err := models.NewSportClubDBRow(row)
		if err != nil {
			slog.Warn("Ошибка создания спортивного клуба", slog.String("error", err.Error()))
			continue
		}

		result = append(result, sportClub)
	}

	return result, nil
}
