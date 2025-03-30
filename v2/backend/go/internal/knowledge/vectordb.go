// backend/go/internal/knowledge/memory_vectordb.go
package knowledge

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"sage-ai-v2/pkg/logger"
	"sort"
	"sync"
)

// MemoryVectorDB is a simple in-memory implementation of VectorDB
type MemoryVectorDB struct {
	items       map[string]TrainingItem
	mu          sync.RWMutex
	persistPath string
}

// CreateMemoryVectorDB initializes a new in-memory vector database
func CreateMemoryVectorDB(persistPath string) (*MemoryVectorDB, error) {
	db := &MemoryVectorDB{
		items:       make(map[string]TrainingItem),
		persistPath: persistPath,
	}
	
	// Create directory if needed
	if persistPath != "" {
		dir := filepath.Dir(persistPath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create directory: %w", err)
		}
		
		// Try to load persisted data
		if err := db.loadFromDisk(); err != nil {
			logger.InfoLogger.Printf("No existing data or error loading: %v", err)
		}
	}
	
	return db, nil
}

// AddTrainingItem adds a new item to the knowledge base
func (db *MemoryVectorDB) AddTrainingItem(ctx context.Context, item *TrainingItem) error {
	// Generate a simple embedding
	item.Embedding = generateSimpleEmbedding(item.Content)
	
	// Store in memory
	db.mu.Lock()
	db.items[item.ID] = *item
	db.mu.Unlock()
	
	// Persist to disk if path is set
	if db.persistPath != "" {
		if err := db.persistToDisk(); err != nil {
			return fmt.Errorf("failed to persist data: %w", err)
		}
	}
	
	return nil
}

// FindSimilar finds training items similar to the given query
func (db *MemoryVectorDB) FindSimilar(ctx context.Context, query string, itemType string, limit int) ([]TrainingItem, error) {
	// Generate embedding for query
	queryEmbedding := generateSimpleEmbedding(query)
	
	// Calculate similarity with all items
	type itemWithSimilarity struct {
		item       TrainingItem
		similarity float64
	}
	
	var similarities []itemWithSimilarity
	
	db.mu.RLock()
	for _, item := range db.items {
		// Skip if type doesn't match (when specified)
		if itemType != "" && item.Type != itemType {
			continue
		}
		
		similarity := cosineSimilarity(queryEmbedding, item.Embedding)
		similarities = append(similarities, itemWithSimilarity{item, similarity})
	}
	db.mu.RUnlock()
	
	// Sort by similarity (descending)
	sort.Slice(similarities, func(i, j int) bool {
		return similarities[i].similarity > similarities[j].similarity
	})
	
	// Take top N results
	result := make([]TrainingItem, 0, limit)
	for i := 0; i < limit && i < len(similarities); i++ {
		result = append(result, similarities[i].item)
	}
	
	return result, nil
}

// Close cleans up resources
func (db *MemoryVectorDB) Close() {
	// Persist data before closing
	if db.persistPath != "" {
		if err := db.persistToDisk(); err != nil {
			logger.ErrorLogger.Printf("Failed to persist data on close: %v", err)
		}
	}
}

// persistToDisk saves the database to disk
func (db *MemoryVectorDB) persistToDisk() error {
	db.mu.RLock()
	defer db.mu.RUnlock()
	
	data, err := json.MarshalIndent(db.items, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal items: %w", err)
	}
	
	if err := os.WriteFile(db.persistPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}
	
	logger.InfoLogger.Printf("Persisted %d items to %s", len(db.items), db.persistPath)
	return nil
}

// loadFromDisk loads the database from disk
func (db *MemoryVectorDB) loadFromDisk() error {
	// Check if file exists
	if _, err := os.Stat(db.persistPath); os.IsNotExist(err) {
		return fmt.Errorf("persist file does not exist")
	}
	
	data, err := os.ReadFile(db.persistPath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}
	
	var items map[string]TrainingItem
	if err := json.Unmarshal(data, &items); err != nil {
		return fmt.Errorf("failed to unmarshal items: %w", err)
	}
	
	db.mu.Lock()
	db.items = items
	db.mu.Unlock()
	
	logger.InfoLogger.Printf("Loaded %d items from %s", len(items), db.persistPath)
	return nil
}

// generateSimpleEmbedding creates a simple embedding for the given text
// This is just a mock implementation - in a real system, you'd use a proper
// embedding model or service
func generateSimpleEmbedding(text string) []float32 {
	// Create a simple "embedding" based on character frequencies
	embedding := make([]float32, 128) // Small dimension for memory efficiency
	
	if text == "" {
		return embedding
	}
	
	// Count character frequencies
	freqs := make(map[byte]int)
	for i := 0; i < len(text); i++ {
		freqs[text[i]]++
	}
	
	// Normalize and assign to embedding
	for char, count := range freqs {
		idx := int(char) % len(embedding)
		embedding[idx] = float32(count) / float32(len(text))
	}
	
	return embedding
}

// cosineSimilarity calculates the cosine similarity between two vectors
func cosineSimilarity(a, b []float32) float64 {
	// Handle different lengths (just compare common dimensions)
	minLen := len(a)
	if len(b) < minLen {
		minLen = len(b)
	}
	
	var dotProduct float64
	var normA float64
	var normB float64
	
	for i := 0; i < minLen; i++ {
		dotProduct += float64(a[i] * b[i])
		normA += float64(a[i] * a[i])
		normB += float64(b[i] * b[i])
	}
	
	// Handle zero vectors
	if normA == 0 || normB == 0 {
		return 0
	}
	
	return dotProduct / (math.Sqrt(normA) * math.Sqrt(normB))
}

func (db *MemoryVectorDB) ListAll(ctx context.Context, itemType string) ([]TrainingItem, error) {
    db.mu.RLock()
    defer db.mu.RUnlock()
    
    var items []TrainingItem
    for _, item := range db.items {
        if itemType == "" || item.Type == itemType {
            items = append(items, item)
        }
    }
    
    return items, nil
}