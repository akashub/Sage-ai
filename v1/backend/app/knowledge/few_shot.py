# v1/backend/app/knowledge/few_shot.py
import asyncio
from typing import List, Dict, Any, Optional
import json
from pathlib import Path
import uuid
from pydantic import BaseModel
from .vector_store import VectorStore, QueryExample
from langchain_openai import OpenAIEmbeddings

class FewShotExample(BaseModel):
    """Model for few-shot examples"""
    id: str
    natural_query: str
    sql_query: str
    schema_context: Optional[str]
    metadata: Dict[str, Any]

class FewShotManager:
    def __init__(self, 
                 examples_file: str,
                 embedding_model: Optional[OpenAIEmbeddings] = None):
        """
        Initialize with both JSON examples and vector store
        """
        self.examples_file = Path('app/data') / examples_file
        self.curated_examples: List[FewShotExample] = []
        self.vector_store = VectorStore(embedding_model)
        self.curated_examples = []
        asyncio.create_task(self.load_examples())
        
    async def load_examples(self) -> None:
        """Load examples from JSON and initialize vector store"""
        try:
            with open(self.examples_file, 'r') as f:
                data = json.load(f)
                self.curated_examples = [FewShotExample(**example) for example in data['examples']]
                
                # Add examples to vector store if not already present
                for example in self.curated_examples:
                    await self._ensure_in_vector_store(example)
        except Exception as e:
            print(f"Error loading examples: {e}")
            self.curated_examples = []

    async def _ensure_in_vector_store(self, example: FewShotExample):
        """Ensure example is in vector store"""
        query_example = QueryExample(
            id=example.id,
            natural_query=example.natural_query,
            sql_query=example.sql_query,
            schema_context=example.schema_context,
            metadata=example.metadata
        )
        await self.vector_store.add_example(query_example)

    async def get_relevant_examples(self, 
                                  query: str,
                                  k: int = 3,
                                  use_patterns: bool = True,
                                  use_vector: bool = True) -> List[FewShotExample]:
        """
        Get relevant examples using both pattern matching and vector similarity
        """
        relevant_examples = []

        if use_vector:
            # Get similar examples from vector store
            vector_examples = await self.vector_store.find_similar(query, k=k)
            relevant_examples.extend(vector_examples)

        if use_patterns:
            # Get pattern-matched examples from curated list
            patterns = self._identify_patterns(query)
            pattern_examples = self.get_examples_by_patterns(patterns)
            relevant_examples.extend(pattern_examples)

        # Deduplicate and limit results
        unique_examples = list({ex.id: ex for ex in relevant_examples}.values())
        return unique_examples[:k]

    def _identify_patterns(self, query: str) -> List[str]:
        """Identify patterns in the query"""
        # Get all unique patterns from metadata
        all_patterns = {
            example.metadata.get('pattern')
            for example in self.curated_examples
            if example.metadata.get('pattern')
        }
        
        # Match patterns based on keywords or rules
        matched_patterns = []
        query_lower = query.lower()
        
        for pattern in all_patterns:
            examples = self.get_examples_by_pattern(pattern)
            # If any example's keywords match the query
            for example in examples:
                keywords = example.natural_query.lower().split()
                if any(keyword in query_lower for keyword in keywords):
                    matched_patterns.append(pattern)
                    break
        
        return matched_patterns

    def get_formatted_examples(self, 
                             examples: List[Any],
                             include_metadata: bool = False) -> str:
        """Format examples for prompt injection"""
        formatted = "Few-shot examples:\n\n"
        for i, example in enumerate(examples, 1):
            formatted += f"Example {i}:\n"
            formatted += f"Query: {example.natural_query}\n"
            formatted += f"SQL: {example.sql_query}\n"
            if example.schema_context:
                formatted += f"Schema: {example.schema_context}\n"
            if include_metadata:
                formatted += f"Metadata: {json.dumps(example.metadata, indent=2)}\n"
            formatted += "\n"
        return formatted

    def get_examples_by_pattern(self, pattern: str) -> List[FewShotExample]:
        """Get examples matching a specific pattern"""
        return [
            example for example in self.curated_examples 
            if example.metadata.get('pattern') == pattern
        ]

    def get_examples_by_patterns(self, patterns: List[str]) -> List[FewShotExample]:
        """Get examples matching any of the patterns"""
        return [
            example for example in self.curated_examples 
            if example.metadata.get('pattern') in patterns
        ]

    async def add_example(self, 
                         natural_query: str,
                         sql_query: str,
                         schema_context: Optional[str] = None,
                         metadata: Optional[Dict[str, Any]] = None) -> None:
        """Add new example to both JSON and vector store"""
        example = FewShotExample(
            id=str(uuid.uuid4()),
            natural_query=natural_query,
            sql_query=sql_query,
            schema_context=schema_context,
            metadata=metadata or {}
        )
        
        # Add to curated examples
        self.curated_examples.append(example)
        self._save_examples()
        
        # Add to vector store
        await self._ensure_in_vector_store(example)