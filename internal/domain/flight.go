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

	// Filters
	PriceRange         []int
	StopsRange         []int
	DepartureTimeRange []string
	ArrivalTimeRange   []string
	Airlines           []string
	DurationRange      []int

	// Sort
	SortBy string // "price_asc", "price_desc", "duration_asc", "duration_desc", "departure_asc", "departure_desc", "arrival_asc", "arrival_desc"
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
	Score          float64 // For ranking purposes, the lower the score, the better
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
	ReturnDate    *string
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
