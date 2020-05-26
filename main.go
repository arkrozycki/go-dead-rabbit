package main

import (
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// main
// Because it is main
func main() {
	log.Info().Msg("Go-Dead-Rabbit Starting")

	// get a mail client
	mailClient := GetMailClient(Conf.Notification.Mailgun.Domain, Conf.Notification.Mailgun.ApiKey)

	// get the datastore client
	var datastoreClient DatastoreClientHelper
	mcl, err := mongo.NewClient(options.Client().ApplyURI(Conf.Datastore.Mongodb.Uri))
	if err != nil {
		log.Error().Err(err).Msg("error")
	}
	datastoreClient = &MongoClient{
		cl:      mcl,
		dbname:  Conf.Datastore.Mongodb.Database,
		colname: Conf.Datastore.Mongodb.Collection,
	}

	// init the Listener
	listener := &Listener{
		config: Conf,
		mail:   mailClient,
		ds:     datastoreClient,
	}

	ExecuteListenerWithRetry(listener, &RabbitConnection{}, GetAMQPUrl(Conf)) // turn up queue listener
	SetupApi()                                                                // turn up REST API

	// Run continuously until interrupt
	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt, os.Kill)
	<-sig
}

// setupListenerWithRetry
// Executes the queue listener and controls connection retries
func ExecuteListenerWithRetry(listener *Listener, amqpClient Amqp, connectUri string) {
	dur := 3000                      // the sleep time for retries
	retry := make(chan bool, 1)      // a channel to communicate retries
	disconnect := make(chan bool, 1) // channel for monitoring disconnects and errors
	go ListenerExec(listener, amqpClient, connectUri, retry, disconnect)
	<-retry // wait until disconnect
	log.Debug().Str("retry", strconv.Itoa(dur)).Msg("MAIN: listener disconnected")
	time.Sleep(time.Duration(dur) * time.Millisecond) // pause before reattempt connection
	ExecuteListenerWithRetry(listener, amqpClient, connectUri)

}

// ListenerExec
// Subscribes, monitors and consumes from the amqp client and queue
func ListenerExec(listener *Listener, amqpClient Amqp, connectUri string, retry chan bool, disconnect chan bool) error {
	var err error
	err = listener.Subscribe(amqpClient, connectUri) // create listener, connect and configure
	if err != nil {
		log.Info().Err(err).Msg("error")
		retry <- true
		listener.client.Close() // cleanup
		return err
	}
	log.Debug().Msg("MAIN: listener subscribed success")

	go listener.monitor(disconnect) // listen for disconnects and errors
	go listener.consume()           // listens for new messages and handles them
	<-disconnect                    // stop here until disconnect
	listener.client.Close()         // cleanup
	retry <- true                   // send dur to retry
	return nil

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
