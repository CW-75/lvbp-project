package rest

import "net/http"

type IngestionHandler struct {
	// TODO: inject scorekeeper service
}

func (h *IngestionHandler) ProcessPitch(w http.ResponseWriter, r *http.Request) {
	// TODO: implement REST endpoint for scorekeeper ingestion
}
