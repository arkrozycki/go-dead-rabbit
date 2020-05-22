package main

import (
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
		t.Errorf("actual: %s, expected: %s.", actual, expected)
	}
}

// func TestSetMessages(t *testing.T) {
// 	var client = &RabbitConnection{
// 		channel: &amqp.Channel{},
// 	}
// 	_, err := client.setMessages("")
// 	if err != nil {

// 	}
// }

// TestSubscribe
func TestSubscribe(t *testing.T) {
	var err error
	mockClient := &MockAmqpConnection{}
	// init the Listener
	listener := &Listener{
		config: MockConf,
		mail:   GetMailClient("", ""),
	}

	err = listener.Subscribe(mockClient)
	if err != nil {
		t.Errorf("subscribe failed")
	}
}

// TestListenerConsume
func TestListenerConsume(t *testing.T) {
	var err error
	mockClient := &MockAmqpConnection{}
	// init the Listener
	listener := &Listener{
		config: MockConf,
		mail:   &MockMailClient{},
	}
	err = listener.Subscribe(mockClient)
	if err != nil {
		t.Errorf("subscribe failed")
	}

	err = listener.consume()
	if err != nil {
		t.Errorf("%v", err)
	}

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

	// init the Listener
	listener := &Listener{
		config: MockConf,
		mail:   &MockMailClient{},
	}

	err := listener.handle(msg)
	if err != nil {
		t.Errorf("actual: %v, expected: nil", err)
	}
}
