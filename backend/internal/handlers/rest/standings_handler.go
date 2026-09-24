package rest

import (
	"encoding/json"
	"net/http"

	"lvbp-project/backend/internal/core/ports"
)

type StandingsHandler struct {
	repo ports.StandingsRepository
}

func NewStandingsHandler(repo ports.StandingsRepository) *StandingsHandler {
	return &StandingsHandler{repo: repo}
}

func (h *StandingsHandler) GetStandings(w http.ResponseWriter, r *http.Request) {
	standings, err := h.repo.GetStandings(r.Context())
	if err != nil {
		http.Error(w, "Failed to retrieve standings", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(standings); err != nil {
		http.Error(w, "Failed to encode standings", http.StatusInternalServerError)
	}
}
