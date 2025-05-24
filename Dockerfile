# Build stage
FROM golang:1.21-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git gcc musl-dev

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod ./

# Initialize modules and generate go.sum
RUN go mod init github.com/JumpingMonkey/go-markdown-note-taking-app || true
RUN go get github.com/gin-gonic/gin@v1.9.1
RUN go get github.com/google/uuid@v1.5.0
RUN go get github.com/russross/blackfriday/v2@v2.1.0
RUN go get github.com/stretchr/testify@v1.8.4
RUN go mod tidy && go mod download && go mod verify

# Copy source code
COPY . .

# Clean any stale modules or cached files
RUN go clean -cache -modcache -i -r

# Show environment for debugging
RUN go env

# Build the application with detailed error output
RUN CGO_ENABLED=0 GOOS=linux go build -v -a -installsuffix cgo -tags netgo -ldflags='-w -extldflags "-static"' -o main cmd/server/main.go

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /app/main .

# Copy API documentation
COPY --from=builder /app/api ./api

# Create notes directory
RUN mkdir -p notes

# Expose port
EXPOSE 8080

# Environment variables
ENV PORT=8080
ENV NOTES_DIR=./notes

# Run the application
CMD ["./main"]
