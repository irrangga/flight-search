package aggregator

import (
	"crypto/md5"
	"fmt"
	"time"

	"flight-search/internal/domain"

	"github.com/patrickmn/go-cache"
)

// CacheManager handles caching for flight search results
type CacheManager struct {
	cache *cache.Cache
}

// NewCacheManager creates a new cache manager with default settings
func NewCacheManager() *CacheManager {
	// Default cache: 1 minute expiration, cleanup every 2 minutes
	c := cache.New(1*time.Minute, 2*time.Minute)
	return &CacheManager{
		cache: c,
	}
}

// NewCacheManagerWithTTL creates a new cache manager with custom TTL
func NewCacheManagerWithTTL(defaultExpiration, cleanupInterval time.Duration) *CacheManager {
	c := cache.New(defaultExpiration, cleanupInterval)
	return &CacheManager{
		cache: c,
	}
}

// generateCacheKey creates a unique cache key for the search request
func (cm *CacheManager) generateCacheKey(request domain.SearchRequest) string {
	// Create a string representation of the request for hashing
	returnDate := ""
	if request.ReturnDate != nil {
		returnDate = *request.ReturnDate
	}

	keyData := fmt.Sprintf("%s-%s-%s-%s-%d-%s-%v-%v-%v-%v-%v-%v-%s",
		request.Origin,
		request.Destination,
		request.DepartureDate,
		returnDate,
		request.Passengers,
		request.CabinClass,
		request.PriceRange,
		request.NumberOfStops,
		request.DepartureTimeRange,
		request.ArrivalTimeRange,
		request.Airlines,
		request.DurationRange,
		request.SortBy,
	)

	// Generate MD5 hash of the key data
	hash := md5.Sum([]byte(keyData))
	return fmt.Sprintf("%x", hash)
}

// Get retrieves a cached search result
func (cm *CacheManager) Get(request domain.SearchRequest) (domain.SearchResult, bool) {
	key := cm.generateCacheKey(request)
	if result, found := cm.cache.Get(key); found {
		if searchResult, ok := result.(domain.SearchResult); ok {
			// Mark as cache hit in metadata
			searchResult.Metadata.CacheHit = true
			return searchResult, true
		}
	}
	return domain.SearchResult{}, false
}

// Set stores a search result in cache
func (cm *CacheManager) Set(request domain.SearchRequest, result domain.SearchResult) {
	key := cm.generateCacheKey(request)
	// Mark as cache miss in metadata
	result.Metadata.CacheHit = false
	cm.cache.Set(key, result, cache.DefaultExpiration)
}

// Clear removes all cached items
func (cm *CacheManager) Clear() {
	cm.cache.Flush()
}

// Size returns the number of items in the cache
func (cm *CacheManager) Size() int {
	return cm.cache.ItemCount()
}
