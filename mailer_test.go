package main

import (
	"context"
	"errors"
	"fmt"
	"testing"
)

type MockMailerClient struct {
}

type MockMailer struct {
	client MockMailerClient
}

func (m *MockMailer) NewMessage(from string, subject string, body string, to ...string) *Message {
	msg := &Message{
		subject: subject,
		body:    body,
	}
	return msg
}

func (m *MockMailer) Send(ctx context.Context, message *Message) (string, string, error) {
	fmt.Printf("\n\n === %v ====\n\n", message)
	var err error
	return "", "", err
}

var MockMailClient Mailer

// TestSendMail
func TestSendMail(t *testing.T) {
	MockMailClient = &MockMailer{}
	actual := SendMail(MockMailClient, "subject", "body")
	if actual != nil {
		t.Errorf("actual: %s, expected: %v.", actual, nil)
	}
}

// TestSendMailError
func TestSendMailError(t *testing.T) {
	MockMailClient = &MockMailer{}
	actual := SendMail(MockMailClient, "should error", "body")
	expected := errors.New("send failed")
	if actual != expected {
		t.Errorf("actual: %v, expected: %v", actual, expected)
	}
}
