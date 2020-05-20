package main

import (
	"context"
	"reflect"
	"testing"

	"github.com/mailgun/mailgun-go/v4"
)

type MockMailClient struct{}

func (m *MockMailClient) NewMessage(from string, subject string, body string, to ...string) *mailgun.Message {
	c := mailgun.NewMailgun("", "")
	return c.NewMessage(from, subject, body, to[0])
}

func (m *MockMailClient) Send(ctx context.Context, message *mailgun.Message) (string, string, error) {
	var err error
	return "", "", err
}

var MockMailer MailClient

func GetMockMailClient() MailClient {
	var m MailClient
	m = &MockMailClient{}
	return m
}

// TestSendMail
func TestSendMailSuccess(t *testing.T) {
	MockMailer = &MockMailClient{}

	msg := &Message{
		from:    "from_test",
		to:      "to_test",
		subject: "subject",
		body:    "body",
	}

	_, _, actual := SendMail(MockMailer, msg)
	if actual != nil {
		t.Errorf("actual: %s, expected: %v.", actual, nil)
	}
}

//  TestGetMailClient
func TestGetMailClient(t *testing.T) {
	c := mailgun.NewMailgun("", "")
	m := GetMailClient("", "")
	if reflect.TypeOf(m) != reflect.TypeOf(c) {
		t.Errorf("%s", reflect.TypeOf(m))
	}

}
