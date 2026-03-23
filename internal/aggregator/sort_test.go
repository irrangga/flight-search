package aggregator

import (
	"flight-search/internal/domain"
	"testing"
)

func TestSort_PriceAsc(t *testing.T) {
	agg := &Aggregator{}

	flights := []domain.Flight{
		{Price: domain.Price{Amount: 2000000}},
		{Price: domain.Price{Amount: 1000000}},
		{Price: domain.Price{Amount: 1500000}},
	}

	agg.Sort(flights, "price_asc")

	if flights[0].Price.Amount != 1000000 {
		t.Errorf("Expected first flight price 1000000, got %d", flights[0].Price.Amount)
	}
	if flights[1].Price.Amount != 1500000 {
		t.Errorf("Expected second flight price 1500000, got %d", flights[1].Price.Amount)
	}
	if flights[2].Price.Amount != 2000000 {
		t.Errorf("Expected third flight price 2000000, got %d", flights[2].Price.Amount)
	}
}

func TestSort_PriceDesc(t *testing.T) {
	agg := &Aggregator{}

	flights := []domain.Flight{
		{Price: domain.Price{Amount: 1000000}},
		{Price: domain.Price{Amount: 2000000}},
		{Price: domain.Price{Amount: 1500000}},
	}

	agg.Sort(flights, "price_desc")

	if flights[0].Price.Amount != 2000000 {
		t.Errorf("Expected first flight price 2000000, got %d", flights[0].Price.Amount)
	}
	if flights[1].Price.Amount != 1500000 {
		t.Errorf("Expected second flight price 1500000, got %d", flights[1].Price.Amount)
	}
	if flights[2].Price.Amount != 1000000 {
		t.Errorf("Expected third flight price 1000000, got %d", flights[2].Price.Amount)
	}
}

func TestSort_DurationAsc(t *testing.T) {
	agg := &Aggregator{}

	flights := []domain.Flight{
		{Duration: domain.Duration{TotalMinutes: 150}},
		{Duration: domain.Duration{TotalMinutes: 90}},
		{Duration: domain.Duration{TotalMinutes: 120}},
	}

	agg.Sort(flights, "duration_asc")

	if flights[0].Duration.TotalMinutes != 90 {
		t.Errorf("Expected first flight duration 90, got %d", flights[0].Duration.TotalMinutes)
	}
	if flights[1].Duration.TotalMinutes != 120 {
		t.Errorf("Expected second flight duration 120, got %d", flights[1].Duration.TotalMinutes)
	}
	if flights[2].Duration.TotalMinutes != 150 {
		t.Errorf("Expected third flight duration 150, got %d", flights[2].Duration.TotalMinutes)
	}
}

func TestSort_DurationDesc(t *testing.T) {
	agg := &Aggregator{}

	flights := []domain.Flight{
		{Duration: domain.Duration{TotalMinutes: 90}},
		{Duration: domain.Duration{TotalMinutes: 150}},
		{Duration: domain.Duration{TotalMinutes: 120}},
	}

	agg.Sort(flights, "duration_desc")

	if flights[0].Duration.TotalMinutes != 150 {
		t.Errorf("Expected first flight duration 150, got %d", flights[0].Duration.TotalMinutes)
	}
	if flights[1].Duration.TotalMinutes != 120 {
		t.Errorf("Expected second flight duration 120, got %d", flights[1].Duration.TotalMinutes)
	}
	if flights[2].Duration.TotalMinutes != 90 {
		t.Errorf("Expected third flight duration 90, got %d", flights[2].Duration.TotalMinutes)
	}
}

func TestSort_DefaultByScore(t *testing.T) {
	agg := &Aggregator{}

	flights := []domain.Flight{
		{Score: 200.0},
		{Score: 100.0},
		{Score: 150.0},
	}

	agg.Sort(flights, "")

	if flights[0].Score != 100.0 {
		t.Errorf("Expected first flight score 100.0, got %f", flights[0].Score)
	}
	if flights[1].Score != 150.0 {
		t.Errorf("Expected second flight score 150.0, got %f", flights[1].Score)
	}
	if flights[2].Score != 200.0 {
		t.Errorf("Expected third flight score 200.0, got %f", flights[2].Score)
	}
}

func TestSort_DepartureAsc(t *testing.T) {
	agg := &Aggregator{}

	flights := []domain.Flight{
		{
			Departure: domain.FlightEndpoint{
				Timestamp: 200,
			},
		},
		{
			Departure: domain.FlightEndpoint{
				Timestamp: 100,
			},
		},
		{
			Departure: domain.FlightEndpoint{
				Timestamp: 150,
			},
		},
	}

	agg.Sort(flights, "departure_asc")

	if flights[0].Departure.Timestamp != 100 {
		t.Errorf("Expected first flight departure timestamp 100, got %d", flights[0].Departure.Timestamp)
	}
	if flights[1].Departure.Timestamp != 150 {
		t.Errorf("Expected second flight departure timestamp 150, got %d", flights[1].Departure.Timestamp)
	}
	if flights[2].Departure.Timestamp != 200 {
		t.Errorf("Expected third flight departure timestamp 200, got %d", flights[2].Departure.Timestamp)
	}
}

func TestSort_ArrivalDesc(t *testing.T) {
	agg := &Aggregator{}

	flights := []domain.Flight{
		{
			Arrival: domain.FlightEndpoint{
				Timestamp: 100,
			},
		},
		{
			Arrival: domain.FlightEndpoint{
				Timestamp: 200,
			},
		},
		{
			Arrival: domain.FlightEndpoint{
				Timestamp: 150,
			},
		},
	}

	agg.Sort(flights, "arrival_desc")

	if flights[0].Arrival.Timestamp != 200 {
		t.Errorf("Expected first flight arrival timestamp 200, got %d", flights[0].Arrival.Timestamp)
	}
	if flights[1].Arrival.Timestamp != 150 {
		t.Errorf("Expected second flight arrival timestamp 150, got %d", flights[1].Arrival.Timestamp)
	}
	if flights[2].Arrival.Timestamp != 100 {
		t.Errorf("Expected third flight arrival timestamp 100, got %d", flights[2].Arrival.Timestamp)
	}
}

func TestSort_DepartureDesc(t *testing.T) {
	agg := &Aggregator{}

	flights := []domain.Flight{
		{
			Departure: domain.FlightEndpoint{
				Timestamp: 100,
			},
		},
		{
			Departure: domain.FlightEndpoint{
				Timestamp: 200,
			},
		},
		{
			Departure: domain.FlightEndpoint{
				Timestamp: 150,
			},
		},
	}

	agg.Sort(flights, "departure_desc")

	if flights[0].Departure.Timestamp != 200 {
		t.Errorf("Expected first flight departure timestamp 200, got %d", flights[0].Departure.Timestamp)
	}
	if flights[1].Departure.Timestamp != 150 {
		t.Errorf("Expected second flight departure timestamp 150, got %d", flights[1].Departure.Timestamp)
	}
	if flights[2].Departure.Timestamp != 100 {
		t.Errorf("Expected third flight departure timestamp 100, got %d", flights[2].Departure.Timestamp)
	}
}

func TestSort_ArrivalAsc(t *testing.T) {
	agg := &Aggregator{}

	flights := []domain.Flight{
		{
			Arrival: domain.FlightEndpoint{
				Timestamp: 200,
			},
		},
		{
			Arrival: domain.FlightEndpoint{
				Timestamp: 100,
			},
		},
		{
			Arrival: domain.FlightEndpoint{
				Timestamp: 150,
			},
		},
	}

	agg.Sort(flights, "arrival_asc")

	if flights[0].Arrival.Timestamp != 100 {
		t.Errorf("Expected first flight arrival timestamp 100, got %d", flights[0].Arrival.Timestamp)
	}
	if flights[1].Arrival.Timestamp != 150 {
		t.Errorf("Expected second flight arrival timestamp 150, got %d", flights[1].Arrival.Timestamp)
	}
	if flights[2].Arrival.Timestamp != 200 {
		t.Errorf("Expected third flight arrival timestamp 200, got %d", flights[2].Arrival.Timestamp)
	}
}
