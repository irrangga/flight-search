package aggregator

import (
	"fmt"
	"sync"
	"time"

	"flight-search/internal/domain"
	"flight-search/internal/normalizer"
	"flight-search/internal/provider"
)

// Aggregator orchestrates flight search across providers
type Aggregator struct {
	providers []provider.Provider
	cache     *CacheManager
}

// NewAggregator creates a new flight aggregator
func NewAggregator() *Aggregator {
	return &Aggregator{
		providers: []provider.Provider{
			provider.NewGarudaIndonesiaProvider(),
			provider.NewLionAirProvider(),
			provider.NewBatikAirProvider(),
			provider.NewAirAsiaProvider(),
		},
		cache: NewCacheManager(),
	}
}

// NewTestAggregator creates a new flight aggregator for testing (with rate limiting disabled)
func NewTestAggregator() *Aggregator {
	return &Aggregator{
		providers: []provider.Provider{
			provider.NewGarudaIndonesiaProviderForTest(),
			provider.NewLionAirProviderForTest(),
			provider.NewBatikAirProviderForTest(),
			provider.NewAirAsiaProviderForTest(),
		},
		cache: NewCacheManager(),
	}
}

// Search performs parallel flight search across all providers
func (a *Aggregator) Search(request domain.SearchRequest) (domain.SearchResult, error) {
	// Check cache first
	if cachedResult, found := a.cache.Get(request); found {
		cachedResult.Metadata.CacheHit = true
		return cachedResult, nil
	}

	start := time.Now()

	// Channel to collect results
	resultsChan := make(chan domain.ProviderResult, len(a.providers))

	// Launch goroutines for parallel execution
	var wg sync.WaitGroup
	for _, p := range a.providers {
		wg.Add(1)
		go func(prov provider.Provider) {
			defer wg.Done()
			result := prov.SearchFlights(request)
			resultsChan <- result
		}(p)
	}

	// Close channel when all goroutines are done
	go func() {
		wg.Wait()
		close(resultsChan)
	}()

	// Collect results
	var successfulResults []domain.ProviderResult
	var failedProviders []error

	for result := range resultsChan {
		if result.Success {
			successfulResults = append(successfulResults, result)
		} else {
			failedProviders = append(failedProviders, fmt.Errorf("%s: %v", result.Provider, result.Error))
		}
	}

	// If all providers failed, return error
	if len(successfulResults) == 0 {
		return domain.SearchResult{}, fmt.Errorf("all providers failed: %v", failedProviders)
	}

	// Normalize flights from all providers
	normalizedFlights := normalizer.NormalizeFlightsFromResults(successfulResults)

	// Apply filters
	normalizedFlights = a.Filter(normalizedFlights, request)

	// Apply sorting
	a.Sort(normalizedFlights, request.SortBy)

	// Create response
	searchTimeMs := time.Since(start).Milliseconds()

	result := domain.SearchResult{
		SearchCriteria: domain.SearchCriteria{
			Origin:        request.Origin,
			Destination:   request.Destination,
			DepartureDate: request.DepartureDate,
			ReturnDate:    request.ReturnDate,
			Passengers:    request.Passengers,
			CabinClass:    request.CabinClass,
		},
		Metadata: domain.Metadata{
			TotalResults:       len(normalizedFlights),
			ProvidersQueried:   len(a.providers),
			ProvidersSucceeded: len(successfulResults),
			ProvidersFailed:    len(failedProviders),
			SearchTimeMs:       searchTimeMs,
			CacheHit:           false,
		},
		Flights: normalizedFlights,
	}

	// Cache the result
	a.cache.Set(request, result)

	return result, nil
}
