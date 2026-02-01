# Makefile for Contest Service Testing

.PHONY: help test test-unit test-integration test-e2e test-all coverage lint benchmark clean

# Default target
help:
	@echo "Contest Service - Test Commands"
	@echo "================================"
	@echo "make test              - Run unit tests only"
	@echo "make test-unit         - Run unit tests with coverage"
	@echo "make test-integration  - Run integration tests (requires Docker)"
	@echo "make test-e2e          - Run end-to-end tests"
	@echo "make test-all          - Run all tests"
	@echo "make coverage          - Generate coverage report"
	@echo "make coverage-html     - Generate HTML coverage report"
	@echo "make lint              - Run linters"
	@echo "make benchmark         - Run performance benchmarks"
	@echo "make clean             - Clean test artifacts"
	@echo "make deps              - Install test dependencies"

# Install dependencies
deps:
	@echo "Installing test dependencies..."
	go mod download
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/rakyll/hey@latest
	go install golang.org/x/tools/cmd/cover@latest

# Run unit tests only (fast)
test:
	@echo "Running unit tests..."
	go test -short -v ./...

# Run unit tests with coverage
test-unit:
	@echo "Running unit tests with coverage..."
	go test -short -race -coverprofile=coverage.out -covermode=atomic ./src/...
	@echo "Coverage: $$(go tool cover -func=coverage.out | grep total | awk '{print $$3}')"

# Run integration tests (requires Docker)
test-integration:
	@echo "Running integration tests..."
	@echo "Note: Docker must be running for test containers"
	go test -v -race -tags=integration ./src/tests/integration/...

# Run E2E tests
test-e2e:
	@echo "Running E2E tests..."
	@echo "Starting services with docker-compose..."
	docker-compose up -d
	@sleep 5
	go test -v ./src/tests/e2e/...
	@echo "Stopping services..."
	docker-compose down

# Run all tests
test-all: test-unit test-integration
	@echo "All tests completed!"

# Generate coverage report
coverage:
	@echo "Generating coverage report..."
	go test -short -coverprofile=coverage.out ./src/...
	go tool cover -func=coverage.out
	@echo "\nTotal Coverage: $$(go tool cover -func=coverage.out | grep total | awk '{print $$3}')"

# Generate HTML coverage report
coverage-html: coverage
	@echo "Generating HTML coverage report..."
	go tool cover -html=coverage.out -o coverage.html
	@echo "Opening coverage report in browser..."
	@which open > /dev/null && open coverage.html || echo "Open coverage.html in your browser"

# Run linters
lint:
	@echo "Running linters..."
	golangci-lint run --timeout=5m ./...

# Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...
	gofmt -s -w .

# Run benchmarks
benchmark:
	@echo "Running benchmarks..."
	go test -bench=. -benchmem -benchtime=5s ./... | tee benchmark.txt
	@echo "Benchmark results saved to benchmark.txt"

# Run specific benchmark
benchmark-scoreboard:
	@echo "Running scoreboard calculation benchmark..."
	go test -bench=BenchmarkScoreboardCalculation -benchmem -benchtime=10s ./src/service/

# Check for race conditions
race:
	@echo "Running race detector..."
	go test -race ./...

# Run tests with verbose output
test-verbose:
	@echo "Running tests with verbose output..."
	go test -v ./...

# Run tests multiple times to detect flakiness
test-flaky:
	@echo "Running tests 10 times to detect flakiness..."
	go test ./... -count=10 -failfast

# Generate mocks (if using mockgen)
mocks:
	@echo "Generating mocks..."
	go generate ./...

# Clean test artifacts
clean:
	@echo "Cleaning test artifacts..."
	rm -f coverage.out coverage.html benchmark.txt
	rm -rf tmp/
	go clean -testcache

# Watch mode - run tests on file changes (requires entr or gow)
watch:
	@echo "Watching for changes..."
	@which gow > /dev/null && gow test ./... || echo "Install gow: go install github.com/mitranim/gow@latest"

# Quick validation before commit
pre-commit: fmt lint test-unit
	@echo "Pre-commit checks passed!"

# CI simulation - run what CI will run
ci: lint test-unit test-integration
	@echo "CI simulation completed!"

# Security scan
security:
	@echo "Running security scan..."
	gosec -quiet ./...

# Test specific package
test-package:
	@read -p "Enter package path (e.g., ./src/domain/entity): " pkg; \
	go test -v $$pkg

# Show test coverage by package
coverage-by-package:
	@echo "Coverage by package:"
	@go test -short -coverprofile=coverage.out ./src/... > /dev/null 2>&1
	@go tool cover -func=coverage.out | grep -v "total"
