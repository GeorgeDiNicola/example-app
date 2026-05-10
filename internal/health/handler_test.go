package health

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler(t *testing.T) {
	// mock HTTP request to pass to the handler
	req, err := http.NewRequest("GET", "/health", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder (acts like an http.ResponseWriter) to record the response
	// and convert the handler func to an http.Handler
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Handler)

	// Call the handler directly
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expectedHeader := "application/json"
	if contentType := rr.Header().Get("Content-Type"); contentType != expectedHeader {
		t.Errorf("handler returned wrong content type: got %v want %v", contentType, expectedHeader)
	}

	var response HealthResponse
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatalf("could not decode JSON response: %v", err)
	}

	if response.Status != "up" {
		t.Errorf("handler returned unexpected status: got %v want %v", response.Status, "up")
	}

	if response.Timestamp.IsZero() {
		t.Error("expected a non-zero timestamp")
	}

}
