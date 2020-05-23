package main

import (
	"testing"
)

const succeed = "\u2713"
const failed = "\u2717"

func TestSetupApi(t *testing.T) {
	SetupApi()
	t.Logf("%s API setup", succeed)
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
		t.Errorf("%s %v", failed, err)
	}
	t.Logf("%s Listener executed", succeed)
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
	if err == nil {
		t.Errorf("%s %v", failed, err)
	}
	t.Logf("%s Listener dial should fail", succeed)
}
