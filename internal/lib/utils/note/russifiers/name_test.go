package russifiers_test

import (
	dictio "judo/internal/io/excel/dict"
	filesutils "judo/internal/lib/utils/files"
	"judo/internal/lib/utils/note/russifiers"
	"judo/internal/models"
	"judo/internal/services/dict"
	"judo/internal/services/parse"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestJudokaRussifier_Russify(t *testing.T) {
	judokaFilePath, err := filesutils.GetRootFilePath("#SPRAVOCHNIK.xlsx")
	require.NoError(t, err)

	reader, err := dictio.NewReader(judokaFilePath)
	require.NoError(t, err)

	dictService := dict.NewDictService(reader, slog.Default())

	judokas, err := dictService.ParseJudokas()
	require.NoError(t, err)
	require.NotEmpty(t, judokas)

	judokaNames := models.JudokaRowsToNames(judokas)
	judokaRussifier := russifiers.NewJudokaRussifier(judokaNames)

	tournamentFilePath, err := filesutils.GetRootFilePath("USSR_tours.xlsx")
	require.NoError(t, err)

	parseService, err := parse.NewParseService([]string{tournamentFilePath})
	require.NoError(t, err)

	tournamentSheets, err := parseService.ParseTournaments()
	require.NoError(t, err)

	cnt := 0
l:
	for _, sheet := range tournamentSheets {
		for _, tournament := range sheet {
			for _, judokas := range tournament.WeightCategories {
				for _, judoka := range judokas {
					if cnt > 10 {
						break l
					}
					cnt++
					got := judokaRussifier.Russify(judoka.FirstName, judoka.LastName)
					t.Logf("%s %s -> %s %s", judoka.LastName, judoka.FirstName, got[1], got[0])
				}
			}
		}
	}
}
