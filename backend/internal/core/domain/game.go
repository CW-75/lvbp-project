package domain

import (
	"time"

	"github.com/google/uuid"
)

type GameStatus string

const (
	GameStatusScheduled  GameStatus = "SCHEDULED"
	GameStatusInProgress GameStatus = "IN_PROGRESS"
	GameStatusFinal      GameStatus = "FINAL"
	GameStatusPostponed  GameStatus = "POSTPONED"
	GameStatusSuspended  GameStatus = "SUSPENDED"
)

type Game struct {
	ID            uuid.UUID
	HomeTeamID    string
	AwayTeamID    string
	Status        GameStatus
	HomeScore     int
	AwayScore     int
	CurrentInning int
	IsTopInning   bool
	StartTime     time.Time
	CreatedAt     time.Time
}
