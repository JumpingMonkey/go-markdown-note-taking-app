package handlers

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/JumpingMonkey/go-markdown-note-taking-app/internal/models"
	"github.com/JumpingMonkey/go-markdown-note-taking-app/internal/services/grammar"
	"github.com/JumpingMonkey/go-markdown-note-taking-app/internal/services/markdown"
	"github.com/JumpingMonkey/go-markdown-note-taking-app/internal/services/storage"
	"github.com/JumpingMonkey/go-markdown-note-taking-app/internal/utils/testutils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTest(t *testing.T) (
	*NotesHandler,
	*gin.Engine,
	*storage.FileStorage,
	*markdown.Service,
	*grammar.Service,
	func(),
) {
	// Create temporary directory for testing
	tempDir, err := testutils.CreateTempMarkdownFile(t, "# Test Note\nTest content.")
	require.NoError(t, err)

	// Create services
	storageService := storage.NewFileStorage(tempDir)
	markdownService := markdown.Service{}
	grammarService := grammar.Service{}

	// Create handler
	handler := NewNotesHandler(storageService, &markdownService, &grammarService)

	// Setup router
	router := testutils.SetupRouter()
	
	// Register routes
	api := router.Group("/api/v1")
	notes := api.Group("/notes")
	notes.POST("", handler.CreateNote)
	notes.GET("", handler.ListNotes)
	notes.GET("/:id", handler.GetNote)
	notes.GET("/:id/html", handler.GetNoteHTML)
	notes.DELETE("/:id", handler.DeleteNote)
	notes.POST("/check-grammar", handler.CheckGrammar)

	return handler, router, storageService, &markdownService, &grammarService, func() {
		// Cleanup
		// Remove temporary directory
	}
}

func TestCreateNote(t *testing.T) {
	_, router, _, _, _, cleanup := setupTest(t)
	defer cleanup()

	// Create request payload
	payload := models.CreateNoteRequest{
		Title:   "Test Note",
		Content: "# Test Note\n\nThis is a test note.",
	}

	body := testutils.CreateJSONRequest(t, payload)
	
	// Perform request
	w := testutils.PerformRequest(router, http.MethodPost, "/api/v1/notes", body)

	// Assert response
	assert.Equal(t, http.StatusCreated, w.Code)

	// Parse response
	var response models.Note
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	// Validate response
	assert.NotEmpty(t, response.ID)
	assert.Equal(t, payload.Title, response.Title)
	assert.Equal(t, payload.Content, response.Content)
	assert.NotZero(t, response.CreatedAt)
	assert.NotZero(t, response.UpdatedAt)
}

func TestCheckGrammar(t *testing.T) {
	_, router, _, _, _, cleanup := setupTest(t)
	defer cleanup()

	// Create request payload
	payload := models.CheckGrammarRequest{
		Content: "this needs grammar checking.",
	}

	body := testutils.CreateJSONRequest(t, payload)
	
	// Perform request
	w := testutils.PerformRequest(router, http.MethodPost, "/api/v1/notes/check-grammar", body)

	// Assert response
	assert.Equal(t, http.StatusOK, w.Code)

	// Parse response
	var response models.GrammarCheckResult
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	// Validate response
	assert.GreaterOrEqual(t, len(response.Issues), 1) // Should have at least one issue
	assert.GreaterOrEqual(t, response.Score, 0.0)
	assert.LessOrEqual(t, response.Score, 100.0)
}
