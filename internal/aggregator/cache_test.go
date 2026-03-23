package aggregator

import (
	"testing"
	"time"

	"flight-search/internal/domain"
)

func TestCacheManager_BasicCaching(t *testing.T) {
	cm := NewCacheManager()

	request := domain.SearchRequest{
		Origin:        "CGK",
		Destination:   "DPS",
		DepartureDate: "2025-12-15",
		Passengers:    1,
		CabinClass:    "economy",
	}

	result := domain.SearchResult{
		SearchCriteria: domain.SearchCriteria{
			Origin:        "CGK",
			Destination:   "DPS",
			DepartureDate: "2025-12-15",
			Passengers:    1,
			CabinClass:    "economy",
		},
		Metadata: domain.Metadata{
			TotalResults:       5,
			ProvidersQueried:   4,
			ProvidersSucceeded: 4,
			ProvidersFailed:    0,
			SearchTimeMs:       100,
			CacheHit:           false,
		},
		Flights: []domain.Flight{
			{
				ID:             "test-flight-1",
				Provider:       "Test Provider",
				FlightNumber:   "TP001",
				Price:          domain.Price{Amount: 1000000, Currency: "IDR"},
				AvailableSeats: 10,
			},
		},
	}

	// Test cache miss
	_, found := cm.Get(request)
	if found {
		t.Error("Expected cache miss, but got cache hit")
	}

	// Store in cache
	cm.Set(request, result)

	// Test cache hit
	cachedResult, found := cm.Get(request)
	if !found {
		t.Error("Expected cache hit, but got cache miss")
	}

	if cachedResult.Metadata.CacheHit != true {
		t.Error("Expected CacheHit to be true for cached result")
	}

	if cachedResult.Metadata.TotalResults != result.Metadata.TotalResults {
		t.Errorf("Expected total results %d, got %d", result.Metadata.TotalResults, cachedResult.Metadata.TotalResults)
	}

	if len(cachedResult.Flights) != len(result.Flights) {
		t.Errorf("Expected %d flights, got %d", len(result.Flights), len(cachedResult.Flights))
	}
}

func TestCacheManager_DifferentRequests(t *testing.T) {
	cm := NewCacheManager()

	request1 := domain.SearchRequest{
		Origin:        "CGK",
		Destination:   "DPS",
		DepartureDate: "2025-12-15",
		Passengers:    1,
		CabinClass:    "economy",
	}

	request2 := domain.SearchRequest{
		Origin:        "CGK",
		Destination:   "DPS",
		DepartureDate: "2025-12-16", // Different date
		Passengers:    1,
		CabinClass:    "economy",
	}

	result1 := domain.SearchResult{
		Metadata: domain.Metadata{TotalResults: 5},
	}
	result2 := domain.SearchResult{
		Metadata: domain.Metadata{TotalResults: 3},
	}

	// Store both results
	cm.Set(request1, result1)
	cm.Set(request2, result2)

	// Test retrieval
	cached1, found1 := cm.Get(request1)
	cached2, found2 := cm.Get(request2)

	if !found1 || !found2 {
		t.Error("Expected both requests to be found in cache")
	}

	if cached1.Metadata.TotalResults != 5 {
		t.Errorf("Expected result1 total results 5, got %d", cached1.Metadata.TotalResults)
	}

	if cached2.Metadata.TotalResults != 3 {
		t.Errorf("Expected result2 total results 3, got %d", cached2.Metadata.TotalResults)
	}
}

func TestCacheManager_CacheExpiration(t *testing.T) {
	// Create cache with very short expiration (1 millisecond)
	cm := NewCacheManagerWithTTL(1*time.Millisecond, 1*time.Millisecond)

	request := domain.SearchRequest{
		Origin:        "CGK",
		Destination:   "DPS",
		DepartureDate: "2025-12-15",
		Passengers:    1,
		CabinClass:    "economy",
	}

	result := domain.SearchResult{
		Metadata: domain.Metadata{TotalResults: 5},
	}

	// Store in cache
	cm.Set(request, result)

	// Immediately check - should be found
	_, found := cm.Get(request)
	if !found {
		t.Error("Expected cache hit immediately after setting")
	}

	// Wait for expiration
	time.Sleep(10 * time.Millisecond)

	// Check again - should be expired
	_, found = cm.Get(request)
	if found {
		t.Error("Expected cache miss after expiration")
	}
}

func TestCacheManager_CacheKeyGeneration(t *testing.T) {
	cm := NewCacheManager()

	request1 := domain.SearchRequest{
		Origin:        "CGK",
		Destination:   "DPS",
		DepartureDate: "2025-12-15",
		Passengers:    1,
		CabinClass:    "economy",
		PriceRange:    []int{1000000, 2000000},
	}

	request2 := domain.SearchRequest{
		Origin:        "CGK",
		Destination:   "DPS",
		DepartureDate: "2025-12-15",
		Passengers:    1,
		CabinClass:    "economy",
		PriceRange:    []int{1000000, 2000000},
	}

	// Same requests should generate same key
	key1 := cm.generateCacheKey(request1)
	key2 := cm.generateCacheKey(request2)

	if key1 != key2 {
		t.Errorf("Expected same cache keys for identical requests, got %s and %s", key1, key2)
	}

	// Different requests should generate different keys
	request3 := domain.SearchRequest{
		Origin:        "CGK",
		Destination:   "SUB", // Different destination
		DepartureDate: "2025-12-15",
		Passengers:    1,
		CabinClass:    "economy",
		PriceRange:    []int{1000000, 2000000},
	}

	key3 := cm.generateCacheKey(request3)
	if key1 == key3 {
		t.Error("Expected different cache keys for different requests")
	}
}

func TestCacheManager_CacheOperations(t *testing.T) {
	cm := NewCacheManager()

	// Test initial state
	if cm.Size() != 0 {
		t.Errorf("Expected empty cache, got size %d", cm.Size())
	}

	request := domain.SearchRequest{
		Origin:        "CGK",
		Destination:   "DPS",
		DepartureDate: "2025-12-15",
		Passengers:    1,
		CabinClass:    "economy",
	}

	result := domain.SearchResult{
		Metadata: domain.Metadata{TotalResults: 5},
	}

	// Add item
	cm.Set(request, result)
	if cm.Size() != 1 {
		t.Errorf("Expected cache size 1, got %d", cm.Size())
	}

	// Clear cache
	cm.Clear()
	if cm.Size() != 0 {
		t.Errorf("Expected empty cache after clear, got size %d", cm.Size())
	}

	// Verify item is gone
	_, found := cm.Get(request)
	if found {
		t.Error("Expected cache miss after clear")
	}
}

func TestAggregator_CacheIntegration(t *testing.T) {
	agg := NewTestAggregator()

	request := domain.SearchRequest{
		Origin:        "CGK",
		Destination:   "DPS",
		DepartureDate: "2025-12-15",
		Passengers:    1,
		CabinClass:    "economy",
	}

	// First search - should be cache miss
	result1, err := agg.Search(request)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if result1.Metadata.CacheHit {
		t.Error("Expected cache miss on first search")
	}

	// Second search - should be cache hit
	result2, err := agg.Search(request)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if !result2.Metadata.CacheHit {
		t.Error("Expected cache hit on second search")
	}

	// Results should be identical
	if result1.Metadata.TotalResults != result2.Metadata.TotalResults {
		t.Errorf("Expected same total results, got %d and %d", result1.Metadata.TotalResults, result2.Metadata.TotalResults)
	}

	if len(result1.Flights) != len(result2.Flights) {
		t.Errorf("Expected same number of flights, got %d and %d", len(result1.Flights), len(result2.Flights))
	}
}
