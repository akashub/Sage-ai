// backend/go/internal/knowledge/vector_test.go
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

func setupTestVectorDB(t *testing.T) (*MemoryVectorDB, func()) {
	// Create a temporary directory for the test
	tempDir, err := os.MkdirTemp("", "vector-test")
	require.NoError(t, err)

	// Create a test file path
	testPath := filepath.Join(tempDir, "test-vector-db.json")

	// Create in-memory vector DB
	db, err := CreateMemoryVectorDB(testPath)
	require.NoError(t, err)
	require.NotNil(t, db)

	// Return cleanup function
	cleanup := func() {
		db.Close()
		os.RemoveAll(tempDir)
	}

	return db, cleanup
}

func TestCreateMemoryVectorDB(t *testing.T) {
	db, cleanup := setupTestVectorDB(t)
	defer cleanup()

	assert.NotNil(t, db.items)
	assert.NotEmpty(t, db.persistPath)
}

func TestAddTrainingItem(t *testing.T) {
	db, cleanup := setupTestVectorDB(t)
	defer cleanup()
	
	ctx := context.Background()
	
	// Create a test item
	item := &TrainingItem{
		ID:          "test-item-1",
		Type:        "ddl",
		Content:     "CREATE TABLE test (id INT PRIMARY KEY)",
		Description: "Test table schema",
		Metadata: map[string]interface{}{
			"name": "test_table",
		},
		DateAdded: time.Now().Format(time.RFC3339),
	}
	
	// Add the item
	err := db.AddTrainingItem(ctx, item)
	assert.NoError(t, err)
	
	// Verify item was added
	db.mu.RLock()
	storedItem, exists := db.items[item.ID]
	db.mu.RUnlock()
	
	assert.True(t, exists)
	assert.Equal(t, item.ID, storedItem.ID)
	assert.Equal(t, item.Type, storedItem.Type)
	assert.Equal(t, item.Content, storedItem.Content)
	assert.Equal(t, item.Description, storedItem.Description)
	assert.NotEmpty(t, storedItem.Embedding, "Embedding should have been generated")
}

func TestFindSimilar(t *testing.T) {
	db, cleanup := setupTestVectorDB(t)
	defer cleanup()
	
	ctx := context.Background()
	
	// Add multiple test items
	items := []*TrainingItem{
		{
			ID:          "ddl-1",
			Type:        "ddl",
			Content:     "CREATE TABLE users (id INT, name VARCHAR(100))",
			Description: "Users table schema",
			DateAdded:   time.Now().Format(time.RFC3339),
		},
		{
			ID:          "ddl-2",
			Type:        "ddl",
			Content:     "CREATE TABLE products (id INT, name VARCHAR(100), price DECIMAL(10,2))",
			Description: "Products table schema",
			DateAdded:   time.Now().Format(time.RFC3339),
		},
		{
			ID:          "doc-1",
			Type:        "documentation",
			Content:     "User accounts are stored in the users table. Each user has an ID and name.",
			Description: "User documentation",
			DateAdded:   time.Now().Format(time.RFC3339),
		},
	}
	
	for _, item := range items {
		err := db.AddTrainingItem(ctx, item)
		assert.NoError(t, err)
	}
	
	// Test finding similar items with no type filter
	result, err := db.FindSimilar(ctx, "user accounts", "", 5)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(result), 1, "Should find at least one similar item")
	
	// Test finding similar items with type filter
	ddlResults, err := db.FindSimilar(ctx, "user accounts", "ddl", 5)
	assert.NoError(t, err)
	for _, item := range ddlResults {
		assert.Equal(t, "ddl", item.Type)
	}
	
	docResults, err := db.FindSimilar(ctx, "user accounts", "documentation", 5)
	assert.NoError(t, err)
	for _, item := range docResults {
		assert.Equal(t, "documentation", item.Type)
	}
	
	// Test limit
	limitedResults, err := db.FindSimilar(ctx, "user accounts", "", 1)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(limitedResults))
}

func TestListAll(t *testing.T) {
	db, cleanup := setupTestVectorDB(t)
	defer cleanup()
	
	ctx := context.Background()
	
	// Add multiple test items of different types
	items := []*TrainingItem{
		{
			ID:          "ddl-1",
			Type:        "ddl",
			Content:     "CREATE TABLE users (id INT)",
			Description: "Users table schema",
			DateAdded:   time.Now().Format(time.RFC3339),
		},
		{
			ID:          "doc-1",
			Type:        "documentation",
			Content:     "Documentation about users",
			Description: "User documentation",
			DateAdded:   time.Now().Format(time.RFC3339),
		},
		{
			ID:          "qa-1",
			Type:        "question_sql",
			Content:     "How many users are there?",
			Description: "User count query",
			DateAdded:   time.Now().Format(time.RFC3339),
		},
	}
	
	for _, item := range items {
		err := db.AddTrainingItem(ctx, item)
		assert.NoError(t, err)
	}
	
	// Test listing all items
	allItems, err := db.ListAll(ctx, "")
	assert.NoError(t, err)
	assert.Equal(t, len(items), len(allItems))
	
	// Test filtering by type
	ddlItems, err := db.ListAll(ctx, "ddl")
	assert.NoError(t, err)
	assert.Equal(t, 1, len(ddlItems))
	assert.Equal(t, "ddl", ddlItems[0].Type)
	
	docItems, err := db.ListAll(ctx, "documentation")
	assert.NoError(t, err)
	assert.Equal(t, 1, len(docItems))
	assert.Equal(t, "documentation", docItems[0].Type)
	
	qaItems, err := db.ListAll(ctx, "question_sql")
	assert.NoError(t, err)
	assert.Equal(t, 1, len(qaItems))
	assert.Equal(t, "question_sql", qaItems[0].Type)
}

func TestDeleteItem(t *testing.T) {
	db, cleanup := setupTestVectorDB(t)
	defer cleanup()
	
	ctx := context.Background()
	
	// Add a test item
	item := &TrainingItem{
		ID:          "test-delete",
		Type:        "ddl",
		Content:     "CREATE TABLE test (id INT)",
		Description: "Test schema",
		DateAdded:   time.Now().Format(time.RFC3339),
	}
	
	err := db.AddTrainingItem(ctx, item)
	assert.NoError(t, err)
	
	// Verify item exists
	db.mu.RLock()
	_, exists := db.items[item.ID]
	db.mu.RUnlock()
	assert.True(t, exists)
	
	// Delete the item
	err = db.DeleteItem(ctx, item.ID)
	assert.NoError(t, err)
	
	// Verify item is deleted
	db.mu.RLock()
	_, exists = db.items[item.ID]
	db.mu.RUnlock()
	assert.False(t, exists)
	
	// Test deleting non-existent item
	err = db.DeleteItem(ctx, "non-existent-id")
	assert.Error(t, err)
}

func TestPersistToDisk(t *testing.T) {
	// Create a temporary directory for the test
	tempDir, err := os.MkdirTemp("", "persist-test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)
	
	// Create a test file path
	testPath := filepath.Join(tempDir, "persist-test.json")
	
	// Create in-memory vector DB
	db, err := CreateMemoryVectorDB(testPath)
	require.NoError(t, err)
	
	ctx := context.Background()
	
	// Add a test item
	item := &TrainingItem{
		ID:          "persist-test",
		Type:        "ddl",
		Content:     "CREATE TABLE test (id INT)",
		Description: "Test schema",
		DateAdded:   time.Now().Format(time.RFC3339),
	}
	
	err = db.AddTrainingItem(ctx, item)
	assert.NoError(t, err)
	
	// Close the DB to trigger persistence
	db.Close()
	
	// Verify the file exists
	_, err = os.Stat(testPath)
	assert.NoError(t, err)
	
	// Create a new DB from the persisted file
	newDB, err := CreateMemoryVectorDB(testPath)
	assert.NoError(t, err)
	defer newDB.Close()
	
	// Verify the item was loaded
	newDB.mu.RLock()
	loadedItem, exists := newDB.items[item.ID]
	newDB.mu.RUnlock()
	
	assert.True(t, exists)
	assert.Equal(t, item.ID, loadedItem.ID)
	assert.Equal(t, item.Type, loadedItem.Type)
	assert.Equal(t, item.Content, loadedItem.Content)
	assert.Equal(t, item.Description, loadedItem.Description)
}

func TestCosineSimilarity(t *testing.T) {
	_, cleanup := setupTestVectorDB(t)
	defer cleanup()
	
	// Test identical vectors
	v1 := []float32{1.0, 2.0, 3.0}
	similarity := cosineSimilarity(v1, v1)
	assert.Equal(t, 1.0, similarity)
	
	// Test orthogonal vectors
	v2 := []float32{1.0, 0.0, 0.0}
	v3 := []float32{0.0, 1.0, 0.0}
	similarity = cosineSimilarity(v2, v3)
	assert.Equal(t, 0.0, similarity)
	
	// Test vectors of different lengths
	v4 := []float32{1.0, 2.0, 3.0, 4.0}
	v5 := []float32{1.0, 2.0}
	similarity = cosineSimilarity(v4, v5)
	assert.NotEqual(t, 0.0, similarity, "Should calculate similarity for common dimensions")
	
	// Test zero vectors
	v6 := []float32{0.0, 0.0, 0.0}
	similarity = cosineSimilarity(v1, v6)
	assert.Equal(t, 0.0, similarity, "Similarity with zero vector should be 0")
}

func TestGenerateSimpleEmbedding(t *testing.T) {
	// Test with empty text
	embedding1 := generateSimpleEmbedding("")
	assert.NotNil(t, embedding1)
	assert.Equal(t, 128, len(embedding1))
	
	// Test with simple text
	embedding2 := generateSimpleEmbedding("test")
	assert.NotNil(t, embedding2)
	assert.Equal(t, 128, len(embedding2))
	
	// Test vectors for different texts should be different
	embedding3 := generateSimpleEmbedding("completely different text")
	similarity := cosineSimilarity(embedding2, embedding3)
	assert.NotEqual(t, 1.0, similarity, "Different texts should have different embeddings")
	
	// Test vectors for similar texts should have some similarity
	embedding4 := generateSimpleEmbedding("test text")
	embedding5 := generateSimpleEmbedding("text test")
	similarity = cosineSimilarity(embedding4, embedding5)
	assert.Greater(t, similarity, 0.7, "Similar texts should have similar embeddings")
}