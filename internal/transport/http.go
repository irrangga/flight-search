package transport

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"flight-search/internal/aggregator"
	"flight-search/internal/domain"
	"flight-search/internal/mapper"
)

// SearchRequest is the HTTP request DTO
type SearchRequest struct {
	Origin        string  `json:"origin"`
	Destination   string  `json:"destination"`
	DepartureDate string  `json:"departureDate"`
	ReturnDate    *string `json:"returnDate,omitempty"`
	Passengers    int     `json:"passengers"`
	CabinClass    string  `json:"cabinClass"`

	// Filters
	PriceRange         []int    `json:"priceRange,omitempty"`
	NumberOfStops      []int    `json:"numberOfStops,omitempty"`
	DepartureTimeRange []string `json:"departureTimeRange,omitempty"`
	ArrivalTimeRange   []string `json:"arrivalTimeRange,omitempty"`
	Airlines           []string `json:"airlines,omitempty"`
	DurationRange      []int    `json:"durationRange,omitempty"`

	// Sort
	SortBy string `json:"sortBy,omitempty"`
}

// Handler handles HTTP requests
type Handler struct {
	aggregator *aggregator.Aggregator
}

// NewHandler creates a new HTTP handler
func NewHandler(agg *aggregator.Aggregator) *Handler {
	return &Handler{aggregator: agg}
}

func (h *Handler) searchFlightsHandler(c *gin.Context) {
	var req SearchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON request"})
		return
	}

	domainRequest := domain.SearchRequest{
		Origin:        req.Origin,
		Destination:   req.Destination,
		DepartureDate: req.DepartureDate,
		ReturnDate:    req.ReturnDate,
		Passengers:    req.Passengers,
		CabinClass:    req.CabinClass,

		PriceRange:         req.PriceRange,
		NumberOfStops:      req.NumberOfStops,
		DepartureTimeRange: req.DepartureTimeRange,
		ArrivalTimeRange:   req.ArrivalTimeRange,
		Airlines:           req.Airlines,
		DurationRange:      req.DurationRange,

		SortBy: req.SortBy,
	}

	result, err := h.aggregator.Search(domainRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	resp := mapper.ToSearchResponse(result)
	c.JSON(http.StatusOK, resp)
}
