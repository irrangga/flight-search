package aggregator

import (
	"flight-search/internal/domain"
)

// Filter applies filters to flights based on the request
func (a *Aggregator) Filter(flights []domain.Flight, request domain.SearchRequest) []domain.Flight {
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

	// Number of stops - exact values list
	if len(request.NumberOfStops) > 0 {
		found := false
		for _, allowedStops := range request.NumberOfStops {
			if flight.Stops == allowedStops {
				found = true
				break
			}
		}
		if !found {
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
