#!/bin/bash

echo "🔧 Fixing Go dependencies..."

# Clean module cache
go clean -modcache

# Remove go.sum to start fresh
rm -f go.sum

# Tidy the module
go mod tidy

# Download dependencies
go mod download

# Verify dependencies
go mod verify

echo "✅ Dependencies fixed!"
