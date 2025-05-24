package testutils

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

// SetupRouter returns a configured Gin router for testing
func SetupRouter() *gin.Engine {
	// Switch to test mode
	gin.SetMode(gin.TestMode)
	return gin.New()
}

// PerformRequest performs an HTTP request for testing
func PerformRequest(r http.Handler, method, path string, body io.Reader) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

// CreateJSONRequest creates a test request with JSON body
func CreateJSONRequest(t *testing.T, data interface{}) io.Reader {
	jsonData, err := json.Marshal(data)
	require.NoError(t, err)
	return bytes.NewBuffer(jsonData)
}

// CreateTempMarkdownFile creates a temporary markdown file for testing
func CreateTempMarkdownFile(t *testing.T, content string) (string, func()) {
	tempDir, err := os.MkdirTemp("", "test-markdown")
	require.NoError(t, err)

	fileName := filepath.Join(tempDir, "test.md")
	err = os.WriteFile(fileName, []byte(content), 0644)
	require.NoError(t, err)

	return fileName, func() { os.RemoveAll(tempDir) }
}

// CreateMultipartRequest creates a multipart form request for file uploads
func CreateMultipartRequest(t *testing.T, fieldName, filePath string) (*bytes.Buffer, string) {
	file, err := os.Open(filePath)
	require.NoError(t, err)
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(fieldName, filepath.Base(filePath))
	require.NoError(t, err)

	_, err = io.Copy(part, file)
	require.NoError(t, err)
	err = writer.Close()
	require.NoError(t, err)

	return body, writer.FormDataContentType()
}
