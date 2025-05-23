#!/bin/bash

# Create necessary directories
mkdir -p notes

# Start the application
echo "Starting Go Markdown Note-Taking App..."
go run cmd/server/main.go
