package parse

import (
	"strings"

	parseio "judo/internal/io/excel/parse"
	parseutils "judo/internal/lib/utils/parse"
	"judo/internal/models"
)

type ParseService struct {
	Reader *parseio.Reader
}

func NewParseService(fileNames []string) (*ParseService, error) {
	reader, err := parseio.NewReader(fileNames)
	if err != nil {
		return nil, err
	}

	service := &ParseService{
		Reader: reader,
	}

	return service, nil
}

func (ps *ParseService) ParseTournaments() (models.ExсelSheet, error) {
	data, err := ps.Reader.Read()
	if err != nil {
		return nil, err
	}

	result := make(models.ExсelSheet)

	//Проход по всем листам
	for curSheet, rows := range data {
		if len(rows) < 2 {
			continue
		}

		lenTables := parseutils.FindLenTables(rows[1])

		left := 1

		//Проход по всей таблице
		for _, lenCurTable := range lenTables {
			right := left + lenCurTable

			tournaments := createTournament(left, right, lenCurTable, rows)

			if _, exists := result[curSheet]; !exists {
				result[curSheet] = make([]*models.Tournament, 0)
			}
			result[curSheet] = append(result[curSheet], tournaments...)

			left = right + 1
		}
	}

	return result, nil
}

func readTournamentHeader(rows [][]string, i, left, right int) (*models.Tournament, int) {
	var tournament models.Tournament

	for j := 0; j < 4; j++ {
		if i+j >= len(rows) || left >= len(rows[i+j]) {
			continue
		}

		row := rows[i+j]
		var tournamentRow []string

		if right < cap(row) {
			tournamentRow = row[left:right]
		} else {
			tournamentRow = make([]string, 0, cap(row)+(right-cap(row)))
			copy(tournamentRow, row)
			tournamentRow = tournamentRow[left:right]
		}

		if len(tournamentRow) == 0 {
			continue
		}

		switch j {
		case 0:
			tournament.Name = tournamentRow[0]
		case 1:
			tournament.Description = tournamentRow[0]
		case 2:
			tournament.Date = tournamentRow[0]
		case 3:
			tournament.Gender = tournamentRow[0]
		}
	}

	return &tournament, i + 3
}

func createTournament(left, right, lenCurTable int, rows [][]string) []*models.Tournament {
	var tournaments []*models.Tournament
	var tournament *models.Tournament
	var curWeightCategoryName string
	WeightCategories := make(map[string][]*models.Judoka)
	isNewTournament := false

	//Проход по турниру
	for i := 0; i < len(rows); i++ {
		row := rows[i]

		if left > len(row) {
			continue
		}

		var curRow []string

		if right < cap(row) {
			curRow = row[left:right]
		} else {
			curRow = make([]string, 0, cap(row)+(right-cap(row)))
			copy(curRow, row)
			curRow = curRow[left:right]
		}

		//Отсеиваем пустые и другие строки без нужной информации
		if !parseutils.IsValidDataRow(curRow) {
			continue
		}

		if curRow[0] == "_" {
			// Сохраняем предыдущий турнир, если он был инициализирован
			if isNewTournament {
				tournament.WeightCategories = WeightCategories
				tournaments = append(tournaments, tournament)
			}

			// Получаем шапку нового турнира
			i++
			tournament, i = readTournamentHeader(rows, i, left, right)
			WeightCategories = make(map[string][]*models.Judoka)
			isNewTournament = true
			curWeightCategoryName = ""
		} else {
			// fmt.Println(curRow[0])
			if parseutils.ReNum.MatchString(curRow[0]) || strings.Contains(curRow[0], "Open") {

				if len(curRow[0]) > 2 {
					curWeightCategoryName = curRow[0]
					WeightCategories[curWeightCategoryName] = make([]*models.Judoka, 0)
				} else {
					athlete := models.NewJudoka(curRow, lenCurTable)
					WeightCategories[curWeightCategoryName] = append(WeightCategories[curWeightCategoryName], athlete)
				}
			} else {
				continue
			}
		}

	}

	// Сохраняем последний турнир
	if isNewTournament {
		tournament.WeightCategories = WeightCategories
		tournaments = append(tournaments, tournament)
	}

	return tournaments
}
