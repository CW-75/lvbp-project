package rest

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"lvbp-project/backend/internal/core/domain"
	"lvbp-project/backend/internal/core/ports"
)

type IngestionHandler struct {
	scorekeeperService ports.ScorekeeperService
}

func NewIngestionHandler(scorekeeperService ports.ScorekeeperService) *IngestionHandler {
	return &IngestionHandler{scorekeeperService: scorekeeperService}
}

type PitchInput struct {
	PitchResult string   `json:"pitchResult"`
	CoordinateX *float64 `json:"coordinateX,omitempty"`
	CoordinateY *float64 `json:"coordinateY,omitempty"`
	VelocityMPH *float64 `json:"velocityMPH,omitempty"`
}

type PitchResponse struct {
	Game    *domain.Game   `json:"game"`
	AtBat   *domain.AtBat  `json:"atBat"`
	Events  []string       `json:"events,omitempty"`
}

func (h *IngestionHandler) ProcessPitch(w http.ResponseWriter, r *http.Request) {
	gameIDStr := chi.URLParam(r, "gameId")
	gameID, err := uuid.Parse(gameIDStr)
	if err != nil {
		http.Error(w, "Invalid game ID format", http.StatusBadRequest)
		return
	}

	var input PitchInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	if !isValidPitchResult(input.PitchResult) {
		http.Error(w, "Invalid pitchResult value", http.StatusBadRequest)
		return
	}

	pitch := &domain.Pitch{
		PitchResult: input.PitchResult,
		CoordinateX: derefFloat64(input.CoordinateX, 0),
		CoordinateY: derefFloat64(input.CoordinateY, 0),
		VelocityMPH: input.VelocityMPH,
	}

	game, atBat, err := h.scorekeeperService.RecordPitch(r.Context(), gameID, pitch)
	if err != nil {
		switch err.Error() {
		case "game is not in progress":
			http.Error(w, err.Error(), http.StatusConflict)
		case "current at-bat is already completed":
			http.Error(w, err.Error(), http.StatusConflict)
		default:
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	response := PitchResponse{
		Game:  game,
		AtBat: atBat,
		Events: []string{"PitchThrown"},
	}
	if atBat.Result != nil {
		response.Events = append(response.Events, "AtBatCompleted")
		if atBat.OutsRecorded > 0 {
			response.Events = append(response.Events, "InningChanged")
		}
		if atBat.RunsScored > 0 {
			response.Events = append(response.Events, "ScoreUpdated")
		}
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func isValidPitchResult(result string) bool {
	validResults := map[string]bool{
		"Ball":            true,
		"Called Strike":   true,
		"Swinging Strike": true,
		"Foul Ball":       true,
		"In Play (Out)":   true,
		"In Play (Hit)":   true,
	}
	return validResults[result]
}

func derefFloat64(ptr *float64, defaultVal float64) float64 {
	if ptr != nil {
		return *ptr
	}
	return defaultVal
}