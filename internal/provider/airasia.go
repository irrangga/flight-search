package provider

import (
	"fmt"
	"time"

	"flight-search/internal/domain"
)

// AirAsiaProvider implements the AirAsia API
type AirAsiaProvider struct {
	BaseProvider
}

func NewAirAsiaProvider() *AirAsiaProvider {
	return &AirAsiaProvider{
		BaseProvider: BaseProvider{
			name:        "AirAsia",
			delayMin:    50 * time.Millisecond,
			delayMax:    150 * time.Millisecond,
			failureRate: 0.1,
		},
	}
}

func (p *AirAsiaProvider) SearchFlights(request domain.SearchRequest) domain.ProviderResult {
	start := time.Now()
	p.simulateDelay()

	if p.simulateFailure() {
		return domain.ProviderResult{
			Provider:     p.Name(),
			Success:      false,
			ResponseTime: time.Since(start),
			Error:        fmt.Errorf("API timeout"),
		}
	}

	data, err := p.loadMockData("airasia_search_response.json")
	if err != nil {
		return domain.ProviderResult{
			Provider:     p.Name(),
			Success:      false,
			ResponseTime: time.Since(start),
			Error:        err,
		}
	}

	return domain.ProviderResult{
		Provider:     p.Name(),
		Success:      true,
		Data:         data,
		ResponseTime: time.Since(start),
	}
}
