package tests

import (
	"bytes"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/JumpingMonkey/go-markdown-note-taking-app/internal/api/routes"
	"github.com/JumpingMonkey/go-markdown-note-taking-app/internal/models"
	"github.com/JumpingMonkey/go-markdown-note-taking-app/internal/services/grammar"
	"github.com/JumpingMonkey/go-markdown-note-taking-app/internal/services/markdown"
	"github.com/JumpingMonkey/go-markdown-note-taking-app/internal/services/storage"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestServer(t *testing.T) (*gin.Engine, string, func()) {
	// Create temporary directory for testing
	tempDir, err := os.MkdirTemp("", "integration-test")
	require.NoError(t, err)

	// Create services
	storageService := storage.NewFileStorage(tempDir)
	markdownService := markdown.NewService()
	grammarService := grammar.NewService()

	// Setup router
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	routes.Setup(router, storageService, markdownService, grammarService)

	return router, tempDir, func() {
		os.RemoveAll(tempDir)
	}
}

func TestFullApiFlow(t *testing.T) {
	router, _, cleanup := setupTestServer(t)
	defer cleanup()

	// Create a new note
	createPayload := models.CreateNoteRequest{
		Title:   "Integration Test Note",
		Content: "# Integration Test\n\nThis is a test of the full API flow.",
	}

	jsonData, err := json.Marshal(createPayload)
	require.NoError(t, err)

	createw := httptest.NewRecorder()
	createReq, _ := http.NewRequest(http.MethodPost, "/api/v1/notes", bytes.NewBuffer(jsonData))
	createReq.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(createw, createReq)

	assert.Equal(t, http.StatusCreated, createw.Code)

	// Parse created note to get ID
	var createdNote models.Note
	err = json.Unmarshal(createw.Body.Bytes(), &createdNote)
	require.NoError(t, err)
	assert.NotEmpty(t, createdNote.ID)

	// Get the note by ID
	getw := httptest.NewRecorder()
	getReq, _ := http.NewRequest(http.MethodGet, "/api/v1/notes/"+createdNote.ID, nil)
	router.ServeHTTP(getw, getReq)

	assert.Equal(t, http.StatusOK, getw.Code)

	var retrievedNote models.Note
	err = json.Unmarshal(getw.Body.Bytes(), &retrievedNote)
	require.NoError(t, err)
	assert.Equal(t, createdNote.ID, retrievedNote.ID)
	assert.Equal(t, createdNote.Title, retrievedNote.Title)
	assert.Equal(t, createdNote.Content, retrievedNote.Content)

	// Get HTML rendering
	htmlw := httptest.NewRecorder()
	htmlReq, _ := http.NewRequest(http.MethodGet, "/api/v1/notes/"+createdNote.ID+"/html", nil)
	router.ServeHTTP(htmlw, htmlReq)

	assert.Equal(t, http.StatusOK, htmlw.Code)
	assert.Contains(t, htmlw.Body.String(), "<h1>Integration Test</h1>")

	// Check grammar
	grammarPayload := models.CheckGrammarRequest{
		Content: "this needs grammar checking",
	}

	grammarData, err := json.Marshal(grammarPayload)
	require.NoError(t, err)

	grammarw := httptest.NewRecorder()
	grammarReq, _ := http.NewRequest(http.MethodPost, "/api/v1/notes/check-grammar", bytes.NewBuffer(grammarData))
	grammarReq.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(grammarw, grammarReq)

	assert.Equal(t, http.StatusOK, grammarw.Code)

	var grammarResult models.GrammarCheckResult
	err = json.Unmarshal(grammarw.Body.Bytes(), &grammarResult)
	require.NoError(t, err)
	assert.NotEmpty(t, grammarResult.Issues)

	// List all notes
	listw := httptest.NewRecorder()
	listReq, _ := http.NewRequest(http.MethodGet, "/api/v1/notes", nil)
	router.ServeHTTP(listw, listReq)

	assert.Equal(t, http.StatusOK, listw.Code)

	var notesList []models.NoteMetadata
	err = json.Unmarshal(listw.Body.Bytes(), &notesList)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(notesList), 1)

	// Delete the note
	deletew := httptest.NewRecorder()
	deleteReq, _ := http.NewRequest(http.MethodDelete, "/api/v1/notes/"+createdNote.ID, nil)
	router.ServeHTTP(deletew, deleteReq)

	assert.Equal(t, http.StatusOK, deletew.Code)

	// Verify note is deleted
	getw2 := httptest.NewRecorder()
	getReq2, _ := http.NewRequest(http.MethodGet, "/api/v1/notes/"+createdNote.ID, nil)
	router.ServeHTTP(getw2, getReq2)

	assert.Equal(t, http.StatusNotFound, getw2.Code)
}

func TestUploadMarkdownFile(t *testing.T) {
	router, tempDir, cleanup := setupTestServer(t)
	defer cleanup()

	// Create a test markdown file
	testFilePath := filepath.Join(tempDir, "test-upload.md")
	testContent := "# Test Upload\n\nThis is a test file for upload functionality."
	err := os.WriteFile(testFilePath, []byte(testContent), 0644)
	require.NoError(t, err)

	// Create multipart form
	var b bytes.Buffer
	writer := multipart.NewWriter(&b)
	part, err := writer.CreateFormFile("file", "test-upload.md")
	require.NoError(t, err)

	fileContents, err := os.ReadFile(testFilePath)
	require.NoError(t, err)
	_, err = part.Write(fileContents)
	require.NoError(t, err)

	err = writer.Close()
	require.NoError(t, err)

	// Upload file
	uploadw := httptest.NewRecorder()
	uploadReq, _ := http.NewRequest(http.MethodPost, "/api/v1/notes/upload", &b)
	uploadReq.Header.Set("Content-Type", writer.FormDataContentType())
	router.ServeHTTP(uploadw, uploadReq)

	assert.Equal(t, http.StatusCreated, uploadw.Code)

	// Parse response to get note ID
	var uploadedNote models.Note
	err = json.Unmarshal(uploadw.Body.Bytes(), &uploadedNote)
	require.NoError(t, err)
	assert.NotEmpty(t, uploadedNote.ID)
	assert.Equal(t, "test-upload", uploadedNote.Title)
	assert.Equal(t, testContent, uploadedNote.Content)
}
