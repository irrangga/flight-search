package normalizer

import (
	"fmt"
	"math"
	"strings"
	"time"

	"flight-search/internal/domain"
)

// NormalizeAirAsia normalizes AirAsia flight data
func NormalizeAirAsia(flight map[string]interface{}) (domain.Flight, error) {
	flightCode, _ := flight["flight_code"].(string)
	airline, _ := flight["airline"].(string)
	fromAirport, _ := flight["from_airport"].(string)
	toAirport, _ := flight["to_airport"].(string)

	// Parse departure
	depTime, _ := flight["depart_time"].(string)
	depDt, err := time.Parse(time.RFC3339, depTime)
	if err != nil {
		return domain.Flight{}, fmt.Errorf("invalid departure time: %w", err)
	}

	// Parse arrival
	arrTime, _ := flight["arrive_time"].(string)
	arrDt, err := time.Parse(time.RFC3339, arrTime)
	if err != nil {
		return domain.Flight{}, fmt.Errorf("invalid arrival time: %w", err)
	}

	// Duration
	durationHours, _ := flight["duration_hours"].(float64)
	durationMinutes := int(math.Round(durationHours * 60))
	hours := int(durationHours)
	minutes := int(math.Round((durationHours - float64(hours)) * 60))
	var durationFormatted string
	if hours > 0 {
		durationFormatted = fmt.Sprintf("%dh %dm", hours, minutes)
	} else {
		durationFormatted = fmt.Sprintf("%dm", minutes)
	}

	// Stops
	stops := 0
	if directFlight, ok := flight["direct_flight"].(bool); ok && !directFlight {
		if stopsData, ok := flight["stops"].([]interface{}); ok {
			stops = len(stopsData)
		}
	}

	// Price
	var priceIDR float64
	if p, ok := flight["price_idr"]; ok {
		switch v := p.(type) {
		case float64:
			priceIDR = v
		case float32:
			priceIDR = float64(v)
		case int:
			priceIDR = float64(v)
		case int64:
			priceIDR = float64(v)
		case int32:
			priceIDR = float64(v)
		}
	}

	// Seats
	seats, _ := flight["seats"].(float64)

	// Cabin class
	cabinClass, _ := flight["cabin_class"].(string)

	// Baggage
	baggageNote, _ := flight["baggage_note"].(string)
	var carryOn, checked string
	if strings.Contains(baggageNote, "Cabin baggage only") {
		carryOn = "Cabin baggage only"
		checked = "Additional fee"
	}

	depCity := resolveCity("", fromAirport)
	arrCity := resolveCity("", toAirport)

	return domain.Flight{
		ID:           fmt.Sprintf("%s_%s", flightCode, "AirAsia"),
		Provider:     "AirAsia",
		Airline:      domain.Airline{Name: airline, Code: "QZ"},
		FlightNumber: flightCode,
		Departure: domain.FlightEndpoint{
			Airport:   fromAirport,
			City:      depCity,
			Datetime:  depTime,
			Timestamp: depDt.Unix(),
		},
		Arrival: domain.FlightEndpoint{
			Airport:   toAirport,
			City:      arrCity,
			Datetime:  arrTime,
			Timestamp: arrDt.Unix(),
		},
		Duration: domain.Duration{
			TotalMinutes: durationMinutes,
			Formatted:    durationFormatted,
		},
		Stops:          stops,
		Price:          domain.Price{Amount: int(priceIDR), Currency: "IDR"},
		AvailableSeats: int(seats),
		CabinClass:     cabinClass,
		Aircraft:       nil,
		Amenities:      []string{},
		Baggage:        domain.Baggage{CarryOn: carryOn, Checked: checked},
	}, nil
}
