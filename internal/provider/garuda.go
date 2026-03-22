package provider

import (
	"fmt"
	"time"

	"flight-search/internal/domain"
)

// GarudaIndonesiaProvider implements the Garuda Indonesia API
type GarudaIndonesiaProvider struct {
	BaseProvider
}

func NewGarudaIndonesiaProvider() *GarudaIndonesiaProvider {
	return &GarudaIndonesiaProvider{
		BaseProvider: BaseProvider{
			name:     "Garuda Indonesia",
			delayMin: 50 * time.Millisecond,
			delayMax: 100 * time.Millisecond,
		},
	}
}

func (p *GarudaIndonesiaProvider) SearchFlights(request domain.SearchRequest) domain.ProviderResult {
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

	data, err := p.loadMockData("garuda_indonesia_search_response.json")
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
