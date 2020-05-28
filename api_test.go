package main

import (
	"net/http"
	"testing"
)

type MockResponseWriter struct {
	code int
	body []byte
}

func (rw *MockResponseWriter) Header() http.Header {
	return http.Header{}
}
func (rw *MockResponseWriter) Write(w []byte) (int, error) {
	rw.body = w
	return 0, nil
}
func (rw *MockResponseWriter) WriteHeader(code int) {
	rw.code = code
}

func TestNotFound(t *testing.T) {
	s := Server{Conf, nil}
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
	s := Server{Conf, &MockMongoClientHelper{}}
	rw := &MockResponseWriter{}
	s.ServeHTTP(rw, &http.Request{
		Method:     "GET",
		RequestURI: "/stats",
	})

	if rw.code != 200 {
		t.Errorf("%s did not return 200", failed)
	}

	if string(rw.body) != `{"totalDocs":10}` {
		t.Errorf("%s returned body does not match, got: %v", failed, string(rw.body))
	}

	t.Logf("%s GET /stats", succeed)
}
