package russifiers_test

import (
	"fmt"
	judoio "judo/internal/io/excel/judoka_parse"
	filesutils "judo/internal/lib/utils/files"
	"judo/internal/lib/utils/note/russifiers"
	"judo/internal/models"
	"judo/internal/services/parse"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJudokaRussifier_Russify(t *testing.T) {
	filePath, err := filesutils.GetRootFilePath("#JUDOKA.xlsx")
	require.NoError(t, err)

	reader, err := judoio.NewReader(filePath)
	require.NoError(t, err)

	judokaService := parse.NewJudokaService(reader, slog.Default())

	judokas, err := judokaService.Parse()
	require.NoError(t, err)

	judokaNames := models.JudokaRowsToNames(judokas)
	judokaRussifier := russifiers.NewJudokaRussifier(judokaNames)

	parseService, err := parse.NewParseService([]string{"USSR_tours"})
	require.NoError(t, err)

	tournamentSheets, err := parseService.ParseTournaments()
	require.NoError(t, err)

	for _, sheet := range tournamentSheets {
		for i, tournament := range sheet {
			if i > 5 {
				break
			}
			for _, judokas := range tournament.WeightCategories {
				for _, judoka := range judokas {
					fullName := fmt.Sprintf("%s %s", judoka.LastName, judoka.FirstName)
					t.Run(fullName, func(t *testing.T) {
						got := judokaRussifier.Russify(fullName)
						assert.Equal(t, got[0], judoka.FirstName)
						assert.Equal(t, got[1], judoka.LastName)
					})
				}
			}
		}
	}
}
