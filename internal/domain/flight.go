package domain

import "time"

// SearchRequest represents a flight search request
type SearchRequest struct {
	Origin        string
	Destination   string
	DepartureDate string
	ReturnDate    *string
	Passengers    int
	CabinClass    string
}

// ProviderResult represents the result from a single provider
type ProviderResult struct {
	Provider     string
	Success      bool
	Data         map[string]interface{}
	ResponseTime time.Duration
	Error        error
}

// Flight is the unified flight domain model
type Flight struct {
	ID             string
	Provider       string
	Airline        Airline
	FlightNumber   string
	Departure      FlightEndpoint
	Arrival        FlightEndpoint
	Duration       Duration
	Stops          int
	Price          Price
	AvailableSeats int
	CabinClass     string
	Aircraft       *string
	Amenities      []string
	Baggage        Baggage
}

// Airline represents airline information
type Airline struct {
	Name string
	Code string
}

// FlightEndpoint represents departure or arrival information
type FlightEndpoint struct {
	Airport   string
	City      string
	Datetime  string
	Timestamp int64
}

// Duration represents flight duration
type Duration struct {
	TotalMinutes int
	Formatted    string
}

// Price represents flight price
type Price struct {
	Amount   int
	Currency string
}

// Baggage represents baggage information
type Baggage struct {
	CarryOn string
	Checked string
}

// SearchResult contains the search response data
type SearchResult struct {
	SearchCriteria SearchCriteria
	Metadata       Metadata
	Flights        []Flight
}

// SearchCriteria represents the search criteria used
type SearchCriteria struct {
	Origin        string
	Destination   string
	DepartureDate string
	Passengers    int
	CabinClass    string
}

// Metadata represents search metadata
type Metadata struct {
	TotalResults       int
	ProvidersQueried   int
	ProvidersSucceeded int
	ProvidersFailed    int
	SearchTimeMs       int64
	CacheHit           bool
}
