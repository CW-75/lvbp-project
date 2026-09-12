package rest

import "net/http"

type StandingsHandler struct {
	// TODO: inject standings query service or repo
}

func (h *StandingsHandler) GetStandings(w http.ResponseWriter, r *http.Request) {
	// TODO: implement REST endpoint
}
