# API Documentation Guide

## Overview

This document provides detailed information about the Go Markdown Note-Taking API, including endpoints, request/response formats, and usage examples.

## Base URL

```
Development: http://localhost:8080/api/v1
Production: https://api.yourdomain.com/api/v1
```

## Authentication

Currently, the API does not require authentication. Future versions will implement JWT-based authentication.

## Common Headers

```http
Content-Type: application/json
Accept: application/json
```

## API Endpoints

### 1. Health Check

Check if the service is running properly.

**Endpoint:** `GET /health`

**Response:**
```json
{
  "status": "ok"
}
```

### 2. Create Note

Create a new markdown note.

**Endpoint:** `POST /api/v1/notes`

**Request Body:**
```json
{
  "title": "My First Note",
  "content": "# Hello World\n\nThis is my first markdown note."
}
```

**Response (201 Created):**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "title": "My First Note",
  "content": "# Hello World\n\nThis is my first markdown note.",
  "created_at": "2024-01-20T10:30:00Z",
  "updated_at": "2024-01-20T10:30:00Z"
}
```

**Error Response (400 Bad Request):**
```json
{
  "error": "title and content are required"
}
```

### 3. List Notes

Get a list of all saved notes (metadata only).

**Endpoint:** `GET /api/v1/notes`

**Response (200 OK):**
```json
[
  {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "title": "My First Note",
    "created_at": "2024-01-20T10:30:00Z",
    "updated_at": "2024-01-20T10:30:00Z"
  },
  {
    "id": "650e8400-e29b-41d4-a716-446655440001",
    "title": "Another Note",
    "created_at": "2024-01-20T11:00:00Z",
    "updated_at": "2024-01-20T11:00:00Z"
  }
]
```

### 4. Get Note

Retrieve a specific note by ID.

**Endpoint:** `GET /api/v1/notes/{id}`

**Parameters:**
- `id` (path parameter): UUID of the note

**Response (200 OK):**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "title": "My First Note",
  "content": "# Hello World\n\nThis is my first markdown note.",
  "created_at": "2024-01-20T10:30:00Z",
  "updated_at": "2024-01-20T10:30:00Z"
}
```

**Error Response (404 Not Found):**
```json
{
  "error": "Note not found"
}
```

### 5. Get Note as HTML

Retrieve a note rendered as HTML.

**Endpoint:** `GET /api/v1/notes/{id}/html`

**Parameters:**
- `id` (path parameter): UUID of the note

**Response (200 OK):**
```html
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>My First Note</title>
    <style>/* CSS styles */</style>
</head>
<body>
    <div class="content">
        <h1>My First Note</h1>
        <h1>Hello World</h1>
        <p>This is my first markdown note.</p>
    </div>
</body>
</html>
```

### 6. Delete Note

Delete a note by ID.

**Endpoint:** `DELETE /api/v1/notes/{id}`

**Parameters:**
- `id` (path parameter): UUID of the note

**Response (200 OK):**
```json
{
  "message": "Note deleted successfully"
}
```

### 7. Upload Markdown File

Upload a markdown file to create a new note.

**Endpoint:** `POST /api/v1/notes/upload`

**Request:** Multipart form data
- `file`: The markdown file to upload

**cURL Example:**
```bash
curl -X POST http://localhost:8080/api/v1/notes/upload \
  -F "file=@mynote.md"
```

**Response (201 Created):**
```json
{
  "id": "750e8400-e29b-41d4-a716-446655440002",
  "title": "mynote",
  "content": "# File content here...",
  "created_at": "2024-01-20T12:00:00Z",
  "updated_at": "2024-01-20T12:00:00Z"
}
```

**Error Response (400 Bad Request):**
```json
{
  "error": "Only markdown files are allowed"
}
```

### 8. Check Grammar

Check the grammar of provided text. The service can detect multiple issues, including capitalization, punctuation, and multiple spaces.

**Endpoint:** `POST /api/v1/notes/check-grammar`

**Request Body:**
```json
{
  "content": "this is a test sentence  with multiple  spaces"
}
```

**Response (200 OK):**
```json
{
  "issues": [
    {
      "message": "Sentence should start with a capital letter",
      "offset": 0,
      "length": 1,
      "replacement": "T",
      "type": "capitalization"
    },
    {
      "message": "Multiple consecutive spaces detected",
      "offset": 19,
      "length": 2,
      "replacement": " ",
      "type": "spacing"
    },
    {
      "message": "Multiple consecutive spaces detected",
      "offset": 31,
      "length": 2,
      "replacement": " ",
      "type": "spacing"
    },
    {
      "message": "Sentence should end with proper punctuation",
      "offset": 38,
      "length": 0,
      "type": "punctuation"
    }
  ],
  "score": 75.0
}
```

## Error Handling

All error responses follow the same format:

```json
{
  "error": "Error message describing what went wrong"
}
```

### Common HTTP Status Codes

- `200 OK`: Request succeeded
- `201 Created`: Resource created successfully
- `400 Bad Request`: Invalid request data
- `404 Not Found`: Resource not found
- `500 Internal Server Error`: Server error

## Rate Limiting

Currently, there is no rate limiting. Future versions will implement rate limiting with the following headers:
- `X-RateLimit-Limit`: Maximum requests per hour
- `X-RateLimit-Remaining`: Remaining requests
- `X-RateLimit-Reset`: Time when limit resets

## Examples

### Creating a Note with cURL

```bash
curl -X POST http://localhost:8080/api/v1/notes \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Meeting Notes",
    "content": "# Meeting Notes\n\n- Discussed project timeline\n- Assigned tasks"
  }'
```

### Getting All Notes

```bash
curl http://localhost:8080/api/v1/notes
```

### Checking Grammar

```bash
curl -X POST http://localhost:8080/api/v1/notes/check-grammar \
  -H "Content-Type: application/json" \
  -d '{
    "content": "this sentence need grammar check"
  }'
```

## SDK Examples

### Go Client Example

```go
package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "net/http"
)

type Note struct {
    ID        string `json:"id"`
    Title     string `json:"title"`
    Content   string `json:"content"`
    CreatedAt string `json:"created_at"`
    UpdatedAt string `json:"updated_at"`
}

func createNote(title, content string) (*Note, error) {
    payload := map[string]string{
        "title":   title,
        "content": content,
    }
    
    jsonData, _ := json.Marshal(payload)
    resp, err := http.Post(
        "http://localhost:8080/api/v1/notes",
        "application/json",
        bytes.NewBuffer(jsonData),
    )
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    
    var note Note
    json.NewDecoder(resp.Body).Decode(&note)
    return &note, nil
}
```

### JavaScript/Node.js Example

```javascript
const axios = require('axios');

async function createNote(title, content) {
    try {
        const response = await axios.post('http://localhost:8080/api/v1/notes', {
            title: title,
            content: content
        });
        return response.data;
    } catch (error) {
        console.error('Error creating note:', error);
        throw error;
    }
}

// Usage
createNote('My Note', '# Hello\n\nThis is my note content')
    .then(note => console.log('Created note:', note))
    .catch(err => console.error('Error:', err));
```

### Python Example

```python
import requests
import json

def create_note(title, content):
    url = "http://localhost:8080/api/v1/notes"
    payload = {
        "title": title,
        "content": content
    }
    headers = {"Content-Type": "application/json"}
    
    response = requests.post(url, json=payload, headers=headers)
    
    if response.status_code == 201:
        return response.json()
    else:
        raise Exception(f"Error: {response.json()['error']}")

# Usage
try:
    note = create_note("Python Note", "# Created from Python\n\nThis is cool!")
    print(f"Created note with ID: {note['id']}")
except Exception as e:
    print(f"Failed to create note: {e}")
```

## OpenAPI Specification

The complete OpenAPI specification is available at:
- YAML: `/api/v1/docs/openapi.yaml`
- Swagger UI: `/api/v1/docs`

You can use the Swagger UI for interactive API exploration and testing.

## Postman Collection

A Postman collection is available for easy API testing. Import the collection from:
`docs/postman_collection.json`

## Support

For API support, please:
1. Check this documentation
2. Review the OpenAPI specification
3. Open an issue on GitHub: https://github.com/JumpingMonkey/go-markdown-note-taking-app/issues
