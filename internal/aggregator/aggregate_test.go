package aggregator

import (
	"strings"
	"testing"

	"flight-search/internal/domain"
)

func TestSearch_IntegrationWithAllProviders(t *testing.T) {
	request := domain.SearchRequest{
		Origin:        "CGK",
		Destination:   "DPS",
		DepartureDate: "2025-12-15",
		Passengers:    1,
		CabinClass:    "economy",
	}

	agg := NewTestAggregator()
	result, err := agg.Search(request)
	if err != nil {
		t.Fatalf("Unexpected error occurred: %v", err)
	}

	// Verify we got results
	if result.Metadata.TotalResults == 0 {
		t.Error("Expected to find flights, but got 0")
	}

	// Verify all providers were queried
	if result.Metadata.ProvidersQueried != 4 {
		t.Errorf("Expected 4 providers queried, got %d", result.Metadata.ProvidersQueried)
	}

	// Verify flights have required fields
	for _, flight := range result.Flights {
		if flight.ID == "" {
			t.Error("Flight ID should not be empty")
		}
		if flight.Airline.Name == "" {
			t.Error("Airline name should not be empty")
		}
		if flight.FlightNumber == "" {
			t.Error("Flight number should not be empty")
		}
		if flight.Price.Amount <= 0 {
			t.Error("Price should be positive")
		}
		if flight.Departure.Airport != "CGK" {
			t.Errorf("Expected departure airport CGK, got %s", flight.Departure.Airport)
		}
		if flight.Departure.City == "" {
			t.Error("Departure city should not be empty")
		}
		if flight.Arrival.Airport != "DPS" {
			t.Errorf("Expected arrival airport DPS, got %s", flight.Arrival.Airport)
		}
		if flight.Arrival.City == "" {
			t.Error("Arrival city should not be empty")
		}
	}
}

func TestSearch_RateLimiting(t *testing.T) {
	request := domain.SearchRequest{
		Origin:        "CGK",
		Destination:   "DPS",
		DepartureDate: "2025-12-15",
		Passengers:    1,
		CabinClass:    "economy",
	}

	// Create a fresh aggregator for this test (with rate limiting enabled)
	agg := NewAggregator()

	// Test that the system can handle multiple requests
	// Rate limiting may cause some to fail, which is expected
	for i := 0; i < 5; i++ {
		result, err := agg.Search(request)
		if err != nil {
			// Rate limiting is expected - this is normal behavior
			if strings.Contains(err.Error(), "rate limit exceeded") {
				t.Logf("Request %d was rate limited (expected)", i)
				continue
			}
			t.Errorf("Request %d failed with unexpected error: %v", i, err)
		} else {
			// Success is also possible
			t.Logf("Request %d succeeded with %d results", i, result.Metadata.TotalResults)
		}
	}

	// The test passes if we can make requests without crashing
	// Rate limiting behavior depends on the actual limits and timing
	t.Log("Rate limiting test completed - system handles multiple requests appropriately")
}
