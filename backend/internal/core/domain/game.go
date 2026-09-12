package domain

type Game struct {
	ID            int
	Status        string // Scheduled, Live, Final
	HomeTeamID    int
	AwayTeamID    int
	HomeScore     int
	AwayScore     int
	CurrentInning int
	IsTopInning   bool
}
