package normalizer

import (
	"fmt"
	"strings"
	"time"

	"flight-search/internal/domain"
)

// NormalizeLionAir normalizes Lion Air flight data
func NormalizeLionAir(flight map[string]interface{}) (domain.Flight, error) {
	id, _ := flight["id"].(string)

	// Carrier
	carrier, _ := flight["carrier"].(map[string]interface{})
	carrierName, _ := carrier["name"].(string)
	carrierCode, _ := carrier["iata"].(string)

	// Route
	route, _ := flight["route"].(map[string]interface{})
	from, _ := route["from"].(map[string]interface{})
	to, _ := route["to"].(map[string]interface{})
	fromCode, _ := from["code"].(string)
	fromCity, _ := from["city"].(string)
	toCode, _ := to["code"].(string)
	toCity, _ := to["city"].(string)

	// Schedule
	schedule, _ := flight["schedule"].(map[string]interface{})
	depTime, _ := schedule["departure"].(string)
	arrTime, _ := schedule["arrival"].(string)

	// Parse times
	depDt, err := time.Parse("2006-01-02T15:04:05", depTime)
	if err != nil {
		return domain.Flight{}, fmt.Errorf("invalid departure time: %w", err)
	}
	arrDt, err := time.Parse("2006-01-02T15:04:05", arrTime)
	if err != nil {
		return domain.Flight{}, fmt.Errorf("invalid arrival time: %w", err)
	}

	// Flight time
	flightTime, _ := flight["flight_time"].(float64)
	hours := int(flightTime / 60)
	minutes := int(flightTime) % 60
	var durationFormatted string
	if hours > 0 {
		durationFormatted = fmt.Sprintf("%dh %dm", hours, minutes)
	} else {
		durationFormatted = fmt.Sprintf("%dm", minutes)
	}

	// Stops
	stops := 0
	if isDirect, ok := flight["is_direct"].(bool); ok && !isDirect {
		stops = 1
	}

	// Price
	pricing, _ := flight["pricing"].(map[string]interface{})
	var total float64
	if t, ok := pricing["total"]; ok {
		switch v := t.(type) {
		case float64:
			total = v
		case float32:
			total = float64(v)
		case int:
			total = float64(v)
		case int64:
			total = float64(v)
		case int32:
			total = float64(v)
		}
	}
	currency, _ := pricing["currency"].(string)

	// Seats
	seatsLeft, _ := flight["seats_left"].(float64)

	// Cabin class
	fareType, _ := pricing["fare_type"].(string)
	cabinClass := strings.ToLower(fareType)

	// Aircraft
	var aircraft *string
	if planeType, ok := flight["plane_type"].(string); ok && planeType != "" {
		aircraft = &planeType
	}

	// Amenities
	var amenities []string
	if services, ok := flight["services"].(map[string]interface{}); ok {
		if wifi, ok := services["wifi_available"].(bool); ok && wifi {
			amenities = append(amenities, "wifi")
		}
		if meals, ok := services["meals_included"].(bool); ok && meals {
			amenities = append(amenities, "meal")
		}
	}

	// Baggage
	var carryOn, checked string
	if services, ok := flight["services"].(map[string]interface{}); ok {
		if baggage, ok := services["baggage_allowance"].(map[string]interface{}); ok {
			if co, ok := baggage["cabin"].(string); ok {
				carryOn = co
			}
			if ch, ok := baggage["hold"].(string); ok {
				checked = ch
			}
		}
	}

	depCity := resolveCity(fromCity, fromCode)
	arrCity := resolveCity(toCity, toCode)

	return domain.Flight{
		ID:           fmt.Sprintf("%s_%s", id, "Lion Air"),
		Provider:     "Lion Air",
		Airline:      domain.Airline{Name: carrierName, Code: carrierCode},
		FlightNumber: id,
		Departure: domain.FlightEndpoint{
			Airport:   fromCode,
			City:      depCity,
			Datetime:  depTime + "+07:00",
			Timestamp: depDt.Unix(),
		},
		Arrival: domain.FlightEndpoint{
			Airport:   toCode,
			City:      arrCity,
			Datetime:  arrTime + "+08:00",
			Timestamp: arrDt.Unix(),
		},
		Duration: domain.Duration{
			TotalMinutes: int(flightTime),
			Formatted:    durationFormatted,
		},
		Stops:          stops,
		Price:          domain.Price{Amount: int(total), Currency: currency},
		AvailableSeats: int(seatsLeft),
		CabinClass:     cabinClass,
		Aircraft:       aircraft,
		Amenities:      amenities,
		Baggage:        domain.Baggage{CarryOn: carryOn, Checked: checked},
	}, nil
}
