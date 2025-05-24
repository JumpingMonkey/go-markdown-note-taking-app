# Simple single-stage build
FROM golang:1.21-alpine

# Install build dependencies
RUN apk add --no-cache git gcc musl-dev ca-certificates

# Set working directory
WORKDIR /app

# Copy the entire application
COPY . .

# Initialize module with explicit version-pinned dependencies
RUN go get github.com/gin-gonic/gin@v1.9.1
RUN go get github.com/google/uuid@v1.5.0
RUN go get github.com/russross/blackfriday/v2@v2.1.0
RUN go get github.com/stretchr/testify@v1.8.4

# Fix go.sum and download dependencies
RUN go mod tidy

# Build the application
RUN go build -o main cmd/server/main.go

# Create notes directory
RUN mkdir -p notes

# Expose port
EXPOSE 8080

# Environment variables
ENV PORT=8080
ENV NOTES_DIR=./notes

# Run the application
CMD ["/app/main"]
