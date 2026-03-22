package provider

import (
	"fmt"
	"time"

	"flight-search/internal/domain"
)

// BatikAirProvider implements the Batik Air API
type BatikAirProvider struct {
	BaseProvider
}

func NewBatikAirProvider() *BatikAirProvider {
	return &BatikAirProvider{
		BaseProvider: BaseProvider{
			name:     "Batik Air",
			delayMin: 200 * time.Millisecond,
			delayMax: 400 * time.Millisecond,
		},
	}
}

func (p *BatikAirProvider) SearchFlights(request domain.SearchRequest) domain.ProviderResult {
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

	data, err := p.loadMockData("batik_air_search_response.json")
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
