package ports

import (
	"context"

	"github.com/google/uuid"
	"lvbp-project/backend/internal/core/domain"
)

type GameRepository interface {
	GetGameByID(ctx context.Context, id uuid.UUID) (*domain.Game, error)
	ListActiveGames(ctx context.Context) ([]*domain.Game, error)
	CreateGame(ctx context.Context, game *domain.Game) (*domain.Game, error)
	UpdateGame(ctx context.Context, game *domain.Game) (*domain.Game, error)
	
	CreateAtBat(ctx context.Context, atBat *domain.AtBat) (*domain.AtBat, error)
	UpdateAtBatResult(ctx context.Context, atBat *domain.AtBat) (*domain.AtBat, error)
	GetCurrentAtBat(ctx context.Context, gameID uuid.UUID) (*domain.AtBat, error)
	
	CreatePitch(ctx context.Context, pitch *domain.Pitch) (*domain.Pitch, error)
	GetPitchesForAtBat(ctx context.Context, atBatID uuid.UUID) ([]*domain.Pitch, error)
	
	GetOutsForInning(ctx context.Context, gameID uuid.UUID, inning int, isTopInning bool) (int, error)
}

type StandingsRepository interface {
	GetStandings(ctx context.Context) ([]*domain.TeamStandings, error)
}
