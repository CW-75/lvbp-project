package sse

import "net/http"

type StreamHandler struct {
	// TODO: inject redis pub/sub subscriber
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

	// TODO: implement event loop and flush
	_ = flusher
}
