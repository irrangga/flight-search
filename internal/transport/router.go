package transport

import "github.com/gin-gonic/gin"

// SetupRouter configures routes with Gin
func SetupRouter(handler *Handler) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(gin.Logger())
	r.HandleMethodNotAllowed = true

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	r.POST("/api/search-flights", handler.searchFlightsHandler)
	return r
}
