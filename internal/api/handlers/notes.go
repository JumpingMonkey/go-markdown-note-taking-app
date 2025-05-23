package handlers

import (
	"net/http"
	"strings"

	"github.com/JumpingMonkey/go-markdown-note-taking-app/internal/models"
	"github.com/JumpingMonkey/go-markdown-note-taking-app/internal/services/grammar"
	"github.com/JumpingMonkey/go-markdown-note-taking-app/internal/services/markdown"
	"github.com/JumpingMonkey/go-markdown-note-taking-app/internal/services/storage"
	"github.com/gin-gonic/gin"
)

// NotesHandler handles note-related requests
type NotesHandler struct {
	storage  storage.Storage
	markdown *markdown.Service
	grammar  *grammar.Service
}

// NewNotesHandler creates a new notes handler
func NewNotesHandler(storage storage.Storage, markdown *markdown.Service, grammar *grammar.Service) *NotesHandler {
	return &NotesHandler{
		storage:  storage,
		markdown: markdown,
		grammar:  grammar,
	}
}

// CreateNote handles creating a new note
func (h *NotesHandler) CreateNote(c *gin.Context) {
	var req models.CreateNoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}

	note := &models.Note{
		Title:   req.Title,
		Content: req.Content,
	}

	if err := h.storage.Save(note); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to save note"})
		return
	}

	c.JSON(http.StatusCreated, note)
}

// ListNotes handles listing all notes
func (h *NotesHandler) ListNotes(c *gin.Context) {
	notes, err := h.storage.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to list notes"})
		return
	}

	c.JSON(http.StatusOK, notes)
}

// GetNote handles getting a specific note
func (h *NotesHandler) GetNote(c *gin.Context) {
	id := c.Param("id")
	
	note, err := h.storage.Get(id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "Note not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to get note"})
		return
	}

	c.JSON(http.StatusOK, note)
}

// GetNoteHTML handles getting a note rendered as HTML
func (h *NotesHandler) GetNoteHTML(c *gin.Context) {
	id := c.Param("id")
	
	note, err := h.storage.Get(id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "Note not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to get note"})
		return
	}

	html := h.markdown.ToHTML(note.Content)
	
	// Return as HTML content
	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(http.StatusOK, `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>%s</title>
    <style>
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, sans-serif;
            line-height: 1.6;
            color: #333;
            max-width: 800px;
            margin: 0 auto;
            padding: 20px;
            background-color: #f5f5f5;
        }
        .content {
            background-color: white;
            padding: 30px;
            border-radius: 8px;
            box-shadow: 0 2px 4px rgba(0,0,0,0.1);
        }
        pre {
            background-color: #f4f4f4;
            padding: 10px;
            border-radius: 4px;
            overflow-x: auto;
        }
        code {
            background-color: #f4f4f4;
            padding: 2px 4px;
            border-radius: 3px;
        }
        blockquote {
            border-left: 4px solid #ddd;
            margin: 0;
            padding-left: 20px;
            color: #666;
        }
    </style>
</head>
<body>
    <div class="content">
        <h1>%s</h1>
        %s
    </div>
</body>
</html>
`, note.Title, note.Title, html)
}

// DeleteNote handles deleting a note
func (h *NotesHandler) DeleteNote(c *gin.Context) {
	id := c.Param("id")
	
	if err := h.storage.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to delete note"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Note deleted successfully"})
}

// CheckGrammar handles grammar checking
func (h *NotesHandler) CheckGrammar(c *gin.Context) {
	var req models.CheckGrammarRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}

	result, err := h.grammar.Check(req.Content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to check grammar"})
		return
	}

	c.JSON(http.StatusOK, result)
}

// UploadNote handles file upload
func (h *NotesHandler) UploadNote(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Failed to get file"})
		return
	}
	defer file.Close()

	// Check if it's a markdown file
	if !strings.HasSuffix(strings.ToLower(header.Filename), ".md") {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Only markdown files are allowed"})
		return
	}

	// Save the uploaded file
	fileStorage, ok := h.storage.(*storage.FileStorage)
	if !ok {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Storage type not supported"})
		return
	}

	note, err := fileStorage.SaveUploadedFile(file, header.Filename)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to save uploaded file"})
		return
	}

	c.JSON(http.StatusCreated, note)
}
