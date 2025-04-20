// backend/go/internal/knowledge/manager_test.go
package knowledge

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestKnowledgeManager(t *testing.T) (*KnowledgeManager, func()) {
	// Create a temporary directory for the test
	tempDir, err := os.MkdirTemp("", "knowledge-test")
	require.NoError(t, err)

	// Create a mock vector DB
	vectorDB, err := CreateMemoryVectorDB(filepath.Join(tempDir, "vector_store.json"))
	require.NoError(t, err)

	// Create knowledge manager
	km := CreateKnowledgeManager(vectorDB, tempDir)
	require.NotNil(t, km)

	// Return cleanup function
	cleanup := func() {
		os.RemoveAll(tempDir)
	}

	return km, cleanup
}

func TestCreateKnowledgeManager(t *testing.T) {
	km, cleanup := setupTestKnowledgeManager(t)
	defer cleanup()

	assert.NotNil(t, km.vectorDB)
	assert.Equal(t, km.fileStoragePath, km.fileStoragePath)
}

func TestStoreFile(t *testing.T) {
	km, cleanup := setupTestKnowledgeManager(t)
	defer cleanup()

	// Test storing a file
	content := []byte("Test content for file storage")
	filePath, err := km.StoreFile("test.txt", content)
	
	assert.NoError(t, err)
	assert.NotEmpty(t, filePath)
	
	// Verify file exists
	storedContent, err := os.ReadFile(filePath)
	assert.NoError(t, err)
	assert.Equal(t, content, storedContent)
}

func TestAddDDLSchema(t *testing.T) {
	km, cleanup := setupTestKnowledgeManager(t)
	defer cleanup()
	
	ctx := context.Background()
	schemaName := "test_schema"
	schemaContent := "CREATE TABLE users (id INT PRIMARY KEY, name VARCHAR(100))"
	description := "Test schema for users table"
	
	// Add DDL schema
	err := km.AddDDLSchema(ctx, schemaName, schemaContent, description)
	assert.NoError(t, err)
	
	// Retrieve and verify
	result, err := km.RetrieveRelevantKnowledge(ctx, "users table")
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(result.DDLSchemas), 1)
	
	// Verify content
	found := false
	for _, schema := range result.DDLSchemas {
		if schema.Description == description {
			assert.Equal(t, schemaContent, schema.Content)
			found = true
			break
		}
	}
	assert.True(t, found, "Added schema was not found in retrieval results")
}

func TestAddDocumentation(t *testing.T) {
	km, cleanup := setupTestKnowledgeManager(t)
	defer cleanup()
	
	ctx := context.Background()
	title := "User Registration Process"
	content := "When a user registers, we validate email and store their information in the users table."
	tags := []string{"users", "registration"}
	
	// Add documentation
	err := km.AddDocumentation(ctx, title, content, tags)
	assert.NoError(t, err)
	
	// Retrieve and verify
	result, err := km.RetrieveRelevantKnowledge(ctx, "user registration")
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(result.Documentation), 1)
	
	// Verify content
	found := false
	for _, doc := range result.Documentation {
		if doc.Description == title {
			assert.Equal(t, content, doc.Content)
			found = true
			break
		}
	}
	assert.True(t, found, "Added documentation was not found in retrieval results")
}

func TestAddQuestionSQLPair(t *testing.T) {
	km, cleanup := setupTestKnowledgeManager(t)
	defer cleanup()
	
	ctx := context.Background()
	
	// Create a question-SQL pair
	pair := QuestionSQLPair{
		Question:    "How many users registered last month?",
		SQL:         "SELECT COUNT(*) FROM users WHERE created_at >= DATE_SUB(NOW(), INTERVAL 1 MONTH)",
		Description: "Monthly user registration count",
		Tags:        []string{"users", "registration", "reporting"},
		DateAdded:   time.Now().Format(time.RFC3339),
		Verified:    true,
	}
	
	// Add the pair
	err := km.AddQuestionSQLPair(ctx, pair)
	assert.NoError(t, err)
	
	// Retrieve and verify
	result, err := km.RetrieveRelevantKnowledge(ctx, "users registered last month")
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(result.QuestionSQLPairs), 1)
	
	// Verify content
	found := false
	for _, resultPair := range result.QuestionSQLPairs {
		if resultPair.Description == pair.Description {
			assert.Equal(t, pair.Question, resultPair.Question)
			assert.Equal(t, pair.SQL, resultPair.SQL)
			found = true
			break
		}
	}
	assert.True(t, found, "Added question-SQL pair was not found in retrieval results")
}

func TestListTrainingData(t *testing.T) {
	km, cleanup := setupTestKnowledgeManager(t)
	defer cleanup()
	
	ctx := context.Background()
	
	// Add some test data
	err := km.AddDDLSchema(ctx, "users_schema", "CREATE TABLE users (id INT)", "Users table schema")
	assert.NoError(t, err)
	
	err = km.AddDocumentation(ctx, "User Guide", "User documentation content", []string{"user"})
	assert.NoError(t, err)
	
	// List all training data
	items, err := km.ListTrainingData(ctx, "")
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(items), 2, "Should list at least 2 items")
	
	// Test filtering by type
	ddlItems, err := km.ListTrainingData(ctx, "ddl")
	assert.NoError(t, err)
	for _, item := range ddlItems {
		assert.Equal(t, "ddl", item["type"])
	}
	
	docItems, err := km.ListTrainingData(ctx, "documentation")
	assert.NoError(t, err)
	for _, item := range docItems {
		assert.Equal(t, "documentation", item["type"])
	}
}

func TestDeleteTrainingItem(t *testing.T) {
	km, cleanup := setupTestKnowledgeManager(t)
	defer cleanup()
	
	ctx := context.Background()
	
	// Add a test item
	err := km.AddDDLSchema(ctx, "temp_schema", "CREATE TABLE temp (id INT)", "Temporary schema")
	assert.NoError(t, err)
	
	// List items to get the ID
	items, err := km.ListTrainingData(ctx, "ddl")
	assert.NoError(t, err)
	assert.NotEmpty(t, items)
	
	// Get the ID of the first item
	id, ok := items[0]["id"].(string)
	assert.True(t, ok, "ID should be a string")
	assert.NotEmpty(t, id)
	
	// Delete the item
	err = km.DeleteTrainingItem(ctx, id)
	assert.NoError(t, err)
	
	// Verify it's deleted
	items, err = km.ListTrainingData(ctx, "")
	assert.NoError(t, err)
	
	// Check the item is no longer in the list
	for _, item := range items {
		itemID, _ := item["id"].(string)
		assert.NotEqual(t, id, itemID, "Item should have been deleted")
	}
}

func TestGetTrainingItem(t *testing.T) {
	km, cleanup := setupTestKnowledgeManager(t)
	defer cleanup()
	
	ctx := context.Background()
	
	// Add a test item
	description := "Test item for retrieval"
	content := "This is test content for retrieval testing"
	err := km.AddDocumentation(ctx, description, content, []string{"test"})
	assert.NoError(t, err)
	
	// List items to get the ID
	items, err := km.ListTrainingData(ctx, "documentation")
	assert.NoError(t, err)
	assert.NotEmpty(t, items)
	
	// Get the ID of the first item
	id, ok := items[0]["id"].(string)
	assert.True(t, ok, "ID should be a string")
	assert.NotEmpty(t, id)
	
	// Get the full item
	item, err := km.GetTrainingItem(ctx, id)
	assert.NoError(t, err)
	assert.NotNil(t, item)
	assert.Equal(t, id, item.ID)
	assert.Equal(t, "documentation", item.Type)
	assert.Equal(t, description, item.Description)
	assert.Equal(t, content, item.Content)
}

func TestRetrieveRelevantKnowledge(t *testing.T) {
	km, cleanup := setupTestKnowledgeManager(t)
	defer cleanup()
	
	ctx := context.Background()
	
	// Add various types of knowledge
	err := km.AddDDLSchema(ctx, "products_schema", 
		"CREATE TABLE products (id INT, name VARCHAR(100), price DECIMAL(10,2))", 
		"Products table schema")
	assert.NoError(t, err)
	
	err = km.AddDocumentation(ctx, "Product Pricing", 
		"Product prices are set based on cost plus markup. Special discounts are available.", 
		[]string{"products", "pricing"})
	assert.NoError(t, err)
	
	pair := QuestionSQLPair{
		Question:    "What are the top 5 most expensive products?",
		SQL:         "SELECT name, price FROM products ORDER BY price DESC LIMIT 5",
		Description: "Top expensive products query",
		Tags:        []string{"products", "pricing", "reporting"},
		DateAdded:   time.Now().Format(time.RFC3339),
		Verified:    true,
	}
	err = km.AddQuestionSQLPair(ctx, pair)
	assert.NoError(t, err)
	
	// Test retrieval with a query that should match all three types
	result, err := km.RetrieveRelevantKnowledge(ctx, "expensive product prices")
	assert.NoError(t, err)
	
	// Should find all types of knowledge
	assert.NotEmpty(t, result.DDLSchemas, "Should find DDL schemas")
	assert.NotEmpty(t, result.Documentation, "Should find documentation")
	assert.NotEmpty(t, result.QuestionSQLPairs, "Should find question-SQL pairs")
	
	// Test retrieval with more specific query
	result, err = km.RetrieveRelevantKnowledge(ctx, "top 5 most expensive products")
	assert.NoError(t, err)
	
	// Should strongly match the question-SQL pair
	assert.NotEmpty(t, result.QuestionSQLPairs, "Should find question-SQL pairs")
	assert.Equal(t, pair.Question, result.QuestionSQLPairs[0].Question)
	assert.Equal(t, pair.SQL, result.QuestionSQLPairs[0].SQL)
}