package routes

import (
	"github.com/JumpingMonkey/go-markdown-note-taking-app/internal/api/handlers"
	"github.com/JumpingMonkey/go-markdown-note-taking-app/internal/api/middleware"
	"github.com/JumpingMonkey/go-markdown-note-taking-app/internal/services/grammar"
	"github.com/JumpingMonkey/go-markdown-note-taking-app/internal/services/markdown"
	"github.com/JumpingMonkey/go-markdown-note-taking-app/internal/services/storage"
	"github.com/gin-gonic/gin"
)

// Setup configures all routes
func Setup(router *gin.Engine, storage storage.Storage, markdown *markdown.Service, grammar *grammar.Service) {
	// Apply global middleware
	router.Use(middleware.Logger())
	router.Use(middleware.CORS())
	
	// Create handlers
	notesHandler := handlers.NewNotesHandler(storage, markdown, grammar)

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Notes routes
		notes := v1.Group("/notes")
		{
			notes.POST("", notesHandler.CreateNote)
			notes.GET("", notesHandler.ListNotes)
			notes.GET("/:id", notesHandler.GetNote)
			notes.GET("/:id/html", notesHandler.GetNoteHTML)
			notes.DELETE("/:id", notesHandler.DeleteNote)
			notes.POST("/upload", notesHandler.UploadNote)
			notes.POST("/check-grammar", notesHandler.CheckGrammar)
		}

		// Documentation routes
		v1.GET("/docs", serveSwaggerUI)
		v1.GET("/docs/openapi.yaml", serveOpenAPISpec)
	}

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
}

// serveSwaggerUI serves the Swagger UI HTML
func serveSwaggerUI(c *gin.Context) {
	html := `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Go Markdown Note-Taking API Documentation</title>
    <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@4.5.0/swagger-ui.css">
</head>
<body>
    <div id="swagger-ui"></div>
    <script src="https://unpkg.com/swagger-ui-dist@4.5.0/swagger-ui-bundle.js"></script>
    <script>
        window.onload = function() {
            window.ui = SwaggerUIBundle({
                url: "/api/v1/docs/openapi.yaml",
                dom_id: '#swagger-ui',
                presets: [
                    SwaggerUIBundle.presets.apis,
                    SwaggerUIBundle.SwaggerUIStandalonePreset
                ],
                layout: "BaseLayout"
            });
        }
    </script>
</body>
</html>
`
	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(200, html)
}

// serveOpenAPISpec serves the OpenAPI specification
func serveOpenAPISpec(c *gin.Context) {
	c.File("./api/openapi.yaml")
}
