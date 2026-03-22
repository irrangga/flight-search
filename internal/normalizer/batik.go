package normalizer

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"flight-search/internal/domain"
)

// NormalizeBatikAir normalizes Batik Air flight data
func NormalizeBatikAir(flight map[string]interface{}) (domain.Flight, error) {
	flightNumber, _ := flight["flightNumber"].(string)
	airlineName, _ := flight["airlineName"].(string)
	airlineIATA, _ := flight["airlineIATA"].(string)
	origin, _ := flight["origin"].(string)
	destination, _ := flight["destination"].(string)

	// Parse departure
	depTimeStr, _ := flight["departureDateTime"].(string)
	re := regexp.MustCompile(`(\+\d{2})(\d{2})$`)
	depTimeStr = re.ReplaceAllString(depTimeStr, "$1:$2")
	depDt, err := time.Parse(time.RFC3339, depTimeStr)
	if err != nil {
		return domain.Flight{}, fmt.Errorf("invalid departure time: %w", err)
	}

	// Parse arrival
	arrTimeStr, _ := flight["arrivalDateTime"].(string)
	arrTimeStr = re.ReplaceAllString(arrTimeStr, "$1:$2")
	arrDt, err := time.Parse(time.RFC3339, arrTimeStr)
	if err != nil {
		return domain.Flight{}, fmt.Errorf("invalid arrival time: %w", err)
	}

	// Travel time
	travelTime, _ := flight["travelTime"].(string)
	hours := 0
	minutes := 0
	if hMatch := regexp.MustCompile(`(\d+)h`).FindStringSubmatch(travelTime); len(hMatch) > 1 {
		hours, _ = strconv.Atoi(hMatch[1])
	}
	if mMatch := regexp.MustCompile(`(\d+)m`).FindStringSubmatch(travelTime); len(mMatch) > 1 {
		minutes, _ = strconv.Atoi(mMatch[1])
	}
	durationMinutes := hours*60 + minutes

	// Stops
	numberOfStops, _ := flight["numberOfStops"].(float64)

	// Price
	fare, _ := flight["fare"].(map[string]interface{})
	totalPrice, _ := fare["totalPrice"].(float64)
	currencyCode, _ := fare["currencyCode"].(string)

	// Seats
	seatsAvailable, _ := flight["seatsAvailable"].(float64)

	// Aircraft
	var aircraft *string
	if aircraftModel, ok := flight["aircraftModel"].(string); ok && aircraftModel != "" {
		aircraft = &aircraftModel
	}

	// Amenities
	var amenities []string
	if onboardServices, ok := flight["onboardServices"].([]interface{}); ok {
		for _, s := range onboardServices {
			if str, ok := s.(string); ok {
				amenities = append(amenities, str)
			}
		}
	}

	// Baggage
	baggageInfo, _ := flight["baggageInfo"].(string)
	var carryOn, checked string
	if strings.Contains(baggageInfo, "7kg") {
		carryOn = "7kg cabin"
	}
	if strings.Contains(baggageInfo, "20kg") {
		checked = "20kg checked"
	}

	return domain.Flight{
		ID:           fmt.Sprintf("%s_%s", flightNumber, "Batik Air"),
		Provider:     "Batik Air",
		Airline:      domain.Airline{Name: airlineName, Code: airlineIATA},
		FlightNumber: flightNumber,
		Departure: domain.FlightEndpoint{
			Airport:   origin,
			City:      "",
			Datetime:  depTimeStr,
			Timestamp: depDt.Unix(),
		},
		Arrival: domain.FlightEndpoint{
			Airport:   destination,
			City:      "",
			Datetime:  arrTimeStr,
			Timestamp: arrDt.Unix(),
		},
		Duration: domain.Duration{
			TotalMinutes: durationMinutes,
			Formatted:    travelTime,
		},
		Stops:          int(numberOfStops),
		Price:          domain.Price{Amount: int(totalPrice), Currency: currencyCode},
		AvailableSeats: int(seatsAvailable),
		CabinClass:     "economy",
		Aircraft:       aircraft,
		Amenities:      amenities,
		Baggage:        domain.Baggage{CarryOn: carryOn, Checked: checked},
	}, nil
}
