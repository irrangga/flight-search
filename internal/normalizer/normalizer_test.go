package normalizer

import (
	"testing"
)

func TestResolveCityFromAirportCode(t *testing.T) {
	city := resolveCity("", "DPS")
	if city != "Denpasar" {
		t.Fatalf("expected city Denpasar, got %q", city)
	}

	city = resolveCity("", "CGK")
	if city != "Jakarta" {
		t.Fatalf("expected city Jakarta, got %q", city)
	}

	city = resolveCity("SomeCity", "DPS")
	if city != "SomeCity" {
		t.Fatalf("expected city SomeCity (override), got %q", city)
	}
}

func TestNormalizeAirAsia(t *testing.T) {
	flight := map[string]interface{}{
		"flight_code":    "QZ520",
		"airline":        "AirAsia",
		"from_airport":   "CGK",
		"to_airport":     "DPS",
		"depart_time":    "2025-12-15T04:45:00+07:00",
		"arrive_time":    "2025-12-15T07:25:00+08:00",
		"duration_hours": 1.67,
		"direct_flight":  true,
		"price_idr":      650000,
		"seats":          67,
		"cabin_class":    "economy",
		"baggage_note":   "Cabin baggage only, checked bags additional fee",
	}

	f, err := NormalizeAirAsia(flight)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if f.Departure.City != "Jakarta" {
		t.Errorf("expected departure city Jakarta, got %q", f.Departure.City)
	}
	if f.Arrival.City != "Denpasar" {
		t.Errorf("expected arrival city Denpasar, got %q", f.Arrival.City)
	}
	if f.Price.Amount != 650000 {
		t.Errorf("expected price 650000, got %d", f.Price.Amount)
	}
	if f.Price.Currency != "IDR" {
		t.Errorf("expected currency IDR, got %q", f.Price.Currency)
	}
	if f.Provider != "AirAsia" {
		t.Errorf("expected provider AirAsia, got %q", f.Provider)
	}
	if f.CabinClass != "economy" {
		t.Errorf("expected cabin class economy, got %q", f.CabinClass)
	}
	if f.Airline.Name != "AirAsia" {
		t.Errorf("expected airline AirAsia, got %q", f.Airline.Name)
	}
}

func TestNormalizeBatikAir(t *testing.T) {
	flight := map[string]interface{}{
		"flightNumber":      "ID6514",
		"airlineName":       "Batik Air",
		"airlineIATA":       "ID",
		"origin":            "CGK",
		"destination":       "DPS",
		"departureDateTime": "2025-12-15T07:15:00+0700",
		"arrivalDateTime":   "2025-12-15T10:00:00+0800",
		"travelTime":        "1h 45m",
		"numberOfStops":     0,
		"fare": map[string]interface{}{
			"basePrice":    980000,
			"taxes":        120000,
			"totalPrice":   1100000,
			"currencyCode": "IDR",
			"class":        "Y",
		},
		"seatsAvailable": 32,
		"aircraftModel":  "Airbus A320",
		"baggageInfo":    "7kg cabin, 20kg checked",
		"onboardServices": []any{
			"Snack",
			"Beverage",
		},
	}

	f, err := NormalizeBatikAir(flight)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if f.Departure.City != "Jakarta" {
		t.Errorf("expected departure city Jakarta, got %q", f.Departure.City)
	}
	if f.Arrival.City != "Denpasar" {
		t.Errorf("expected arrival city Denpasar, got %q", f.Arrival.City)
	}
	if f.Price.Amount != 1100000 {
		t.Errorf("expected price 1100000, got %d", f.Price.Amount)
	}
	if f.Price.Currency != "IDR" {
		t.Errorf("expected currency IDR, got %q", f.Price.Currency)
	}
	if f.Provider != "Batik Air" {
		t.Errorf("expected provider Batik Air, got %q", f.Provider)
	}
	if f.CabinClass != "economy" {
		t.Errorf("expected cabin class economy, got %q", f.CabinClass)
	}
	if f.Airline.Name != "Batik Air" {
		t.Errorf("expected airline Batik Air, got %q", f.Airline.Name)
	}
	if len(f.Amenities) != 2 {
		t.Errorf("expected amenities length is 2, got %v", f.Amenities)
	}
	for i, amenity := range f.Amenities {
		if i == 0 && amenity != "snack" {
			t.Errorf("expected snack amenity")
		}
		if i == 1 && amenity != "beverage" {
			t.Errorf("expected beverage amenity")
		}
	}
}

func TestNormalizeGarudaIndonesia(t *testing.T) {
	flight := map[string]interface{}{
		"flight_id":    "GA400",
		"airline":      "Garuda Indonesia",
		"airline_code": "GA",
		"departure": map[string]interface{}{
			"airport":  "CGK",
			"city":     "Jakarta",
			"time":     "2025-12-15T06:00:00+07:00",
			"terminal": "3",
		},
		"arrival": map[string]interface{}{
			"airport":  "DPS",
			"city":     "Denpasar",
			"time":     "2025-12-15T08:50:00+08:00",
			"terminal": "I",
		},
		"duration_minutes": 110,
		"stops":            0,
		"aircraft":         "Boeing 737-800",
		"price": map[string]interface{}{
			"amount":   1250000,
			"currency": "IDR",
		},
		"available_seats": 28,
		"fare_class":      "economy",
		"baggage": map[string]interface{}{
			"carry_on": 1,
			"checked":  2,
		},
		"amenities": []any{
			"wifi",
			"meal",
			"entertainment",
		},
	}

	f, err := NormalizeGarudaIndonesia(flight)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if f.Departure.City != "Jakarta" {
		t.Errorf("expected departure city Jakarta, got %q", f.Departure.City)
	}
	if f.Arrival.City != "Denpasar" {
		t.Errorf("expected arrival city Denpasar, got %q", f.Arrival.City)
	}
	if f.Price.Amount != 1250000 {
		t.Errorf("expected price 1250000, got %d", f.Price.Amount)
	}
	if f.Price.Currency != "IDR" {
		t.Errorf("expected currency IDR, got %q", f.Price.Currency)
	}
	if f.Provider != "Garuda Indonesia" {
		t.Errorf("expected provider Garuda Indonesia, got %q", f.Provider)
	}
	if f.CabinClass != "economy" {
		t.Errorf("expected cabin class economy, got %q", f.CabinClass)
	}
	if f.Airline.Name != "Garuda Indonesia" {
		t.Errorf("expected airline Garuda Indonesia, got %q", f.Airline.Name)
	}
	if len(f.Amenities) != 3 {
		t.Errorf("expected amenities length is 3, got %v", f.Amenities)
	}
	for i, amenity := range f.Amenities {
		if i == 0 && amenity != "wifi" {
			t.Errorf("expected wifi amenity")
		}
		if i == 1 && amenity != "meal" {
			t.Errorf("expected meal amenity")
		}
		if i == 2 && amenity != "entertainment" {
			t.Errorf("expected entertainment amenity")
		}
	}
}

func TestNormalizeLionAir(t *testing.T) {
	flight := map[string]interface{}{
		"id": "JT740",
		"carrier": map[string]interface{}{
			"name": "Lion Air",
			"iata": "JT",
		},
		"route": map[string]interface{}{
			"from": map[string]interface{}{
				"code": "CGK",
				"name": "Soekarno-Hatta International",
				"city": "Jakarta",
			},
			"to": map[string]interface{}{
				"code": "DPS",
				"name": "Ngurah Rai International",
				"city": "Denpasar",
			},
		},
		"schedule": map[string]interface{}{
			"departure":          "2025-12-15T05:30:00",
			"departure_timezone": "Asia/Jakarta",
			"arrival":            "2025-12-15T08:15:00",
			"arrival_timezone":   "Asia/Makassar",
		},
		"flight_time": 105,
		"is_direct":   true,
		"pricing": map[string]interface{}{
			"total":     950000,
			"currency":  "IDR",
			"fare_type": "ECONOMY",
		},
		"seats_left": 45,
		"plane_type": "Boeing 737-900ER",
		"services": map[string]interface{}{
			"wifi_available": true,
			"meals_included": true,
			"baggage_allowance": map[string]interface{}{
				"cabin": "7 kg",
				"hold":  "20 kg",
			},
		},
	}

	f, err := NormalizeLionAir(flight)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if f.Departure.City != "Jakarta" {
		t.Errorf("expected departure city Jakarta, got %q", f.Departure.City)
	}
	if f.Arrival.City != "Denpasar" {
		t.Errorf("expected arrival city Denpasar, got %q", f.Arrival.City)
	}
	if f.Price.Amount != 950000 {
		t.Errorf("expected price 950000, got %d", f.Price.Amount)
	}
	if f.Price.Currency != "IDR" {
		t.Errorf("expected currency IDR, got %q", f.Price.Currency)
	}
	if f.Provider != "Lion Air" {
		t.Errorf("expected provider Lion Air, got %q", f.Provider)
	}
	if f.CabinClass != "economy" {
		t.Errorf("expected cabin class economy, got %q", f.CabinClass)
	}
	if f.Airline.Name != "Lion Air" {
		t.Errorf("expected airline Lion Air, got %q", f.Airline.Name)
	}
	if len(f.Amenities) != 2 {
		t.Errorf("expected amenities length is 2, got %v", f.Amenities)
	}
	for i, amenity := range f.Amenities {
		if i == 0 && amenity != "wifi" {
			t.Errorf("expected wifi amenity")
		}
		if i == 1 && amenity != "meal" {
			t.Errorf("expected meal amenity")
		}
	}
}
