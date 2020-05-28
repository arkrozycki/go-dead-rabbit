package main

import (
	"net/http"
	"strconv"

	"github.com/rs/zerolog/log"
)

type Api interface {
	ServeHTTP(http.ResponseWriter, *http.Request)
}

type Server struct {
	config Config
	ds     DatastoreClientHelper
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == "GET" && r.RequestURI == "/stats":
		GetStatsHandler(w, r, s.ds)
	default:
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(``))
	}
}

func GetStatsHandler(w http.ResponseWriter, r *http.Request, ds DatastoreClientHelper) {

	count, err := ds.Count()

	if err != nil {
		log.Info().Err(err).Msg("error")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"totalDocs":` + strconv.FormatInt(count, 10) + `}`))
}
