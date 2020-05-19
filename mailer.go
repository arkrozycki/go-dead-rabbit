package main

import (
	"context"
	"time"

	"github.com/mailgun/mailgun-go/v4"
	"github.com/rs/zerolog/log"
)

type Mailer interface {
	NewMessage(string, string, string, ...string) *Message
	Send(context.Context, *Message) (string, string, error)
}

type Message interface{}

type MailgunMessage struct {
	content *mailgun.Message
}

type MailgunMailer struct {
	client *mailgun.MailgunImpl
}

func (m *MailgunMailer) NewMessage(from string, subject string, body string, to ...string) *Message {
	// var c Message

	content := m.client.NewMessage(
		Conf.Notification.Mailgun.From,
		subject,
		body,
		Conf.Notification.Mailgun.To,
	)
	c := &Message{content}
	return c
}

func (m *MailgunMailer) Send(ctx context.Context, message *Message) (string, string, error) {
	resp, id, err := m.client.Send(ctx, message.content)
	return resp, id, err
}

// send
// Email something to someone/thing
func SendMail(client Mailer, subject string, body string) error {
	message := client.NewMessage(
		Conf.Notification.Mailgun.From,
		subject,
		body,
		Conf.Notification.Mailgun.To,
	)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	// Send the message with a 10 second timeout
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
