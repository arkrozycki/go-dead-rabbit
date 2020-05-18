package main

import (
	"context"
	"time"

	"github.com/mailgun/mailgun-go/v4"
	"github.com/rs/zerolog/log"
)

type Mailer interface {
	send(string, string, []byte) error
}

type Mail struct {
	config Config
}

// send
// Email something to someone/thing
func (m *Mail) send(subject string, body string, attachment []byte) error {
	mg := mailgun.NewMailgun(m.config.Notification.Mailgun.Domain, m.config.Notification.Mailgun.ApiKey)
	// The message object allows you to add attachments and Bcc recipients
	message := mg.NewMessage(
		m.config.Notification.Mailgun.From,
		subject,
		body,
		Conf.Notification.Mailgun.To)

	message.AddBufferAttachment("something.b64", attachment)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// Send the message with a 10 second timeout
	resp, id, err := mg.Send(ctx, message)

	if err != nil {
		log.Error().Err(err)
	}

	log.Debug().
		Str("ID", id).
		Str("Resp", resp).
		Msgf("MAILER")

	return err
}
