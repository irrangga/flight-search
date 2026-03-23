package provider

import (
	"fmt"
	"time"

	"flight-search/internal/domain"

	"golang.org/x/time/rate"
)

// BatikAirProvider implements the Batik Air API
type BatikAirProvider struct {
	BaseProvider
}

func NewBatikAirProvider() *BatikAirProvider {
	return &BatikAirProvider{
		BaseProvider: BaseProvider{
			name:             "Batik Air",
			delayMin:         200 * time.Millisecond,
			delayMax:         400 * time.Millisecond,
			rateLimiter:      rate.NewLimiter(0.5, 1), // 0.5 requests per second, burst of 1
			disableRateLimit: false,
		},
	}
}

func NewBatikAirProviderForTest() *BatikAirProvider {
	return &BatikAirProvider{
		BaseProvider: BaseProvider{
			name:             "Batik Air",
			delayMin:         200 * time.Millisecond,
			delayMax:         400 * time.Millisecond,
			rateLimiter:      rate.NewLimiter(0.5, 1), // 0.5 requests per second, burst of 1
			disableRateLimit: true,
		},
	}
}

func (p *BatikAirProvider) SearchFlights(request domain.SearchRequest) domain.ProviderResult {
	start := time.Now()

	// Wait for rate limit allowance
	if err := p.waitForRateLimit(); err != nil {
		return domain.ProviderResult{
			Provider:     p.Name(),
			Success:      false,
			ResponseTime: time.Since(start),
			Error:        err,
		}
	}

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
