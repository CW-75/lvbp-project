package rest

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"lvbp-project/backend/internal/core/domain"
)

type mockGameRepo struct{}

func (m *mockGameRepo) GetGameByID(ctx context.Context, id uuid.UUID) (*domain.Game, error) {
	return &domain.Game{
		ID: id,
	}, nil
}
func (m *mockGameRepo) ListActiveGames(ctx context.Context) ([]*domain.Game, error) { return nil, nil }
func (m *mockGameRepo) CreateGame(ctx context.Context, game *domain.Game) (*domain.Game, error) { return nil, nil }
func (m *mockGameRepo) UpdateGame(ctx context.Context, game *domain.Game) (*domain.Game, error) { return nil, nil }
func (m *mockGameRepo) CreateAtBat(ctx context.Context, atBat *domain.AtBat) (*domain.AtBat, error) { return nil, nil }
func (m *mockGameRepo) UpdateAtBatResult(ctx context.Context, atBat *domain.AtBat) (*domain.AtBat, error) { return nil, nil }
func (m *mockGameRepo) GetCurrentAtBat(ctx context.Context, gameID uuid.UUID) (*domain.AtBat, error) { return nil, nil }
func (m *mockGameRepo) CreatePitch(ctx context.Context, pitch *domain.Pitch) (*domain.Pitch, error) { return nil, nil }
func (m *mockGameRepo) GetPitchesForAtBat(ctx context.Context, atBatID uuid.UUID) ([]*domain.Pitch, error) { return nil, nil }
func (m *mockGameRepo) GetOutsForInning(ctx context.Context, gameID uuid.UUID, inning int, isTopInning bool) (int, error) { return 0, nil }

func TestGetBoxscore(t *testing.T) {
	repo := &mockGameRepo{}
	handler := NewBoxscoreHandler(repo)

	gameID := uuid.New()
	req, err := http.NewRequest("GET", "/games/"+gameID.String()+"/boxscore", nil)
	assert.NoError(t, err)

	// Utilizar chi para simular el URL param
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("gameId", gameID.String())
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	rr := httptest.NewRecorder()
	handler.GetBoxscore(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Contains(t, rr.Body.String(), gameID.String())
}
