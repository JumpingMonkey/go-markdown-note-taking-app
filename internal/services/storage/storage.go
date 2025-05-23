package storage

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/JumpingMonkey/go-markdown-note-taking-app/internal/models"
	"github.com/google/uuid"
)

// Storage defines the interface for note storage
type Storage interface {
	Save(note *models.Note) error
	Get(id string) (*models.Note, error)
	List() ([]*models.NoteMetadata, error)
	Delete(id string) error
}

// FileStorage implements file-based storage for notes
type FileStorage struct {
	baseDir string
}

// NewFileStorage creates a new file storage instance
func NewFileStorage(baseDir string) *FileStorage {
	// Ensure the directory exists
	os.MkdirAll(baseDir, 0755)
	return &FileStorage{
		baseDir: baseDir,
	}
}

// Save saves a note to the file system
func (fs *FileStorage) Save(note *models.Note) error {
	if note.ID == "" {
		note.ID = uuid.New().String()
		note.CreatedAt = time.Now()
	}
	note.UpdatedAt = time.Now()

	// Create metadata file
	metadataPath := filepath.Join(fs.baseDir, note.ID+".json")
	metadata := map[string]interface{}{
		"id":         note.ID,
		"title":      note.Title,
		"created_at": note.CreatedAt,
		"updated_at": note.UpdatedAt,
	}

	metadataJSON, err := json.MarshalIndent(metadata, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}

	if err := os.WriteFile(metadataPath, metadataJSON, 0644); err != nil {
		return fmt.Errorf("failed to write metadata: %w", err)
	}

	// Create markdown file
	markdownPath := filepath.Join(fs.baseDir, note.ID+".md")
	if err := os.WriteFile(markdownPath, []byte(note.Content), 0644); err != nil {
		return fmt.Errorf("failed to write markdown: %w", err)
	}

	return nil
}

// Get retrieves a note by ID
func (fs *FileStorage) Get(id string) (*models.Note, error) {
	// Read metadata
	metadataPath := filepath.Join(fs.baseDir, id+".json")
	metadataData, err := os.ReadFile(metadataPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("note not found")
		}
		return nil, fmt.Errorf("failed to read metadata: %w", err)
	}

	var metadata map[string]interface{}
	if err := json.Unmarshal(metadataData, &metadata); err != nil {
		return nil, fmt.Errorf("failed to unmarshal metadata: %w", err)
	}

	// Read content
	markdownPath := filepath.Join(fs.baseDir, id+".md")
	content, err := os.ReadFile(markdownPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read content: %w", err)
	}

	// Parse timestamps
	createdAt, _ := time.Parse(time.RFC3339, metadata["created_at"].(string))
	updatedAt, _ := time.Parse(time.RFC3339, metadata["updated_at"].(string))

	return &models.Note{
		ID:        metadata["id"].(string),
		Title:     metadata["title"].(string),
		Content:   string(content),
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}, nil
}

// List returns all note metadata
func (fs *FileStorage) List() ([]*models.NoteMetadata, error) {
	files, err := os.ReadDir(fs.baseDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory: %w", err)
	}

	var notes []*models.NoteMetadata
	for _, file := range files {
		if filepath.Ext(file.Name()) == ".json" {
			metadataPath := filepath.Join(fs.baseDir, file.Name())
			metadataData, err := os.ReadFile(metadataPath)
			if err != nil {
				continue
			}

			var metadata map[string]interface{}
			if err := json.Unmarshal(metadataData, &metadata); err != nil {
				continue
			}

			createdAt, _ := time.Parse(time.RFC3339, metadata["created_at"].(string))
			updatedAt, _ := time.Parse(time.RFC3339, metadata["updated_at"].(string))

			notes = append(notes, &models.NoteMetadata{
				ID:        metadata["id"].(string),
				Title:     metadata["title"].(string),
				CreatedAt: createdAt,
				UpdatedAt: updatedAt,
			})
		}
	}

	return notes, nil
}

// Delete removes a note
func (fs *FileStorage) Delete(id string) error {
	metadataPath := filepath.Join(fs.baseDir, id+".json")
	markdownPath := filepath.Join(fs.baseDir, id+".md")

	if err := os.Remove(metadataPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to remove metadata: %w", err)
	}

	if err := os.Remove(markdownPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to remove markdown: %w", err)
	}

	return nil
}

// SaveUploadedFile saves an uploaded file
func (fs *FileStorage) SaveUploadedFile(reader io.Reader, filename string) (*models.Note, error) {
	content, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	// Extract title from filename (remove extension)
	title := filename
	if ext := filepath.Ext(filename); ext != "" {
		title = filename[:len(filename)-len(ext)]
	}

	note := &models.Note{
		Title:   title,
		Content: string(content),
	}

	if err := fs.Save(note); err != nil {
		return nil, err
	}

	return note, nil
}
