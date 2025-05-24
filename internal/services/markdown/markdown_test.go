package markdown

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMarkdownService_ToHTML(t *testing.T) {
	tests := []struct {
		name     string
		markdown string
		expected string
	}{
		{
			name:     "heading",
			markdown: "# Heading",
			expected: "<h1 id=\"heading\">Heading</h1>\n",
		},
		{
			name:     "paragraph",
			markdown: "This is a paragraph",
			expected: "<p>This is a paragraph</p>\n",
		},
		{
			name:     "bold",
			markdown: "**Bold text**",
			expected: "<p><strong>Bold text</strong></p>\n",
		},
		{
			name:     "italic",
			markdown: "*Italic text*",
			expected: "<p><em>Italic text</em></p>\n",
		},
		{
			name:     "link",
			markdown: "[Link](https://example.com)",
			expected: "<p><a href=\"https://example.com\">Link</a></p>\n",
		},
		{
			name:     "code block",
			markdown: "```\nCode block\n```",
			expected: "<pre><code>Code block\n</code></pre>\n",
		},
		{
			name:     "combined",
			markdown: "# Title\n\nParagraph with **bold** and *italic*\n\n- List item 1\n- List item 2",
			expected: "<h1 id=\"title\">Title</h1>\n\n<p>Paragraph with <strong>bold</strong> and <em>italic</em></p>\n\n<ul>\n<li>List item 1</li>\n<li>List item 2</li>\n</ul>\n",
		},
	}

	service := NewService()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			html := service.ToHTML(tt.markdown)
			html = strings.ReplaceAll(html, "\r\n", "\n") // normalize line endings
			assert.Equal(t, tt.expected, html)
		})
	}
}

func TestMarkdownService_Validate(t *testing.T) {
	service := NewService()

	// Test valid markdown
	err := service.Validate("# Valid Markdown")
	assert.NoError(t, err)

	// Test empty markdown (should still be valid)
	err = service.Validate("")
	assert.NoError(t, err)
}
