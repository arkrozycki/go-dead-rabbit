package main

import (
	"testing"
)

func TestStart(t *testing.T) {
	s := &Server{}
	err := s.start()
	if err != nil {
		t.Errorf("%s error should be nil", failed)
	}
	t.Logf("%s API started", succeed)
}
