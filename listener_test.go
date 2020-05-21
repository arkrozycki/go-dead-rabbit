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

func TestRabbitConnectionDail(t *testing.T) {
	var amqpClient = &RabbitConnection{}
	amqpClient.dial("")
}

// import (
// 	"testing"

// 	"github.com/streadway/amqp"
// )

// var simpleListener = &Listener{
// 	config: MockConf,
// 	mail:   GetMockMailClient(),
// 	conn:   GetAmqpConnection(),
// }

// // TestSubscribe
// func TestSubscribe(t *testing.T) {

// }

// // TestAmqpMessageHandler
// func TestAmqpMessageHandler(t *testing.T) {
// 	msg := amqp.Delivery{
// 		CorrelationId: "none",
// 		ReplyTo:       "reply",
// 		MessageId:     "123",
// 		ContentType:   "binary",
// 		RoutingKey:    "evt.bstock.all",
// 		Exchange:      "fake",
// 		Headers:       amqp.Table{},
// 		Body:          []byte("fake data"),
// 	}

// 	err := simpleListener.amqpMessageHandler(msg)
// 	if err != nil {
// 		t.Errorf("actual: %v, expected: nil", err)
// 	}
// }

// type MockAmqpConnection struct{}
// type MockAmqpChannel struct{}

// func (m *MockAmqpConnection) Channel() (AmqpChannel, error) {
// 	var err error
// 	c := &MockAmqpChannel{}
// 	return c, err
// }

// func (c *MockAmqpChannel) Qos(count int, size int, global bool) error {
// 	var err error
// 	return err
// }

// func TestGetAMQPChannel(t *testing.T) {
// 	var conn = &MockAmqpConnection{}
// 	_, err := GetAMQPChannel(conn)
// 	if err != nil {
// 		t.Errorf("amqp connection failed %v\n", err)
// 	}

// 	// if reflect.TypeOf(conn) != reflect.TypeOf(connection) {
// 	// 	t.Errorf("%s", reflect.TypeOf(conn))
// 	// }
// }

// TestGetAMQPConnection
// func TestGetAMQPConnection(t *testing.T) {
// 	var connection = &amqp.Connection{}
// 	uri := GetAMQPUrl(MockConf)
// 	conn, err := GetAMQPConnection(uri)
// 	if err != nil {
// 		t.Errorf("amqp connection failed %v\n", err)
// 	}

// 	if reflect.TypeOf(conn) != reflect.TypeOf(connection) {
// 		t.Errorf("%s", reflect.TypeOf(conn))
// 	}

// }

// // TestConnect
// func TestConnect(t *testing.T) {}

// // TestConfigure
// func TestConfigure(t *testing.T) {}

// // TestConsume
// func TestConsume(t *testing.T) {

// }
