package provider

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"flight-search/internal/domain"
)

// Provider defines the interface for flight providers
type Provider interface {
	Name() string
	SearchFlights(request domain.SearchRequest) domain.ProviderResult
}

// BaseProvider provides common functionality for all providers
type BaseProvider struct {
	name        string
	delayMin    time.Duration
	delayMax    time.Duration
	failureRate float64
}

func (p *BaseProvider) Name() string {
	return p.name
}

func (p *BaseProvider) simulateDelay() {
	if p.delayMax > 0 {
		delay := p.delayMin + time.Duration(rand.Int63n(int64(p.delayMax-p.delayMin)))
		time.Sleep(delay)
	}
}

func (p *BaseProvider) simulateFailure() bool {
	return rand.Float64() < p.failureRate
}

func (p *BaseProvider) loadMockData(filename string) (map[string]interface{}, error) {
	// Get the directory of this file (provider.go)
	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		return nil, fmt.Errorf("failed to get current file path")
	}
	dir := filepath.Dir(currentFile)

	// Build path relative to internal/mock-api
	path := filepath.Join(dir, "../mock-api", filename)

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read mock data file %s: %w", path, err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to parse JSON from %s: %w", path, err)
	}

	return result, nil
}
