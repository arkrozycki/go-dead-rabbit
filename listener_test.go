package main

import (
	"testing"

	"github.com/streadway/amqp"
)

var simpleListener = &Listener{
	config: MockConf,
	mail:   GetMockMailClient(),
}

// TestSubscribe
func TestSubscribe(t *testing.T) {

}

// TestAmqpMessageHandler
func TestAmqpMessageHandler(t *testing.T) {
	msg := amqp.Delivery{
		CorrelationId: "none",
		ReplyTo:       "reply",
		MessageId:     "123",
		ContentType:   "binary",
		RoutingKey:    "evt.bstock.all",
		Exchange:      "fake",
		Headers:       amqp.Table{},
		Body:          []byte("fake data"),
	}

	err := simpleListener.amqpMessageHandler(msg)
	if err != nil {
		t.Errorf("actual: %v, expected: nil", err)
	}
}

// TestAmqpUrl
func TestAmqpUrl(t *testing.T) {
	expected := "amqp://tester:password@testServer:5672/vhost"
	actual := simpleListener.amqpUrl()

	if actual != expected {
		t.Errorf("actual: %s, expected: %s.", actual, expected)
	}
}

// TestConnect
func TestConnect(t *testing.T) {}

// TestConfigure
func TestConfigure(t *testing.T) {}

// TestConsume
func TestConsume(t *testing.T) {}
