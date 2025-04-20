// backend/go/tests/knowledge/mock_vectordb.go
package knowledge

import (
	"context"
	"sage-ai-v2/internal/knowledge"
)

// MockVectorDB provides a testable implementation of the VectorDB interface
type MockVectorDB struct {
	items                  []knowledge.TrainingItem
	addItemCalled          bool
	findSimilarCalled      bool
	listAllCalled          bool
	deleteItemCalled       bool
	lastAddedItem          *knowledge.TrainingItem
	lastFindQuery          string
	lastFindItemType       string
	lastFindLimit          int
	lastDeleteID           string
	findSimilarResults     []knowledge.TrainingItem
	findSimilarError       error
	listAllResults         []knowledge.TrainingItem
	listAllError           error
	addItemError           error
	deleteItemError        error
}

// AddTrainingItem adds a training item to the mock database
func (m *MockVectorDB) AddTrainingItem(ctx context.Context, item *knowledge.TrainingItem) error {
	m.addItemCalled = true
	m.lastAddedItem = item
	
	if m.addItemError != nil {
		return m.addItemError
	}
	
	m.items = append(m.items, *item)
	return nil
}

// FindSimilar finds similar items to the query in the mock database
func (m *MockVectorDB) FindSimilar(ctx context.Context, query string, itemType string, limit int) ([]knowledge.TrainingItem, error) {
	m.findSimilarCalled = true
	m.lastFindQuery = query
	m.lastFindItemType = itemType
	m.lastFindLimit = limit
	
	if m.findSimilarError != nil {
		return nil, m.findSimilarError
	}
	
	// Return predefined results or filter items based on type
	if m.findSimilarResults != nil {
		return m.findSimilarResults, nil
	}
	
	var results []knowledge.TrainingItem
	for _, item := range m.items {
		if itemType == "" || item.Type == itemType {
			results = append(results, item)
			if len(results) >= limit {
				break
			}
		}
	}
	
	return results, nil
}

// ListAll lists all items in the mock database
func (m *MockVectorDB) ListAll(ctx context.Context, dataType string) ([]knowledge.TrainingItem, error) {
	m.listAllCalled = true
	
	if m.listAllError != nil {
		return nil, m.listAllError
	}
	
	// Return predefined results or filter items based on type
	if m.listAllResults != nil {
		return m.listAllResults, nil
	}
	
	var results []knowledge.TrainingItem
	for _, item := range m.items {
		if dataType == "" || item.Type == dataType {
			results = append(results, item)
		}
	}
	
	return results, nil
}

// DeleteItem deletes an item from the mock database
func (m *MockVectorDB) DeleteItem(ctx context.Context, id string) error {
	m.deleteItemCalled = true
	m.lastDeleteID = id
	
	if m.deleteItemError != nil {
		return m.deleteItemError
	}
	
	var filteredItems []knowledge.TrainingItem
	for _, item := range m.items {
		if item.ID != id {
			filteredItems = append(filteredItems, item)
		}
	}
	
	m.items = filteredItems
	return nil
}

// Close is a no-op for the mock database
func (m *MockVectorDB) Close() {
	// No-op for testing
}