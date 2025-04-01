// // backend/go/tests/knowledge/manager_test.go
// package knowledge_test

// import (
// 	"context"
// 	"os"
// 	"sage-ai-v2/internal/knowledge"
// 	"testing"
// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/require"
// )

// // TestKnowledgeManagerWithMockDB tests the knowledge manager with a mock vector database
// func TestKnowledgeManagerWithMockDB(t *testing.T) {
// 	// Create temporary directory for testing
// 	tempDir, err := os.MkdirTemp("", "knowledge-test")
// 	require.NoError(t, err)
// 	defer os.RemoveAll(tempDir)

// 	// Create mock vector DB
// 	mockDB := &MockVectorDB{
// 		items: []knowledge.TrainingItem{},
// 	}

// 	// Create knowledge manager with mock DB
// 	km := knowledge.CreateKnowledgeManager(mockDB, tempDir)
// 	require.NotNil(t, km)

// 	// Test AddDDLSchema
// 	t.Run("AddDDLSchema", func(t *testing.T) {
// 		ctx := context.Background()
// 		schemaName := "test_schema"
// 		schemaContent := "CREATE TABLE users (id INT PRIMARY KEY, name VARCHAR(100))"
// 		description := "Test schema for users table"
		
// 		err := km.AddDDLSchema(ctx, schemaName, schemaContent, description)
// 		assert.NoError(t, err)
// 		assert.True(t, mockDB.addItemCalled)
// 		assert.Equal(t, "ddl", mockDB.lastAddedItem.Type)
// 		assert.Equal(t, schemaContent, mockDB.lastAddedItem.Content)
// 		assert.Equal(t, description, mockDB.lastAddedItem.Description)
// 	})

// 	// Test AddDocumentation
// 	t.Run("AddDocumentation", func(t *testing.T) {
// 		ctx := context.Background()
// 		title := "User Registration Process"
// 		content := "When a user registers, we validate email and store their information in the users table."
// 		tags := []string{"users", "registration"}
		
// 		err := km.AddDocumentation(ctx, title, content, tags)
// 		assert.NoError(t, err)
// 		assert.True(t, mockDB.addItemCalled)
// 		assert.Equal(t, "documentation", mockDB.lastAddedItem.Type)
// 		assert.Equal(t, content, mockDB.lastAddedItem.Content)
// 		assert.Equal(t, title, mockDB.lastAddedItem.Description)
		
// 		// Verify tags are in metadata
// 		metadataTags, ok := mockDB.lastAddedItem.Metadata["tags"].([]string)
// 		assert.True(t, ok)
// 		assert.Equal(t, tags, metadataTags)
// 	})

// 	// Test AddQuestionSQLPair
// 	t.Run("AddQuestionSQLPair", func(t *testing.T) {
// 		ctx := context.Background()
// 		pair := knowledge.QuestionSQLPair{
// 			Question:    "How many users registered last month?",
// 			SQL:         "SELECT COUNT(*) FROM users WHERE created_at >= DATE_SUB(NOW(), INTERVAL 1 MONTH)",
// 			Description: "Monthly user registration count",
// 			Tags:        []string{"users", "registration", "reporting"},
// 			DateAdded:   "2025-03-01T00:00:00Z",
// 			Verified:    true,
// 		}
		
// 		err := km.AddQuestionSQLPair(ctx, pair)
// 		assert.NoError(t, err)
// 		assert.True(t, mockDB.addItemCalled)
// 		assert.Equal(t, "question_sql", mockDB.lastAddedItem.Type)
// 		assert.Equal(t, pair.Question, mockDB.lastAddedItem.Content)
// 		assert.Equal(t, pair.Description, mockDB.lastAddedItem.Description)
// 	})

// 	// Test RetrieveRelevantKnowledge
// 	t.Run("RetrieveRelevantKnowledge", func(t *testing.T) {
// 		ctx := context.Background()
		
// 		// Set up mock results
// 		mockDB.findSimilarResults = []knowledge.TrainingItem{
// 			{
// 				ID:          "ddl_123",
// 				Type:        "ddl",
// 				Content:     "CREATE TABLE users (id INT PRIMARY KEY, name VARCHAR(100))",
// 				Description: "Users table schema",
// 			},
// 			{
// 				ID:          "doc_456",
// 				Type:        "documentation",
// 				Content:     "User documentation for registration process",
// 				Description: "User Registration",
// 			},
// 		}
		
// 		// Also prepare a question-SQL pair (these are handled specially in the manager)
// 		pairItem := knowledge.TrainingItem{
// 			ID:          "qa_789",
// 			Type:        "question_sql",
// 			Content:     "How many users registered last month?",
// 			Description: "Monthly user registration count",
// 			Metadata: map[string]interface{}{
// 				"pair": `{"question":"How many users registered last month?","sql":"SELECT COUNT(*) FROM users WHERE created_at >= DATE_SUB(NOW(), INTERVAL 1 MONTH)","description":"Monthly user registration count"}`,
// 			},
// 		}
// 		mockDB.items = append(mockDB.items, pairItem)
		
// 		// Query for relevant knowledge
// 		result, err := km.RetrieveRelevantKnowledge(ctx, "users registered last month")
// 		assert.NoError(t, err)
// 		assert.True(t, mockDB.findSimilarCalled)
		
// 		// Verify results
// 		assert.Len(t, result.DDLSchemas, 1)
// 		assert.Equal(t, "CREATE TABLE users (id INT PRIMARY KEY, name VARCHAR(100))", result.DDLSchemas[0].Content)
		
// 		assert.Len(t, result.Documentation, 1)
// 		assert.Equal(t, "User documentation for registration process", result.Documentation[0].Content)
		
// 		// Reset the mock for question-SQL pairs specific query
// 		mockDB.findSimilarCalled = false
// 		mockDB.findSimilarResults = []knowledge.TrainingItem{pairItem}
		
// 		// RetrieveRelevantKnowledge will call FindSimilar for each type, so test with a different query
// 		result, err = km.RetrieveRelevantKnowledge(ctx, "monthly registration")
// 		assert.NoError(t, err)
// 		assert.True(t, mockDB.findSimilarCalled)
		
// 		// Check that the question-SQL pair was processed correctly
// 		assert.Len(t, result.QuestionSQLPairs, 1)
// 		assert.Equal(t, "How many users registered last month?", result.QuestionSQLPairs[0].Question)
// 		assert.Equal(t, "SELECT COUNT(*) FROM users WHERE created_at >= DATE_SUB(NOW(), INTERVAL 1 MONTH)", result.QuestionSQLPairs[0].SQL)
// 	})

// 	// Test ListTrainingData
// 	t.Run("ListTrainingData", func(t *testing.T) {
// 		ctx := context.Background()
		
// 		// Set up mock results
// 		mockDB.listAllResults = []knowledge.TrainingItem{
// 			{
// 				ID:          "ddl_123",
// 				Type:        "ddl",
// 				Content:     "CREATE TABLE users (id INT PRIMARY KEY, name VARCHAR(100))",
// 				Description: "Users table schema",
// 				DateAdded:   "2025-03-01T00:00:00Z",
// 			},
// 			{
// 				ID:          "doc_456",
// 				Type:        "documentation",
// 				Content:     "User documentation for registration process",
// 				Description: "User Registration",
// 				DateAdded:   "2025-03-02T00:00:00Z",
// 			},
// 		}
		
// 		// Get all training data
// 		items, err := km.ListTrainingData(ctx, "")
// 		assert.NoError(t, err)
// 		assert.True(t, mockDB.listAllCalled)
// 		assert.Len(t, items, 2)
		
// 		// Check the format of the returned items
// 		assert.Equal(t, "ddl_123", items[0]["id"])
// 		assert.Equal(t, "ddl", items[0]["type"])
// 		assert.Equal(t, "Users table schema", items[0]["description"])
// 		assert.Equal(t, "2025-03-01T00:00:00Z", items[0]["date_added"])
		
// 		// Reset mock and test with type filter
// 		mockDB.listAllCalled = false
// 		mockDB.listAllResults = []knowledge.TrainingItem{
// 			{
// 				ID:          "ddl_123",
// 				Type:        "ddl",
// 				Content:     "CREATE TABLE users (id INT PRIMARY KEY, name VARCHAR(100))",
// 				Description: "Users table schema",
// 				DateAdded:   "2025-03-01T00:00:00Z",
// 			},
// 		}
		
// 		// Get only DDL schemas
// 		items, err = km.ListTrainingData(ctx, "ddl")
// 		assert.NoError(t, err)
// 		assert.True(t, mockDB.listAllCalled)
// 		assert.Len(t, items, 1)
// 		assert.Equal(t, "ddl", items[0]["type"])
// 	})

// 	// Test DeleteTrainingItem
// 	t.Run("DeleteTrainingItem", func(t *testing.T) {
// 		ctx := context.Background()
		
// 		// Delete an item
// 		err := km.DeleteTrainingItem(ctx, "ddl_123")
// 		assert.NoError(t, err)
// 		assert.True(t, mockDB.deleteItemCalled)
// 		assert.Equal(t, "ddl_123", mockDB.lastDeleteID)
		
// 		// Test with error
// 		mockDB.deleteItemError = assert.AnError
// 		err = km.DeleteTrainingItem(ctx, "ddl_456")
// 		assert.Error(t, err)
// 		assert.Equal(t, assert.AnError, err)
// 		assert.True(t, mockDB.deleteItemCalled)
// 		assert.Equal(t, "ddl_456", mockDB.lastDeleteID)
// 	})

// 	// Test GetTrainingItem
// 	t.Run("GetTrainingItem", func(t *testing.T) {
// 		ctx := context.Background()
		
// 		// Set up mock data
// 		testItem := knowledge.TrainingItem{
// 			ID:          "ddl_123",
// 			Type:        "ddl",
// 			Content:     "CREATE TABLE users (id INT PRIMARY KEY, name VARCHAR(100))",
// 			Description: "Users table schema",
// 			DateAdded:   "2025-03-01T00:00:00Z",
// 		}
// 		mockDB.items = []knowledge.TrainingItem{testItem}
		
// 		// Reset the mock and prepare list results to simulate finding the item
// 		mockDB.listAllCalled = false
// 		mockDB.listAllResults = []knowledge.TrainingItem{testItem}
		
// 		// Get the training item
// 		item, err := km.GetTrainingItem(ctx, "ddl_123")
// 		assert.NoError(t, err)
// 		assert.True(t, mockDB.listAllCalled)
// 		assert.Equal(t, "ddl_123", item.ID)
// 		assert.Equal(t, "ddl", item.Type)
// 		assert.Equal(t, "CREATE TABLE users (id INT PRIMARY KEY, name VARCHAR(100))", item.Content)
// 		assert.Equal(t, "Users table schema", item.Description)
		
// 		// Test with non-existent item
// 		mockDB.listAllResults = []knowledge.TrainingItem{}
// 		item, err = km.GetTrainingItem(ctx, "non_existent")
// 		assert.Error(t, err)
// 		assert.Nil(t, item)
// 	})
// }

// backend/go/tests/knowledge/manager_test.go
package knowledge

import (
	"context"
	"os"
	"sage-ai-v2/internal/knowledge"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestKnowledgeManagerWithMockDB tests the knowledge manager with a mock vector database
func TestKnowledgeManagerWithMockDB(t *testing.T) {
	// Create temporary directory for testing
	tempDir, err := os.MkdirTemp("", "knowledge-test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create mock vector DB
	mockDB := &MockVectorDB{
		items: []knowledge.TrainingItem{},
	}

	// Create knowledge manager with mock DB
	km := knowledge.CreateKnowledgeManager(mockDB, tempDir)
	require.NotNil(t, km)

	// Test AddDDLSchema
	t.Run("AddDDLSchema", func(t *testing.T) {
		ctx := context.Background()
		schemaName := "test_schema"
		schemaContent := "CREATE TABLE users (id INT PRIMARY KEY, name VARCHAR(100))"
		description := "Test schema for users table"
		
		err := km.AddDDLSchema(ctx, schemaName, schemaContent, description)
		assert.NoError(t, err)
		assert.True(t, mockDB.addItemCalled)
		assert.Equal(t, "ddl", mockDB.lastAddedItem.Type)
		assert.Equal(t, schemaContent, mockDB.lastAddedItem.Content)
		assert.Equal(t, description, mockDB.lastAddedItem.Description)
	})

	// Test AddDocumentation
	t.Run("AddDocumentation", func(t *testing.T) {
		ctx := context.Background()
		title := "User Registration Process"
		content := "When a user registers, we validate email and store their information in the users table."
		tags := []string{"users", "registration"}
		
		err := km.AddDocumentation(ctx, title, content, tags)
		assert.NoError(t, err)
		assert.True(t, mockDB.addItemCalled)
		assert.Equal(t, "documentation", mockDB.lastAddedItem.Type)
		assert.Equal(t, content, mockDB.lastAddedItem.Content)
		assert.Equal(t, title, mockDB.lastAddedItem.Description)
		
		// Verify tags are in metadata
		metadataTags, ok := mockDB.lastAddedItem.Metadata["tags"].([]string)
		assert.True(t, ok)
		assert.Equal(t, tags, metadataTags)
	})

	// Test AddQuestionSQLPair
	t.Run("AddQuestionSQLPair", func(t *testing.T) {
		ctx := context.Background()
		pair := knowledge.QuestionSQLPair{
			Question:    "How many users registered last month?",
			SQL:         "SELECT COUNT(*) FROM users WHERE created_at >= DATE_SUB(NOW(), INTERVAL 1 MONTH)",
			Description: "Monthly user registration count",
			Tags:        []string{"users", "registration", "reporting"},
			DateAdded:   "2025-03-01T00:00:00Z",
			Verified:    true,
		}
		
		err := km.AddQuestionSQLPair(ctx, pair)
		assert.NoError(t, err)
		assert.True(t, mockDB.addItemCalled)
		assert.Equal(t, "question_sql", mockDB.lastAddedItem.Type)
		assert.Equal(t, pair.Question, mockDB.lastAddedItem.Content)
		assert.Equal(t, pair.Description, mockDB.lastAddedItem.Description)
	})

	// Test RetrieveRelevantKnowledge
	t.Run("RetrieveRelevantKnowledge", func(t *testing.T) {
		ctx := context.Background()
		
		// Set up mock results
		mockDB.findSimilarResults = []knowledge.TrainingItem{
			{
				ID:          "ddl_123",
				Type:        "ddl",
				Content:     "CREATE TABLE users (id INT PRIMARY KEY, name VARCHAR(100))",
				Description: "Users table schema",
			},
			{
				ID:          "doc_456",
				Type:        "documentation",
				Content:     "User documentation for registration process",
				Description: "User Registration",
			},
		}
		
		// Also prepare a question-SQL pair (these are handled specially in the manager)
		pairItem := knowledge.TrainingItem{
			ID:          "qa_789",
			Type:        "question_sql",
			Content:     "How many users registered last month?",
			Description: "Monthly user registration count",
			Metadata: map[string]interface{}{
				"pair": `{"question":"How many users registered last month?","sql":"SELECT COUNT(*) FROM users WHERE created_at >= DATE_SUB(NOW(), INTERVAL 1 MONTH)","description":"Monthly user registration count"}`,
			},
		}
		mockDB.items = append(mockDB.items, pairItem)
		
		// Query for relevant knowledge
		result, err := km.RetrieveRelevantKnowledge(ctx, "users registered last month")
		assert.NoError(t, err)
		assert.True(t, mockDB.findSimilarCalled)
		
		// Verify results
		assert.Len(t, result.DDLSchemas, 1)
		assert.Equal(t, "CREATE TABLE users (id INT PRIMARY KEY, name VARCHAR(100))", result.DDLSchemas[0].Content)
		
		assert.Len(t, result.Documentation, 1)
		assert.Equal(t, "User documentation for registration process", result.Documentation[0].Content)
		
		// Reset the mock for question-SQL pairs specific query
		mockDB.findSimilarCalled = false
		mockDB.findSimilarResults = []knowledge.TrainingItem{pairItem}
		
		// RetrieveRelevantKnowledge will call FindSimilar for each type, so test with a different query
		result, err = km.RetrieveRelevantKnowledge(ctx, "monthly registration")
		assert.NoError(t, err)
		assert.True(t, mockDB.findSimilarCalled)
		
		// Check that the question-SQL pair was processed correctly
		assert.Len(t, result.QuestionSQLPairs, 1)
		assert.Equal(t, "How many users registered last month?", result.QuestionSQLPairs[0].Question)
		assert.Equal(t, "SELECT COUNT(*) FROM users WHERE created_at >= DATE_SUB(NOW(), INTERVAL 1 MONTH)", result.QuestionSQLPairs[0].SQL)
	})

	// Test ListTrainingData
	t.Run("ListTrainingData", func(t *testing.T) {
		ctx := context.Background()
		
		// Set up mock results
		mockDB.listAllResults = []knowledge.TrainingItem{
			{
				ID:          "ddl_123",
				Type:        "ddl",
				Content:     "CREATE TABLE users (id INT PRIMARY KEY, name VARCHAR(100))",
				Description: "Users table schema",
				DateAdded:   "2025-03-01T00:00:00Z",
			},
			{
				ID:          "doc_456",
				Type:        "documentation",
				Content:     "User documentation for registration process",
				Description: "User Registration",
				DateAdded:   "2025-03-02T00:00:00Z",
			},
		}
		
		// Get all training data
		items, err := km.ListTrainingData(ctx, "")
		assert.NoError(t, err)
		assert.True(t, mockDB.listAllCalled)
		assert.Len(t, items, 2)
		
		// Check the format of the returned items
		assert.Equal(t, "ddl_123", items[0]["id"])
		assert.Equal(t, "ddl", items[0]["type"])
		assert.Equal(t, "Users table schema", items[0]["description"])
		assert.Equal(t, "2025-03-01T00:00:00Z", items[0]["date_added"])
		
		// Reset mock and test with type filter
		mockDB.listAllCalled = false
		mockDB.listAllResults = []knowledge.TrainingItem{
			{
				ID:          "ddl_123",
				Type:        "ddl",
				Content:     "CREATE TABLE users (id INT PRIMARY KEY, name VARCHAR(100))",
				Description: "Users table schema",
				DateAdded:   "2025-03-01T00:00:00Z",
			},
		}
		
		// Get only DDL schemas
		items, err = km.ListTrainingData(ctx, "ddl")
		assert.NoError(t, err)
		assert.True(t, mockDB.listAllCalled)
		assert.Len(t, items, 1)
		assert.Equal(t, "ddl", items[0]["type"])
	})

	// Test DeleteTrainingItem
	t.Run("DeleteTrainingItem", func(t *testing.T) {
		ctx := context.Background()
		
		// Delete an item
		err := km.DeleteTrainingItem(ctx, "ddl_123")
		assert.NoError(t, err)
		assert.True(t, mockDB.deleteItemCalled)
		assert.Equal(t, "ddl_123", mockDB.lastDeleteID)
		
		// Test with error
		mockDB.deleteItemError = assert.AnError
		err = km.DeleteTrainingItem(ctx, "ddl_456")
		assert.Error(t, err)
		assert.Equal(t, assert.AnError, err)
		assert.True(t, mockDB.deleteItemCalled)
		assert.Equal(t, "ddl_456", mockDB.lastDeleteID)
	})

	// Test GetTrainingItem
	t.Run("GetTrainingItem", func(t *testing.T) {
		ctx := context.Background()
		
		// Set up mock data
		testItem := knowledge.TrainingItem{
			ID:          "ddl_123",
			Type:        "ddl",
			Content:     "CREATE TABLE users (id INT PRIMARY KEY, name VARCHAR(100))",
			Description: "Users table schema",
			DateAdded:   "2025-03-01T00:00:00Z",
		}
		mockDB.items = []knowledge.TrainingItem{testItem}
		
		// Reset the mock and prepare list results to simulate finding the item
		mockDB.listAllCalled = false
		mockDB.listAllResults = []knowledge.TrainingItem{testItem}
		
		// Get the training item
		item, err := km.GetTrainingItem(ctx, "ddl_123")
		assert.NoError(t, err)
		assert.True(t, mockDB.listAllCalled)
		assert.Equal(t, "ddl_123", item.ID)
		assert.Equal(t, "ddl", item.Type)
		assert.Equal(t, "CREATE TABLE users (id INT PRIMARY KEY, name VARCHAR(100))", item.Content)
		assert.Equal(t, "Users table schema", item.Description)
		
		// Test with non-existent item
		mockDB.listAllResults = []knowledge.TrainingItem{}
		item, err = km.GetTrainingItem(ctx, "non_existent")
		assert.Error(t, err)
		assert.Nil(t, item)
	})
}