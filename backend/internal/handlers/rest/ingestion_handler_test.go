package rest

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"lvbp-project/backend/internal/core/domain"
)

type mockScorekeeperService struct {
	mock.Mock
}

func (m *mockScorekeeperService) RecordPitch(ctx context.Context, gameID uuid.UUID, pitch *domain.Pitch) (*domain.Game, *domain.AtBat, error) {
	args := m.Called(ctx, gameID, pitch)
	if args.Get(0) == nil && args.Get(1) == nil {
		return nil, nil, args.Error(2)
	}
	return args.Get(0).(*domain.Game), args.Get(1).(*domain.AtBat), args.Error(2)
}

func (m *mockScorekeeperService) StartAtBat(ctx context.Context, gameID uuid.UUID, batterID string, pitcherID string) (*domain.AtBat, error) {
	return nil, nil
}

func (m *mockScorekeeperService) RecordAtBatResult(ctx context.Context, gameID uuid.UUID, result string) (*domain.Game, *domain.AtBat, error) {
	return nil, nil, nil
}

func (m *mockScorekeeperService) ScoreRuns(ctx context.Context, gameID uuid.UUID, runs int) (*domain.Game, error) {
	return nil, nil
}

func newTestRequest(method, path string, body []byte) *http.Request {
	req := httptest.NewRequest(method, path, nil)
	if body != nil {
		req = httptest.NewRequest(method, path, nil)
		req.Body = &testBody{body: body}
	}
	return req
}

type testBody struct {
	body []byte
	pos  int
}

func (t *testBody) Read(p []byte) (n int, err error) {
	if t.pos >= len(t.body) {
		return 0, nil
	}
	n = copy(p, t.body[t.pos:])
	t.pos += n
	return n, nil
}

func (t *testBody) Close() error { return nil }

func setupHandler(gameID uuid.UUID) (*IngestionHandler, *mockScorekeeperService, *http.Request) {
	mockSvc := new(mockScorekeeperService)
	handler := NewIngestionHandler(mockSvc)

	req := httptest.NewRequest("POST", "/games/"+gameID.String()+"/pitches", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("gameId", gameID.String())
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	return handler, mockSvc, req
}

func TestProcessPitch_Success_Strikeout(t *testing.T) {
	gameID := uuid.New()
	handler, mockSvc, req := setupHandler(gameID)

	game := &domain.Game{ID: gameID, Status: domain.GameStatusInProgress, CurrentInning: 1, IsTopInning: true}
	atBat := &domain.AtBat{ID: uuid.New(), GameID: gameID, Result: strPtr("Strikeout"), OutsRecorded: 1}

	mockSvc.On("RecordPitch", mock.Anything, gameID, mock.MatchedBy(func(p *domain.Pitch) bool {
		return p.PitchResult == "Called Strike"
	})).Return(game, atBat, nil)

	body, _ := json.Marshal(PitchInput{PitchResult: "Called Strike"})
	req.Body = &testBody{body: body}
	rr := httptest.NewRecorder()

	handler.ProcessPitch(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Contains(t, rr.Body.String(), "Strikeout")
	assert.Contains(t, rr.Body.String(), "PitchThrown")
	assert.Contains(t, rr.Body.String(), "AtBatCompleted")
	assert.Contains(t, rr.Body.String(), "InningChanged")

	mockSvc.AssertExpectations(t)
}

func TestProcessPitch_Success_BaseOnBalls(t *testing.T) {
	gameID := uuid.New()
	handler, mockSvc, req := setupHandler(gameID)

	game := &domain.Game{ID: gameID, Status: domain.GameStatusInProgress, CurrentInning: 1, IsTopInning: true}
	atBat := &domain.AtBat{ID: uuid.New(), GameID: gameID, Result: strPtr("Base on Balls"), OutsRecorded: 0}

	mockSvc.On("RecordPitch", mock.Anything, gameID, mock.MatchedBy(func(p *domain.Pitch) bool {
		return p.PitchResult == "Ball"
	})).Return(game, atBat, nil)

	body, _ := json.Marshal(PitchInput{PitchResult: "Ball"})
	req.Body = &testBody{body: body}
	rr := httptest.NewRecorder()

	handler.ProcessPitch(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Contains(t, rr.Body.String(), "Base on Balls")
	assert.Contains(t, rr.Body.String(), "PitchThrown")
	assert.Contains(t, rr.Body.String(), "AtBatCompleted")
	assert.NotContains(t, rr.Body.String(), "InningChanged")

	mockSvc.AssertExpectations(t)
}

func TestProcessPitch_Success_InPlayHit(t *testing.T) {
	gameID := uuid.New()
	handler, mockSvc, req := setupHandler(gameID)

	game := &domain.Game{ID: gameID, Status: domain.GameStatusInProgress, CurrentInning: 1, IsTopInning: true}
	atBat := &domain.AtBat{ID: uuid.New(), GameID: gameID, Result: strPtr("Hit"), OutsRecorded: 0, RunsScored: 1}

	mockSvc.On("RecordPitch", mock.Anything, gameID, mock.MatchedBy(func(p *domain.Pitch) bool {
		return p.PitchResult == "In Play (Hit)"
	})).Return(game, atBat, nil)

	body, _ := json.Marshal(PitchInput{PitchResult: "In Play (Hit)"})
	req.Body = &testBody{body: body}
	rr := httptest.NewRecorder()

	handler.ProcessPitch(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Contains(t, rr.Body.String(), "Hit")
	assert.Contains(t, rr.Body.String(), "ScoreUpdated")

	mockSvc.AssertExpectations(t)
}

func TestProcessPitch_InvalidGameID(t *testing.T) {
	handler := NewIngestionHandler(new(mockScorekeeperService))
	req := httptest.NewRequest("POST", "/games/invalid-id/pitches", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("gameId", "invalid-id")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	body, _ := json.Marshal(PitchInput{PitchResult: "Ball"})
	req.Body = &testBody{body: body}
	rr := httptest.NewRecorder()

	handler.ProcessPitch(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Contains(t, rr.Body.String(), "Invalid game ID")
}

func TestProcessPitch_InvalidJSON(t *testing.T) {
	gameID := uuid.New()
	handler, _, req := setupHandler(gameID)
	req.Body = &testBody{body: []byte("invalid json")}
	rr := httptest.NewRecorder()

	handler.ProcessPitch(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Contains(t, rr.Body.String(), "Invalid JSON")
}

func TestProcessPitch_InvalidPitchResult(t *testing.T) {
	gameID := uuid.New()
	handler, _, req := setupHandler(gameID)

	body, _ := json.Marshal(PitchInput{PitchResult: "Invalid Result"})
	req.Body = &testBody{body: body}
	rr := httptest.NewRecorder()

	handler.ProcessPitch(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Contains(t, rr.Body.String(), "Invalid pitchResult")
}

func TestProcessPitch_GameNotInProgress(t *testing.T) {
	gameID := uuid.New()
	handler, mockSvc, req := setupHandler(gameID)

	mockSvc.On("RecordPitch", mock.Anything, gameID, mock.Anything).Return(nil, nil, assert.AnError)

	body, _ := json.Marshal(PitchInput{PitchResult: "Ball"})
	req.Body = &testBody{body: body}
	rr := httptest.NewRecorder()

	handler.ProcessPitch(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	mockSvc.AssertExpectations(t)
}

func TestProcessPitch_GameConflict(t *testing.T) {
	gameID := uuid.New()
	handler, mockSvc, req := setupHandler(gameID)

	mockSvc.On("RecordPitch", mock.Anything, gameID, mock.Anything).Return(nil, nil, assert.AnError)

	body, _ := json.Marshal(PitchInput{PitchResult: "Ball"})
	req.Body = &testBody{body: body}
	rr := httptest.NewRecorder()

	handler.ProcessPitch(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestProcessPitch_AtBatCompleted(t *testing.T) {
	gameID := uuid.New()
	handler, mockSvc, req := setupHandler(gameID)

	game := &domain.Game{ID: gameID, Status: domain.GameStatusInProgress}
	atBat := &domain.AtBat{ID: uuid.New(), GameID: gameID, Result: strPtr("Strikeout"), OutsRecorded: 1}

	mockSvc.On("RecordPitch", mock.Anything, gameID, mock.Anything).Return(game, atBat, nil)

	body, _ := json.Marshal(PitchInput{PitchResult: "Called Strike"})
	req.Body = &testBody{body: body}
	rr := httptest.NewRecorder()

	handler.ProcessPitch(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	var resp PitchResponse
	require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &resp))
	assert.Equal(t, "Strikeout", *resp.AtBat.Result)
	assert.Equal(t, 1, resp.AtBat.OutsRecorded)
	assert.Contains(t, resp.Events, "PitchThrown")
	assert.Contains(t, resp.Events, "AtBatCompleted")
	assert.Contains(t, resp.Events, "InningChanged")

	mockSvc.AssertExpectations(t)
}

func TestProcessPitch_WithOptionalFields(t *testing.T) {
	gameID := uuid.New()
	handler, mockSvc, req := setupHandler(gameID)

	game := &domain.Game{ID: gameID, Status: domain.GameStatusInProgress}
	atBat := &domain.AtBat{ID: uuid.New(), GameID: gameID, Result: strPtr("In Play (Hit)")}

	capturedPitch := new(domain.Pitch)
	mockSvc.On("RecordPitch", mock.Anything, gameID, mock.MatchedBy(func(p *domain.Pitch) bool {
		capturedPitch = p
		return true
	})).Return(game, atBat, nil)

	coordX := 0.75
	coordY := -0.25
	vel := 98.5
	body, _ := json.Marshal(PitchInput{
		PitchResult: "In Play (Hit)",
		CoordinateX: &coordX,
		CoordinateY: &coordY,
		VelocityMPH: &vel,
	})
	req.Body = &testBody{body: body}
	rr := httptest.NewRecorder()

	handler.ProcessPitch(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, 0.75, capturedPitch.CoordinateX)
	assert.Equal(t, -0.25, capturedPitch.CoordinateY)
	assert.Equal(t, 98.5, *capturedPitch.VelocityMPH)

	mockSvc.AssertExpectations(t)
}

func strPtr(s string) *string {
	return &s
}