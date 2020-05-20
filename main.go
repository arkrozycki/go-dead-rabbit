package main

import (
	"os"
	"os/signal"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/streadway/amqp"
)

// main
// Because its main
func main() {
	log.Info().Msg("Go-Dead-Rabbit Starting")

	// setupMailer()            // turn up the mailer
	SetupListenerWithRetry() // turn up queue listener
	SetupApi()               // turn up REST API

	// Run continuously until interrupt
	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt, os.Kill)
	<-sig
}

// func setupMailer() {
// 	// MailClient := NewMailgunClient(Conf)
// 	// MailClient = &MailgunMailer{Conf, mailer} // global mailer object
// }

// setupListenerWithRetry
// Runs the queue listener and controls connection retries
func SetupListenerWithRetry() {
	dur := 1000             // the sleep time for retries
	retry := make(chan int) // a channel to communicate retries

	// init the Listener
	listener := &Listener{
		config: Conf,
		mail:   GetMailClient(Conf.Notification.Mailgun.Domain, Conf.Notification.Mailgun.ApiKey),
	}

	for {
		go func() {
			var conn AmqpConnection
			var client AmqpChannel
			// conn, err := GetAMQPConnection(GetAMQPUrl(Conf))
			conn, err := amqp.Dial(GetAMQPUrl(Conf))
			if err != nil {
				log.Info().Msg("MAIN: Listener not connected")
			}
			client, err = GetAMQPChannel(conn)
			if err != nil {
				log.Info().Msg("MAIN: Channel not connected")
			}
			err = listener.Subscribe(client)
			// err := listener.subscribe(retry)
			// if err != nil {
			// 	log.Info().Msg("MAIN: Listener not connected")
			// 	retry <- dur
			// }
		}()

		dur = <-retry // wait until disconnect

		log.Debug().Msgf("MAIN: Retry listener connection in %vms", dur)
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
func SetupApi() {
	api := &Server{Conf}
	err := api.start()
	if err != nil {
		log.Error().Err(err)
	}
}
