package debug

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestErrorHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/debug/error", nil)
	recorder := httptest.NewRecorder()

	handler := http.HandlerFunc(ErrorHandler)
	handler.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusInternalServerError {
		t.Fatalf("unexpected status code: got %d want %d", recorder.Code, http.StatusInternalServerError)
	}
}
