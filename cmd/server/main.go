package main

import (
	"log"
	"os"

	"github.com/JumpingMonkey/go-markdown-note-taking-app/internal/api/routes"
	"github.com/JumpingMonkey/go-markdown-note-taking-app/internal/config"
	"github.com/JumpingMonkey/go-markdown-note-taking-app/internal/services/grammar"
	"github.com/JumpingMonkey/go-markdown-note-taking-app/internal/services/markdown"
	"github.com/JumpingMonkey/go-markdown-note-taking-app/internal/services/storage"
	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize services
	storageService := storage.NewFileStorage(cfg.NotesDir)
	markdownService := markdown.NewService()
	grammarService := grammar.NewService()

	// Initialize Gin router
	router := gin.Default()

	// Setup routes
	routes.Setup(router, storageService, markdownService, grammarService)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting server on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
