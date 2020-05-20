package main

import (
	"context"
	"errors"
	"testing"
)

type MockMailerClient struct {
}

type MockMailer struct {
	client MockMailerClient
}

func (m *MockMailer) Send(ctx context.Context, message *Message) (string, string, error) {
	var err error
	if message.subject == "should error" {
		return "", "", errors.New("send failed")
	}

	return "", "", err
}

var MockMailClient Mailer

// TestSendMail
func TestSendMailSuccess(t *testing.T) {
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
