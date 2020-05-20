package main

import (
	"testing"
)

func TestStart(t *testing.T) {
	s := &Server{}
	err := s.start()
	if err != nil {
		t.Errorf("error should be nil")
	}
}
