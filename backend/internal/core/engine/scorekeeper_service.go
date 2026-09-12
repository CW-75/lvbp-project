package engine

import (
	"lvbp-project/backend/internal/core/domain"
	"lvbp-project/backend/internal/core/ports"
)

type ScorekeeperService struct {
	gameRepo ports.GameRepository
	eventPub ports.EventPublisher
}

func NewScorekeeperService(gr ports.GameRepository, ep ports.EventPublisher) *ScorekeeperService {
	return &ScorekeeperService{
		gameRepo: gr,
		eventPub: ep,
	}
}

// ProcessPitch handles the logic for a new pitch inputted by the scorekeeper.
func (s *ScorekeeperService) ProcessPitch(pitch *domain.Pitch) error {
	// TODO: Implement baseball logic (4 balls -> walk, 3 strikes -> out)
	return nil
}
