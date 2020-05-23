package main

import "testing"

var MockConf = Config{
	Connection: ConnectionConfig{
		Server:   "testServer",
		Port:     "5672",
		Vhost:    "vhost",
		User:     "tester",
		Password: "password",
	},
}

func TestInit(t *testing.T) {
	t.Logf("%s Test init-ed", succeed)
}
