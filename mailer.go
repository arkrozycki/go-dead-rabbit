package main

import (
	"context"
	"time"

	"github.com/mailgun/mailgun-go/v4"
)

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

func GetMailClient(domain string, apiKey string) MailClient {
	var m MailClient
	m = mailgun.NewMailgun(domain, apiKey)
	return m
}

func SendMail(client MailClient, message *Message) (string, string, error) {
	msg := client.NewMessage(
		message.from,
		message.subject,
		message.body,
		message.to,
	)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// Send the message with a 10 second timeout
	resp, id, err := client.Send(ctx, msg)
	return resp, id, err
}
