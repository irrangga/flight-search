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
	}

	result := h.aggregator.Search(domainRequest)
	resp := mapper.ToSearchResponse(result)
	c.JSON(http.StatusOK, resp)
}

// SetupRouter configures routes with Gin
func SetupRouter(handler *Handler) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())
	r.HandleMethodNotAllowed = true
	r.POST("/api/search-flights", handler.searchFlightsHandler)
	return r
}
