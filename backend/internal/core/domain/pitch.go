package domain

type Pitch struct {
	ID          int
	AtBatID     int
	CoordX      float64
	CoordY      float64
	PitchResult string // Ball, Called Strike, Swinging Strike, Foul, In Play, Hit By Pitch
	CountBefore string // e.g. "0-0-0"
	VelocityMPH float64
}

type AtBat struct {
	ID           int
	GameID       int
	BatterID     int
	PitcherID    int
	Result       string
	RunsScored   int
	OutsRecorded int
}
