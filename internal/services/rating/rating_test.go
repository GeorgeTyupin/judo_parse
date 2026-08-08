package rating

import (
	"log/slog"
	"testing"

	dictio "judo/internal/io/excel/dict"
	filesutils "judo/internal/lib/utils/files"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

const dictFileName string = "#SPRAVOCHNIK.xlsx"

type RatingServiceTestSuite struct {
	service *RatingService
	suite.Suite
}

func (s *RatingServiceTestSuite) SetupSuite() {
	filePath, err := filesutils.GetRootFilePath(dictFileName)
	require.NoError(s.T(), err)

	reader, err := dictio.NewReader(filePath)
	require.NoError(s.T(), err)

	s.service, err = NewRatingService(reader, slog.Default())
	require.NoError(s.T(), err)
	require.NotNil(s.T(), s.service)
}

func (s *RatingServiceTestSuite) TestRatingsParsing() {
	ratings, err := s.service.GetRatings()
	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), ratings)

	i := 0
	for tourType, places := range ratings {
		require.NotEmpty(s.T(), places)

		if i >= 10 {
			break
		}
		i++

		s.T().Logf("Тип турнира №%d: %v\n Места и рейтинг - %v\n\n", i, tourType, places)
	}
}

func TestRatingServiceSuite(t *testing.T) {
	suite.Run(t, new(RatingServiceTestSuite))
}
