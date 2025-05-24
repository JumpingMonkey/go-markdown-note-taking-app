package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// EnsureDirectory ensures that the specified directory exists
func EnsureDirectory(path string) error {
	return os.MkdirAll(path, 0755)
}

// IsMarkdownFile checks if the filename has a markdown extension
func IsMarkdownFile(filename string) bool {
	extensions := []string{".md", ".markdown"}
	extension := strings.ToLower(filepath.Ext(filename))
	
	for _, validExt := range extensions {
		if extension == validExt {
			return true
		}
	}
	
	return false
}

// SanitizeFilename removes potentially dangerous characters from a filename
func SanitizeFilename(filename string) string {
	// Replace all non-alphanumeric characters except for safe ones
	safe := []byte{'_', '-', '.', ' '}
	result := ""
	
	for _, char := range filename {
		if (char >= 'a' && char <= 'z') || 
		   (char >= 'A' && char <= 'Z') || 
		   (char >= '0' && char <= '9') {
			result += string(char)
			continue
		}
		
		isSafe := false
		for _, safeChar := range safe {
			if byte(char) == safeChar {
				isSafe = true
				break
			}
		}
		
		if isSafe {
			result += string(char)
		}
	}
	
	// Prevent empty filenames
	if result == "" {
		return "untitled"
	}
	
	return result
}

// GetFileSize returns the size of a file in bytes
func GetFileSize(path string) (int64, error) {
	info, err := os.Stat(path)
	if err != nil {
		return 0, fmt.Errorf("failed to get file size: %w", err)
	}
	
	return info.Size(), nil
}
