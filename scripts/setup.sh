#!/bin/bash

# Clean setup script for the Go project 19
echo "🧹 Cleaning up existing module files..."
rm -f go.sum

echo "📦 Initializing Go module..."
go mod init product-catalog-sorting

echo "📥 Adding dependencies..."
go get go.uber.org/zap@v1.26.0
go get go.uber.org/zap/zapcore@v1.26.0
go get github.com/pkg/errors@v0.9.1
go get github.com/stretchr/testify@v1.8.4

echo "🔧 Tidying module..."
go mod tidy

echo "✅ Setup complete! You can now run: go run ./cmd"

go test ./... -v 
go test ./test/unit -count=1 -v
 go test ./test/unit       
 go run ./cmd  
 find . -name '*_test.go'
grep -r "func Test" .
grep -r "suite.Run" .
go test -cover ./...
go test ./test/integration
go test -coverprofile=coverage.out ./...

