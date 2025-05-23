package storage

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/JumpingMonkey/go-markdown-note-taking-app/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFileStorage_Save(t *testing.T) {
	// Create temporary directory for testing
	tempDir, err := os.MkdirTemp("", "notes_test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create storage instance
	storage := NewFileStorage(tempDir)

	// Test saving a new note
	note := &models.Note{
		Title:   "Test Note",
		Content: "# Test Content\n\nThis is a test.",
	}

	err = storage.Save(note)
	assert.NoError(t, err)
	assert.NotEmpty(t, note.ID)
	assert.NotZero(t, note.CreatedAt)
	assert.NotZero(t, note.UpdatedAt)

	// Verify files were created
	metadataPath := filepath.Join(tempDir, note.ID+".json")
	markdownPath := filepath.Join(tempDir, note.ID+".md")

	assert.FileExists(t, metadataPath)
	assert.FileExists(t, markdownPath)
}

func TestFileStorage_Get(t *testing.T) {
	// Create temporary directory for testing
	tempDir, err := os.MkdirTemp("", "notes_test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create storage instance
	storage := NewFileStorage(tempDir)

	// Save a note first
	originalNote := &models.Note{
		Title:   "Test Note",
		Content: "# Test Content\n\nThis is a test.",
	}
	err = storage.Save(originalNote)
	require.NoError(t, err)

	// Test getting the note
	retrievedNote, err := storage.Get(originalNote.ID)
	assert.NoError(t, err)
	assert.Equal(t, originalNote.ID, retrievedNote.ID)
	assert.Equal(t, originalNote.Title, retrievedNote.Title)
	assert.Equal(t, originalNote.Content, retrievedNote.Content)

	// Test getting non-existent note
	_, err = storage.Get("non-existent-id")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestFileStorage_List(t *testing.T) {
	// Create temporary directory for testing
	tempDir, err := os.MkdirTemp("", "notes_test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create storage instance
	storage := NewFileStorage(tempDir)

	// Save multiple notes
	notes := []*models.Note{
		{Title: "Note 1", Content: "Content 1"},
		{Title: "Note 2", Content: "Content 2"},
		{Title: "Note 3", Content: "Content 3"},
	}

	for _, note := range notes {
		err = storage.Save(note)
		require.NoError(t, err)
	}

	// Test listing notes
	list, err := storage.List()
	assert.NoError(t, err)
	assert.Len(t, list, 3)

	// Verify all notes are in the list
	titles := make(map[string]bool)
	for _, meta := range list {
		titles[meta.Title] = true
	}

	assert.True(t, titles["Note 1"])
	assert.True(t, titles["Note 2"])
	assert.True(t, titles["Note 3"])
}

func TestFileStorage_Delete(t *testing.T) {
	// Create temporary directory for testing
	tempDir, err := os.MkdirTemp("", "notes_test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create storage instance
	storage := NewFileStorage(tempDir)

	// Save a note
	note := &models.Note{
		Title:   "Test Note",
		Content: "# Test Content",
	}
	err = storage.Save(note)
	require.NoError(t, err)

	// Delete the note
	err = storage.Delete(note.ID)
	assert.NoError(t, err)

	// Verify files were deleted
	metadataPath := filepath.Join(tempDir, note.ID+".json")
	markdownPath := filepath.Join(tempDir, note.ID+".md")

	assert.NoFileExists(t, metadataPath)
	assert.NoFileExists(t, markdownPath)

	// Test deleting non-existent note (should not error)
	err = storage.Delete("non-existent-id")
	assert.NoError(t, err)
}
