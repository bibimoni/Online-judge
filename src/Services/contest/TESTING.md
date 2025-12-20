# Testing Strategy for Online Judge Platform

## Overview
This document outlines the comprehensive testing strategy for the online judge platform, with a focus on the contest service.

## Test Pyramid

```
        /\
       /E2E\          ~5% - Critical user flows
      /------\
     /Integr-\        ~15% - Service interactions
    /----------\
   /   Unit    \      ~80% - Business logic
  /--------------\
```

## Test Categories

### 1. Unit Tests (80% of test suite)
**Location**: `src/Services/contest/src/*/`*_test.go`
**Purpose**: Test individual functions and methods in isolation

#### Coverage Targets
- Domain entities: 90%+
- Services: 85%+
- Use cases: 90%+
- Utilities: 80%+

#### Examples
- `contest_test.go` - Domain validation, state transitions
- `contest_service_test.go` - Business logic with mocks
- `scoreboard_calculator_test.go` - Scoring algorithms

#### Running Unit Tests
```bash
# Run all unit tests
cd src/Services/contest
go test ./... -v

# Run with coverage
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out

# Run specific test
go test ./src/domain/entity -run TestContest_Validation

# Run tests in short mode (skip integration)
go test ./... -short
```

### 2. Integration Tests (15% of test suite)
**Location**: `src/Services/contest/src/tests/integration/`
**Purpose**: Test interactions between components with real dependencies

#### Test Containers
We use [Testcontainers](https://golang.testcontainers.org/) for spinning up real databases:
- MongoDB for data persistence
- Redis for caching and pub/sub

#### Examples
- Contest creation → Database → Retrieval
- Submission ingestion → Scoreboard update
- Rejudge flow with database updates

#### Running Integration Tests
```bash
# Requires Docker running
cd src/Services/contest
go test ./src/tests/integration/... -v

# Skip integration tests
go test ./... -short
```

### 3. E2E Tests (5% of test suite)
**Location**: `src/Services/contest/src/tests/e2e/`
**Purpose**: Test complete user workflows across services

#### Critical Paths
1. **Contest Registration Flow**
   - User registers → Contest appears → Can submit
2. **Submission to Scoreboard Flow**
   - Submit solution → Judge evaluates → Scoreboard updates
3. **Rejudge Flow**
   - Admin triggers rejudge → Submissions re-evaluated → Scoreboard recalculated

#### Running E2E Tests
```bash
# Start all services with docker-compose
docker-compose up -d

# Run E2E tests
cd src/Services/contest
go test ./src/tests/e2e/... -v

# Clean up
docker-compose down
```

## Test Fixtures & Utilities

### Builders Pattern
Use fluent builders for creating test data:

```go
contest := testutil.NewContestBuilder().
    WithName("ICPC Regional").
    AsICPC().
    AsOngoing().
    WithRated(true).
    Build()
```

### Mock Implementations
Located in `src/testutil/mocks/`:
- `MockContestRepository`
- `MockContestSubmissionRepository`
- `MockScoreboardCalculator`
- `MockRatingPolicy`

### Test Data
Use deterministic test data:
- Fixed timestamps (avoid `time.Now()` in assertions)
- Predictable UUIDs for specific scenarios
- Consistent usernames, problem IDs

## CI/CD Integration

### GitHub Actions Workflow
File: `.github/workflows/contest-service-ci.yml`

**Triggers**:
- Push to `main` or `develop`
- Pull requests

**Jobs**:
1. **Test** - Run all tests with multiple Go versions
2. **Security** - Gosec security scanning
3. **Build** - Docker image build and push
4. **Benchmark** - Performance regression detection

### Pre-commit Hooks
```bash
# Install pre-commit
pip install pre-commit

# Install hooks
pre-commit install

# Run manually
pre-commit run --all-files
```

## Coverage Requirements

### Minimum Coverage Targets
- Overall: 80%
- Critical paths (contest creation, submission ingestion): 95%
- Domain entities: 90%
- Services: 85%

### Coverage Reports
```bash
# Generate HTML coverage report
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html
open coverage.html

# View coverage in terminal
go test ./... -cover

# Detailed coverage by package
go test ./... -coverprofile=coverage.out
go tool cover -func=coverage.out
```

## Performance Testing

### Benchmarks
Located in `*_test.go` files with `Benchmark` prefix

```bash
# Run benchmarks
go test -bench=. -benchmem ./...

# Run specific benchmark
go test -bench=BenchmarkScoreboardCalculation -benchtime=10s

# Compare benchmarks
go test -bench=. -benchmem ./... > old.txt
# Make changes
go test -bench=. -benchmem ./... > new.txt
benchcmp old.txt new.txt
```

### Load Testing
Use `hey` or `k6` for API load testing:

```bash
# Install hey
go install github.com/rakyll/hey@latest

# Load test scoreboard endpoint
hey -n 10000 -c 100 -m GET http://localhost:8001/api/v1/contest/123/scoreboard
```

## Test Data Management

### Database Seeding
For integration tests:
```go
func seedDatabase(t *testing.T, db *mongo.Database) {
    // Insert test contests
    // Insert test users
    // Insert test submissions
}
```

### Cleanup Strategy
```go
t.Cleanup(func() {
    // Clean up resources
    db.Drop(ctx)
    container.Terminate(ctx)
})
```

## Best Practices

### 1. Test Naming
```go
// Pattern: Test<FunctionName>_<Scenario>_<ExpectedBehavior>
func TestContest_Create_ValidData_ReturnsID(t *testing.T)
func TestContest_Create_InvalidDates_ReturnsError(t *testing.T)
```

### 2. Table-Driven Tests
```go
tests := []struct {
    name     string
    input    Contest
    wantErr  bool
}{
    {"valid contest", validContest, false},
    {"invalid dates", invalidContest, true},
}
```

### 3. Arrange-Act-Assert Pattern
```go
// Arrange
contest := testutil.NewContestBuilder().Build()

// Act
result, err := service.CreateContest(ctx, contest)

// Assert
assert.NoError(t, err)
assert.NotEmpty(t, result.ID)
```

### 4. Avoid Test Interdependence
- Each test should be independent
- Use `t.Parallel()` when possible
- Clean up after each test

### 5. Mock External Dependencies
- Don't call real APIs in tests
- Use mocks for repositories
- Use test containers for databases

### 6. Test Edge Cases
- Empty inputs
- Nil values
- Boundary conditions
- Concurrent access
- Network failures
- Database errors

## Continuous Improvement

### Flaky Test Detection
```bash
# Run tests multiple times to detect flakiness
go test ./... -count=10 -failfast
```

### Mutation Testing
Consider using [go-mutesting](https://github.com/zimmski/go-mutesting) to ensure tests catch bugs

### Coverage Trends
Track coverage over time in CI/CD:
- Fail if coverage drops below threshold
- Comment coverage changes on PRs

## Resources

- [Go Testing Documentation](https://golang.org/pkg/testing/)
- [Testcontainers for Go](https://golang.testcontainers.org/)
- [Testify - Testing toolkit](https://github.com/stretchr/testify)
- [gomock - Mocking framework](https://github.com/golang/mock)

## Quick Reference

```bash
# Run all tests
go test ./...

# Run with coverage
go test ./... -cover

# Run integration tests only
go test -tags=integration ./src/tests/integration/...

# Skip integration tests
go test ./... -short

# Run specific test
go test ./src/domain/entity -run TestContest_Validation

# Run benchmarks
go test -bench=. ./...

# Race detection
go test -race ./...

# Verbose output
go test -v ./...

# Generate mocks (if using mockgen)
go generate ./...
```
