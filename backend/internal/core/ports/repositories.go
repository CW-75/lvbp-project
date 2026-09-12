package ports

import "lvbp-project/backend/internal/core/domain"

type GameRepository interface {
	GetGameByID(id int) (*domain.Game, error)
	UpdateGame(game *domain.Game) error
}

type StandingsRepository interface {
	GetStandings() ([]*domain.TeamStandings, error)
}
