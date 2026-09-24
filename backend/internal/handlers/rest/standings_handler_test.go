package rest

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"lvbp-project/backend/internal/core/domain"
)

type mockStandingsRepo struct{}

func (m *mockStandingsRepo) GetStandings(ctx context.Context) ([]*domain.TeamStandings, error) {
	return []*domain.TeamStandings{
		{TeamID: 1, Won: 10, Lost: 5, Pct: 0.667},
	}, nil
}

func TestGetStandings(t *testing.T) {
	repo := &mockStandingsRepo{}
	handler := NewStandingsHandler(repo)

	req, err := http.NewRequest("GET", "/standings", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler.GetStandings(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Contains(t, rr.Body.String(), "0.667")
}
