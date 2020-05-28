package main

import (
	"net/http"
)

type Api interface {
	ServeHTTP(http.ResponseWriter, *http.Request)
}

type Server struct {
	config Config
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// log.Trace().Msgf("%s", r.RequestURI)

	switch {
	case r.Method == "GET" && r.RequestURI == "/stats":
		GetStatsHandler(w, r)
	default:
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(``))
	}
}

func GetStatsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"totalDocs": "100000000"}`))
}
