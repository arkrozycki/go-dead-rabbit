package main

import (
	"testing"
)

type MockMailer struct {
	config Config
}

func TestSend(t *testing.T) {

	var conf = Config{
		Notification: NotificationConfig{
			Mailgun: MailgunConfig{
				ApiKey:  "none",
				BaseUrl: "http://test.test",
				Domain:  "tester",
				From:    "a",
				To:      "b",
			},
		},
	}

	var MailClientMock = &MockMailer{conf}

	err := MailClientMock.send("unit test", "body of test")
	if err != nil {
		t.Errorf("actual: %v, expected: nil", err)
	}
}
