# Makefile for Gateway Service

.PHONY: help test test-unit coverage lint clean

help:
	@echo "Gateway Service - Test Commands"
	@echo "==============================="
	@echo "make test         - Run unit tests"
	@echo "make coverage     - Generate coverage report"
	@echo "make lint         - Run linters"
	@echo "make clean        - Clean test artifacts"

test:
	@echo "Running unit tests..."
	go test -v ./...

test-unit:
	@echo "Running unit tests with coverage..."
	go test -race -coverprofile=coverage.out ./src/...

coverage:
	@echo "Generating coverage report..."
	go test -coverprofile=coverage.out ./src/...
	go tool cover -func=coverage.out

coverage-html:
	go test -coverprofile=coverage.out ./src/...
	go tool cover -html=coverage.out -o coverage.html

lint:
	@echo "Running linters..."
	golangci-lint run --timeout=5m ./...

clean:
	@echo "Cleaning test artifacts..."
	rm -f coverage.out coverage.html
	go clean -testcache
