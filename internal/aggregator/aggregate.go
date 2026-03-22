package aggregator

import (
	"fmt"
	"sort"
	"sync"
	"time"

	"flight-search/internal/domain"
	"flight-search/internal/normalizer"
	"flight-search/internal/provider"
)

// Aggregator orchestrates flight search across providers
type Aggregator struct {
	providers []provider.Provider
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
	}
}

// Search performs parallel flight search across all providers
func (a *Aggregator) Search(request domain.SearchRequest) domain.SearchResult {
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
	var failedProviders []string

	for result := range resultsChan {
		if result.Success {
			successfulResults = append(successfulResults, result)
		} else {
			failedProviders = append(failedProviders, result.Provider)
		}
	}

	// Normalize flights from all providers
	var normalizedFlights []domain.Flight
	for _, result := range successfulResults {
		flights := a.extractFlightsFromResponse(result.Data, result.Provider)
		for _, flight := range flights {
			normalized, err := normalizer.NormalizeFlight(flight, result.Provider)
			if err != nil {
				fmt.Printf("Error normalizing flight from %s: %v\n", result.Provider, err)
				continue
			}
			normalizedFlights = append(normalizedFlights, normalized)
		}
	}

	// Sort by price (lowest first) as default
	sort.Slice(normalizedFlights, func(i, j int) bool {
		return normalizedFlights[i].Price.Amount < normalizedFlights[j].Price.Amount
	})

	// Create response
	searchTimeMs := time.Since(start).Milliseconds()

	return domain.SearchResult{
		SearchCriteria: domain.SearchCriteria{
			Origin:        request.Origin,
			Destination:   request.Destination,
			DepartureDate: request.DepartureDate,
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
}

// extractFlightsFromResponse extracts flight data from provider response
func (a *Aggregator) extractFlightsFromResponse(data map[string]interface{}, provider string) []map[string]interface{} {
	switch provider {
	case "Garuda Indonesia":
		if flights, ok := data["flights"].([]interface{}); ok {
			result := make([]map[string]interface{}, len(flights))
			for i, f := range flights {
				if flight, ok := f.(map[string]interface{}); ok {
					result[i] = flight
				}
			}
			return result
		}
	case "Lion Air":
		if dataSection, ok := data["data"].(map[string]interface{}); ok {
			if flights, ok := dataSection["available_flights"].([]interface{}); ok {
				result := make([]map[string]interface{}, len(flights))
				for i, f := range flights {
					if flight, ok := f.(map[string]interface{}); ok {
						result[i] = flight
					}
				}
				return result
			}
		}
	case "Batik Air":
		if results, ok := data["results"].([]interface{}); ok {
			result := make([]map[string]interface{}, len(results))
			for i, r := range results {
				if flight, ok := r.(map[string]interface{}); ok {
					result[i] = flight
				}
			}
			return result
		}
	case "AirAsia":
		if flights, ok := data["flights"].([]interface{}); ok {
			result := make([]map[string]interface{}, len(flights))
			for i, f := range flights {
				if flight, ok := f.(map[string]interface{}); ok {
					result[i] = flight
				}
			}
			return result
		}
	}
	return []map[string]interface{}{}
}
