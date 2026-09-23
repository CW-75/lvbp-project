package domain

import (
	"time"

	"github.com/google/uuid"
)

type AtBat struct {
	ID           uuid.UUID
	GameID       uuid.UUID
	Inning       int
	IsTopInning  bool
	BatterID     string
	PitcherID    string
	Result       *string // Can be nil initially
	RunsScored   int
	OutsRecorded int
	CreatedAt    time.Time
}
