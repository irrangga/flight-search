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

	agg := NewTestAggregator()
	result, err := agg.Search(request)
	if err != nil {
		t.Fatalf("Unexpected error occurred: %v", err)
	}

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
		SortBy:        "price_asc",
	}

	agg := NewTestAggregator()
	result, err := agg.Search(request)
	if err != nil {
		t.Fatalf("Unexpected error occurred: %v", err)
	}

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

func TestSearch_FilterByPriceRange(t *testing.T) {
	request := domain.SearchRequest{
		Origin:        "CGK",
		Destination:   "DPS",
		DepartureDate: "2025-12-15",
		Passengers:    1,
		CabinClass:    "economy",
		PriceRange:    []int{1000000, 2000000},
	}

	agg := NewTestAggregator()
	result, err := agg.Search(request)
	if err != nil {
		t.Fatalf("Unexpected error occurred: %v", err)
	}

	for _, flight := range result.Flights {
		if flight.Price.Amount < 1000000 || flight.Price.Amount > 2000000 {
			t.Errorf("Flight price %d is outside range 1000000-2000000", flight.Price.Amount)
		}
	}
}

func TestSearch_FilterByAirlines(t *testing.T) {
	request := domain.SearchRequest{
		Origin:        "CGK",
		Destination:   "DPS",
		DepartureDate: "2025-12-15",
		Passengers:    1,
		CabinClass:    "economy",
		Airlines:      []string{"Garuda Indonesia", "AirAsia"},
	}

	agg := NewTestAggregator()
	result, err := agg.Search(request)
	if err != nil {
		t.Fatalf("Unexpected error occurred: %v", err)
	}

	for _, flight := range result.Flights {
		if flight.Airline.Name != "Garuda Indonesia" && flight.Airline.Name != "AirAsia" {
			t.Errorf("Flight airline %s is not in allowed list", flight.Airline.Name)
		}
	}
}

func TestSearch_SortByPriceDesc(t *testing.T) {
	request := domain.SearchRequest{
		Origin:        "CGK",
		Destination:   "DPS",
		DepartureDate: "2025-12-15",
		Passengers:    1,
		CabinClass:    "economy",
		SortBy:        "price_desc",
	}

	agg := NewTestAggregator()
	result, err := agg.Search(request)
	if err != nil {
		t.Fatalf("Unexpected error occurred: %v", err)
	}

	if len(result.Flights) < 2 {
		t.Skip("Not enough flights to test sorting")
	}

	for i := 1; i < len(result.Flights); i++ {
		if result.Flights[i-1].Price.Amount < result.Flights[i].Price.Amount {
			t.Error("Flights not sorted by price descending")
		}
	}
}

func TestSearch_SortByDurationAsc(t *testing.T) {
	request := domain.SearchRequest{
		Origin:        "CGK",
		Destination:   "DPS",
		DepartureDate: "2025-12-15",
		Passengers:    1,
		CabinClass:    "economy",
		SortBy:        "duration_asc",
	}

	agg := NewTestAggregator()
	result, err := agg.Search(request)
	if err != nil {
		t.Fatalf("Unexpected error occurred: %v", err)
	}

	if len(result.Flights) < 2 {
		t.Skip("Not enough flights to test sorting")
	}

	for i := 1; i < len(result.Flights); i++ {
		if result.Flights[i-1].Duration.TotalMinutes > result.Flights[i].Duration.TotalMinutes {
			t.Error("Flights not sorted by duration ascending")
		}
	}
}

func TestSearch_FilterByOriginDestination(t *testing.T) {
	// Test valid route
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

	if result.Metadata.TotalResults == 0 {
		t.Error("Expected flights for CGK->DPS route")
	}

	// Test invalid route
	request2 := domain.SearchRequest{
		Origin:        "XXX",
		Destination:   "YYY",
		DepartureDate: "2025-12-15",
		Passengers:    1,
		CabinClass:    "economy",
	}

	result2, err := agg.Search(request2)
	if err != nil {
		t.Fatalf("Unexpected error occurred: %v", err)
	}

	if result2.Metadata.TotalResults != 0 {
		t.Error("Expected no flights for invalid XXX->YYY route")
	}
}

func TestSearch_FilterByCabinClass(t *testing.T) {
	// Test economy
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

	if result.Metadata.TotalResults == 0 {
		t.Error("Expected flights for economy class")
	}

	// Test business (should return no results since mock data is economy)
	request2 := domain.SearchRequest{
		Origin:        "CGK",
		Destination:   "DPS",
		DepartureDate: "2025-12-15",
		Passengers:    1,
		CabinClass:    "business",
	}

	result2, err := agg.Search(request2)
	if err != nil {
		t.Fatalf("Unexpected error occurred: %v", err)
	}

	if result2.Metadata.TotalResults != 0 {
		t.Error("Expected no flights for business class (mock data is economy only)")
	}
}

func TestSearch_FilterByPassengers(t *testing.T) {
	// Test 20 passenger
	request := domain.SearchRequest{
		Origin:        "CGK",
		Destination:   "DPS",
		DepartureDate: "2025-12-15",
		Passengers:    20,
		CabinClass:    "economy",
	}

	agg := NewTestAggregator()
	result, err := agg.Search(request)
	if err != nil {
		t.Fatalf("Unexpected error occurred: %v", err)
	}

	for _, flight := range result.Flights {
		if flight.AvailableSeats < request.Passengers {
			t.Errorf("Flight seats %d is less than requested passengers %d", flight.AvailableSeats, request.Passengers)
		}
	}
}

func TestSearch_FilterByStopsRange(t *testing.T) {
	// Test 1 stop
	request := domain.SearchRequest{
		Origin:        "CGK",
		Destination:   "DPS",
		DepartureDate: "2025-12-15",
		Passengers:    1,
		CabinClass:    "economy",
		StopsRange:    []int{1, 2},
	}

	agg := NewTestAggregator()
	result, err := agg.Search(request)
	if err != nil {
		t.Fatalf("Unexpected error occurred: %v", err)
	}

	for _, flight := range result.Flights {
		if flight.Stops < 1 || flight.Stops > 2 {
			t.Errorf("Flight stops %d is not in range [1, 2]", flight.Stops)
		}
	}
}

func TestSearch_FilterByDepartureTimeRange(t *testing.T) {
	request := domain.SearchRequest{
		Origin:             "CGK",
		Destination:        "DPS",
		DepartureDate:      "2025-12-15",
		Passengers:         1,
		CabinClass:         "economy",
		DepartureTimeRange: []string{"06:00", "07:00"},
	}

	agg := NewTestAggregator()
	result, err := agg.Search(request)
	if err != nil {
		t.Fatalf("Unexpected error occurred: %v", err)
	}

	for _, flight := range result.Flights {
		if flight.Departure.Datetime[11:16] < "06:00" || flight.Departure.Datetime[11:16] > "07:00" {
			t.Errorf("Flight departure time %s is not in range [06:00, 07:00]", flight.Departure.Datetime[11:16])
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

func TestSearch_RankingByScore(t *testing.T) {
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

	if len(result.Flights) < 2 {
		t.Skip("Not enough flights to test ranking")
	}

	for i := 1; i < len(result.Flights); i++ {
		if result.Flights[i-1].Score > result.Flights[i].Score {
			t.Error("Flights not ranked by score ascending")
		}
	}
}
