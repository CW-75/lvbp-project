package services_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"lvbp-project/backend/internal/core/domain"
	"lvbp-project/backend/internal/core/services"
)

// MockGameRepository
type MockGameRepository struct {
	mock.Mock
}

func (m *MockGameRepository) GetGameByID(ctx context.Context, id uuid.UUID) (*domain.Game, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Game), args.Error(1)
}

func (m *MockGameRepository) ListActiveGames(ctx context.Context) ([]*domain.Game, error) { return nil, nil }
func (m *MockGameRepository) CreateGame(ctx context.Context, game *domain.Game) (*domain.Game, error) { return nil, nil }

func (m *MockGameRepository) UpdateGame(ctx context.Context, game *domain.Game) (*domain.Game, error) {
	args := m.Called(ctx, game)
	return args.Get(0).(*domain.Game), args.Error(1)
}

func (m *MockGameRepository) CreateAtBat(ctx context.Context, atBat *domain.AtBat) (*domain.AtBat, error) {
	args := m.Called(ctx, atBat)
	return args.Get(0).(*domain.AtBat), args.Error(1)
}

func (m *MockGameRepository) UpdateAtBatResult(ctx context.Context, atBat *domain.AtBat) (*domain.AtBat, error) {
	args := m.Called(ctx, atBat)
	return args.Get(0).(*domain.AtBat), args.Error(1)
}

func (m *MockGameRepository) GetCurrentAtBat(ctx context.Context, gameID uuid.UUID) (*domain.AtBat, error) {
	args := m.Called(ctx, gameID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.AtBat), args.Error(1)
}

func (m *MockGameRepository) CreatePitch(ctx context.Context, pitch *domain.Pitch) (*domain.Pitch, error) {
	args := m.Called(ctx, pitch)
	return args.Get(0).(*domain.Pitch), args.Error(1)
}

func (m *MockGameRepository) GetPitchesForAtBat(ctx context.Context, atBatID uuid.UUID) ([]*domain.Pitch, error) {
	args := m.Called(ctx, atBatID)
	return args.Get(0).([]*domain.Pitch), args.Error(1)
}

func (m *MockGameRepository) GetOutsForInning(ctx context.Context, gameID uuid.UUID, inning int, isTopInning bool) (int, error) {
	args := m.Called(ctx, gameID, inning, isTopInning)
	return args.Int(0), args.Error(1)
}

// MockEventBus
type MockEventBus struct {
	mock.Mock
}

func (m *MockEventBus) PublishEvent(topic string, payload interface{}) error {
	args := m.Called(topic, payload)
	return args.Error(0)
}
func (m *MockEventBus) Subscribe(topic string) (<-chan string, error) { return nil, nil }

// Tests
func TestScorekeeperService_RecordPitch_Strikeout(t *testing.T) {
	gameRepo := new(MockGameRepository)
	eventBus := new(MockEventBus)
	svc := services.NewScorekeeperService(gameRepo, eventBus)
	ctx := context.Background()
	gameID := uuid.New()
	atBatID := uuid.New()

	game := &domain.Game{
		ID:            gameID,
		Status:        domain.GameStatusInProgress,
		CurrentInning: 1,
		IsTopInning:   true,
	}
	atBat := &domain.AtBat{
		ID:     atBatID,
		GameID: gameID,
	}

	// We simulate that the batter already has 2 strikes
	pitches := []*domain.Pitch{
		{PitchResult: "Called Strike"},
		{PitchResult: "Swinging Strike"},
	}

	gameRepo.On("GetGameByID", ctx, gameID).Return(game, nil)
	gameRepo.On("GetCurrentAtBat", ctx, gameID).Return(atBat, nil)
	gameRepo.On("GetPitchesForAtBat", ctx, atBatID).Return(pitches, nil)
	gameRepo.On("GetOutsForInning", ctx, gameID, 1, true).Return(0, nil)
	
	newPitch := &domain.Pitch{PitchResult: "Called Strike"} // This should be strike 3
	
	// CreatePitch will be called
	gameRepo.On("CreatePitch", ctx, mock.MatchedBy(func(p *domain.Pitch) bool {
		return p.StrikesBefore == 2 && p.BallsBefore == 0 && p.OutsBefore == 0
	})).Return(newPitch, nil)

	// Since it's strike 3, it should UpdateAtBatResult with Strikeout
	gameRepo.On("UpdateAtBatResult", ctx, mock.MatchedBy(func(ab *domain.AtBat) bool {
		return ab.Result != nil && *ab.Result == "Strikeout" && ab.OutsRecorded == 1
	})).Return(atBat, nil)

	// Events should be published
	eventBus.On("PublishEvent", "game:"+gameID.String(), mock.Anything).Return(nil)

	resGame, resAtBat, err := svc.RecordPitch(ctx, gameID, newPitch)
	assert.NoError(t, err)
	assert.NotNil(t, resGame)
	assert.Equal(t, "Strikeout", *resAtBat.Result)
	assert.Equal(t, 1, resAtBat.OutsRecorded)
	
	gameRepo.AssertExpectations(t)
	eventBus.AssertExpectations(t)
}

func TestScorekeeperService_RecordPitch_BaseOnBalls(t *testing.T) {
	gameRepo := new(MockGameRepository)
	eventBus := new(MockEventBus)
	svc := services.NewScorekeeperService(gameRepo, eventBus)
	ctx := context.Background()
	gameID := uuid.New()
	atBatID := uuid.New()

	gameRepo.On("GetGameByID", ctx, gameID).Return(&domain.Game{ID: gameID, Status: domain.GameStatusInProgress, CurrentInning: 1}, nil)
	gameRepo.On("GetCurrentAtBat", ctx, gameID).Return(&domain.AtBat{ID: atBatID, GameID: gameID}, nil)
	gameRepo.On("GetPitchesForAtBat", ctx, atBatID).Return([]*domain.Pitch{
		{PitchResult: "Ball"}, {PitchResult: "Ball"}, {PitchResult: "Ball"},
	}, nil)
	gameRepo.On("GetOutsForInning", ctx, gameID, 1, false).Return(0, nil)
	
	gameRepo.On("CreatePitch", ctx, mock.Anything).Return(&domain.Pitch{PitchResult: "Ball"}, nil)
	
	resultStr := "Base on Balls"
	gameRepo.On("UpdateAtBatResult", ctx, mock.MatchedBy(func(ab *domain.AtBat) bool {
		return ab.Result != nil && *ab.Result == "Base on Balls" && ab.OutsRecorded == 0
	})).Return(&domain.AtBat{Result: &resultStr}, nil)

	eventBus.On("PublishEvent", mock.Anything, mock.Anything).Return(nil)

	_, resAtBat, err := svc.RecordPitch(ctx, gameID, &domain.Pitch{PitchResult: "Ball"})
	
	assert.NoError(t, err)
	assert.Equal(t, "Base on Balls", *resAtBat.Result)
}

func TestScorekeeperService_RecordPitch_InningTransition(t *testing.T) {
	gameRepo := new(MockGameRepository)
	eventBus := new(MockEventBus)
	svc := services.NewScorekeeperService(gameRepo, eventBus)
	ctx := context.Background()
	gameID := uuid.New()

	// Top of the 1st inning, already 2 outs
	game := &domain.Game{ID: gameID, Status: domain.GameStatusInProgress, CurrentInning: 1, IsTopInning: true}
	atBat := &domain.AtBat{ID: uuid.New(), GameID: gameID}

	gameRepo.On("GetGameByID", ctx, gameID).Return(game, nil)
	gameRepo.On("GetCurrentAtBat", ctx, gameID).Return(atBat, nil)
	gameRepo.On("GetPitchesForAtBat", ctx, atBat.ID).Return([]*domain.Pitch{}, nil)
	gameRepo.On("GetOutsForInning", ctx, gameID, 1, true).Return(2, nil) // 2 Outs previously!
	
	gameRepo.On("CreatePitch", ctx, mock.Anything).Return(&domain.Pitch{PitchResult: "In Play (Out)"}, nil)
	
	gameRepo.On("UpdateAtBatResult", ctx, mock.Anything).Return(atBat, nil)
	
	// Expect UpdateGame to be called to flip the inning!
	gameRepo.On("UpdateGame", ctx, mock.MatchedBy(func(g *domain.Game) bool {
		return g.CurrentInning == 1 && g.IsTopInning == false // Flipped to bottom of the 1st
	})).Return(game, nil)

	eventBus.On("PublishEvent", mock.Anything, mock.Anything).Return(nil)

	resGame, _, err := svc.RecordPitch(ctx, gameID, &domain.Pitch{PitchResult: "In Play (Out)"})
	
	assert.NoError(t, err)
	assert.False(t, resGame.IsTopInning) // Top changed to bottom
	
	gameRepo.AssertExpectations(t)
}

func TestScorekeeperService_ScoreRuns(t *testing.T) {
	gameRepo := new(MockGameRepository)
	eventBus := new(MockEventBus)
	svc := services.NewScorekeeperService(gameRepo, eventBus)
	ctx := context.Background()
	gameID := uuid.New()

	game := &domain.Game{ID: gameID, Status: domain.GameStatusInProgress, CurrentInning: 1, IsTopInning: true, AwayScore: 0}
	atBat := &domain.AtBat{ID: uuid.New(), GameID: gameID, RunsScored: 0}

	gameRepo.On("GetGameByID", ctx, gameID).Return(game, nil)
	gameRepo.On("UpdateGame", ctx, mock.MatchedBy(func(g *domain.Game) bool {
		return g.AwayScore == 2 // 2 runs scored by Away Team (since it's Top Inning)
	})).Return(game, nil)
	
	gameRepo.On("GetCurrentAtBat", ctx, gameID).Return(atBat, nil)
	gameRepo.On("UpdateAtBatResult", ctx, mock.MatchedBy(func(ab *domain.AtBat) bool {
		return ab.RunsScored == 2
	})).Return(atBat, nil)

	eventBus.On("PublishEvent", mock.Anything, mock.Anything).Return(nil)

	resGame, err := svc.ScoreRuns(ctx, gameID, 2)
	assert.NoError(t, err)
	assert.Equal(t, 2, resGame.AwayScore)
	
	gameRepo.AssertExpectations(t)
}
