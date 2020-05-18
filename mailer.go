package main

import (
	"context"
	"time"

	"github.com/mailgun/mailgun-go/v4"
	"github.com/rs/zerolog/log"
)

type Mailer interface {
	send(string, string) error
}

type Mail struct {
	config Config
	client *mailgun.MailgunImpl
}

func (m *Mail) getClient() *mailgun.MailgunImpl {
	if m.client == nil {
		m.client = mailgun.NewMailgun(m.config.Notification.Mailgun.Domain, m.config.Notification.Mailgun.ApiKey)
	}

	return m.client
}

func (m *Mail) Message(from string, to string, subject string, body string) *mailgun.Message {
	// The message object allows you to add attachments and Bcc recipients
	return m.client.NewMessage(
		from,
		subject,
		body,
		to)
}

// send
// Email something to someone/thing
func (m *Mail) send(subject string, body string) error {
	mg := m.getClient()
	message := m.Message(
		Conf.Notification.Mailgun.From,
		Conf.Notification.Mailgun.To,
		subject,
		body)

	// message.AddBufferAttachment(attachmentFilename, attachment)

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
		Msgf("MAILER:")

	return err
}
