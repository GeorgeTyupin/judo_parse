package rating

import (
	"log/slog"
	"testing"

	dictio "judo/internal/io/excel/dict"
	filesutils "judo/internal/lib/utils/files"

	"github.com/stretchr/testify/require"
)

const dictFileName string = "#SPRAVOCHNIK.xlsx"

func TestRatingsParsing(t *testing.T) {
	filePath, err := filesutils.GetRootFilePath(dictFileName)
	require.NoError(t, err)

	reader, err := dictio.NewReader(filePath)
	require.NoError(t, err)

	service, err := NewRatingService(reader, slog.Default())
	require.NoError(t, err)
	require.NotNil(t, service)

	ratings, err := service.GetRatings()
	require.NoError(t, err)
	require.NotEmpty(t, ratings)

	i := 0
	for tourType, places := range ratings {
		require.NotEmpty(t, places)

		if i >= 10 {
			break
		}
		i++
		

		t.Logf("Тип турнира №%d: %v\n Места и рейтинг - %v\n\n", i, tourType, places)
	}
}