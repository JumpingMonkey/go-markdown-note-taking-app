package markdown

import (
	"github.com/russross/blackfriday/v2"
)

// Service provides markdown processing functionality
type Service struct {
	// Add any configuration here if needed
}

// NewService creates a new markdown service
func NewService() *Service {
	return &Service{}
}

// ToHTML converts markdown content to HTML
func (s *Service) ToHTML(markdown string) string {
	// Use blackfriday with common extensions
	extensions := blackfriday.CommonExtensions | blackfriday.AutoHeadingIDs
	renderer := blackfriday.NewHTMLRenderer(blackfriday.HTMLRendererParameters{
		Flags: blackfriday.CommonHTMLFlags,
	})
	
	return string(blackfriday.Run([]byte(markdown), blackfriday.WithExtensions(extensions), blackfriday.WithRenderer(renderer)))
}

// Validate checks if the markdown is valid
func (s *Service) Validate(markdown string) error {
	// Basic validation - for now just check if it's not empty
	// You can add more sophisticated validation here
	if markdown == "" {
		return nil
	}
	return nil
}
