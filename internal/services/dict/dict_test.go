package dict

import (
	"log/slog"
	"testing"

	dictio "judo/internal/io/excel/dict"
	filesutils "judo/internal/lib/utils/files"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

const dictFileName string = "#SPRAVOCHNIK.xlsx"

type DictServiceTestSuite struct {
	service *DictService
	suite.Suite
}

func (s *DictServiceTestSuite) SetupSuite() {
	filePath, err := filesutils.GetRootFilePath(dictFileName)
	require.NoError(s.T(), err)

	reader, err := dictio.NewReader(filePath)
	require.NoError(s.T(), err)

	s.service = NewDictService(reader, slog.Default())
}

func (s *DictServiceTestSuite) TestJudokasParsing() {
	judokas, err := s.service.ParseJudokas()
	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), judokas)

	for i, got := range judokas {
		if i >= 10 {
			break
		}

		s.T().Logf("Дзюдоист №%d:\n Английское имя - %v \n Английская фамилия - %v \n Русское имя - %v \n Русская фамилия - %v\n\n",
			i+1,
			got.FirstName,
			got.LastName,
			*got.FirstNameRus,
			*got.LastNameRus,
		)
	}

}

func (s *DictServiceTestSuite) TestCitiesParsing() {
	cities, err := s.service.ParseCities()
	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), cities)

	for i, got := range cities {
		if i >= 10 {
			break
		}

		s.T().Logf("Город №%d:\n Английское название - %v \n Русское название - %v \n Городской округ - %v \n Республика - %v\n\n",
			i+1,
			got.Name,
			*got.NameRus,
			got.RepublicNameEng,
			got.RepublicNameRus,
		)
	}
}

func (s *DictServiceTestSuite) TestCountriesParsing() {
	countries, err := s.service.ParseCountries()
	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), countries)

	for i, got := range countries {
		if i >= 10 {
			break
		}

		s.T().Logf("Страна №%d:\n Английское название - %v \n Английская аббревиатура - %v \n Русское название - %v \n\n",
			i+1,
			got.Name,
			*got.ISOCode,
			*got.NameRus,
		)
	}

}

func (s *DictServiceTestSuite) TestSportClubsParsing() {
	sportClubs, err := s.service.ParseSportClubs()
	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), sportClubs)

	for i, got := range sportClubs {
		if i >= 10 {
			break
		}

		s.T().Logf("Спортивный клуб №%d:\n Английское название - %v \n Русское название - %v \n\n",
			i+1,
			got.Name,
			*got.NameRus,
		)
	}
}

func TestDictServiceSuite(t *testing.T) {
	suite.Run(t, new(DictServiceTestSuite))
}
