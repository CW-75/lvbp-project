package repositories

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	
	"lvbp-project/backend/internal/core/domain"
	"lvbp-project/backend/internal/infrastructure/db"
)

type GameRepository struct {
	queries *db.Queries
	db      *pgxpool.Pool
}

func NewGameRepository(dbPool *pgxpool.Pool) *GameRepository {
	return &GameRepository{
		queries: db.New(dbPool),
		db:      dbPool,
	}
}

func (r *GameRepository) GetGameByID(ctx context.Context, id uuid.UUID) (*domain.Game, error) {
	row, err := r.queries.GetGame(ctx, id)
	if err != nil {
		return nil, err
	}
	
	return &domain.Game{
		ID:            row.ID,
		HomeTeamID:    row.HomeTeamID.String,
		AwayTeamID:    row.AwayTeamID.String,
		Status:        domain.GameStatus(row.Status.GameStatus),
		HomeScore:     int(row.HomeScore.Int32),
		AwayScore:     int(row.AwayScore.Int32),
		CurrentInning: int(row.CurrentInning.Int32),
		IsTopInning:   row.IsTopInning.Bool,
		StartTime:     row.StartTime,
		CreatedAt:     row.CreatedAt.Time,
	}, nil
}

func (r *GameRepository) ListActiveGames(ctx context.Context) ([]*domain.Game, error) {
	rows, err := r.queries.ListActiveGames(ctx)
	if err != nil {
		return nil, err
	}
	
	games := make([]*domain.Game, len(rows))
	for i, row := range rows {
		games[i] = &domain.Game{
			ID:            row.ID,
			HomeTeamID:    row.HomeTeamID.String,
			AwayTeamID:    row.AwayTeamID.String,
			Status:        domain.GameStatus(row.Status.GameStatus),
			HomeScore:     int(row.HomeScore.Int32),
			AwayScore:     int(row.AwayScore.Int32),
			CurrentInning: int(row.CurrentInning.Int32),
			IsTopInning:   row.IsTopInning.Bool,
			StartTime:     row.StartTime,
			CreatedAt:     row.CreatedAt.Time,
		}
	}
	return games, nil
}

func (r *GameRepository) CreateGame(ctx context.Context, game *domain.Game) (*domain.Game, error) {
	row, err := r.queries.CreateGame(ctx, db.CreateGameParams{
		HomeTeamID: pgtype.Text{String: game.HomeTeamID, Valid: true},
		AwayTeamID: pgtype.Text{String: game.AwayTeamID, Valid: true},
		StartTime:  game.StartTime,
	})
	if err != nil {
		return nil, err
	}
	
	game.ID = row.ID
	game.Status = domain.GameStatus(row.Status.GameStatus)
	game.HomeScore = int(row.HomeScore.Int32)
	game.AwayScore = int(row.AwayScore.Int32)
	game.CurrentInning = int(row.CurrentInning.Int32)
	game.IsTopInning = row.IsTopInning.Bool
	game.CreatedAt = row.CreatedAt.Time
	
	return game, nil
}

func (r *GameRepository) UpdateGame(ctx context.Context, game *domain.Game) (*domain.Game, error) {
	row, err := r.queries.UpdateGame(ctx, db.UpdateGameParams{
		ID:            game.ID,
		Status:        db.NullGameStatus{GameStatus: db.GameStatus(game.Status), Valid: true},
		HomeScore:     pgtype.Int4{Int32: int32(game.HomeScore), Valid: true},
		AwayScore:     pgtype.Int4{Int32: int32(game.AwayScore), Valid: true},
		CurrentInning: pgtype.Int4{Int32: int32(game.CurrentInning), Valid: true},
		IsTopInning:   pgtype.Bool{Bool: game.IsTopInning, Valid: true},
	})
	if err != nil {
		return nil, err
	}
	
	game.Status = domain.GameStatus(row.Status.GameStatus)
	return game, nil
}

// Add dummy implementations to satisfy interface for now
func (r *GameRepository) CreateAtBat(ctx context.Context, atBat *domain.AtBat) (*domain.AtBat, error) { return nil, nil }
func (r *GameRepository) UpdateAtBatResult(ctx context.Context, atBat *domain.AtBat) (*domain.AtBat, error) { return nil, nil }
func (r *GameRepository) GetCurrentAtBat(ctx context.Context, gameID uuid.UUID) (*domain.AtBat, error) { return nil, nil }
func (r *GameRepository) CreatePitch(ctx context.Context, pitch *domain.Pitch) (*domain.Pitch, error) { return nil, nil }
func (r *GameRepository) GetPitchesForAtBat(ctx context.Context, atBatID uuid.UUID) ([]*domain.Pitch, error) { return nil, nil }
