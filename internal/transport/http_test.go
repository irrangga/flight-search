package transport

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"flight-search/internal/aggregator"
)

func newTestRouter() *gin.Engine {
	agg := aggregator.NewAggregator()
	h := NewHandler(agg)
	return SetupRouter(h)
}

func TestSearchFlights_ValidRequest(t *testing.T) {
	router := newTestRouter()

	requestBody := SearchRequest{
		Origin:        "CGK",
		Destination:   "DPS",
		DepartureDate: "2025-12-15",
		Passengers:    1,
		CabinClass:    "economy",
	}

	body, _ := json.Marshal(requestBody)
	req := httptest.NewRequest("POST", "/api/search-flights", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", w.Code)
	}

	var resp map[string]interface{}
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("invalid response JSON: %v", err)
	}

	if _, ok := resp["metadata"]; !ok {
		t.Error("Expected metadata in response")
	}
}

func TestSearchFlights_InvalidMethod(t *testing.T) {
	router := newTestRouter()

	req := httptest.NewRequest("GET", "/api/search-flights", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status 405, got %d", w.Code)
	}
}

func TestSearchFlights_InvalidJSON(t *testing.T) {
	router := newTestRouter()

	req := httptest.NewRequest("POST", "/api/search-flights", bytes.NewReader([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}
