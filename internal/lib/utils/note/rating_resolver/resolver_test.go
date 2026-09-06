package rating_resolver_test

import (
	"log/slog"
	"testing"

	dictio "judo/internal/io/excel/dict"
	filesutils "judo/internal/lib/utils/files"
	ratingResolver "judo/internal/lib/utils/note/rating_resolver"
	"judo/internal/services/rating"

	"github.com/stretchr/testify/require"
)

const dictFileName string = "#SPRAVOCHNIK.xlsx"

func TestRatingResolver_Resolve(t *testing.T) {
	filePath, err := filesutils.GetRootFilePath(dictFileName)
	require.NoError(t, err)

	reader, err := dictio.NewReader(filePath)
	require.NoError(t, err)

	service, err := rating.NewRatingService(reader, slog.Default())
	require.NoError(t, err)
	require.NotNil(t, service)

	ratings, err := service.GetRatings()
	require.NoError(t, err)
	require.NotEmpty(t, ratings)

	resolver := ratingResolver.NewRatingResolver(ratings)

	tests := []struct {
		name     string
		tourType string
		rank     string
		want     string
	}{
		{"Olympic Games, 1 место", "Olympic Games", "1", "10000"},
		{"Olympic Games, 3 место", "Olympic Games", "3", "4000"},
		{"European Championships, 1 место", "European Championships", "1", "3510"},
		{"European Championships, 7 место", "European Championships", "7", "570"},
		{"National Team Cup, 2 место", "National Team Cup", "2", "450"},
		{"National Championships, 5 место", "National Championships", "5", "200"},
		{"неизвестный тип турнира", "Несуществующий турнир", "1", ""},
		{"известный тип турнира, несуществующее место", "Olympic Games", "4", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := resolver.Resolve(tt.tourType, tt.rank)
			require.Equal(t, tt.want, got)
		})
	}
}
