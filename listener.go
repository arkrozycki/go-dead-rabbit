package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/streadway/amqp"
)

// Listener
// Struct for the listener
type Listener struct {
	config  Config
	conn    *amqp.Connection
	channel *amqp.Channel
	queue   amqp.Queue
	mail    MailClient
}

// subscribe
// Subscribes to a queue based on configuration.
// Executes a connect, opens channel, consumes.
// Handles disconnects via the NotifyError channel
func (l *Listener) subscribe(retry chan<- int) error {
	log.Debug().Msg("LISTENER: Starting up")

	// connect to amqp
	notify, err := l.connect()
	if err != nil {
		return err
	}

	// configure queue
	err = l.configure()
	if err != nil {
		return err
	}

	// start consuming messages
	msgs, err := l.consume()
	if err != nil {
		return err
	}

	defer l.channel.Close()
	defer l.conn.Close()

	// channel for monitoring disconnects and errors
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
			log.Debug().Msgf("LISTENER: Msg received - %s", msg.MessageId)
			err := l.amqpMessageHandler(msg)
			if err != nil {
				log.Error().Err(err)
				msg.Ack(false)
				continue
			}

			msg.Ack(false) // message acknowledgement
		}
	}()

	<-disconnect  // stop here until disconnect
	retry <- 1000 // reset the retry duration
	return errors.New("disconnected")
}

// amqpMessageHandler
// Processes the incoming messages
func (l *Listener) amqpMessageHandler(message amqp.Delivery) error {
	var err error

	headers, err := json.Marshal(message)

	log.Debug().
		RawJSON("Headers", headers).
		Str("CorrelationId", message.CorrelationId).
		Str("ReplyTo", message.ReplyTo).
		Str("MessageId", message.MessageId).
		Str("ContentType", message.ContentType).
		Str("RoutingKey", message.RoutingKey).
		Str("Exchange", message.Exchange).
		Msg("LISTENER:")

	var prettyJSON bytes.Buffer
	err = json.Indent(&prettyJSON, headers, "", "\t")
	if err != nil {
		return err
	}

	// attachmentName := fmt.Sprintf("msg_%s.rabbit", message.Headers["proto"])
	msg := &Message{
		from:    Conf.Notification.Mailgun.From,
		to:      Conf.Notification.Mailgun.To,
		subject: message.RoutingKey,
		body:    string(prettyJSON.Bytes()),
	}
	resp, id, err := SendMail(l.mail, msg)
	log.Debug().Str("ID", id).Str("Resp", resp).Msgf("MAILER:")
	return err
}

// amqpUrl
// Generates a connection string from config
func (l *Listener) amqpUrl() string {
	return fmt.Sprintf("amqp://%s:%s@%s:%s/%s",
		l.config.Connection.User,
		l.config.Connection.Password,
		l.config.Connection.Server,
		l.config.Connection.Port,
		l.config.Connection.Vhost)
}

// connect
// Establishes a connection to the AMQP host,
// retrieves a channel,
// configures the Qos on the channel,
// returns the amqp.NotifyClose channel for detecting errors and disconnects
func (l *Listener) connect() (chan *amqp.Error, error) {
	var err error
	// connect to rabbitmq
	l.conn, err = amqp.Dial(l.amqpUrl())
	if err != nil {
		return nil, err
	}
	log.Info().Msgf("LISTENER: Connected with %s", l.amqpUrl())

	// open channel from connection
	l.channel, err = l.conn.Channel()
	if err != nil {
		return nil, err
	}
	log.Debug().Msg("LISTENER: Open channel success")

	// set the QoS for the channel
	// we will be pulling only 1 at a time for simplicity
	err = l.channel.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)

	// capture the error disconnect channel
	notify := l.conn.NotifyClose(make(chan *amqp.Error)) //error channel

	return notify, err
}

// configure
// Passively checks that the queue in configuration exists.
// Note, the app does not create the queue nor the bindings, those
// need to be preemptively configured.
func (l *Listener) configure() error {
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
	log.Debug().Msgf("LISTENER: Queue exists - %s", l.config.Listener.Queue.Name)

	return err
}

// consume
// Consumes messages from the queue,
// returns the channel used for receiving messages.
func (l *Listener) consume() (<-chan amqp.Delivery, error) {
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
