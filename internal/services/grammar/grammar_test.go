package grammar

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGrammarService_Check(t *testing.T) {
	tests := []struct {
		name          string
		text          string
		expectedCount int
		expectedTypes []string
	}{
		{
			name:          "correct text",
			text:          "This is a correct sentence.",
			expectedCount: 0,
			expectedTypes: []string{},
		},
		{
			name:          "missing capitalization",
			text:          "this should start with a capital letter.",
			expectedCount: 1,
			expectedTypes: []string{"capitalization"},
		},
		{
			name:          "missing period",
			text:          "This sentence is missing a period",
			expectedCount: 1,
			expectedTypes: []string{"punctuation"},
		},
		{
			name:          "multiple spaces",
			text:          "This has  multiple  spaces.",
			expectedCount: 2,
			expectedTypes: []string{"spacing", "spacing"},
		},
		{
			name:          "multiple issues",
			text:          "this has multiple issues  and no period",
			expectedCount: 3,
			expectedTypes: []string{"capitalization", "spacing", "punctuation"},
		},
	}

	service := NewService()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := service.Check(tt.text)
			require.NoError(t, err)
			require.NotNil(t, result)
			
			assert.Len(t, result.Issues, tt.expectedCount)
			
			if tt.expectedCount > 0 {
				types := make([]string, len(result.Issues))
				for i, issue := range result.Issues {
					types[i] = issue.Type
				}
				
				// Just check that each expected type appears in the result
				// (not necessarily in the same order)
				for _, expectedType := range tt.expectedTypes {
					found := false
					for _, actualType := range types {
						if expectedType == actualType {
							found = true
							break
						}
					}
					assert.True(t, found, "Expected issue type %s not found", expectedType)
				}
			}
			
			// Score should be between 0 and 100
			assert.GreaterOrEqual(t, result.Score, 0.0)
			assert.LessOrEqual(t, result.Score, 100.0)
		})
	}
}
