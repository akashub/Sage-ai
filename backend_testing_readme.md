# Testing the Sage AI Backend

This directory contains test files for the Sage AI backend, which uses Go for the main server and Python for the LLM service.

## Overview

The tests are organized by component:

- `tests/orchestrator/`: Tests for the orchestration layer
- `tests/api/`: Tests for API handlers
- `tests/knowledge/`: Tests for the knowledge management system

## Prerequisites

- Go 1.19 or later
- The `testify` package for assertions: `go get github.com/stretchr/testify`

## Running Tests

You can use the provided Makefile to run tests for specific components or all tests at once.

### Run All Tests

```bash
make test
```

### Run Specific Component Tests

```bash
# Test orchestrator
make test-orchestrator

# Test API handlers
make test-api

# Test knowledge management
make test-knowledge

# Test chat functionality
make test-chat

# Test file upload functionality
make test-upload
```

### Run Tests with Coverage

```bash
make coverage
```

This will generate a coverage report and open it in your browser.

## Test Structure

### Mock Implementation Approach

We use a custom mocking approach to avoid compatibility issues with interface types:

1. Instead of trying to create mock implementations of interfaces like `*llm.Bridge` that are used by the real code, we create stand-alone mock orchestrators that implement the same public methods.

2. We test components in isolation using these mocks to avoid complex dependencies.

### Test Helpers

The `tests/api/helpers.go` file contains helper functions and types for testing:

- `TestServer`: A test HTTP server with mocked dependencies
- `MockOrchestrator`: A mock implementation of the orchestrator
- `MockKnowledgeStore`: A simplified in-memory store for testing
- `SetupTestEnvironment`: Prepares a test environment with all necessary mocks

## Adding New Tests

When adding new tests, follow these guidelines:

1. Create test files in the appropriate directory (`tests/orchestrator/`, `tests/api/`, or `tests/knowledge/`).
2. Use the existing mock implementations or create new ones as needed.
3. Follow the pattern of existing tests to maintain consistency.
4. Make sure to clean up any resources created during tests (e.g., temporary files).

## Troubleshooting

### Common Issues

1. **Type Compatibility Errors**: If you see errors like "cannot use mockClient as *llm.Bridge value", use our approach of creating standalone mocks instead of trying to implement the exact interfaces.

2. **Missing Dependencies**: Make sure you have all required Go packages:

   ```bash
   go mod tidy
   ```

3. **Database Errors**: By default, tests use in-memory structures and don't require a real database.

### Debugging Tests

To run tests with verbose output:

```bash
go test -v ./tests/...
```

To run a specific test:

```bash
go test -v ./tests/api/... -run TestQueryHandler
```

## Notes for Writing Good Tests

1. **Isolate tests**: Each test should run independently of others.
2. **Clean up resources**: Use `defer` to clean up resources created during tests.
3. **Use meaningful assertions**: Assertions should verify specific behaviors, not just check if code runs.
4. **Mock external dependencies**: Don't rely on external services in unit tests.
5. **Test error cases**: Don't just test the happy path; also test error handling.


# Test Suite Documentation for Sage AI

This document provides a comprehensive overview of the test suites implemented for the Sage AI backend system. These tests ensure the reliability, correctness, and robustness of the system's components.

## 1. Orchestrator Tests

Located in: `tests/orchestrator/orchestrator_test.go`

The orchestrator is a central component that coordinates the flow of data processing through the various system components. These tests verify that the orchestrator correctly manages this process.

### Test Cases:

#### `TestOrchestratorProcessQueryWithCustomMock`
- **Purpose**: Tests the basic query processing flow using a custom mocked implementation
- **What it verifies**:
  - The orchestrator correctly processes a simple query
  - The generated SQL matches the expected output
  - No errors occur during normal processing

#### `TestOrchestratorProcessQueryWithKnowledge`
- **Purpose**: Tests the knowledge-enhanced query flow
- **What it verifies**:
  - The orchestrator correctly incorporates knowledge base information
  - DDL schemas from the knowledge base are properly utilized
  - Knowledge context is properly populated in the response

#### `TestNewSession`
- **Purpose**: Tests the session creation functionality
- **What it verifies**:
  - The NewSession method is called correctly
  - A new session can be initialized

## 2. Knowledge Manager Tests

Located in: `tests/knowledge/manager_test.go`

The knowledge manager is responsible for storing, retrieving, and managing the system's knowledge base, which includes DDL schemas, documentation, and question-SQL pairs.

### Test Cases:

#### `TestKnowledgeManagerWithMockDB`
- **Purpose**: Tests the knowledge manager with a mock vector database
- **Contains subtests for**:

##### `AddDDLSchema`
- **Purpose**: Tests adding DDL schemas to the knowledge base
- **What it verifies**:
  - The schema is correctly added to the database
  - The content, type, and description are stored correctly

##### `AddDocumentation`
- **Purpose**: Tests adding documentation to the knowledge base
- **What it verifies**:
  - Documentation is correctly stored
  - Metadata like tags are properly associated

##### `AddQuestionSQLPair`
- **Purpose**: Tests adding question-SQL pairs to the knowledge base
- **What it verifies**:
  - The question and SQL are correctly stored
  - The pair can be retrieved for future use

##### `RetrieveRelevantKnowledge`
- **Purpose**: Tests retrieving knowledge relevant to a query
- **What it verifies**:
  - The system can find schemas, documentation, and question-SQL pairs related to a query
  - The vector similarity search works correctly

##### `ListTrainingData`
- **Purpose**: Tests listing all training data or filtered by type
- **What it verifies**:
  - All data can be retrieved
  - Filtering by type works correctly

##### `DeleteTrainingItem`
- **Purpose**: Tests deleting items from the knowledge base
- **What it verifies**:
  - Items can be deleted by ID
  - Error handling works correctly

##### `GetTrainingItem`
- **Purpose**: Tests retrieving a specific item by ID
- **What it verifies**:
  - Items can be retrieved by ID
  - Error handling for non-existent items

## 3. API Tests

Located in: `tests/api/`

These tests verify that the HTTP API endpoints function correctly, handling requests and responses as expected.

### Chat API Tests

Located in: `tests/api/chat_test.go`

#### `TestChatStore`
- **Purpose**: Tests the basic functionality of the chat store
- **What it verifies**:
  - Chats can be created, retrieved, updated, and deleted
  - Chat messages and metadata are stored correctly

#### `TestChatRoutes`
- **Purpose**: Tests the chat API routes
- **Contains subtests for**:
  - `CreateChat`: Tests creating a new chat
  - `GetChats`: Tests retrieving a list of chats
  - `GetChat`: Tests retrieving a specific chat by ID
  - `UpdateChat`: Tests updating an existing chat
  - `DeleteChat`: Tests deleting a chat

### Handler Tests

Located in: `tests/api/handlers_test.go`

#### `TestQueryHandlerDirectly`
- **Purpose**: Tests the query handler directly with HTTP requests
- **What it verifies**:
  - The handler correctly processes query requests
  - Responses contain SQL, results, and other expected fields

#### `TestUploadHandler`
- **Purpose**: Tests the file upload endpoint
- **What it verifies**:
  - CSV files can be uploaded
  - Headers are correctly extracted
  - Success responses are properly formatted

## 4. Mock Implementations

### MockVectorDB

Located in: `tests/knowledge/mock_vectordb.go`

This mock implementation provides a testable version of the vector database, which is used for storing and retrieving knowledge items.

**Features**:
- Tracks method calls for verification
- Stores items in memory
- Allows predefined results to be returned
- Simulates errors for testing error handling

### MockOrchestrator

Located in: `tests/orchestrator/orchestrator_test.go`

This mock implementation provides a testable version of the orchestrator, which coordinates the flow of data processing.

**Features**:
- Allows customizing behavior through function fields
- Tracks method calls for verification
- Returns predefined results for testing different scenarios

## How Tests Work Together

The test suite is designed to test components at different levels of abstraction:

1. **Unit Tests**: Test individual functions and methods in isolation
2. **Integration Tests**: Test interactions between components
3. **API Tests**: Test HTTP endpoints and responses

By using mock implementations, the tests can focus on specific components without requiring the entire system to be operational. This allows for more targeted testing and easier identification of issues.

## Running the Tests

Tests can be run using the Go test command:

```bash
# Run all tests
go test ./tests/...

# Run specific test suites
go test ./tests/orchestrator/...
go test ./tests/knowledge/...
go test ./tests/api/...