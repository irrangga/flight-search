package aggregator

import (
	"flight-search/internal/domain"
	"sort"
)

// Sort sorts the flights based on the sortBy parameter
func (a *Aggregator) Sort(flights []domain.Flight, sortBy string) {
	switch sortBy {
	case "price_asc":
		sort.Slice(flights, func(i, j int) bool {
			return flights[i].Price.Amount < flights[j].Price.Amount
		})
	case "price_desc":
		sort.Slice(flights, func(i, j int) bool {
			return flights[i].Price.Amount > flights[j].Price.Amount
		})
	case "duration_asc":
		sort.Slice(flights, func(i, j int) bool {
			return flights[i].Duration.TotalMinutes < flights[j].Duration.TotalMinutes
		})
	case "duration_desc":
		sort.Slice(flights, func(i, j int) bool {
			return flights[i].Duration.TotalMinutes > flights[j].Duration.TotalMinutes
		})
	case "departure_asc":
		sort.Slice(flights, func(i, j int) bool {
			return flights[i].Departure.Timestamp < flights[j].Departure.Timestamp
		})
	case "departure_desc":
		sort.Slice(flights, func(i, j int) bool {
			return flights[i].Departure.Timestamp > flights[j].Departure.Timestamp
		})
	case "arrival_asc":
		sort.Slice(flights, func(i, j int) bool {
			return flights[i].Arrival.Timestamp < flights[j].Arrival.Timestamp
		})
	case "arrival_desc":
		sort.Slice(flights, func(i, j int) bool {
			return flights[i].Arrival.Timestamp > flights[j].Arrival.Timestamp
		})
	default:
		// Default: based on score, the lower the score, the better
		sort.Slice(flights, func(i, j int) bool {
			return flights[i].Score < flights[j].Score
		})
	}
}
