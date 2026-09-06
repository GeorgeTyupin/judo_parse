package rating

import (
	"fmt"
	"log/slog"
	"strconv"
	"strings"

	dictio "judo/internal/io/excel/dict"
	ratingResolver "judo/internal/lib/utils/note/rating_resolver"
)

const (
	ratingSheetKey = "RATING"

	// Reader.Read() уже отрезает первую строку листа (коэффициенты мест),
	// поэтому здесь заголовок — это первая оставшаяся строка.
	headerRowIdx    = 0
	firstDataRowIdx = 1

	tourTypeColIdx   = 2 // колонка C — Tour type
	firstPlaceColIdx = 3 // колонка D — начало колонок с местами (1, 2, 3, 5, 7...)
)

type RatingService struct {
	logger *slog.Logger
	rows   [][]string
}

func NewRatingService(reader *dictio.Reader, logger *slog.Logger) (*RatingService, error) {
	data, err := reader.Read()
	if err != nil {
		return nil, fmt.Errorf("ошибка чтения данных: %w", err)
	}

	return &RatingService{
		logger: logger,
		rows:   data[ratingSheetKey],
	}, nil
}

func (s RatingService) GetRatings() (ratingResolver.TournamentRatings, error) {
	result := make(ratingResolver.TournamentRatings)

	if len(s.rows) <= firstDataRowIdx {
		return result, nil
	}

	places := parsePlaces(s.rows[headerRowIdx])

	for _, row := range s.rows[firstDataRowIdx:] {
		if len(row) <= tourTypeColIdx {
			s.logger.Debug("Пропуск строки рейтинга", slog.Int("длина", len(row)))
			continue
		}

		tourType := strings.TrimSpace(row[tourTypeColIdx])
		if tourType == "" {
			continue
		}

		if result[tourType] == nil {
			result[tourType] = make(map[string]int)
		}

		for i, place := range places {
			colIdx := firstPlaceColIdx + i
			if colIdx >= len(row) {
				break
			}

			value := strings.TrimSpace(row[colIdx])
			if value == "" {
				continue
			}

			rating, err := strconv.Atoi(value)
			if err != nil {
				s.logger.Warn("Ошибка чтения рейтинга", slog.String("значение", value), slog.String("error", err.Error()))
				continue
			}

			result[tourType][place] = rating
		}
	}

	return result, nil
}

// parsePlaces читает номера мест (1, 2, 3, 5, 7...) из строки заголовков начиная с
// firstPlaceColIdx и останавливается на первой пустой ячейке — правее в этом же
// листе идёт вторая, отдельная таблица с такими же заголовками.
func parsePlaces(header []string) []string {
	places := make([]string, 0)
	for i := firstPlaceColIdx; i < len(header); i++ {
		place := strings.TrimSpace(header[i])
		if place == "" {
			break
		}
		places = append(places, place)
	}

	return places
}
