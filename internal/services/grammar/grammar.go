package grammar

import (
	"strings"
	"unicode"

	"github.com/JumpingMonkey/go-markdown-note-taking-app/internal/models"
)

// Service provides grammar checking functionality
type Service struct {
	// In a real implementation, this might integrate with LanguageTool API
}

// NewService creates a new grammar service
func NewService() *Service {
	return &Service{}
}

// Check performs grammar checking on the provided text
func (s *Service) Check(text string) (*models.GrammarCheckResult, error) {
	// This is a simple implementation. In production, you would integrate
	// with a proper grammar checking service like LanguageTool API
	
	issues := []models.GrammarIssue{}
	
	// Simple checks for demonstration
	
	// Check for multiple spaces
	if strings.Contains(text, "  ") {
		idx := strings.Index(text, "  ")
		issues = append(issues, models.GrammarIssue{
			Message:     "Multiple consecutive spaces found",
			Offset:      idx,
			Length:      2,
			Replacement: " ",
			Type:        "spacing",
		})
	}
	
	// Check for missing capital letter at start
	if len(text) > 0 && unicode.IsLower(rune(text[0])) {
		issues = append(issues, models.GrammarIssue{
			Message:     "Sentence should start with a capital letter",
			Offset:      0,
			Length:      1,
			Replacement: strings.ToUpper(string(text[0])),
			Type:        "capitalization",
		})
	}
	
	// Check for missing period at end
	trimmed := strings.TrimSpace(text)
	if len(trimmed) > 0 && !strings.HasSuffix(trimmed, ".") && !strings.HasSuffix(trimmed, "!") && !strings.HasSuffix(trimmed, "?") {
		issues = append(issues, models.GrammarIssue{
			Message: "Sentence should end with proper punctuation",
			Offset:  len(trimmed),
			Length:  0,
			Type:    "punctuation",
		})
	}
	
	// Calculate a simple score (100 - 10 points per issue)
	score := 100.0 - float64(len(issues)*10)
	if score < 0 {
		score = 0
	}
	
	return &models.GrammarCheckResult{
		Issues: issues,
		Score:  score,
	}, nil
}

// For a production implementation, you would add methods like:
// - Integration with LanguageTool API
// - More sophisticated grammar rules
// - Support for different languages
// - Caching of results
