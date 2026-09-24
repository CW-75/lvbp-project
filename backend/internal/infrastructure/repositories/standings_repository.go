package repositories

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"lvbp-project/backend/internal/core/domain"
	"lvbp-project/backend/internal/infrastructure/db"
)

type StandingsRepository struct {
	queries *db.Queries
	db      *pgxpool.Pool
}

func NewStandingsRepository(dbPool *pgxpool.Pool) *StandingsRepository {
	return &StandingsRepository{
		queries: db.New(dbPool),
		db:      dbPool,
	}
}

func (r *StandingsRepository) GetStandings(ctx context.Context) ([]*domain.TeamStandings, error) {
	rows, err := r.queries.GetStandings(ctx)
	if err != nil {
		return nil, err
	}

	standings := make([]*domain.TeamStandings, len(rows))
	for i, row := range rows {
		standings[i] = &domain.TeamStandings{
			TeamID:      int(row.TeamID.Int32),
			GamesPlayed: int(row.GamesPlayed.Int32),
			Won:         int(row.Won.Int32),
			Lost:        int(row.Lost.Int32),
			Pct:         float64(row.Pct.Float64),
			GamesBehind: float64(row.GamesBehind.Float64),
		}
	}
	return standings, nil
}