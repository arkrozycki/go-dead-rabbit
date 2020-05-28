package main

import (
	"errors"
	"testing"

	"github.com/streadway/amqp"
)

type (
	MockAmqpConnection struct {
		connection *MockConnection
		channel    *MockChannel
	}

	MockConnection struct{}
	MockChannel    struct{}
)

func (m *MockAmqpConnection) dial(uri string) error {
	m.connection = &MockConnection{}
	if uri == "error" {
		return errors.New("no dial tone")
	}
	return nil
}

func (m *MockAmqpConnection) getChannel() error {
	m.channel = &MockChannel{}
	return nil
}

func (m *MockAmqpConnection) Close()                                        {}
func (m *MockAmqpConnection) setQos(count int, size int, global bool) error { return nil }
func (m *MockAmqpConnection) setQueue(string) error                         { return nil }
func (m *MockAmqpConnection) setNotifyCloseChannel(chan *amqp.Error) chan *amqp.Error {
	return make(chan *amqp.Error)
}
func (m *MockAmqpConnection) setMessages(q string) (<-chan amqp.Delivery, error) {
	return make(<-chan amqp.Delivery), nil
}

// TestGetAMQPUrl
func TestGetAMQPUrl(t *testing.T) {
	expected := "amqp://tester:password@testServer:5672/vhost"
	actual := GetAMQPUrl(MockConf)

	if actual != expected {
		t.Errorf("%s actual: %s, expected: %s.", failed, actual, expected)
	}
	t.Logf("%s AMQP Uri correct", succeed)
}

// TestSubscribe
func TestSubscribe(t *testing.T) {
	var err error
	mockClient := &MockAmqpConnection{}
	listener := &Listener{
		config: MockConf,
		mail:   GetMailClient("", ""),
	}

	err = listener.Subscribe(mockClient, "")
	if err != nil {
		t.Errorf("%s Subscribe failed", failed)
	}
	t.Logf("%s Subscribe success", succeed)
}

// // TestAmqpMessageHandler
func TestListenerHandle(t *testing.T) {
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

	// mc := &MockDatastoreClient{}
	// init the Listener
	listener := &Listener{
		config: MockConf,
		mail:   &MockMailClient{},
		ds: &MockMongoClientHelper{
			cl:      &MockDatastoreClient{},
			dbname:  "testdb",
			colname: "testCollection",
		},
	}

	err := listener.handle(msg)
	if err != nil {
		t.Errorf("%s actual: %v, expected: nil", failed, err)
	}
	t.Logf("%s Message handled", succeed)
}
