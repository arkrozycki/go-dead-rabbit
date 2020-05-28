package main

import (
	"net/http"
	"testing"
)

type MockResponseWriter struct {
	code int
}

func (rw *MockResponseWriter) Header() http.Header {
	return http.Header{}
}
func (rw *MockResponseWriter) Write(w []byte) (int, error) {
	return 0, nil
}
func (rw *MockResponseWriter) WriteHeader(code int) {
	rw.code = code
}

func TestNotFound(t *testing.T) {
	s := Server{Conf}
	rw := &MockResponseWriter{}
	s.ServeHTTP(rw, &http.Request{
		Method:     "GET",
		RequestURI: "/w3qr",
	})

	if rw.code != 404 {
		t.Errorf("%s GET /w3qr expected: 404, got: %d", failed, rw.code)
	}

	t.Logf("%s GET /w3qr 404", succeed)
}

func TestGetStatsHandlerSuccess(t *testing.T) {
	s := Server{Conf}
	rw := &MockResponseWriter{}
	s.ServeHTTP(rw, &http.Request{
		Method:     "GET",
		RequestURI: "/stats",
	})

	if rw.code != 200 {
		t.Errorf("%s did not return 200", failed)
	}

	t.Logf("%s GET /stats", succeed)
}
