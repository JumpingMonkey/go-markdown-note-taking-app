package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

// Logger is a middleware that logs request details
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method

		// Process request
		c.Next()

		// Calculate latency
		latency := time.Since(start)

		// Log details
		statusCode := c.Writer.Status()
		log.Printf("%s | %3d | %v | %s | %s", 
			method, 
			statusCode, 
			latency, 
			path, 
			c.ClientIP(),
		)
	}
}
