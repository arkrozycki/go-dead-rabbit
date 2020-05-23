package main

import (
	"testing"
)

func TestSetupApi(t *testing.T) {
	SetupApi()
}

func TestListenerExec(t *testing.T) {
	mockClient := &MockAmqpConnection{}
	listener := &Listener{
		config: MockConf,
		mail:   &MockMailClient{},
	}
	retry := make(chan bool, 1)
	disconnect := make(chan bool, 1)
	disconnect <- true
	err := ListenerExec(listener, mockClient, "", retry, disconnect)
	if err != nil {
		t.Errorf("%v", err)
	}
}

func TestListenerExecDialErr(t *testing.T) {
	mockClient := &MockAmqpConnection{}
	listener := &Listener{
		config: MockConf,
		mail:   &MockMailClient{},
	}
	retry := make(chan bool, 1)
	disconnect := make(chan bool, 1)
	disconnect <- true
	err := ListenerExec(listener, mockClient, "error", retry, disconnect)
	if err != nil {
		t.Errorf("%v", err)
	}
}
