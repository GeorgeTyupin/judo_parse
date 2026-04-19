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
	"Ordzhonikidzeabad", "Vahdat",
	"Ordzhonikidze", "Vladikavkaz",
	"Ustinov", "Izhevsk",
	"Kuibyshev", "Samara",
	"Zhdanov", "Mariupol",
	"Voroshilovgrad", "Lugansk",
	"Gorkiy", "Nizhniy Novgorod",
	"Kaliningrad", "Kaliningrad",
	"Kalinin", "Tver",
)

func NormalizeCityName(s string) string {
	return cityReplacer.Replace(s)
}
