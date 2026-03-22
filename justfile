# go-output - Build and test commands
# https://github.com/casey/just

# Default recipe - show help
default:
    @just --list

# Build all packages
build:
    go build ./...

# Run all tests
test:
    go test ./...

# Run tests with verbose output
test-v:
    go test -v ./...

# Run tests with coverage
test-cover:
    go test -cover ./...

# Run linting
lint:
    golangci-lint run --fix ./...

# Format code
fmt:
    go fmt ./...

# Tidy dependencies
tidy:
    go mod tidy

# Verify build and tests
verify: build test lint
    @echo "All checks passed!"

# Run example
run-example:
    go run ./examples/basic/main.go

# Clean build artifacts
clean:
    rm -f coverage.out
    go clean

# Show dependency graph
deps:
    go mod graph
