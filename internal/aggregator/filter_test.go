package aggregator

import (
	"flight-search/internal/domain"
	"testing"
)

func TestFilter_BasicSearchCriteria(t *testing.T) {
	agg := &Aggregator{}

	request := domain.SearchRequest{
		Origin:        "CGK",
		Destination:   "DPS",
		DepartureDate: "2025-12-15",
		Passengers:    1,
		CabinClass:    "economy",
	}

	flights := []domain.Flight{
		{
			Departure: domain.FlightEndpoint{
				Airport:  "CGK",
				Datetime: "2025-12-15T06:00:00+07:00",
			},
			Arrival: domain.FlightEndpoint{
				Airport:  "DPS",
				Datetime: "2025-12-15T08:50:00+08:00",
			},
			CabinClass:     "economy",
			AvailableSeats: 10,
		},
		{
			Departure: domain.FlightEndpoint{
				Airport:  "CGK",
				Datetime: "2025-12-15T07:00:00+07:00",
			},
			Arrival: domain.FlightEndpoint{
				Airport:  "SUB",
				Datetime: "2025-12-15T08:30:00+07:00",
			},
			CabinClass:     "economy",
			AvailableSeats: 5,
		},
	}

	filtered := agg.Filter(flights, request)

	if len(filtered) != 1 {
		t.Errorf("Expected 1 flight, got %d", len(filtered))
	}

	if filtered[0].Arrival.Airport != "DPS" {
		t.Errorf("Expected destination DPS, got %s", filtered[0].Arrival.Airport)
	}
}

func TestFilter_PriceRange(t *testing.T) {
	agg := &Aggregator{}

	request := domain.SearchRequest{
		Origin:      "CGK",
		Destination: "DPS",
		Passengers:  1,
		PriceRange:  []int{1000000, 2000000},
	}

	flights := []domain.Flight{
		{
			Departure:      domain.FlightEndpoint{Airport: "CGK"},
			Arrival:        domain.FlightEndpoint{Airport: "DPS"},
			Price:          domain.Price{Amount: 1500000},
			AvailableSeats: 5,
		},
		{
			Departure:      domain.FlightEndpoint{Airport: "CGK"},
			Arrival:        domain.FlightEndpoint{Airport: "DPS"},
			Price:          domain.Price{Amount: 2500000},
			AvailableSeats: 5,
		},
	}

	filtered := agg.Filter(flights, request)

	if len(filtered) != 1 {
		t.Errorf("Expected 1 flight, got %d", len(filtered))
	}

	if filtered[0].Price.Amount != 1500000 {
		t.Errorf("Expected price 1500000, got %d", filtered[0].Price.Amount)
	}
}

func TestFilter_Airlines(t *testing.T) {
	agg := &Aggregator{}

	request := domain.SearchRequest{
		Origin:      "CGK",
		Destination: "DPS",
		Passengers:  1,
		Airlines:    []string{"Garuda Indonesia", "AirAsia"},
	}

	flights := []domain.Flight{
		{
			Departure:      domain.FlightEndpoint{Airport: "CGK"},
			Arrival:        domain.FlightEndpoint{Airport: "DPS"},
			Airline:        domain.Airline{Name: "Garuda Indonesia"},
			AvailableSeats: 5,
		},
		{
			Departure:      domain.FlightEndpoint{Airport: "CGK"},
			Arrival:        domain.FlightEndpoint{Airport: "DPS"},
			Airline:        domain.Airline{Name: "Lion Air"},
			AvailableSeats: 5,
		},
	}

	filtered := agg.Filter(flights, request)

	if len(filtered) != 1 {
		t.Errorf("Expected 1 flight, got %d", len(filtered))
	}

	if filtered[0].Airline.Name != "Garuda Indonesia" {
		t.Errorf("Expected airline Garuda Indonesia, got %s", filtered[0].Airline.Name)
	}
}

func TestFilter_NumberOfStops(t *testing.T) {
	agg := &Aggregator{}

	request := domain.SearchRequest{
		Origin:        "CGK",
		Destination:   "DPS",
		Passengers:    1,
		NumberOfStops: []int{1, 2},
	}

	flights := []domain.Flight{
		{
			Departure:      domain.FlightEndpoint{Airport: "CGK"},
			Arrival:        domain.FlightEndpoint{Airport: "DPS"},
			Stops:          1,
			AvailableSeats: 5,
		},
		{
			Departure:      domain.FlightEndpoint{Airport: "CGK"},
			Arrival:        domain.FlightEndpoint{Airport: "DPS"},
			Stops:          2,
			AvailableSeats: 5,
		},
		{
			Departure:      domain.FlightEndpoint{Airport: "CGK"},
			Arrival:        domain.FlightEndpoint{Airport: "DPS"},
			Stops:          0,
			AvailableSeats: 5,
		},
		{
			Departure:      domain.FlightEndpoint{Airport: "CGK"},
			Arrival:        domain.FlightEndpoint{Airport: "DPS"},
			Stops:          3,
			AvailableSeats: 5,
		},
	}

	filtered := agg.Filter(flights, request)

	if len(filtered) != 2 {
		t.Errorf("Expected 2 flights, got %d", len(filtered))
	}

	// Should include flights with 1 and 2 stops
	expectedStops := map[int]bool{1: false, 2: false}
	for _, flight := range filtered {
		if flight.Stops == 1 {
			expectedStops[1] = true
		}
		if flight.Stops == 2 {
			expectedStops[2] = true
		}
	}

	if !expectedStops[1] || !expectedStops[2] {
		t.Error("Expected flights with 1 and 2 stops to be included")
	}
}

func TestFilter_DepartureTimeRange(t *testing.T) {
	agg := &Aggregator{}

	request := domain.SearchRequest{
		Origin:             "CGK",
		Destination:        "DPS",
		Passengers:         1,
		DepartureTimeRange: []string{"06:00", "07:00"},
	}

	flights := []domain.Flight{
		{
			Departure: domain.FlightEndpoint{
				Airport:  "CGK",
				Datetime: "2025-12-15T06:30:00+07:00",
			},
			Arrival:        domain.FlightEndpoint{Airport: "DPS"},
			AvailableSeats: 5,
		},
		{
			Departure: domain.FlightEndpoint{
				Airport:  "CGK",
				Datetime: "2025-12-15T08:00:00+07:00",
			},
			Arrival:        domain.FlightEndpoint{Airport: "DPS"},
			AvailableSeats: 5,
		},
	}

	filtered := agg.Filter(flights, request)

	if len(filtered) != 1 {
		t.Errorf("Expected 1 flight, got %d", len(filtered))
	}

	if filtered[0].Departure.Datetime[11:16] != "06:30" {
		t.Errorf("Expected departure time 06:30, got %s", filtered[0].Departure.Datetime[11:16])
	}
}

func TestFilter_DurationRange(t *testing.T) {
	agg := &Aggregator{}

	request := domain.SearchRequest{
		Origin:        "CGK",
		Destination:   "DPS",
		Passengers:    1,
		DurationRange: []int{90, 150},
	}

	flights := []domain.Flight{
		{
			Departure:      domain.FlightEndpoint{Airport: "CGK"},
			Arrival:        domain.FlightEndpoint{Airport: "DPS"},
			Duration:       domain.Duration{TotalMinutes: 120},
			AvailableSeats: 5,
		},
		{
			Departure:      domain.FlightEndpoint{Airport: "CGK"},
			Arrival:        domain.FlightEndpoint{Airport: "DPS"},
			Duration:       domain.Duration{TotalMinutes: 200},
			AvailableSeats: 5,
		},
	}

	filtered := agg.Filter(flights, request)

	if len(filtered) != 1 {
		t.Errorf("Expected 1 flight, got %d", len(filtered))
	}

	if filtered[0].Duration.TotalMinutes != 120 {
		t.Errorf("Expected duration 120, got %d", filtered[0].Duration.TotalMinutes)
	}
}

func TestFilter_ArrivalTimeRange(t *testing.T) {
	agg := &Aggregator{}

	request := domain.SearchRequest{
		Origin:           "CGK",
		Destination:      "DPS",
		Passengers:       1,
		ArrivalTimeRange: []string{"08:00", "09:00"},
	}

	flights := []domain.Flight{
		{
			Departure: domain.FlightEndpoint{Airport: "CGK"},
			Arrival: domain.FlightEndpoint{
				Airport:  "DPS",
				Datetime: "2025-12-15T08:30:00+08:00",
			},
			AvailableSeats: 5,
		},
		{
			Departure: domain.FlightEndpoint{Airport: "CGK"},
			Arrival: domain.FlightEndpoint{
				Airport:  "DPS",
				Datetime: "2025-12-15T10:00:00+08:00",
			},
			AvailableSeats: 5,
		},
	}

	filtered := agg.Filter(flights, request)

	if len(filtered) != 1 {
		t.Errorf("Expected 1 flight, got %d", len(filtered))
	}

	if filtered[0].Arrival.Datetime[11:16] != "08:30" {
		t.Errorf("Expected arrival time 08:30, got %s", filtered[0].Arrival.Datetime[11:16])
	}
}

func TestFilter_CabinClass(t *testing.T) {
	agg := &Aggregator{}

	request := domain.SearchRequest{
		Origin:      "CGK",
		Destination: "DPS",
		Passengers:  1,
		CabinClass:  "economy",
	}

	flights := []domain.Flight{
		{
			Departure:      domain.FlightEndpoint{Airport: "CGK"},
			Arrival:        domain.FlightEndpoint{Airport: "DPS"},
			CabinClass:     "economy",
			AvailableSeats: 5,
		},
		{
			Departure:      domain.FlightEndpoint{Airport: "CGK"},
			Arrival:        domain.FlightEndpoint{Airport: "DPS"},
			CabinClass:     "business",
			AvailableSeats: 5,
		},
	}

	filtered := agg.Filter(flights, request)

	if len(filtered) != 1 {
		t.Errorf("Expected 1 flight, got %d", len(filtered))
	}

	if filtered[0].CabinClass != "economy" {
		t.Errorf("Expected cabin class economy, got %s", filtered[0].CabinClass)
	}
}

func TestFilter_Passengers(t *testing.T) {
	agg := &Aggregator{}

	request := domain.SearchRequest{
		Origin:      "CGK",
		Destination: "DPS",
		Passengers:  5,
	}

	flights := []domain.Flight{
		{
			Departure:      domain.FlightEndpoint{Airport: "CGK"},
			Arrival:        domain.FlightEndpoint{Airport: "DPS"},
			AvailableSeats: 10,
		},
		{
			Departure:      domain.FlightEndpoint{Airport: "CGK"},
			Arrival:        domain.FlightEndpoint{Airport: "DPS"},
			AvailableSeats: 3,
		},
	}

	filtered := agg.Filter(flights, request)

	if len(filtered) != 1 {
		t.Errorf("Expected 1 flight, got %d", len(filtered))
	}

	if filtered[0].AvailableSeats != 10 {
		t.Errorf("Expected available seats 10, got %d", filtered[0].AvailableSeats)
	}
}
