package replacers

import "strings"

var cityReplacer = strings.NewReplacer(
	"Brezhnev", "Naberezhnye Chelny",
	"Leningrad", "St. Petersburg",
	"Frunze", "Bishkek",
	"Kubyshev", "Samara",
	"Sverdlovsk", "Ekaterinburg",
	"Gorky", "Nizhniy Novgorod",
	"Andropov", "Rybinsk",
	"Ordzhonikidze", "Vladikavkaz",
	"Ustinov", "Izhevsk",
)

func NormalizeCityName(s string) string {
	return cityReplacer.Replace(s)
}
