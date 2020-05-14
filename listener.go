package main

import "github.com/rs/zerolog/log"

func init() {
	log.Trace().Msg("listener init")
}

func RunListener() {
	log.Debug().Msg("listener run")
}
