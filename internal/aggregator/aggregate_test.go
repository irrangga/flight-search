package aggregator

import (
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

	agg := NewAggregator()
	result := agg.Search(request)

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
	}
}

func TestSearch_FlightsAreNormalized(t *testing.T) {
	request := domain.SearchRequest{
		Origin:        "CGK",
		Destination:   "DPS",
		DepartureDate: "2025-12-15",
		Passengers:    1,
		CabinClass:    "economy",
	}

	agg := NewAggregator()
	result := agg.Search(request)

	if len(result.Flights) == 0 {
		t.Skip("No flights found, skipping normalization test")
	}

	// Verify flights are from different providers
	providersFound := make(map[string]bool)
	for _, flight := range result.Flights {
		providersFound[flight.Provider] = true
	}

	expectedProviders := map[string]bool{
		"Garuda Indonesia": true,
		"Lion Air":         true,
		"Batik Air":        true,
		"AirAsia":          true,
	}

	for provider := range expectedProviders {
		if !providersFound[provider] {
			t.Logf("Warning: Provider %s not found in results (may have failed)", provider)
		}
	}
}

func TestSearch_FlightsSortedByPrice(t *testing.T) {
	request := domain.SearchRequest{
		Origin:        "CGK",
		Destination:   "DPS",
		DepartureDate: "2025-12-15",
		Passengers:    1,
		CabinClass:    "economy",
	}

	agg := NewAggregator()
	result := agg.Search(request)

	if len(result.Flights) < 2 {
		t.Skip("Not enough flights to test sorting")
	}

	// Verify flights are sorted by price ascending
	for i := 1; i < len(result.Flights); i++ {
		if result.Flights[i-1].Price.Amount > result.Flights[i].Price.Amount {
			t.Error("Flights not sorted by price ascending")
		}
	}
}
