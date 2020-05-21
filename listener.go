package main

import (
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/streadway/amqp"
)

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
func (r *RabbitConnection) dial(uri string) error {
	var err error
	r.connection, err = amqp.Dial(uri)
	return err
}

// channel
func (r *RabbitConnection) getChannel() error {
	var err error
	r.channel, err = r.connection.Channel()
	return err
}

// Close
func (r *RabbitConnection) Close() {
	r.channel.Close()
	r.connection.Close()
}

// setQos
// Qos controls how many messages or how many bytes the server will try to keep on the network for consumers before receiving delivery acks. The intent of Qos is to make sure the network buffers stay full between the server and client
func (r *RabbitConnection) setQos(count int, size int, global bool) error {
	err := r.channel.Qos(
		count,  // prefetch count
		size,   // prefetch size
		global, // global
	)
	return err
}

// setNotifyCloseChannel
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

func (r *RabbitConnection) setMessages(qName string) (<-chan amqp.Delivery, error) {
	msgs, err := r.channel.Consume(
		qName, // queue
		"",    // consumer
		false, // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)

	return msgs, err
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

// Listener
// Struct for the listener
type Listener struct {
	config   Config
	mail     MailClient
	client   Amqp
	errChan  chan *amqp.Error
	messages <-chan amqp.Delivery
}

// subscribe
// Subscribes to a queue based on configuration.
// Executes a connect, opens channel, consumes.
// Handles disconnects via the NotifyError channel
func (l *Listener) Subscribe(client Amqp) error {
	var err error
	l.client = client
	err = l.connect()
	return err
}

// connect
// Establishes a connection to the AMQP host,
// retrieves a channel,
func (l *Listener) connect() error {
	var err error

	// client connect
	err = l.client.dial(GetAMQPUrl(Conf))
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
	err = l.configure(1, 0, false)
	if err != nil {
		return err
	}

	return nil
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

	// capture the error disconnect channel
	// l.errorChan = l.client.getNotifyClose()

	// configure the queue
	err = l.client.setQueue(l.config.Listener.Queue.Name)
	if err != nil {
		return err
	}

	return nil
}

// consume
// Consumes messages from the queue,
// returns the channel used for receiving messages.
func (l *Listener) consume() error {
	var err error
	// start consuming messages
	log.Debug().Msg(l.config.Listener.Queue.Name)
	l.messages, err = l.client.setMessages(l.config.Listener.Queue.Name)

	// msgs, err := l.channel.Consume(
	// 	l.config.Listener.Queue.Name, // queue
	// 	"",                           // consumer
	// 	false,                        // auto-ack
	// 	false,                        // exclusive
	// 	false,                        // no-local
	// 	false,                        // no-wait
	// 	nil,                          // args
	// )

	return err
}

/*


// subscribe
// Subscribes to a queue based on configuration.
// Executes a connect, opens channel, consumes.
// Handles disconnects via the NotifyError channel
func (l *Listener) subscribe(retry chan<- int) error {

	// start consuming messages
	msgs, err := l.consume()
	if err != nil {
		return err
	}

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


*/
