package middleware

import (
	"bytes"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestLoggingMiddleware(t *testing.T) {
	// Setup a buffer to capture logs
	var buf bytes.Buffer
	testLogger := slog.New(slog.NewTextHandler(&buf, nil))

	// Create the middleware with a mock handler
	handlerWasCalled := false
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlerWasCalled = true
	})

	// Create a mock GET request and recorder
	req := httptest.NewRequest("GET", "/test-path", nil)
	rr := httptest.NewRecorder()

	// Save a copy of the previous logger and switch to the test logger
	previousLogger := slog.Default()
	slog.SetDefault(testLogger)
	defer slog.SetDefault(previousLogger) // Restore after test

	// Execute the middleware using the test logger that will write out to the buffer
	Logging(nextHandler).ServeHTTP(rr, req)

	// Assert the inner handler was actually executed
	if !handlerWasCalled {
		t.Error("loggingMiddleware did not call the next handler")
	}

	// Assert the log returned the path
	if !strings.Contains(buf.String(), "/test-path") {
		t.Errorf("Log did not contain the request path. Got: %s", buf.String())
	}
}
