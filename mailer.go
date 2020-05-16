package main

import (
	"context"
	"time"

	"github.com/mailgun/mailgun-go/v4"
	"github.com/rs/zerolog/log"
)

type Mailer struct {
	config   Config
	Provider *mailgun.MailgunImpl
}

func (m *Mailer) _init() {
	m.Provider = mailgun.NewMailgun(Conf.Notification.Mailgun.Domain, Conf.Notification.Mailgun.ApiKey)
}

// send
// Email something to someone/thing
func (m *Mailer) send(subject string, body string, attachment []byte) {

	// log.Debug().Msg("MAILER")

	// The message object allows you to add attachments and Bcc recipients
	message := m.Provider.NewMessage(
		Conf.Notification.Mailgun.From,
		subject,
		body,
		Conf.Notification.Mailgun.To)

	message.AddBufferAttachment("something.b64", attachment)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// Send the message with a 10 second timeout
	resp, id, err := m.Provider.Send(ctx, message)

	if err != nil {
		log.Error().Err(err)
	}

	log.Debug().
		Str("ID", id).
		Str("Resp", resp).
		Msgf("MAILER")
}
