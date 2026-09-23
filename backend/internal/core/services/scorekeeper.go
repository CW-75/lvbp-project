package services

import (
	"context"
	"errors"

	"github.com/google/uuid"

	"lvbp-project/backend/internal/core/domain"
	"lvbp-project/backend/internal/core/ports"
)

type scorekeeperService struct {
	gameRepo ports.GameRepository
	eventBus ports.EventBus
}

func NewScorekeeperService(gameRepo ports.GameRepository, eventBus ports.EventBus) ports.ScorekeeperService {
	return &scorekeeperService{
		gameRepo: gameRepo,
		eventBus: eventBus,
	}
}

func (s *scorekeeperService) StartAtBat(ctx context.Context, gameID uuid.UUID, batterID string, pitcherID string) (*domain.AtBat, error) {
	game, err := s.gameRepo.GetGameByID(ctx, gameID)
	if err != nil {
		return nil, err
	}

	if game.Status != domain.GameStatusInProgress {
		return nil, errors.New("game is not in progress")
	}

	atBat := &domain.AtBat{
		GameID:      gameID,
		Inning:      game.CurrentInning,
		IsTopInning: game.IsTopInning,
		BatterID:    batterID,
		PitcherID:   pitcherID,
	}

	atBat, err = s.gameRepo.CreateAtBat(ctx, atBat)
	if err != nil {
		return nil, err
	}

	// Publish event
	_ = s.eventBus.PublishEvent("game:"+gameID.String(), map[string]interface{}{
		"type":  "AtBatStarted",
		"atBat": atBat,
	})

	return atBat, nil
}

func (s *scorekeeperService) RecordPitch(ctx context.Context, gameID uuid.UUID, pitch *domain.Pitch) (*domain.Game, *domain.AtBat, error) {
	game, err := s.gameRepo.GetGameByID(ctx, gameID)
	if err != nil {
		return nil, nil, err
	}
	
	if game.Status != domain.GameStatusInProgress {
		return nil, nil, errors.New("game is not in progress")
	}

	atBat, err := s.gameRepo.GetCurrentAtBat(ctx, gameID)
	if err != nil {
		return nil, nil, err
	}
	if atBat.Result != nil {
		return nil, nil, errors.New("current at-bat is already completed")
	}

	pitches, err := s.gameRepo.GetPitchesForAtBat(ctx, atBat.ID)
	if err != nil {
		return nil, nil, err
	}

	outs, err := s.gameRepo.GetOutsForInning(ctx, game.ID, game.CurrentInning, game.IsTopInning)
	if err != nil {
		return nil, nil, err
	}

	balls, strikes := calculateCount(pitches)

	pitch.AtBatID = atBat.ID
	pitch.BallsBefore = balls
	pitch.StrikesBefore = strikes
	pitch.OutsBefore = outs
	pitch.PitchNumber = len(pitches) + 1

	pitch, err = s.gameRepo.CreatePitch(ctx, pitch)
	if err != nil {
		return nil, nil, err
	}

	_ = s.eventBus.PublishEvent("game:"+gameID.String(), map[string]interface{}{
		"type":  "PitchThrown",
		"pitch": pitch,
	})

	// Evaluate Baseball Rules
	newBalls, newStrikes := balls, strikes
	atBatResult := ""
	outsRecorded := 0
	isCompleted := false

	switch pitch.PitchResult {
	case "Ball":
		newBalls++
		if newBalls == 4 {
			atBatResult = "Base on Balls"
			isCompleted = true
		}
	case "Called Strike", "Swinging Strike":
		newStrikes++
		if newStrikes == 3 {
			atBatResult = "Strikeout"
			outsRecorded = 1
			isCompleted = true
		}
	case "Foul Ball":
		if newStrikes < 2 {
			newStrikes++
		}
	case "In Play (Out)":
		atBatResult = "Out"
		outsRecorded = 1
		isCompleted = true
	case "In Play (Hit)":
		atBatResult = "Hit"
		isCompleted = true
	}

	if isCompleted {
		atBat.Result = &atBatResult
		atBat.OutsRecorded = outsRecorded
		atBat, err = s.gameRepo.UpdateAtBatResult(ctx, atBat)
		if err != nil {
			return nil, nil, err
		}

		_ = s.eventBus.PublishEvent("game:"+gameID.String(), map[string]interface{}{
			"type":  "AtBatCompleted",
			"atBat": atBat,
		})

		outs += outsRecorded
		if outs >= 3 {
			// Inning change
			if game.IsTopInning {
				game.IsTopInning = false
			} else {
				game.IsTopInning = true
				game.CurrentInning++
			}
			game, err = s.gameRepo.UpdateGame(ctx, game)
			if err != nil {
				return nil, nil, err
			}

			_ = s.eventBus.PublishEvent("game:"+gameID.String(), map[string]interface{}{
				"type":  "InningChanged",
				"game":  game,
			})
		}
	}

	return game, atBat, nil
}

func (s *scorekeeperService) RecordAtBatResult(ctx context.Context, gameID uuid.UUID, result string) (*domain.Game, *domain.AtBat, error) {
	// For manual scorekeeping overrides (e.g., intentionally walking a batter without pitches, or double plays)
	// Simplified for this phase.
	return nil, nil, errors.New("not implemented yet")
}

func (s *scorekeeperService) ScoreRuns(ctx context.Context, gameID uuid.UUID, runs int) (*domain.Game, error) {
	if runs <= 0 {
		return nil, errors.New("runs must be greater than 0")
	}

	game, err := s.gameRepo.GetGameByID(ctx, gameID)
	if err != nil {
		return nil, err
	}

	if game.Status != domain.GameStatusInProgress {
		return nil, errors.New("game is not in progress")
	}

	// Update game score depending on who is batting
	if game.IsTopInning {
		game.AwayScore += runs // Top of inning = Away team bats
	} else {
		game.HomeScore += runs // Bottom of inning = Home team bats
	}

	game, err = s.gameRepo.UpdateGame(ctx, game)
	if err != nil {
		return nil, err
	}

	// Update the AtBat to record RBI/Runs Scored on the play
	atBat, err := s.gameRepo.GetCurrentAtBat(ctx, gameID)
	if err == nil && atBat != nil && atBat.Result == nil {
		atBat.RunsScored += runs
		_, _ = s.gameRepo.UpdateAtBatResult(ctx, atBat) // We ignore error here, game score is priority
	}

	_ = s.eventBus.PublishEvent("game:"+gameID.String(), map[string]interface{}{
		"type": "ScoreUpdated",
		"game": game,
		"runs_added": runs,
	})

	return game, nil
}

func calculateCount(pitches []*domain.Pitch) (balls int, strikes int) {
	for _, p := range pitches {
		switch p.PitchResult {
		case "Ball":
			balls++
		case "Called Strike", "Swinging Strike":
			strikes++
		case "Foul Ball":
			if strikes < 2 {
				strikes++
			}
		}
	}
	return balls, strikes
}
