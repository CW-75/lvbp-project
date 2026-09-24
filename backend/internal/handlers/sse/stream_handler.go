package sse

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	redisPkg "lvbp-project/backend/internal/pkg/redis"
)

type StreamHandler struct {
	redisClient *redisPkg.Client
}

func NewStreamHandler(redisClient *redisPkg.Client) *StreamHandler {
	return &StreamHandler{redisClient: redisClient}
}

func (h *StreamHandler) StreamGameEvents(w http.ResponseWriter, r *http.Request) {
	// NFR-3: Disable buffering and caching for SSE
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("X-Accel-Buffering", "no")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	gameID := chi.URLParam(r, "gameId")
	if gameID == "" {
		http.Error(w, "gameId is required", http.StatusBadRequest)
		return
	}

	channel := fmt.Sprintf("game:%s:events", gameID)
	pubsub := h.redisClient.GetDB().Subscribe(r.Context(), channel)
	defer pubsub.Close()

	// Ensure connection is established
	_, err := pubsub.Receive(r.Context())
	if err != nil {
		http.Error(w, "Failed to subscribe to events", http.StatusInternalServerError)
		return
	}

	// Send an initial connected event
	fmt.Fprintf(w, "event: connected\ndata: {\"status\":\"connected\",\"gameId\":\"%s\"}\n\n", gameID)
	flusher.Flush()

	ch := pubsub.Channel()

	for {
		select {
		case <-r.Context().Done():
			// Client disconnected
			return
		case msg, ok := <-ch:
			if !ok {
				return
			}
			// msg.Payload is a JSON string of the event
			fmt.Fprintf(w, "data: %s\n\n", msg.Payload)
			flusher.Flush()
		}
	}
}
