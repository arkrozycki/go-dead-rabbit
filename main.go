package main

import (
	"os"
	"os/signal"
	"time"

	"github.com/rs/zerolog/log"
)

// main
// Because its main
func main() {
	log.Info().Msg("Go-Dead-Rabbit Starting")

	listenWithRetry() // turn up queue listener
	api()             // turn up REST API

	// Run continuously until interrupt
	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt, os.Kill)
	<-sig
}

// listenWithRetry
// Runs the queue listener and controls connection retries
func listenWithRetry() {
	dur := 1000
	listener := &Listener{
		config: Conf,
	}

	retry := make(chan int)

	for {
		go func() {
			err := listener.subscribe(retry)
			if err != nil {
				log.Info().Msg("LISTENER not connected")
				retry <- dur
			}
		}()

		dur = <-retry // wait until disconnect

		log.Debug().Msgf("LISTENER retry connection in %vms", dur)
		time.Sleep(time.Duration(dur) * time.Millisecond) // pause before reattempt connection
		if dur < 128001 {
			dur = dur * 2 // increase duration
		} else {
			dur = 180000 // fix duration at specific interval e.g. 3mins
		}
	}
}

// api
// Runs the RESTful API
func api() {

}
