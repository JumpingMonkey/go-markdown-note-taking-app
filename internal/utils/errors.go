package utils

import (
	"github.com/JumpingMonkey/go-markdown-note-taking-app/internal/models"
	"github.com/gin-gonic/gin"
)

// ErrorResponse standardizes error responses across the API
func ErrorResponse(c *gin.Context, status int, message string) {
	c.JSON(status, models.ErrorResponse{Error: message})
}

// ErrorResponseWithDetails returns a detailed error response
func ErrorResponseWithDetails(c *gin.Context, status int, message string, details interface{}) {
	c.JSON(status, gin.H{
		"error":   message,
		"details": details,
	})
}
