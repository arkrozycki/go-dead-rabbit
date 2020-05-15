package main

import (
	"errors"
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/streadway/amqp"
)

// Listener
type Listener struct {
	config  Config
	conn    *amqp.Connection
	channel *amqp.Channel
	queue   amqp.Queue
}

// start
func (l *Listener) subscribe() error {
	log.Debug().Msg("listener starting up")
	// connect to amqp
	notify, err := l._connect()
	if err != nil {
		return err
	}

	// configure queue
	err = l._configure()
	if err != nil {
		return err
	}

	// start consuming messages
	msgs, err := l._consume()
	if err != nil {
		return err
	}

	defer l.channel.Close()
	defer l.conn.Close()

	disconnect := make(chan bool)
	// connection error monitoring
	go func() {
		for e := range notify {
			log.Error().Msgf("%v", e)
			disconnect <- true
		}
	}()

	// queue message processing
	go func() {
		for msg := range msgs {
			// listen on channel for new messages
			log.Debug().Msgf("LISTENER MSG: %s", msg.Body)
			msg.Ack(false)
		}
	}()

	<-disconnect // stop here until disconnect
	return errors.New("disconnected")
}

// getAMQPURL
func (l *Listener) getAMQPURL() string {
	return fmt.Sprintf("amqp://%s:%s@%s:%s/%s",
		l.config.Connection.User,
		l.config.Connection.Password,
		l.config.Connection.Server,
		l.config.Connection.Port,
		l.config.Connection.Vhost)
}

// _connect
func (l *Listener) _connect() (chan *amqp.Error, error) {
	var err error
	// connect to rabbitmq
	l.conn, err = amqp.Dial(l.getAMQPURL())
	if err != nil {
		return nil, err
	}
	log.Info().Msgf("LISTENER connected with %s", l.getAMQPURL())

	// open channel from connection
	l.channel, err = l.conn.Channel()
	if err != nil {
		return nil, err
	}
	log.Debug().Msg("LISTENER open channel success")

	// set the QoS for the channel
	// we will be pulling only 1 at a time for simplicity
	err = l.channel.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)

	notify := l.conn.NotifyClose(make(chan *amqp.Error)) //error channel

	return notify, err
}

// _configure
func (l *Listener) _configure() error {
	var err error
	// Connect to an existing queue, will throw if not exist
	l.queue, err = l.channel.QueueDeclarePassive(
		l.config.Listener.Queue.Name, // name
		true,                         // durable
		false,                        // delete when unused
		false,                        // exclusive
		false,                        // no-wait
		nil,                          // arguments
	)
	if err != nil {
		return err
	}
	log.Debug().Msgf("LISTENER queue exists %s", l.config.Listener.Queue.Name)

	return err
}

// _consume
func (l *Listener) _consume() (<-chan amqp.Delivery, error) {
	// start consuming messages
	msgs, err := l.channel.Consume(
		l.config.Listener.Queue.Name, // queue
		"",                           // consumer
		false,                        // auto-ack
		false,                        // exclusive
		false,                        // no-local
		false,                        // no-wait
		nil,                          // args
	)

	return msgs, err
}
