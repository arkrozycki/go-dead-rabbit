package main

import (
	"os"
	"os/signal"

	"github.com/rs/zerolog/log"
)

// main
func main() {
	log.Info().Msg("Go-Dead-Rabbit Starting")

	// listener := &Listener{
	// 	config: Conf,
	// }
	// go func() {
	// 	err := listener.subscribe()
	// 	log.Error().Msgf("%v", err)
	// }()
	listenWithRetry()

	go RunAPI()

	// Run continuously until interrupt
	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt, os.Kill)
	<-sig
}

func listenWithRetry() {
	retry := 0
	listener := &Listener{
		config: Conf,
	}

	go func() {
		err := listener.subscribe()
	}()

}
