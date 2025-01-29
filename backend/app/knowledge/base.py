# # # # app/knowledge/base.py
# # # from typing import List, Dict, Any, Optional
# # # from pydantic import BaseModel

# # # class Example(BaseModel):
# # #     natural_query: str
# # #     sql_query: str
# # #     schema_context: Optional[str]
# # #     metadata: Optional[Dict[str, Any]]

# # # class KnowledgeBase:
# # #     def __init__(self):
# # #         # Initialize with some default examples
# # #         self.examples: List[Example] = [
# # #             Example(
# # #                 natural_query="Show me all users who signed up last month",
# # #                 sql_query="SELECT * FROM users WHERE DATE_TRUNC('month', created_at) = DATE_TRUNC('month', CURRENT_DATE - INTERVAL '1 month')",
# # #                 schema_context="users(id, name, email, created_at)",
# # #                 metadata={"type": "temporal_query", "difficulty": "medium"}
# # #             ),
# # #             Example(
# # #                 natural_query="Find the total sales by category in descending order",
# # #                 sql_query="SELECT category, SUM(amount) as total_sales FROM sales GROUP BY category ORDER BY total_sales DESC",
# # #                 schema_context="sales(id, category, amount, date)",
# # #                 metadata={"type": "aggregation", "difficulty": "easy"}
# # #             ),
# # #             # Add more examples here
# # #         ]
    
# # #     def add_example(self, example: Example):
# # #         """Add a new example to the knowledge base"""
# # #         self.examples.append(example)
    
# # #     def get_similar_examples(self, query: str, schema: Optional[str] = None, k: int = 3) -> List[Example]:
# # #         """Retrieve k most similar examples using semantic similarity"""
# # #         # TODO: Implement semantic similarity using embeddings
# # #         return self.examples[:k]  # Placeholder implementation

# # # # app/knowledge/few_shot.py
# # # from typing import List, Optional
# # # from langchain_openai import ChatOpenAI
# # # from .base import Example, KnowledgeBase

# # # class FewShotLearner:
# # #     def __init__(self, llm: ChatOpenAI, knowledge_base: KnowledgeBase):
# # #         self.llm = llm
# # #         self.knowledge_base = knowledge_base

# # #     def _format_examples(self, examples: List[Example]) -> str:
# # #         """Format examples for prompt"""
# # #         formatted = ""
# # #         for i, ex in enumerate(examples, 1):
# # #             formatted += f"Example {i}:\n"
# # #             formatted += f"Natural Query: {ex.natural_query}\n"
# # #             formatted += f"Schema Context: {ex.schema_context}\n"
# # #             formatted += f"SQL Query: {ex.sql_query}\n\n"
# # #         return formatted

# # #     async def generate_sql(self, query: str, schema: Optional[str] = None) -> str:
# # #         """Generate SQL using few-shot learning"""
# # #         # Get similar examples
# # #         similar_examples = self.knowledge_base.get_similar_examples(query, schema)
        
# # #         # Format prompt with examples
# # #         prompt = f"""Given the following examples and schema, convert the natural language query to SQL.

# # # {self._format_examples(similar_examples)}

# # # Current Schema: {schema or 'Not provided'}
# # # Natural Query: {query}

# # # Generate a SQL query following the patterns shown in the examples above.
# # # SQL Query:"""

# # #         # Generate SQL using LLM
# # #         response = await self.llm.ainvoke(prompt)
# # #         return response.content.strip()

# # # # app/knowledge/schema_manager.py
# # # from typing import Dict, Optional
# # # import json

# # # class SchemaManager:
# # #     def __init__(self):
# # #         self.schemas: Dict[str, dict] = {}
    
# # #     def add_schema(self, name: str, schema_dict: dict):
# # #         """Add or update a schema"""
# # #         self.schemas[name] = schema_dict
    
# # #     def get_schema(self, name: str) -> Optional[dict]:
# # #         """Retrieve a schema by name"""
# # #         return self.schemas.get(name)
    
# # #     def format_schema_for_prompt(self, name: str) -> str:
# # #         """Format schema for inclusion in prompts"""
# # #         schema = self.get_schema(name)
# # #         if not schema:
# # #             return "No schema available"
            
# # #         formatted = "Database Schema:\n"
# # #         for table, columns in schema.items():
# # #             formatted += f"Table: {table}\n"
# # #             formatted += "Columns:\n"
# # #             for col, details in columns.items():
# # #                 formatted += f"- {col}: {details['type']}"
# # #                 if details.get('description'):
# # #                     formatted += f" ({details['description']})"
# # #                 formatted += "\n"
# # #             formatted += "\n"
# # #         return formatted

# # app/knowledge/base.py
# from typing import List, Dict, Any, Optional
# from pydantic import BaseModel
# from uuid import uuid4
# from datetime import datetime
# from langchain_openai import OpenAIEmbeddings
# from .vector_store import VectorStore, QueryExample

# class QueryPattern(BaseModel):
#     """Model for common query patterns"""
#     pattern_type: str
#     description: str
#     example_query: str
#     sql_template: str
#     metadata: Dict[str, Any] = {}

# class QueryLog(BaseModel):
#     """Model for logging query executions"""
#     query_id: str
#     timestamp: datetime
#     natural_query: str
#     generated_sql: str
#     execution_time: float
#     success: bool
#     error: Optional[str] = None

# class KnowledgeBase:
#     def __init__(self, 
#                  embedding_model: Optional[OpenAIEmbeddings] = None,
#                  vector_store: Optional[VectorStore] = None):
#         """Initialize knowledge base with vector store and common patterns"""
#         self.vector_store = vector_store or VectorStore(embedding_model)
#         self.patterns = self._initialize_patterns()
#         self.query_logs: List[QueryLog] = []
#         self._initialize_default_examples()

#     def _initialize_patterns(self) -> Dict[str, QueryPattern]:
#         """Initialize common query patterns"""
#         return {
#             "temporal": QueryPattern(
#                 pattern_type="temporal",
#                 description="Time-based queries",
#                 example_query="Show sales from last month",
#                 sql_template="SELECT * FROM {table} WHERE DATE_TRUNC('month', {date_column}) = DATE_TRUNC('month', CURRENT_DATE - INTERVAL '1 month')",
#                 metadata={"complexity": "medium"}
#             ),
#             "aggregation": QueryPattern(
#                 pattern_type="aggregation",
#                 description="Aggregate calculations",
#                 example_query="Calculate total sales by category",
#                 sql_template="SELECT {group_by}, SUM({value_column}) as total FROM {table} GROUP BY {group_by}",
#                 metadata={"complexity": "medium"}
#             ),
#             "join": QueryPattern(
#                 pattern_type="join",
#                 description="Table joins",
#                 example_query="Show customers with their orders",
#                 sql_template="SELECT * FROM {table1} JOIN {table2} ON {table1}.{key1} = {table2}.{key2}",
#                 metadata={"complexity": "hard"}
#             ),
#             "pagination": QueryPattern(
#                 pattern_type="pagination",
#                 description="Paginated results",
#                 example_query="Show first 10 customers",
#                 sql_template="SELECT * FROM {table} ORDER BY {order_column} LIMIT {limit} OFFSET {offset}",
#                 metadata={"complexity": "easy"}
#             ),
#             "search": QueryPattern(
#                 pattern_type="search",
#                 description="Text search",
#                 example_query="Find products containing 'phone'",
#                 sql_template="SELECT * FROM {table} WHERE {column} ILIKE '%{search_term}%'",
#                 metadata={"complexity": "easy"}
#             ),
#             "nested": QueryPattern(
#                 pattern_type="nested",
#                 description="Nested queries",
#                 example_query="Find customers who spent more than average",
#                 sql_template="SELECT * FROM {table} WHERE {value_column} > (SELECT AVG({value_column}) FROM {table})",
#                 metadata={"complexity": "hard"}
#             )
#         }

#     async def _initialize_default_examples(self):
#         """Initialize knowledge base with default examples"""
#         default_examples = [
#             {
#                 "natural_query": "Show me sales from last month",
#                 "sql_query": """
#                 SELECT * FROM sales 
#                 WHERE DATE_TRUNC('month', date) = DATE_TRUNC('month', CURRENT_DATE - INTERVAL '1 month')
#                 """,
#                 "schema_context": "sales(id, amount, date, category)",
#                 "metadata": {
#                     "type": "temporal",
#                     "difficulty": "medium",
#                     "pattern": "temporal"
#                 }
#             },
#             {
#                 "natural_query": "Calculate total revenue by product category",
#                 "sql_query": """
#                 SELECT category, SUM(amount) as total_revenue 
#                 FROM sales 
#                 GROUP BY category 
#                 ORDER BY total_revenue DESC
#                 """,
#                 "schema_context": "sales(id, amount, date, category)",
#                 "metadata": {
#                     "type": "aggregation",
#                     "difficulty": "medium",
#                     "pattern": "aggregation"
#                 }
#             },
#             {
#                 "natural_query": "Show customers with their orders sorted by order date",
#                 "sql_query": """
#                 SELECT c.name, o.order_date, o.amount 
#                 FROM customers c 
#                 JOIN orders o ON c.id = o.customer_id 
#                 ORDER BY o.order_date DESC
#                 """,
#                 "schema_context": """
#                 customers(id, name, email)
#                 orders(id, customer_id, order_date, amount)
#                 """,
#                 "metadata": {
#                     "type": "join",
#                     "difficulty": "hard",
#                     "pattern": "join"
#                 }
#             }
#         ]

#         for example_data in default_examples:
#             example = QueryExample(
#                 id=str(uuid4()),
#                 **example_data
#             )
#             await self.vector_store.add_example(example)

#     async def add_example(self, 
#                          natural_query: str,
#                          sql_query: str,
#                          schema_context: Optional[str] = None,
#                          metadata: Optional[dict] = None) -> bool:
#         """Add new example to knowledge base"""
#         try:
#             example = QueryExample(
#                 id=str(uuid4()),
#                 natural_query=natural_query,
#                 sql_query=sql_query,
#                 schema_context=schema_context,
#                 metadata=metadata or {}
#             )
#             return await self.vector_store.add_example(example)
#         except Exception as e:
#             print(f"Error adding example: {e}")
#             return False

#     async def find_similar_examples(self, 
#                                   query: str,
#                                   k: int = 3) -> List[QueryExample]:
#         """Find similar examples using vector similarity"""
#         return await self.vector_store.find_similar(query, k)

#     def get_pattern_template(self, pattern_type: str) -> Optional[QueryPattern]:
#         """Get pattern template by type"""
#         return self.patterns.get(pattern_type)

#     def identify_patterns(self, query: str) -> List[str]:
#         """Identify potential patterns in a query"""
#         matched_patterns = []
#         for pattern_type, pattern in self.patterns.items():
#             # Simple keyword matching - could be enhanced with more sophisticated matching
#             keywords = pattern.example_query.lower().split()
#             if any(keyword in query.lower() for keyword in keywords):
#                 matched_patterns.append(pattern_type)
#         return matched_patterns

#     def log_query(self, 
#                  natural_query: str,
#                  generated_sql: str,
#                  execution_time: float,
#                  success: bool,
#                  error: Optional[str] = None):
#         """Log query execution for analysis"""
#         log = QueryLog(
#             query_id=str(uuid4()),
#             timestamp=datetime.now(),
#             natural_query=natural_query,
#             generated_sql=generated_sql,
#             execution_time=execution_time,
#             success=success,
#             error=error
#         )
#         self.query_logs.append(log)

#     def get_query_statistics(self) -> Dict[str, Any]:
#         """Get statistics about query executions"""
#         total_queries = len(self.query_logs)
#         if total_queries == 0:
#             return {
#                 "total_queries": 0,
#                 "success_rate": 0,
#                 "average_execution_time": 0
#             }

#         successful_queries = len([log for log in self.query_logs if log.success])
#         total_execution_time = sum(log.execution_time for log in self.query_logs)

#         return {
#             "total_queries": total_queries,
#             "success_rate": successful_queries / total_queries,
#             "average_execution_time": total_execution_time / total_queries
#         }

#     async def find_similar_by_pattern(self, 
#                                     query: str,
#                                     pattern_type: str,
#                                     k: int = 3) -> List[QueryExample]:
#         """Find similar examples matching a specific pattern"""
#         all_similar = await self.find_similar_examples(query, k=k*2)  # Get more examples initially
#         pattern_examples = [
#             ex for ex in all_similar 
#             if ex.metadata.get("pattern") == pattern_type
#         ]
#         return pattern_examples[:k]  # Return top k pattern-matching examples

# # app/knowledge/base.py
# from typing import List, Dict, Any, Optional
# from pydantic import BaseModel
# from uuid import uuid4
# from datetime import datetime
# from langchain_openai import OpenAIEmbeddings
# from .vector_store import VectorStore, QueryExample

# class QueryPattern(BaseModel):
#     """Model for common query patterns"""
#     pattern_type: str
#     description: str
#     example_query: str
#     sql_template: str
#     metadata: Dict[str, Any] = {}

# class QueryLog(BaseModel):
#     """Model for logging query executions"""
#     query_id: str
#     timestamp: datetime
#     natural_query: str
#     generated_sql: str
#     execution_time: float
#     success: bool
#     error: Optional[str] = None

# class KnowledgeBase:
#     def __init__(self, 
#                  embedding_model: Optional[OpenAIEmbeddings] = None,
#                  vector_store: Optional[VectorStore] = None):
#         """Initialize knowledge base with vector store and common patterns"""
#         self.vector_store = vector_store or VectorStore(embedding_model)
#         self.patterns = self._initialize_patterns()
#         self.query_logs: List[QueryLog] = []
#         self._initialize_default_examples()

#     def _initialize_patterns(self) -> Dict[str, QueryPattern]:
#         """Initialize IMDB movie-specific query patterns"""
#         return {
#             "temporal": QueryPattern(
#                 pattern_type="temporal",
#                 description="Time-based movie queries",
#                 example_query="Show movies released in 2022",
#                 sql_template="SELECT * FROM movies WHERE EXTRACT(YEAR FROM release_date) = {year}",
#                 metadata={"complexity": "easy"}
#             ),
#             "rating": QueryPattern(
#                 pattern_type="rating",
#                 description="Rating-based movie queries",
#                 example_query="Find highest rated movies",
#                 sql_template="SELECT * FROM movies WHERE vote_average > {rating} AND vote_count > {min_votes}",
#                 metadata={"complexity": "medium"}
#             ),
#             "financial": QueryPattern(
#                 pattern_type="financial",
#                 description="Budget and revenue queries",
#                 example_query="Show most profitable movies",
#                 sql_template="SELECT *, (revenue - budget) as profit FROM movies WHERE revenue > budget",
#                 metadata={"complexity": "medium"}
#             ),
#             "genre": QueryPattern(
#                 pattern_type="genre",
#                 description="Genre-based queries",
#                 example_query="Find action movies",
#                 sql_template="SELECT * FROM movies WHERE genres LIKE '%{genre}%'",
#                 metadata={"complexity": "easy"}
#             ),
#             "language": QueryPattern(
#                 pattern_type="language",
#                 description="Language-based queries",
#                 example_query="Show English language movies",
#                 sql_template="SELECT * FROM movies WHERE original_language = '{language}'",
#                 metadata={"complexity": "easy"}
#             ),
#             "popularity": QueryPattern(
#                 pattern_type="popularity",
#                 description="Popularity-based queries",
#                 example_query="Show most popular movies",
#                 sql_template="SELECT * FROM movies ORDER BY popularity DESC LIMIT {limit}",
#                 metadata={"complexity": "easy"}
#             ),
#             "advanced_search": QueryPattern(
#                 pattern_type="advanced_search",
#                 description="Complex movie search",
#                 example_query="Find popular action movies with high ratings",
#                 sql_template="SELECT * FROM movies WHERE genres LIKE '%{genre}%' AND vote_average > {rating} ORDER BY popularity DESC",
#                 metadata={"complexity": "hard"}
#             )
#         }

#     async def _initialize_default_examples(self):
#         """Initialize knowledge base with IMDB-specific examples"""
#         default_examples = [
#             {
#                 "natural_query": "Show me movies released in 2022",
#                 "sql_query": """
#                 SELECT title, release_date, vote_average 
#                 FROM movies 
#                 WHERE EXTRACT(YEAR FROM release_date) = 2022 
#                 ORDER BY vote_average DESC
#                 """,
#                 "schema_context": "movies schema",
#                 "metadata": {
#                     "type": "temporal",
#                     "difficulty": "easy",
#                     "pattern": "temporal"
#                 }
#             },
#             {
#                 "natural_query": "What are the highest grossing movies with a budget under 50 million",
#                 "sql_query": """
#                 SELECT title, budget, revenue, (revenue - budget) as profit 
#                 FROM movies 
#                 WHERE budget < 50000000 
#                 AND revenue > 0 
#                 ORDER BY revenue DESC 
#                 LIMIT 10
#                 """,
#                 "schema_context": "movies schema",
#                 "metadata": {
#                     "type": "financial",
#                     "difficulty": "medium",
#                     "pattern": "financial"
#                 }
#             },
#             {
#                 "natural_query": "Find popular action movies with rating above 8",
#                 "sql_query": """
#                 SELECT title, vote_average, popularity 
#                 FROM movies 
#                 WHERE vote_average > 8 
#                 AND genres LIKE '%Action%' 
#                 AND vote_count > 1000
#                 ORDER BY popularity DESC 
#                 LIMIT 10
#                 """,
#                 "schema_context": "movies schema",
#                 "metadata": {
#                     "type": "advanced_search",
#                     "difficulty": "hard",
#                     "pattern": "advanced_search"
#                 }
#             },
#             {
#                 "natural_query": "Show me English language movies released in 2023",
#                 "sql_query": """
#                 SELECT title, release_date, original_language 
#                 FROM movies 
#                 WHERE original_language = 'en' 
#                 AND EXTRACT(YEAR FROM release_date) = 2023 
#                 ORDER BY popularity DESC
#                 """,
#                 "schema_context": "movies schema",
#                 "metadata": {
#                     "type": "language",
#                     "difficulty": "medium",
#                     "pattern": "language"
#                 }
#             }
#         ]

#         for example_data in default_examples:
#             example = QueryExample(
#                 id=str(uuid4()),
#                 **example_data
#             )
#             await self.vector_store.add_example(example)

#     # Rest of the methods remain the same as they are generic enough
#     # Only adding one new method specific to movies:

#     async def add_example(self, 
#                          natural_query: str,
#                          sql_query: str,
#                          schema_context: Optional[str] = None,
#                          metadata: Optional[dict] = None) -> bool:
#         """Add new example to knowledge base"""
#         try:
#             example = QueryExample(
#                 id=str(uuid4()),
#                 natural_query=natural_query,
#                 sql_query=sql_query,
#                 schema_context=schema_context,
#                 metadata=metadata or {}
#             )
#             return await self.vector_store.add_example(example)
#         except Exception as e:
#             print(f"Error adding example: {e}")
#             return False

#     async def find_similar_examples(self, 
#                                   query: str,
#                                   k: int = 3) -> List[QueryExample]:
#         """Find similar examples using vector similarity"""
#         return await self.vector_store.find_similar(query, k)

#     def get_pattern_template(self, pattern_type: str) -> Optional[QueryPattern]:
#         """Get pattern template by type"""
#         return self.patterns.get(pattern_type)

#     def identify_patterns(self, query: str) -> List[str]:
#         """Identify potential patterns in a query"""
#         matched_patterns = []
#         for pattern_type, pattern in self.patterns.items():
#             # Simple keyword matching - could be enhanced with more sophisticated matching
#             keywords = pattern.example_query.lower().split()
#             if any(keyword in query.lower() for keyword in keywords):
#                 matched_patterns.append(pattern_type)
#         return matched_patterns

#     def log_query(self, 
#                  natural_query: str,
#                  generated_sql: str,
#                  execution_time: float,
#                  success: bool,
#                  error: Optional[str] = None):
#         """Log query execution for analysis"""
#         log = QueryLog(
#             query_id=str(uuid4()),
#             timestamp=datetime.now(),
#             natural_query=natural_query,
#             generated_sql=generated_sql,
#             execution_time=execution_time,
#             success=success,
#             error=error
#         )
#         self.query_logs.append(log)

#     def get_query_statistics(self) -> Dict[str, Any]:
#         """Get statistics about query executions"""
#         total_queries = len(self.query_logs)
#         if total_queries == 0:
#             return {
#                 "total_queries": 0,
#                 "success_rate": 0,
#                 "average_execution_time": 0
#             }

#         successful_queries = len([log for log in self.query_logs if log.success])
#         total_execution_time = sum(log.execution_time for log in self.query_logs)

#         return {
#             "total_queries": total_queries,
#             "success_rate": successful_queries / total_queries,
#             "average_execution_time": total_execution_time / total_queries
#         }

#     async def find_similar_by_genre(self, 
#                                   query: str,
#                                   genre: str,
#                                   k: int = 3) -> List[QueryExample]:
#         """Find similar examples for a specific movie genre"""
#         all_similar = await self.find_similar_examples(query, k=k*2)
#         genre_examples = []
        
#         for ex in all_similar:
#             sql_lower = ex.sql_query.lower()
#             if f"genres like '%{genre.lower()}%'" in sql_lower:
#                 genre_examples.append(ex)
        
#         return genre_examples[:k]

# app/knowledge/base.py
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