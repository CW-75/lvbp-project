package ports

import (
	"context"
	"github.com/google/uuid"
	"lvbp-project/backend/internal/core/domain"
)

// ScorekeeperService represents the core domain logic for baseball rules and state engine.
type ScorekeeperService interface {
	// RecordPitch records a pitch for the current at-bat in a game and evaluates baseball rules
	// (balls, strikes, outs, inning transitions).
	RecordPitch(ctx context.Context, gameID uuid.UUID, pitch *domain.Pitch) (*domain.Game, *domain.AtBat, error)
	
	// StartAtBat starts a new at-bat for a game.
	StartAtBat(ctx context.Context, gameID uuid.UUID, batterID string, pitcherID string) (*domain.AtBat, error)
	
	// RecordAtBatResult records a manual at-bat result (e.g. hit, error, out) instead of pitch-by-pitch.
	RecordAtBatResult(ctx context.Context, gameID uuid.UUID, result string) (*domain.Game, *domain.AtBat, error)
}

// EventBus represents the port for publishing domain events
type EventBus interface {
	PublishEvent(topic string, payload interface{}) error
}
