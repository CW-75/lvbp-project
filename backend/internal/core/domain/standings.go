package domain

type TeamStandings struct {
	TeamID       int
	GamesPlayed  int
	Won          int
	Lost         int
	Pct          float64
	GamesBehind  float64
}
