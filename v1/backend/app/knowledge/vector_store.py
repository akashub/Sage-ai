# app/knowledge/vector_store.py
from typing import List, Dict, Any, Optional
from pydantic import BaseModel
import numpy as np
from langchain_openai import OpenAIEmbeddings
import faiss
import json
import os
from datetime import datetime

class QueryExample(BaseModel):
    """Example class with embedding support"""
    id: str
    natural_query: str
    sql_query: str
    schema_context: Optional[str]
    metadata: Dict[str, Any]
    embedding: Optional[List[float]] = None
    created_at: datetime = datetime.now()
    last_used: Optional[datetime] = None
    success_count: int = 0
    failure_count: int = 0

class VectorStore:
    def __init__(self, 
                 embedding_model: Optional[OpenAIEmbeddings] = None,
                 dimension: int = 1536,  # OpenAI embedding dimension
                 store_path: str = 'vector_store'):
        """Initialize vector store"""
        self.embedding_model = embedding_model or OpenAIEmbeddings()
        self.dimension = dimension
        self.store_path = store_path
        self.index = None
        self.examples: Dict[str, QueryExample] = {}
        self.initialize_store()

    def initialize_store(self):
        """Initialize FAISS index and load existing examples"""
        try:
            if not os.path.exists(self.store_path):
                os.makedirs(self.store_path)

            index_path = os.path.join(self.store_path, 'faiss_index.bin')
            examples_path = os.path.join(self.store_path, 'examples.json')

            if os.path.exists(index_path):
                self.index = faiss.read_index(index_path)
                
                with open(examples_path, 'r') as f:
                    examples_data = json.load(f)
                    self.examples = {
                        id_: QueryExample(**{
                            **data,
                            'created_at': datetime.fromisoformat(data['created_at']),
                            'last_used': (datetime.fromisoformat(data['last_used']) 
                                        if data.get('last_used') else None)
                        })
                        for id_, data in examples_data.items()
                    }
            else:
                self.index = faiss.IndexFlatL2(self.dimension)
        except Exception as e:
            print(f"Error initializing vector store: {e}")
            self.index = faiss.IndexFlatL2(self.dimension)

    async def add_example(self, example: QueryExample) -> bool:
        """Add new example with embedding"""
        try:
            # Generate embedding if not provided
            if not example.embedding:
                example.embedding = await self._get_embedding(example.natural_query)

            # Add to FAISS index
            self.index.add(np.array([example.embedding], dtype=np.float32))
            
            # Store example
            self.examples[example.id] = example
            
            # Save updated store
            self._save_store()
            
            return True
        except Exception as e:
            print(f"Error adding example: {e}")
            return False

    async def find_similar(self, 
                          query: str, 
                          k: int = 3,
                          threshold: float = 0.8) -> List[QueryExample]:
        """Find k most similar examples with score threshold"""
        try:
            # Get query embedding
            query_embedding = await self._get_embedding(query)
            
            # Search in FAISS
            D, I = self.index.search(
                np.array([query_embedding], dtype=np.float32), 
                k * 2  # Get more candidates for threshold filtering
            )
            
            # Filter and sort results
            similar_examples = []
            for distance, idx in zip(D[0], I[0]):
                # Convert distance to similarity score (0 to 1)
                similarity = 1 / (1 + distance)
                
                if similarity >= threshold:
                    # Find example with this index
                    for example in self.examples.values():
                        if example.embedding == self.index.reconstruct(idx).tolist():
                            # Update usage statistics
                            example.last_used = datetime.now()
                            similar_examples.append(example)
                            break

            # Save updated usage statistics
            self._save_store()
            
            return similar_examples[:k]  # Return top k after threshold filtering
        except Exception as e:
            print(f"Error finding similar examples: {e}")
            return []

    async def _get_embedding(self, text: str) -> List[float]:
        """Get embedding for text"""
        try:
            embeddings = await self.embedding_model.aembed_documents([text])
            return embeddings[0]
        except Exception as e:
            print(f"Error generating embedding: {e}")
            raise

    def _save_store(self):
        """Save vector store and examples"""
        try:
            # Ensure directory exists
            os.makedirs(self.store_path, exist_ok=True)
            
            # Save FAISS index
            faiss.write_index(self.index, 
                            os.path.join(self.store_path, 'faiss_index.bin'))
            
            # Save examples
            examples_data = {
                id_: {
                    **example.dict(),
                    'created_at': example.created_at.isoformat(),
                    'last_used': example.last_used.isoformat() if example.last_used else None
                }
                for id_, example in self.examples.items()
            }
            
            with open(os.path.join(self.store_path, 'examples.json'), 'w') as f:
                json.dump(examples_data, f, indent=2)
        except Exception as e:
            print(f"Error saving vector store: {e}")

    def update_example_stats(self, example_id: str, success: bool):
        """Update example success/failure statistics"""
        if example_id in self.examples:
            example = self.examples[example_id]
            if success:
                example.success_count += 1
            else:
                example.failure_count += 1
            self._save_store()

    def get_example_stats(self) -> Dict[str, Any]:
        """Get statistics about stored examples"""
        total_examples = len(self.examples)
        if total_examples == 0:
            return {
                "total_examples": 0,
                "avg_success_rate": 0,
                "most_used_patterns": []
            }

        # Calculate statistics
        success_rates = []
        pattern_counts = {}

        for example in self.examples.values():
            total_uses = example.success_count + example.failure_count
            if total_uses > 0:
                success_rates.append(example.success_count / total_uses)

            pattern = example.metadata.get('pattern')
            if pattern:
                pattern_counts[pattern] = pattern_counts.get(pattern, 0) + 1

        # Sort patterns by frequency
        sorted_patterns = sorted(
            pattern_counts.items(), 
            key=lambda x: x[1], 
            reverse=True
        )

        return {
            "total_examples": total_examples,
            "avg_success_rate": sum(success_rates) / len(success_rates) if success_rates else 0,
            "most_used_patterns": sorted_patterns[:5]
        }

    def cleanup_old_examples(self, days_threshold: int = 30):
        """Remove old, unused examples"""
        current_time = datetime.now()
        examples_to_remove = []

        for id_, example in self.examples.items():
            if example.last_used:
                days_since_used = (current_time - example.last_used).days
                if days_since_used > days_threshold:
                    examples_to_remove.append(id_)

        for id_ in examples_to_remove:
            del self.examples[id_]

        # Rebuild index if examples were removed
        if examples_to_remove:
            self._rebuild_index()
            self._save_store()

    def _rebuild_index(self):
        """Rebuild FAISS index from scratch"""
        self.index = faiss.IndexFlatL2(self.dimension)
        for example in self.examples.values():
            if example.embedding:
                self.index.add(np.array([example.embedding], dtype=np.float32))