package main

import (
	"context"
	"time"

	"github.com/mailgun/mailgun-go/v4"
	"github.com/rs/zerolog/log"
)

type Mailer interface {
	// NewMessage(string, string, string, ...string) *Message
	Send(context.Context, *Message) (string, string, error)
}

type Message struct {
	from    string
	to      string
	subject string
	body    string
}

type MailgunMailer struct {
	client *mailgun.MailgunImpl
}

func (m *MailgunMailer) Send(ctx context.Context, message *Message) (string, string, error) {
	msg := m.client.NewMessage(
		message.from,
		message.subject,
		message.body,
		message.to,
	)
	// Send the message with a 10 second timeout
	resp, id, err := m.client.Send(ctx, msg)
	return resp, id, err
}

// send
// Email something to someone/thing
func SendMail(client Mailer, subject string, body string) error {
	message := &Message{
		from:    Conf.Notification.Mailgun.From,
		to:      Conf.Notification.Mailgun.To,
		subject: subject,
		body:    body,
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	resp, id, err := client.Send(ctx, message)
	if err != nil {
		log.Error().Err(err)
	}

	log.Debug().
		Str("ID", id).
		Str("Resp", resp).
		Msgf("MAILER:")

	return err
}
