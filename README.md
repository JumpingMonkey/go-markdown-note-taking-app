# Go Markdown Note-Taking App

A simple note-taking API built with Go that allows users to upload markdown files, check grammar, save notes, and render them as HTML.

## Features

- ✅ Upload and save markdown notes
- ✅ Grammar checking for notes
- ✅ List all saved notes
- ✅ Render markdown notes as HTML
- ✅ RESTful API design
- ✅ Docker support for easy deployment
- ✅ Comprehensive API documentation (OpenAPI/Swagger)

## Project Structure

```
go-markdown-note-taking-app/
├── cmd/
│   └── server/
│       └── main.go           # Application entry point
├── internal/
│   ├── api/
│   │   ├── handlers/         # HTTP handlers
│   │   ├── middleware/       # Middleware functions
│   │   └── routes/           # Route definitions
│   ├── config/               # Configuration management
│   ├── models/               # Data models
│   ├── services/             # Business logic
│   │   ├── grammar/          # Grammar checking service
│   │   ├── markdown/         # Markdown processing service
│   │   └── storage/          # Note storage service
│   └── utils/                # Utility functions
├── api/
│   └── openapi.yaml          # OpenAPI specification
├── docs/                     # Documentation
├── scripts/                  # Build and deployment scripts
├── docker/                   # Docker-related files
├── notes/                    # Directory for stored notes
├── go.mod
├── go.sum
├── Dockerfile
├── docker-compose.yml
└── README.md
```

## Prerequisites

- Go 1.21 or higher
- Docker and Docker Compose (optional, for containerized deployment)
- Git

## Installation

1. Clone the repository:
```bash
git clone https://github.com/JumpingMonkey/go-markdown-note-taking-app.git
cd go-markdown-note-taking-app
```

2. Install dependencies:
```bash
go mod download
```

3. Create the notes directory:
```bash
mkdir -p notes
```

## Running the Application

### Local Development

```bash
go run cmd/server/main.go
```

The server will start on `http://localhost:8080` by default.

### Using Docker

```bash
docker-compose up --build
```

## API Endpoints

### 1. Save a Note
- **POST** `/api/v1/notes`
- **Request Body**: 
  ```json
  {
    "title": "My Note",
    "content": "# Markdown content here"
  }
  ```
- **Response**: Created note with ID

### 2. Check Grammar
- **POST** `/api/v1/notes/check-grammar`
- **Request Body**: 
  ```json
  {
    "content": "Text to check for grammar"
  }
  ```
- **Response**: Grammar check results

### 3. List All Notes
- **GET** `/api/v1/notes`
- **Response**: Array of saved notes

### 4. Get a Specific Note
- **GET** `/api/v1/notes/{id}`
- **Response**: Note details in markdown format

### 5. Get HTML Rendered Note
- **GET** `/api/v1/notes/{id}/html`
- **Response**: Note rendered as HTML

### 6. Upload Markdown File
- **POST** `/api/v1/notes/upload`
- **Request**: Multipart form with markdown file
- **Response**: Created note with ID

## API Documentation

The API documentation is available in OpenAPI format:
- View the OpenAPI specification: `/api/v1/docs/openapi.yaml`
- Interactive Swagger UI: `/api/v1/docs` (when running the server)

## Configuration

The application can be configured using environment variables:

- `PORT`: Server port (default: 8080)
- `NOTES_DIR`: Directory to store notes (default: ./notes)
- `LOG_LEVEL`: Logging level (default: info)

## Development

### Running Tests

```bash
go test ./...
```

### Building the Application

```bash
go build -o bin/server cmd/server/main.go
```

### Code Style

This project follows the standard Go formatting guidelines. Use `go fmt` to format your code:

```bash
go fmt ./...
```

## Docker Support

### Building the Docker Image

```bash
docker build -t go-markdown-note-app .
```

### Running with Docker Compose

```bash
docker-compose up -d
```

This will start the application along with any required services.

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -am 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Acknowledgments

- [Blackfriday](https://github.com/russross/blackfriday) - Markdown processor
- [Gin](https://github.com/gin-gonic/gin) - HTTP web framework
- [LanguageTool](https://languagetool.org/) - Grammar checking
