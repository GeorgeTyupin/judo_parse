package rating_resolver

import (
	"fmt"
)

// TournamentRatings мапа, которая сопоставляет типу турнира другую мапу, где каждому месту присвоен рейтинг
type TournamentRatings map[string]map[string]int
type RatingResolver struct {
	tournamentRatings TournamentRatings
}

func NewRatingResolver(ratings TournamentRatings) RatingResolver {
	return RatingResolver{
		tournamentRatings: ratings,
	}
}

func (r RatingResolver) Resolve(tourType, rank string) string {
	if r.tournamentRatings[tourType] == nil {
		return ""
	}
	res, ok := r.tournamentRatings[tourType][rank]
	if !ok {
		return ""
	}

	return fmt.Sprint(res)
}
