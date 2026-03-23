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

// NewTestAggregator creates a new flight aggregator for testing (with rate limiting disabled)
func NewTestAggregator() *Aggregator {
	return &Aggregator{
		providers: []provider.Provider{
			provider.NewGarudaIndonesiaProviderForTest(),
			provider.NewLionAirProviderForTest(),
			provider.NewBatikAirProviderForTest(),
			provider.NewAirAsiaProviderForTest(),
		},
	}
}

// Search performs parallel flight search across all providers
func (a *Aggregator) Search(request domain.SearchRequest) (domain.SearchResult, error) {
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
	var normalizedFlights []domain.Flight
	for _, result := range successfulResults {
		flights := a.extractFlightsFromResponse(result.Data, result.Provider)
		for _, flight := range flights {
			normalized, err := normalizer.NormalizeFlight(flight, result.Provider)
			if err != nil {
				fmt.Printf("Error normalizing flight from %s: %v\n", result.Provider, err)
				continue
			}
			// Calculate score: price + duration (in minutes) * 1250 to give more weight to price (based on 4 hours delay get Rp300.000,-)
			normalized.Score = float64(normalized.Price.Amount) + float64(normalized.Duration.TotalMinutes*1250)

			normalizedFlights = append(normalizedFlights, normalized)
		}
	}

	// Apply filters
	normalizedFlights = a.applyFilters(normalizedFlights, request)

	// Apply sorting
	a.applySorting(normalizedFlights, request.SortBy)

	// Create response
	searchTimeMs := time.Since(start).Milliseconds()

	return domain.SearchResult{
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
	}, nil
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

// applyFilters applies the filters from the request to the flights
func (a *Aggregator) applyFilters(flights []domain.Flight, request domain.SearchRequest) []domain.Flight {
	var filtered []domain.Flight
	for _, flight := range flights {
		if a.matchesFilters(flight, request) {
			filtered = append(filtered, flight)
		}
	}
	return filtered
}

// matchesFilters checks if a flight matches all the filter criteria
func (a *Aggregator) matchesFilters(flight domain.Flight, request domain.SearchRequest) bool {
	// Basic search criteria
	if flight.Departure.Airport != request.Origin {
		return false
	}
	if flight.Arrival.Airport != request.Destination {
		return false
	}
	if request.DepartureDate != "" && flight.Departure.Datetime[:10] != request.DepartureDate {
		return false
	}
	if request.CabinClass != "" && flight.CabinClass != request.CabinClass {
		return false
	}
	if flight.AvailableSeats < request.Passengers {
		return false
	}

	// Price range
	if len(request.PriceRange) == 2 {
		minPrice := request.PriceRange[0]
		maxPrice := request.PriceRange[1]
		if flight.Price.Amount < minPrice || flight.Price.Amount > maxPrice {
			return false
		}
	}

	// Number of stops range
	if len(request.StopsRange) == 2 {
		minStops := request.StopsRange[0]
		maxStops := request.StopsRange[1]
		if flight.Stops < minStops || flight.Stops > maxStops {
			return false
		}
	}

	// Airlines
	if len(request.Airlines) > 0 {
		found := false
		for _, airline := range request.Airlines {
			if flight.Airline.Name == airline {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	// Duration range
	if len(request.DurationRange) == 2 {
		minDuration := request.DurationRange[0]
		maxDuration := request.DurationRange[1]
		if flight.Duration.TotalMinutes < minDuration || flight.Duration.TotalMinutes > maxDuration {
			return false
		}
	}

	// Departure time range
	if len(request.DepartureTimeRange) == 2 {
		minTime := request.DepartureTimeRange[0]
		maxTime := request.DepartureTimeRange[1]
		if !a.isTimeInRange(flight.Departure.Datetime, minTime, maxTime) {
			return false
		}
	}

	// Arrival time range
	if len(request.ArrivalTimeRange) == 2 {
		minTime := request.ArrivalTimeRange[0]
		maxTime := request.ArrivalTimeRange[1]
		if !a.isTimeInRange(flight.Arrival.Datetime, minTime, maxTime) {
			return false
		}
	}

	return true
}

// isTimeInRange checks if the flight time is within the specified range (HH:MM format)
func (a *Aggregator) isTimeInRange(flightTime, minTime, maxTime string) bool {
	// Extract HH:MM from flightTime (assume it's in RFC3339 format)
	// For simplicity, compare the time part
	flightHourMin := flightTime[11:16] // HH:MM from "2025-12-15T04:45:00+07:00"
	return flightHourMin >= minTime && flightHourMin <= maxTime
}

// applySorting sorts the flights based on the sortBy parameter
func (a *Aggregator) applySorting(flights []domain.Flight, sortBy string) {
	switch sortBy {
	case "price_asc":
		sort.Slice(flights, func(i, j int) bool {
			return flights[i].Price.Amount < flights[j].Price.Amount
		})
	case "price_desc":
		sort.Slice(flights, func(i, j int) bool {
			return flights[i].Price.Amount > flights[j].Price.Amount
		})
	case "duration_asc":
		sort.Slice(flights, func(i, j int) bool {
			return flights[i].Duration.TotalMinutes < flights[j].Duration.TotalMinutes
		})
	case "duration_desc":
		sort.Slice(flights, func(i, j int) bool {
			return flights[i].Duration.TotalMinutes > flights[j].Duration.TotalMinutes
		})
	case "departure_asc":
		sort.Slice(flights, func(i, j int) bool {
			return flights[i].Departure.Timestamp < flights[j].Departure.Timestamp
		})
	case "departure_desc":
		sort.Slice(flights, func(i, j int) bool {
			return flights[i].Departure.Timestamp > flights[j].Departure.Timestamp
		})
	case "arrival_asc":
		sort.Slice(flights, func(i, j int) bool {
			return flights[i].Arrival.Timestamp < flights[j].Arrival.Timestamp
		})
	case "arrival_desc":
		sort.Slice(flights, func(i, j int) bool {
			return flights[i].Arrival.Timestamp > flights[j].Arrival.Timestamp
		})
	default:
		// Default: based on score, the lower the score, the better
		sort.Slice(flights, func(i, j int) bool {
			return flights[i].Score < flights[j].Score
		})
	}
}
