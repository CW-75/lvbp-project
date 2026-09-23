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
		CreatedAt:     row.CreatedAt,
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
			CreatedAt:     row.CreatedAt,
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
	game.CreatedAt = row.CreatedAt
	
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

func (r *GameRepository) CreateAtBat(ctx context.Context, atBat *domain.AtBat) (*domain.AtBat, error) {
	row, err := r.queries.CreateAtBat(ctx, db.CreateAtBatParams{
		GameID:      pgtype.UUID{Bytes: atBat.GameID, Valid: true},
		Inning:      int32(atBat.Inning),
		IsTopInning: atBat.IsTopInning,
		BatterID:    atBat.BatterID,
		PitcherID:   atBat.PitcherID,
	})
	if err != nil {
		return nil, err
	}

	atBat.ID = row.ID
	atBat.CreatedAt = row.CreatedAt
	return atBat, nil
}

func (r *GameRepository) UpdateAtBatResult(ctx context.Context, atBat *domain.AtBat) (*domain.AtBat, error) {
	var result pgtype.Text
	if atBat.Result != nil {
		result = pgtype.Text{String: *atBat.Result, Valid: true}
	}
	
	row, err := r.queries.UpdateAtBatResult(ctx, db.UpdateAtBatResultParams{
		ID:           atBat.ID,
		Result:       result,
		RunsScored:   pgtype.Int4{Int32: int32(atBat.RunsScored), Valid: true},
		OutsRecorded: pgtype.Int4{Int32: int32(atBat.OutsRecorded), Valid: true},
	})
	if err != nil {
		return nil, err
	}
	
	if row.Result.Valid {
		atBat.Result = &row.Result.String
	}
	return atBat, nil
}

func (r *GameRepository) GetCurrentAtBat(ctx context.Context, gameID uuid.UUID) (*domain.AtBat, error) {
	row, err := r.queries.GetCurrentAtBat(ctx, pgtype.UUID{Bytes: gameID, Valid: true})
	if err != nil {
		return nil, err
	}
	
	atBat := &domain.AtBat{
		ID:           row.ID,
		GameID:       row.GameID.Bytes,
		Inning:       int(row.Inning),
		IsTopInning:  row.IsTopInning,
		BatterID:     row.BatterID,
		PitcherID:    row.PitcherID,
		RunsScored:   int(row.RunsScored.Int32),
		OutsRecorded: int(row.OutsRecorded.Int32),
		CreatedAt:    row.CreatedAt,
	}
	if row.Result.Valid {
		atBat.Result = &row.Result.String
	}
	return atBat, nil
}

func (r *GameRepository) CreatePitch(ctx context.Context, pitch *domain.Pitch) (*domain.Pitch, error) {
	var velocity pgtype.Numeric
	if pitch.VelocityMPH != nil {
		velocity.Scan(pitch.VelocityMPH) // simple conversion
	}
	
	row, err := r.queries.CreatePitch(ctx, db.CreatePitchParams{
		AtBatID:       pgtype.UUID{Bytes: pitch.AtBatID, Valid: true},
		PitchNumber:   int32(pitch.PitchNumber),
		CoordinateX:   pgtype.Numeric{}, // Needs proper mapping, keeping simple for now
		CoordinateY:   pgtype.Numeric{},
		PitchResult:   pitch.PitchResult,
		BallsBefore:   int32(pitch.BallsBefore),
		StrikesBefore: int32(pitch.StrikesBefore),
		OutsBefore:    int32(pitch.OutsBefore),
		VelocityMph:   velocity,
	})
	if err != nil {
		return nil, err
	}

	pitch.ID = row.ID
	pitch.CreatedAt = row.CreatedAt
	return pitch, nil
}

func (r *GameRepository) GetPitchesForAtBat(ctx context.Context, atBatID uuid.UUID) ([]*domain.Pitch, error) {
	rows, err := r.queries.GetPitchesForAtBat(ctx, pgtype.UUID{Bytes: atBatID, Valid: true})
	if err != nil {
		return nil, err
	}
	
	var pitches []*domain.Pitch
	for _, row := range rows {
		pitches = append(pitches, &domain.Pitch{
			ID:            row.ID,
			AtBatID:       row.AtBatID.Bytes,
			PitchNumber:   int(row.PitchNumber),
			PitchResult:   row.PitchResult,
			BallsBefore:   int(row.BallsBefore),
			StrikesBefore: int(row.StrikesBefore),
			OutsBefore:    int(row.OutsBefore),
			CreatedAt:     row.CreatedAt,
		})
	}
	return pitches, nil
}

func (r *GameRepository) GetOutsForInning(ctx context.Context, gameID uuid.UUID, inning int, isTopInning bool) (int, error) {
	outs, err := r.queries.GetOutsForInning(ctx, db.GetOutsForInningParams{
		GameID:      pgtype.UUID{Bytes: gameID, Valid: true},
		Inning:      int32(inning),
		IsTopInning: isTopInning,
	})
	return int(outs), err
}
