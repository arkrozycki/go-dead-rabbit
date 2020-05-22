package main

import (
	"context"
	"time"

	"github.com/mailgun/mailgun-go/v4"
	"github.com/rs/zerolog/log"
)

var EMAIL_SEND_TIMEOUT = 30 // seconds

type MailClient interface {
	NewMessage(string, string, string, ...string) *mailgun.Message
	Send(context.Context, *mailgun.Message) (string, string, error)
}

type Message struct {
	from    string
	to      string
	subject string
	body    string
}

// GetMailClient
// returns a MailClient
func GetMailClient(domain string, apiKey string) MailClient {
	var m MailClient
	m = mailgun.NewMailgun(domain, apiKey)
	return m
}

// SendMail
// Sends email with provided MailClient
func SendMail(client MailClient, message *Message) (string, string, error) {
	msg := client.NewMessage(
		message.from,
		message.subject,
		message.body,
		message.to,
	)

	// get context and set timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(EMAIL_SEND_TIMEOUT))
	defer cancel()

	// Send the message with a 10 second timeout
	resp, id, err := client.Send(ctx, msg)
	log.Debug().Str("ID", id).Str("Resp", resp).Msgf("MAILER:")
	return resp, id, err
}
