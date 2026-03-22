package normalizer

import (
	"fmt"
	"time"

	"flight-search/internal/domain"
)

// NormalizeGarudaIndonesia normalizes Garuda Indonesia flight data
func NormalizeGarudaIndonesia(flight map[string]interface{}) (domain.Flight, error) {
	flightID, _ := flight["flight_id"].(string)
	airline, _ := flight["airline"].(string)
	airlineCode, _ := flight["airline_code"].(string)

	// Parse departure
	depData, _ := flight["departure"].(map[string]interface{})
	depTime, _ := depData["time"].(string)
	depDt, err := time.Parse(time.RFC3339, depTime)
	if err != nil {
		return domain.Flight{}, fmt.Errorf("invalid departure time: %w", err)
	}
	depAirport, _ := depData["airport"].(string)
	depCity, _ := depData["city"].(string)

	// Parse arrival
	arrData, _ := flight["arrival"].(map[string]interface{})
	arrTime, _ := arrData["time"].(string)
	arrDt, err := time.Parse(time.RFC3339, arrTime)
	if err != nil {
		return domain.Flight{}, fmt.Errorf("invalid arrival time: %w", err)
	}
	arrAirport, _ := arrData["airport"].(string)
	arrCity, _ := arrData["city"].(string)

	// Duration
	durationMinutes, _ := flight["duration_minutes"].(float64)
	hours := int(durationMinutes / 60)
	minutes := int(durationMinutes) % 60
	var durationFormatted string
	if hours > 0 {
		durationFormatted = fmt.Sprintf("%dh %dm", hours, minutes)
	} else {
		durationFormatted = fmt.Sprintf("%dm", minutes)
	}

	// Stops
	stops := 0
	if stopsVal, ok := flight["stops"].(float64); ok {
		stops = int(stopsVal)
	}

	// Price
	priceData, _ := flight["price"].(map[string]interface{})
	priceAmount, _ := priceData["amount"].(float64)
	currency, _ := priceData["currency"].(string)

	// Available seats
	availableSeats, _ := flight["available_seats"].(float64)

	// Cabin class
	cabinClass, _ := flight["fare_class"].(string)

	// Aircraft
	var aircraft *string
	if aircraftVal, ok := flight["aircraft"].(string); ok && aircraftVal != "" {
		aircraft = &aircraftVal
	}

	// Amenities
	var amenities []string
	if amenitiesVal, ok := flight["amenities"].([]interface{}); ok {
		for _, a := range amenitiesVal {
			if str, ok := a.(string); ok {
				amenities = append(amenities, str)
			}
		}
	}

	// Baggage
	var carryOn, checked string
	if baggage, ok := flight["baggage"].(map[string]interface{}); ok {
		if co, ok := baggage["carry_on"].(float64); ok {
			carryOn = fmt.Sprintf("%.0f piece(s)", co)
		} else if coStr, ok := baggage["carry_on"].(string); ok {
			carryOn = coStr
		}

		if ch, ok := baggage["checked"].(float64); ok {
			checked = fmt.Sprintf("%.0f piece(s)", ch)
		} else if chStr, ok := baggage["checked"].(string); ok {
			checked = chStr
		}
	}

	return domain.Flight{
		ID:           fmt.Sprintf("%s_%s", flightID, "Garuda Indonesia"),
		Provider:     "Garuda Indonesia",
		Airline:      domain.Airline{Name: airline, Code: airlineCode},
		FlightNumber: flightID,
		Departure: domain.FlightEndpoint{
			Airport:   depAirport,
			City:      depCity,
			Datetime:  depTime,
			Timestamp: depDt.Unix(),
		},
		Arrival: domain.FlightEndpoint{
			Airport:   arrAirport,
			City:      arrCity,
			Datetime:  arrTime,
			Timestamp: arrDt.Unix(),
		},
		Duration: domain.Duration{
			TotalMinutes: int(durationMinutes),
			Formatted:    durationFormatted,
		},
		Stops:          stops,
		Price:          domain.Price{Amount: int(priceAmount), Currency: currency},
		AvailableSeats: int(availableSeats),
		CabinClass:     cabinClass,
		Aircraft:       aircraft,
		Amenities:      amenities,
		Baggage:        domain.Baggage{CarryOn: carryOn, Checked: checked},
	}, nil
}
