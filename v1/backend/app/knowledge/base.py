# v1/backend/app/knowledge/base.py
from typing import List, Dict, Any, Optional
from pydantic import BaseModel
from uuid import uuid4
from datetime import datetime
from langchain_openai import OpenAIEmbeddings
from .vector_store import VectorStore, QueryExample
from .few_shot import FewShotManager


class QueryPattern(BaseModel):
    """Model for common query patterns"""
    pattern_type: str
    description: str
    example_query: str
    sql_template: str
    metadata: Dict[str, Any] = {}

class QueryLog(BaseModel):
    """Model for logging query executions"""
    query_id: str
    timestamp: datetime
    natural_query: str
    generated_sql: str
    execution_time: float
    success: bool
    error: Optional[str] = None

class KnowledgeBase:
    def __init__(self, 
                #  examples_file: str,
                 embedding_model: Optional[OpenAIEmbeddings] = None,
                 vector_store: Optional[VectorStore] = None):
        """Initialize knowledge base with vector store and common patterns"""
        self.vector_store = vector_store or VectorStore(embedding_model)
        self.patterns = self._initialize_patterns()
        self.query_logs: List[QueryLog] = []
        # self._initialize_default_examples()
        # self.few_shot_manager = FewShotManager(examples_file)


    def _initialize_patterns(self) -> Dict[str, QueryPattern]:
        """Initialize common query patterns"""
        return {
            "temporal": QueryPattern(
                pattern_type="temporal",
                description="Time-based queries",
                example_query="Show records from last month",
                sql_template="SELECT * FROM {table} WHERE DATE_TRUNC('month', {date_column}) = DATE_TRUNC('month', CURRENT_DATE - INTERVAL '1 month')",
                metadata={"complexity": "medium"}
            ),
            "aggregation": QueryPattern(
                pattern_type="aggregation",
                description="Aggregate calculations",
                example_query="Calculate total by category",
                sql_template="SELECT {group_by}, SUM({value_column}) as total FROM {table} GROUP BY {group_by}",
                metadata={"complexity": "medium"}
            ),
            "search": QueryPattern(
                pattern_type="search",
                description="Text search queries",
                example_query="Find records containing 'keyword'",
                sql_template="SELECT * FROM {table} WHERE {column} ILIKE '%{search_term}%'",
                metadata={"complexity": "easy"}
            ),
            "filter": QueryPattern(
                pattern_type="filter",
                description="Filter-based queries",
                example_query="Show records where value > 100",
                sql_template="SELECT * FROM {table} WHERE {column} {operator} {value}",
                metadata={"complexity": "easy"}
            ),
            "join": QueryPattern(
                pattern_type="join",
                description="Join-based queries",
                example_query="Combine data from two tables",
                sql_template="SELECT * FROM {table1} JOIN {table2} ON {table1}.{key1} = {table2}.{key2}",
                metadata={"complexity": "hard"}
            ),
            "pagination": QueryPattern(
                pattern_type="pagination",
                description="Paginated results",
                example_query="Show first 10 records",
                sql_template="SELECT * FROM {table} ORDER BY {order_column} LIMIT {limit} OFFSET {offset}",
                metadata={"complexity": "easy"}
            )
        }

    async def _initialize_default_examples(self):
        """Initialize with default examples - override in subclass"""
        examples = self.few_shot_manager.get_formatted_examples()

    async def add_example(self, 
                         natural_query: str,
                         sql_query: str,
                         schema_context: Optional[str] = None,
                         metadata: Optional[dict] = None) -> bool:
        """Add new example to knowledge base"""
        try:
            example = QueryExample(
                id=str(uuid4()),
                natural_query=natural_query,
                sql_query=sql_query,
                schema_context=schema_context,
                metadata=metadata or {}
            )
            return await self.vector_store.add_example(example)
        except Exception as e:
            print(f"Error adding example: {e}")
            return False

    # async def find_similar_examples(self, 
    #                               query: str,
    #                               k: int = 3) -> List[QueryExample]:
    #     """Find similar examples using vector similarity"""
    #     return await self.vector_store.find_similar(query, k)

    async def find_similar_examples(self, query: str, k: int = 3):
        return await self.few_shot_manager.get_relevant_examples(
            query, 
            k=k,
            use_patterns=True,
            use_vector=True
        )
    
    def get_pattern_template(self, pattern_type: str) -> Optional[QueryPattern]:
        """Get pattern template by type"""
        return self.patterns.get(pattern_type)

    def identify_patterns(self, query: str) -> List[str]:
        """Identify potential patterns in a query"""
        matched_patterns = []
        for pattern_type, pattern in self.patterns.items():
            # Simple keyword matching - could be enhanced with more sophisticated matching
            keywords = pattern.example_query.lower().split()
            if any(keyword in query.lower() for keyword in keywords):
                matched_patterns.append(pattern_type)
        return matched_patterns

    def log_query(self, 
                 natural_query: str,
                 generated_sql: str,
                 execution_time: float,
                 success: bool,
                 error: Optional[str] = None):
        """Log query execution for analysis"""
        log = QueryLog(
            query_id=str(uuid4()),
            timestamp=datetime.now(),
            natural_query=natural_query,
            generated_sql=generated_sql,
            execution_time=execution_time,
            success=success,
            error=error
        )
        self.query_logs.append(log)

    def get_query_statistics(self) -> Dict[str, Any]:
        """Get statistics about query executions"""
        total_queries = len(self.query_logs)
        if total_queries == 0:
            return {
                "total_queries": 0,
                "success_rate": 0,
                "average_execution_time": 0
            }

        successful_queries = len([log for log in self.query_logs if log.success])
        total_execution_time = sum(log.execution_time for log in self.query_logs)

        return {
            "total_queries": total_queries,
            "success_rate": successful_queries / total_queries,
            "average_execution_time": total_execution_time / total_queries
        }

    async def find_similar_by_pattern(self, 
                                    query: str,
                                    pattern_type: str,
                                    k: int = 3) -> List[QueryExample]:
        """Find similar examples matching a specific pattern"""
        all_similar = await self.find_similar_examples(query, k=k*2)
        pattern_examples = [
            ex for ex in all_similar 
            if ex.metadata.get("pattern") == pattern_type
        ]
        return pattern_examples[:k]