package normalizer

import (
	"flight-search/internal/constant"
	"flight-search/internal/domain"
)

// Normalizer normalizes raw provider data to unified domain model
type Normalizer interface {
	Normalize(flight map[string]interface{}) (domain.Flight, error)
}

func resolveCity(city, airportCode string) string {
	if city != "" {
		return city
	}
	if airportCode == "" {
		return ""
	}
	if lookup, ok := constant.GetAirportCity(airportCode); ok {
		return lookup
	}
	return ""
}

// NormalizeFlight routes normalization to the appropriate provider normalizer
func NormalizeFlight(flight map[string]interface{}, provider string) (domain.Flight, error) {
	switch provider {
	case "Garuda Indonesia":
		return NormalizeGarudaIndonesia(flight)
	case "Lion Air":
		return NormalizeLionAir(flight)
	case "Batik Air":
		return NormalizeBatikAir(flight)
	case "AirAsia":
		return NormalizeAirAsia(flight)
	default:
		return domain.Flight{}, nil
	}
}
