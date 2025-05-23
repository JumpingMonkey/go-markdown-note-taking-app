package models

import "time"

// Note represents a markdown note
type Note struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// NoteMetadata represents note metadata without content
type NoteMetadata struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CreateNoteRequest represents a request to create a new note
type CreateNoteRequest struct {
	Title   string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
}

// CheckGrammarRequest represents a request to check grammar
type CheckGrammarRequest struct {
	Content string `json:"content" binding:"required"`
}

// GrammarCheckResult represents the result of a grammar check
type GrammarCheckResult struct {
	Issues []GrammarIssue `json:"issues"`
	Score  float64        `json:"score"`
}

// GrammarIssue represents a single grammar issue
type GrammarIssue struct {
	Message     string `json:"message"`
	Offset      int    `json:"offset"`
	Length      int    `json:"length"`
	Replacement string `json:"replacement,omitempty"`
	Type        string `json:"type"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error string `json:"error"`
}
