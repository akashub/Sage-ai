// backend/go/internal/knowledge/manager.go
package knowledge

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sage-ai-v2/pkg/logger"
	"strings"
	"time"
)

// KnowledgeManager handles retrieval and management of training data
type KnowledgeManager struct {
	vectorDB        VectorDB
	fileStoragePath string
}

// VectorDB interface defines methods for vector database operations
// type VectorDB interface {
// 	AddTrainingItem(ctx context.Context, item *TrainingItem) error
// 	FindSimilar(ctx context.Context, query string, itemType string, limit int) ([]TrainingItem, error)
// 	ListAll(ctx context.Context, dataType string) ([]TrainingItem, error)
// 	Close()
// }
type VectorDB interface {
    AddTrainingItem(ctx context.Context, item *TrainingItem) error
    FindSimilar(ctx context.Context, query string, itemType string, limit int) ([]TrainingItem, error)
    ListAll(ctx context.Context, dataType string) ([]TrainingItem, error)
    DeleteItem(ctx context.Context, id string) error // Add this method
    Close()
}

// CreateKnowledgeManager initializes a new knowledge manager
func CreateKnowledgeManager(vectorDB VectorDB, storagePath string) *KnowledgeManager {
	// Ensure storage path exists
	if err := os.MkdirAll(storagePath, 0755); err != nil {
		logger.ErrorLogger.Printf("Failed to create storage directory: %v", err)
	}

	return &KnowledgeManager{
		vectorDB:        vectorDB,
		fileStoragePath: storagePath,
	}
}

// StoreFile saves a file to the storage path and returns its local path
func (km *KnowledgeManager) StoreFile(fileName string, content []byte) (string, error) {
	// Create a sanitized filename
	sanitized := strings.ReplaceAll(fileName, " ", "_")
	timestamp := time.Now().UnixNano()
	localPath := filepath.Join(km.fileStoragePath, fmt.Sprintf("%d_%s", timestamp, sanitized))
	
	err := os.WriteFile(localPath, content, 0644)
	if err != nil {
		return "", fmt.Errorf("failed to write file: %w", err)
	}
	
	return localPath, nil
}

// AddDDLSchema adds a DDL schema to the knowledge base
func (km *KnowledgeManager) AddDDLSchema(ctx context.Context, name string, content string, description string) error {
	id := fmt.Sprintf("ddl_%d", time.Now().UnixNano())
	
	item := &TrainingItem{
		ID:          id,
		Type:        "ddl",
		Content:     content,
		Description: description,
		Metadata: map[string]interface{}{
			"name": name,
		},
		DateAdded: time.Now().Format(time.RFC3339),
	}
	
	return km.vectorDB.AddTrainingItem(ctx, item)
}

// AddDocumentation adds documentation to the knowledge base
func (km *KnowledgeManager) AddDocumentation(ctx context.Context, title string, content string, tags []string) error {
	id := fmt.Sprintf("doc_%d", time.Now().UnixNano())
	
	item := &TrainingItem{
		ID:          id,
		Type:        "documentation",
		Content:     content,
		Description: title,
		Metadata: map[string]interface{}{
			"title": title,
			"tags":  tags,
		},
		DateAdded: time.Now().Format(time.RFC3339),
	}
	
	return km.vectorDB.AddTrainingItem(ctx, item)
}

// AddQuestionSQLPair adds a question-SQL pair to the knowledge base
func (km *KnowledgeManager) AddQuestionSQLPair(ctx context.Context, pair QuestionSQLPair) error {
	id := fmt.Sprintf("qa_%d", time.Now().UnixNano())
	
	// Convert the question-SQL pair to JSON
	pairJSON, err := json.Marshal(pair)
	if err != nil {
		return fmt.Errorf("failed to marshal question-SQL pair: %w", err)
	}
	
	item := &TrainingItem{
		ID:          id,
		Type:        "question_sql",
		Content:     pair.Question, // Index by question for semantic search
		Description: pair.Description,
		Metadata: map[string]interface{}{
			"pair":     string(pairJSON),
			"tags":     pair.Tags,
			"sql":      pair.SQL,
			"verified": pair.Verified,
		},
		DateAdded: time.Now().Format(time.RFC3339),
	}
	
	return km.vectorDB.AddTrainingItem(ctx, item)
}

// RetrieveRelevantKnowledge finds relevant knowledge items for a query
func (km *KnowledgeManager) RetrieveRelevantKnowledge(ctx context.Context, query string) (*KnowledgeResult, error) {
	result := &KnowledgeResult{
		DDLSchemas:       []TrainingItem{},
		Documentation:    []TrainingItem{},
		QuestionSQLPairs: []QuestionSQLPair{},
	}
	
	// Find similar DDL schemas
	ddlItems, err := km.vectorDB.FindSimilar(ctx, query, "ddl", 2)
	if err != nil {
		logger.ErrorLogger.Printf("Failed to retrieve DDL schemas: %v", err)
	} else {
		result.DDLSchemas = ddlItems
	}
	
	// Find similar documentation
	docItems, err := km.vectorDB.FindSimilar(ctx, query, "documentation", 3)
	if err != nil {
		logger.ErrorLogger.Printf("Failed to retrieve documentation: %v", err)
	} else {
		result.Documentation = docItems
	}
	
	// Find similar question-SQL pairs
	qaPairs, err := km.vectorDB.FindSimilar(ctx, query, "question_sql", 5)
	if err != nil {
		logger.ErrorLogger.Printf("Failed to retrieve question-SQL pairs: %v", err)
	} else {
		// Convert the training items to question-SQL pairs
		for _, pair := range qaPairs {
			pairJSON, ok := pair.Metadata["pair"].(string)
			if !ok {
				continue
			}
			
			var questionSQLPair QuestionSQLPair
			err := json.Unmarshal([]byte(pairJSON), &questionSQLPair)
			if err != nil {
				logger.ErrorLogger.Printf("Failed to unmarshal question-SQL pair: %v", err)
				continue
			}
			
			result.QuestionSQLPairs = append(result.QuestionSQLPairs, questionSQLPair)
		}
	}
	
	return result, nil
}

// LoadQuestionSQLPairsFromJSON loads question-SQL pairs from a JSON file
func (km *KnowledgeManager) LoadQuestionSQLPairsFromJSON(ctx context.Context, filePath string) (int, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return 0, fmt.Errorf("failed to read file: %w", err)
	}
	
	var pairs []QuestionSQLPair
	err = json.Unmarshal(data, &pairs)
	if err != nil {
		return 0, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}
	
	added := 0
	for _, pair := range pairs {
		if pair.DateAdded == "" {
			pair.DateAdded = time.Now().Format(time.RFC3339)
		}
		
		err := km.AddQuestionSQLPair(ctx, pair)
		if err != nil {
			logger.ErrorLogger.Printf("Failed to add question-SQL pair: %v", err)
			continue
		}
		
		added++
	}
	
	return added, nil
}

// ListTrainingData retrieves all training data items
// func (km *KnowledgeManager) ListTrainingData(ctx context.Context, dataType string) ([]map[string]interface{}, error) {
// 	// For now, we'll return a simple mock implementation
// 	// In a full implementation, you would query your vector database
// 	return []map[string]interface{}{
// 		{
// 			"id":          "ddl_example",
// 			"type":        "ddl",
// 			"description": "Example DDL Schema",
// 			"date_added":  time.Now().Format(time.RFC3339),
// 		},
// 		{
// 			"id":          "doc_example",
// 			"type":        "documentation",
// 			"description": "Example Documentation",
// 			"date_added":  time.Now().Format(time.RFC3339),
// 		},
// 		{
// 			"id":          "qa_example",
// 			"type":        "question_sql",
// 			"description": "Example Question-SQL Pair",
// 			"date_added":  time.Now().Format(time.RFC3339),
// 		},
// 	}, nil
// }
// ListTrainingData retrieves all training data items
func (km *KnowledgeManager) ListTrainingData(ctx context.Context, dataType string) ([]map[string]interface{}, error) {
    logger.InfoLogger.Printf("Listing training data. Filter: %s", dataType)
    
    // Check if VectorDB is available
    if km.vectorDB == nil {
        return nil, fmt.Errorf("vector database is not initialized")
    }
    
    // Get all items from vector database
    items, err := km.vectorDB.ListAll(ctx, dataType)
    if err != nil {
        logger.ErrorLogger.Printf("Failed to list items from vector db: %v", err)
        return nil, err
    }
    
    // Convert to response format
    result := make([]map[string]interface{}, len(items))
    for i, item := range items {
        result[i] = map[string]interface{}{
            "id":          item.ID,
            "type":        item.Type,
            "description": item.Description,
            "date_added":  item.DateAdded,
            // Don't include content to keep response size small
        }
    }
    
    logger.InfoLogger.Printf("Found %d training data items", len(result))
    return result, nil
}

func (km *KnowledgeManager) GetVectorDB() VectorDB {
    return km.vectorDB
}