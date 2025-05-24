#!/bin/bash

# This script runs tests inside a Docker container

echo "Building test container..."
docker build -t go-markdown-note-taking-app-test .

echo "\nRunning tests inside container..."
docker run --rm go-markdown-note-taking-app-test sh -c "go test ./... -v"
