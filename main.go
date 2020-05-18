package main

import (
	"os"
	"os/signal"
	"time"

	"github.com/rs/zerolog/log"
)

var MailClient Mailer

// main
// Because its main
func main() {
	log.Info().Msg("Go-Dead-Rabbit Starting")

	setupMailer()            // turn up the mailer
	setupListenerWithRetry() // turn up queue listener
	setupApi()               // turn up REST API

	// Run continuously until interrupt
	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt, os.Kill)
	<-sig
}

func setupMailer() {
	// Mail._init()
	MailClient = Mail{Conf} // global mailer object
}

// setupListenerWithRetry
// Runs the queue listener and controls connection retries
func setupListenerWithRetry() {
	dur := 1000             // the sleep time for retries
	retry := make(chan int) // a channel to communicate retries

	// init the Listener
	listener := &Listener{
		config: Conf,
	}

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

// setupApi
// Runs the RESTful API
func setupApi() {

}
