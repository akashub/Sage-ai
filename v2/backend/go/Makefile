# backend/go/Makefile
.PHONY: all build test clean test-orchestrator test-api test-knowledge test-chat test-upload coverage

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get

# Main package
MAIN_PACKAGE=./cmd/server

# Binary name
BINARY_NAME=sage-ai-server

# Test flags
TEST_FLAGS=-v -race
COVERAGE_FLAGS=-coverprofile=coverage.out -covermode=atomic

# Build the main application
build:
	$(GOBUILD) -o $(BINARY_NAME) $(MAIN_PACKAGE)

# Clean up
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f coverage.out

# Run all tests
test:
	$(GOTEST) $(TEST_FLAGS) ./...

# Run orchestrator tests
test-orchestrator:
	$(GOTEST) $(TEST_FLAGS) ./tests/orchestrator/...

# Run API tests
test-api:
	$(GOTEST) $(TEST_FLAGS) ./tests/api/...

# Run knowledge tests
test-knowledge:
	$(GOTEST) $(TEST_FLAGS) ./tests/knowledge/...

# Run chat tests
test-chat:
	$(GOTEST) $(TEST_FLAGS) ./tests/api/chat_test.go

# Run upload tests
test-upload:
	$(GOTEST) $(TEST_FLAGS) ./tests/api/upload_test.go

# Run tests with coverage
coverage:
	$(GOTEST) $(COVERAGE_FLAGS) ./...
	go tool cover -html=coverage.out

# Get dependencies
deps:
	$(GOGET) -v ./...

# Run benchmarks
benchmark:
	$(GOTEST) -bench=. ./...

# For CI pipeline
ci-test: test coverage

# Create test directories if they don't exist
create-test-dirs:
	mkdir -p tests/orchestrator
	mkdir -p tests/api
	mkdir -p tests/knowledge

# Initialize go modules if needed
init-modules:
	go mod tidy

# All target - build and test
all: clean deps build test
	@echo "Build and test completed successfully."
	@echo "Run 'make coverage' to generate coverage report."
	@echo "Run 'make clean' to clean up the build."