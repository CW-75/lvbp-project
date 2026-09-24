package rest

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"lvbp-project/backend/internal/core/ports"
)

type BoxscoreHandler struct {
	repo ports.GameRepository
}

func NewBoxscoreHandler(repo ports.GameRepository) *BoxscoreHandler {
	return &BoxscoreHandler{repo: repo}
}

func (h *BoxscoreHandler) GetBoxscore(w http.ResponseWriter, r *http.Request) {
	gameIDStr := chi.URLParam(r, "gameId")
	gameID, err := uuid.Parse(gameIDStr)
	if err != nil {
		http.Error(w, "Invalid game ID format", http.StatusBadRequest)
		return
	}

	game, err := h.repo.GetGameByID(r.Context(), gameID)
	if err != nil {
		http.Error(w, "Game not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(game); err != nil {
		http.Error(w, "Failed to encode boxscore", http.StatusInternalServerError)
	}
}
