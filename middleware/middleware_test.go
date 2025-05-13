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
	var logBuffer bytes.Buffer
	handler := slog.NewTextHandler(&logBuffer, nil)
	logger := slog.New(handler)
	slog.SetDefault(logger)

	mockHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("test response"))
	})

	req := httptest.NewRequest("GET", "/test-path", nil)
	recorder := httptest.NewRecorder()

	middlewareHandler := LoggingMiddleware(mockHandler)
	middlewareHandler.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, recorder.Code)
	}
	if recorder.Body.String() != "test response" {
		t.Errorf("Expected response body %q, got %q", "test response", recorder.Body.String())
	}

	logOutput := logBuffer.String()
	expectedLogEntries := []string{
		"Request started",
		"method=GET",
		"path=/test-path",
		"request processed",
		"duration=",
	}

	for _, entry := range expectedLogEntries {
		if !strings.Contains(logOutput, entry) {
			t.Errorf("Expected log to contain %q, but it didn't.\nLog output: %s", entry, logOutput)
		}
	}
}
