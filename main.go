package main

import (
	"os"
	"os/signal"

	"github.com/rs/zerolog/log"
)

// main
func main() {
	log.Info().Msg("Go-Dead-Rabbit Starting")

	go RunListener()
	go RunAPI()

	// Run continuously until interrupt
	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt, os.Kill)
	<-sig
}
