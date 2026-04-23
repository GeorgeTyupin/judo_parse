package parse

import (
	"fmt"
	"log/slog"
	"testing"

	judoio "judo/internal/io/excel/judoka_parse"
	filesutils "judo/internal/lib/utils/files"

	"github.com/stretchr/testify/require"
)

func TestJudokasParsing(t *testing.T) {
	filePath, err := filesutils.GetRootFilePath("#JUDOKA.xlsx")
	require.NoError(t, err)

	reader, err := judoio.NewReader(filePath)
	require.NoError(t, err)

	service := NewJudokaService(reader, slog.Default())

	judokas, err := service.Parse()
	require.NoError(t, err)
	require.NotEmpty(t, judokas)

	for i, got := range judokas {
		if i >= 10 {
			break
		}

		fmt.Printf("Дзюдоист №%d:\n Английское имя - %v \n Английская фамилия - %v \n Русское имя - %v \n Русская фамилия - %v\n\n",
			i+1,
			got.FirstName,
			got.LastName,
			*got.FirstNameRus,
			*got.LastNameRus,
		)
	}

}
