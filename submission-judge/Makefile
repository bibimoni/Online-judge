# Makefile for Submission-Judge Service

.PHONY: help test test-unit test-integration coverage lint benchmark clean

help:
	@echo "Submission-Judge Service - Test Commands"
	@echo "========================================"
	@echo "make test              - Run unit tests"
	@echo "make test-integration  - Run integration tests"
	@echo "make coverage          - Generate coverage report"
	@echo "make lint              - Run linters"
	@echo "make benchmark         - Run benchmarks"
	@echo "make clean             - Clean test artifacts"

test:
	@echo "Running unit tests..."
	go test -short -v ./...

test-unit:
	@echo "Running unit tests with coverage..."
	go test -short -race -coverprofile=coverage.out ./src/...

test-integration:
	@echo "Running integration tests..."
	go test -v ./src/tests/...

coverage:
	@echo "Generating coverage report..."
	go test -short -coverprofile=coverage.out ./src/...
	go tool cover -func=coverage.out

coverage-html:
	go test -short -coverprofile=coverage.out ./src/...
	go tool cover -html=coverage.out -o coverage.html

lint:
	@echo "Running linters..."
	golangci-lint run --timeout=5m ./...

benchmark:
	@echo "Running benchmarks..."
	go test -bench=. -benchmem ./...

clean:
	@echo "Cleaning test artifacts..."
	rm -f coverage.out coverage.html
	go clean -testcache
