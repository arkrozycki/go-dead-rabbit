package main

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/streadway/amqp"
)

// Listener
// Struct for the listener
type Listener struct {
	config   Config
	mail     MailClient
	ds       DatastoreClientHelper
	client   Amqp
	errChan  chan *amqp.Error
	messages <-chan amqp.Delivery
}

// Amqp
type (
	Amqp interface {
		dial(string) error
		getChannel() error
		Close()
		setQos(int, int, bool) error
		setQueue(string) error
		setNotifyCloseChannel(chan *amqp.Error) chan *amqp.Error
		setMessages(string) (<-chan amqp.Delivery, error)
	}

	RabbitConnection struct {
		connection  *amqp.Connection
		channel     *amqp.Channel
		notifyClose chan *amqp.Error
		queue       amqp.Queue
	}
)

// dial
// connects to the amqp host
func (r *RabbitConnection) dial(uri string) error {
	var err error
	r.connection, err = amqp.Dial(uri)
	return err
}

// channel
// retrieves a channel from the connection
func (r *RabbitConnection) getChannel() error {
	var err error
	r.channel, err = r.connection.Channel()
	return err
}

// Close
// Closes the channel and connection
func (r *RabbitConnection) Close() {
	if r.channel != nil {
		r.channel.Close()
	}
	if r.connection != nil {
		r.connection.Close()
	}
}

// setQos
// Qos controls how many messages or how many bytes the server will try to keep on the network for consumers before receiving delivery acks. The intent of Qos is to make sure the network buffers stay full between the server and client
func (r *RabbitConnection) setQos(count int, size int, global bool) error {
	return r.channel.Qos(
		count,  // prefetch count
		size,   // prefetch size
		global, // global
	)
}

// setNotifyCloseChannel
// returns the error/disconnect channel
func (r *RabbitConnection) setNotifyCloseChannel(ch chan *amqp.Error) chan *amqp.Error {
	return r.connection.NotifyClose(ch)
}

// setQueue
// Connect to an existing queue, will throw if not exist
func (r *RabbitConnection) setQueue(name string) error {
	var err error
	r.queue, err = r.channel.QueueDeclarePassive(
		name,  // name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	return err
}

// setMessages
// Starts consuming messages from the queue
func (r *RabbitConnection) setMessages(qName string) (<-chan amqp.Delivery, error) {
	return r.channel.Consume(
		qName, // queue
		"",    // consumer
		false, // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
}

// amqpUrl
// Generates a connection string from config
func GetAMQPUrl(conf Config) string {
	return fmt.Sprintf("amqp://%s:%s@%s:%s/%s",
		conf.Connection.User,
		conf.Connection.Password,
		conf.Connection.Server,
		conf.Connection.Port,
		conf.Connection.Vhost)
}

// subscribe
// Subscribes to a queue based on configuration.
// Executes a connect, opens channel, consumes.
// Handles disconnects via the NotifyError channel
func (l *Listener) Subscribe(client Amqp, connectUri string) error {
	l.client = client
	return l.connect(connectUri)
}

// connect
// Establishes a connection to the AMQP host,
// retrieves a channel,
func (l *Listener) connect(connectUri string) error {
	var err error

	// client connect
	err = l.client.dial(connectUri)
	if err != nil {
		return err
	}

	// client channel
	err = l.client.getChannel()
	if err != nil {
		return err
	}

	// client error / disconnect channel
	l.errChan = l.client.setNotifyCloseChannel(make(chan *amqp.Error))

	// configure
	return l.configure(1, 0, false)
}

// configure
// configures the Qos on the channel,
// sets the NotifyClose channel for detecting errors and disconnects
func (l *Listener) configure(prefetchCount int, prefetchSize int, global bool) error {
	var err error
	// set the QoS for the channel
	err = l.client.setQos(prefetchCount, prefetchSize, global)
	if err != nil {
		return err
	}

	// configure the queue
	return l.client.setQueue(l.config.Listener.Queue.Name)
}

// consume
// Consumes messages from the queue,
// returns the channel used for receiving messages.
func (l *Listener) consume() error {
	var err error
	log.Debug().Str("queue", l.config.Listener.Queue.Name).Msg("LISTENER: consume")
	messages, err := l.client.setMessages(l.config.Listener.Queue.Name)
	if err != nil {
		log.Info().Err(err).Msg("error")
		return err
	}

	for msg := range messages {
		log.Debug().Msg("new message received")
		err = l.handle(msg)
		if err != nil {
			log.Info().Err(err).Msg("error")
		}
		msg.Ack(false)
	}

	return err
}

// handle
// AMQP message handler:
// 	- Marshals to JSON and prettify
//  - Saves to datastore
// 	- Sends email notification
func (l *Listener) handle(message amqp.Delivery) error {
	messageJSON, err := json.Marshal(message)

	log.Debug().
		RawJSON("Headers", messageJSON).
		Str("CorrelationId", message.CorrelationId).
		Str("ReplyTo", message.ReplyTo).
		Str("MessageId", message.MessageId).
		Str("ContentType", message.ContentType).
		Str("RoutingKey", message.RoutingKey).
		Str("Exchange", message.Exchange).
		Msg("LISTENER:")

	id, err := l.ds.Insert(messageJSON) // persist to storage
	if err != nil {
		return err
	}

	err = l.mailJSON(message.RoutingKey+":"+id, messageJSON) // mailer
	return err
}

// mailJSON
// Formats JSON for pretty
// Sends email to configured destination
func (l *Listener) mailJSON(subject string, body []byte) error {
	var prettyJSON bytes.Buffer
	err := json.Indent(&prettyJSON, body, "", "\t")
	if err != nil {
		return err
	}

	// attachmentName := fmt.Sprintf("msg_%s.rabbit", message.Headers["proto"])
	msg := &Message{
		from:    Conf.Notification.Mailgun.From,
		to:      Conf.Notification.Mailgun.To,
		subject: subject,
		body:    string(prettyJSON.Bytes()),
	}
	_, _, err = SendMail(l.mail, msg)
	return err
}

// monitor
// Monitors the listener for errors and disconnects
func (l *Listener) monitor(disconn chan bool) {
	log.Debug().Msg("LISTENER: listening for disconnects")
	err := <-l.errChan
	log.Error().Err(err).Msg("error")
	disconn <- true
}
