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

	// setupMailer()            // turn up the mailer
	SetupListenerWithRetry() // turn up queue listener
	SetupApi()               // turn up REST API

	// Run continuously until interrupt
	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt, os.Kill)
	<-sig
}

// setupListenerWithRetry
// Runs the queue listener and controls connection retries
func SetupListenerWithRetry() {
	dur := 3000             // the sleep time for retries
	retry := make(chan int) // a channel to communicate retries

	for {
		go func() {
			var err error

			// init the Listener
			listener := &Listener{
				config: Conf,
				mail:   GetMailClient(Conf.Notification.Mailgun.Domain, Conf.Notification.Mailgun.ApiKey),
			}

			amqpClient := &RabbitConnection{} // amqp client

			err = listener.Subscribe(amqpClient) // create listener, connect and configure
			if err != nil {
				log.Info().Err(err).Msg("error")
				retry <- dur
				return
			} else {
				log.Debug().Msg("MAIN: listener subscribed success")

				// // channel for monitoring disconnects and errors
				disconnect := make(chan bool)
				// // connection error monitoring
				go func() {
					// errChan := listener.client.setNotifyCloseChannel(make(chan *amqp.Error))
					log.Debug().Msg("MAIN: listening for disconnects")
					err := <-listener.errChan
					log.Error().Err(err).Msg("error")
					disconnect <- true
				}()

				go func() {
					err := listener.consume()
					if err != nil {
						log.Info().Err(err).Msg("error")
					}
					log.Debug().Msg("MAIN: listening for messages")
					for msg := range listener.messages {
						log.Debug().Msgf("received %v", msg)
					}
				}()

				<-disconnect            // stop here until disconnect
				listener.client.Close() // cleanup
				log.Debug().Msg("MAIN: listener disconnected")
				retry <- dur
			}
		}()

		dur = <-retry // wait until disconnect
		log.Debug().Msgf("MAIN: Retry listener connection in %vms", dur)
		time.Sleep(time.Duration(dur) * time.Millisecond) // pause before reattempt connection
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
