package parse

import (
	"log/slog"
	"testing"

	judoio "judo/internal/io/excel/judoka_parse"
	filesutils "judo/internal/lib/utils/files"
	"judo/internal/models"

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

	testJudokas := []models.JudokaDBRow{
		{
			LastName:     "A. Almeida",
			FirstName:    "Miguel",
			LastNameRus:  new("A. Almeida"),
			FirstNameRus: new("Miguel"),
		},
		{
			LastName:     "Aaltonen",
			FirstName:    "Tiina",
			LastNameRus:  new("Aaltonen"),
			FirstNameRus: new("Tiina"),
		},
		{
			LastName:     "Aamodt",
			FirstName:    "Alexander",
			LastNameRus:  new("Aamodt"),
			FirstNameRus: new("Alexander"),
		},
		{
			LastName:     "Aarts",
			FirstName:    "Monique",
			LastNameRus:  new("Aarts"),
			FirstNameRus: new("Monique"),
		},
		{
			LastName:     "Abadie",
			FirstName:    "",
			LastNameRus:  new("Abadie"),
			FirstNameRus: new(""),
		},
		{
			LastName:     "Abanoz",
			FirstName:    "Salim",
			LastNameRus:  new("Abanoz"),
			FirstNameRus: new("Salim"),
		},
		{
			LastName:     "Abbad",
			FirstName:    "Tahar",
			LastNameRus:  new("Abbad"),
			FirstNameRus: new("Tahar"),
		},
		{
			LastName:     "Abdoune",
			FirstName:    "Hamid",
			LastNameRus:  new("Abdoune"),
			FirstNameRus: new("Hamid"),
		},
		{
			LastName:     "Abe",
			FirstName:    "Takahiro",
			LastNameRus:  new("Abe"),
			FirstNameRus: new("Takahiro"),
		},
		{
			LastName:     "Abe",
			FirstName:    "Yuji",
			LastNameRus:  new("Abe"),
			FirstNameRus: new("Yuji"),
		},
	}

	for i, expected := range testJudokas {
		require.Equal(t, expected.LastName, judokas[i].LastName)
		require.Equal(t, expected.FirstName, judokas[i].FirstName)
		require.Equal(t, expected.LastNameRus, judokas[i].LastNameRus)
		require.Equal(t, expected.FirstNameRus, judokas[i].FirstNameRus)
	}

}
