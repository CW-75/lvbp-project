package domain

import (
	"time"

	"github.com/google/uuid"
)

type Pitch struct {
	ID            uuid.UUID
	AtBatID       uuid.UUID
	PitchNumber   int
	CoordinateX   float64
	CoordinateY   float64
	PitchResult   string
	BallsBefore   int
	StrikesBefore int
	OutsBefore    int
	VelocityMPH   *float64 // Can be nil
	CreatedAt     time.Time
}
