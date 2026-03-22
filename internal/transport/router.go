package transport

import "github.com/gin-gonic/gin"

// SetupRouter configures routes with Gin
func SetupRouter(handler *Handler) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())
	r.HandleMethodNotAllowed = true
	r.POST("/api/search-flights", handler.searchFlightsHandler)
	return r
}
