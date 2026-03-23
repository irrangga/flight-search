package mapper

import (
	"flight-search/internal/domain"
	"flight-search/internal/utils"
)

// SearchResponse is the HTTP response DTO
type SearchResponse struct {
	SearchCriteria SearchCriteria `json:"search_criteria"`
	Metadata       Metadata       `json:"metadata"`
	Flights        []FlightDTO    `json:"flights"`
}

// SearchCriteria represents the search criteria used
type SearchCriteria struct {
	Origin        string  `json:"origin"`
	Destination   string  `json:"destination"`
	DepartureDate string  `json:"departure_date"`
	ReturnDate    *string `json:"return_date,omitempty"`
	Passengers    int     `json:"passengers"`
	CabinClass    string  `json:"cabin_class"`
}

// Metadata represents search metadata
type Metadata struct {
	TotalResults       int   `json:"total_results"`
	ProvidersQueried   int   `json:"providers_queried"`
	ProvidersSucceeded int   `json:"providers_succeeded"`
	ProvidersFailed    int   `json:"providers_failed"`
	SearchTimeMs       int64 `json:"search_time_ms"`
	CacheHit           bool  `json:"cache_hit"`
}

// FlightDTO is the HTTP response flight DTO
type FlightDTO struct {
	ID             string   `json:"id"`
	Provider       string   `json:"provider"`
	Airline        Airline  `json:"airline"`
	FlightNumber   string   `json:"flight_number"`
	Departure      Endpoint `json:"departure"`
	Arrival        Endpoint `json:"arrival"`
	Duration       Duration `json:"duration"`
	Stops          int      `json:"stops"`
	Price          Price    `json:"price"`
	AvailableSeats int      `json:"available_seats"`
	CabinClass     string   `json:"cabin_class"`
	Aircraft       *string  `json:"aircraft"`
	Amenities      []string `json:"amenities"`
	Baggage        Baggage  `json:"baggage"`
}

// Airline DTO
type Airline struct {
	Name string `json:"name"`
	Code string `json:"code"`
}

// Endpoint DTO
type Endpoint struct {
	Airport   string `json:"airport"`
	City      string `json:"city"`
	Datetime  string `json:"datetime"`
	Timestamp int64  `json:"timestamp"`
}

// Duration DTO
type Duration struct {
	TotalMinutes int    `json:"total_minutes"`
	Formatted    string `json:"formatted"`
}

// Price DTO
type Price struct {
	Amount    int    `json:"amount"`
	Currency  string `json:"currency"`
	Formatted string `json:"formatted"`
}

// Baggage DTO
type Baggage struct {
	CarryOn string `json:"carry_on"`
	Checked string `json:"checked"`
}

// ToSearchResponse converts domain SearchResult to HTTP response DTO
func ToSearchResponse(result domain.SearchResult) SearchResponse {
	flightDTOs := make([]FlightDTO, len(result.Flights))
	for i, flight := range result.Flights {
		flightDTOs[i] = toFlightDTO(flight)
	}

	return SearchResponse{
		SearchCriteria: SearchCriteria{
			Origin:        result.SearchCriteria.Origin,
			Destination:   result.SearchCriteria.Destination,
			DepartureDate: result.SearchCriteria.DepartureDate,
			ReturnDate:    result.SearchCriteria.ReturnDate,
			Passengers:    result.SearchCriteria.Passengers,
			CabinClass:    result.SearchCriteria.CabinClass,
		},
		Metadata: Metadata{
			TotalResults:       result.Metadata.TotalResults,
			ProvidersQueried:   result.Metadata.ProvidersQueried,
			ProvidersSucceeded: result.Metadata.ProvidersSucceeded,
			ProvidersFailed:    result.Metadata.ProvidersFailed,
			SearchTimeMs:       result.Metadata.SearchTimeMs,
			CacheHit:           result.Metadata.CacheHit,
		},
		Flights: flightDTOs,
	}
}

func toFlightDTO(flight domain.Flight) FlightDTO {
	return FlightDTO{
		ID:       flight.ID,
		Provider: flight.Provider,
		Airline: Airline{
			Name: flight.Airline.Name,
			Code: flight.Airline.Code,
		},
		FlightNumber: flight.FlightNumber,
		Departure: Endpoint{
			Airport:   flight.Departure.Airport,
			City:      flight.Departure.City,
			Datetime:  flight.Departure.Datetime,
			Timestamp: flight.Departure.Timestamp,
		},
		Arrival: Endpoint{
			Airport:   flight.Arrival.Airport,
			City:      flight.Arrival.City,
			Datetime:  flight.Arrival.Datetime,
			Timestamp: flight.Arrival.Timestamp,
		},
		Duration: Duration{
			TotalMinutes: flight.Duration.TotalMinutes,
			Formatted:    flight.Duration.Formatted,
		},
		Stops: flight.Stops,
		Price: Price{
			Amount:    flight.Price.Amount,
			Currency:  flight.Price.Currency,
			Formatted: utils.FormatCurrency(float64(flight.Price.Amount), flight.Price.Currency),
		},
		AvailableSeats: flight.AvailableSeats,
		CabinClass:     flight.CabinClass,
		Aircraft:       flight.Aircraft,
		Amenities:      flight.Amenities,
		Baggage: Baggage{
			CarryOn: flight.Baggage.CarryOn,
			Checked: flight.Baggage.Checked,
		},
	}
}
