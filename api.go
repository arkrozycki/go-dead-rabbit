package main

import "github.com/rs/zerolog/log"

type Api interface {
	start() error
}

type Server struct {
	config Config
}

func (s *Server) start() error {
	log.Debug().Msg("API: Starting up")
	return nil
}
