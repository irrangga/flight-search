package normalizer

import (
	"flight-search/internal/constant"
	"flight-search/internal/domain"
	"fmt"
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

// normalizeFlight routes normalization to the appropriate provider normalizer
func normalizeFlight(flight map[string]interface{}, provider string) (domain.Flight, error) {
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

// NormalizeFlightsFromResults normalizes flights from all successful provider results
func NormalizeFlightsFromResults(successfulResults []domain.ProviderResult) []domain.Flight {
	var normalizedFlights []domain.Flight

	for _, result := range successfulResults {
		flights := extractFlightsFromResponse(result.Data, result.Provider)
		for _, flight := range flights {
			normalized, err := normalizeFlight(flight, result.Provider)
			if err != nil {
				fmt.Printf("Error normalizing flight from %s: %v\n", result.Provider, err)
				continue
			}
			// Calculate score: price + duration (in minutes) * 1250 to give more weight to price (based on 4 hours delay get Rp300.000,-)
			normalized.Score = float64(normalized.Price.Amount) + float64(normalized.Duration.TotalMinutes*1250)

			normalizedFlights = append(normalizedFlights, normalized)
		}
	}

	return normalizedFlights
}

// extractFlightsFromResponse extracts flight data from provider response
func extractFlightsFromResponse(data map[string]interface{}, provider string) []map[string]interface{} {
	switch provider {
	case "Garuda Indonesia":
		if flights, ok := data["flights"].([]interface{}); ok {
			result := make([]map[string]interface{}, len(flights))
			for i, f := range flights {
				if flight, ok := f.(map[string]interface{}); ok {
					result[i] = flight
				}
			}
			return result
		}
	case "Lion Air":
		if dataSection, ok := data["data"].(map[string]interface{}); ok {
			if flights, ok := dataSection["available_flights"].([]interface{}); ok {
				result := make([]map[string]interface{}, len(flights))
				for i, f := range flights {
					if flight, ok := f.(map[string]interface{}); ok {
						result[i] = flight
					}
				}
				return result
			}
		}
	case "Batik Air":
		if results, ok := data["results"].([]interface{}); ok {
			result := make([]map[string]interface{}, len(results))
			for i, r := range results {
				if flight, ok := r.(map[string]interface{}); ok {
					result[i] = flight
				}
			}
			return result
		}
	case "AirAsia":
		if flights, ok := data["flights"].([]interface{}); ok {
			result := make([]map[string]interface{}, len(flights))
			for i, f := range flights {
				if flight, ok := f.(map[string]interface{}); ok {
					result[i] = flight
				}
			}
			return result
		}
	}
	return []map[string]interface{}{}
}
