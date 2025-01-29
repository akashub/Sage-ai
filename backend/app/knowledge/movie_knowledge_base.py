# app/knowledge/movie_knowledge_base.py
from typing import Dict, List, Any, Optional
from .base import KnowledgeBase, QueryPattern
from langchain_openai import OpenAIEmbeddings
from .vector_store import VectorStore

class MovieKnowledgeBase(KnowledgeBase):
    def __init__(self, 
                 embedding_model: Optional[OpenAIEmbeddings] = None,
                 vector_store: Optional[VectorStore] = None):
        super().__init__(embedding_model=embedding_model, vector_store=vector_store)

    def _initialize_patterns(self) -> Dict[str, QueryPattern]:
        """Initialize movie-specific query patterns"""
        patterns = super()._initialize_patterns()  # Get base patterns
        
        # Add movie-specific patterns
        movie_patterns = {
            "rating": QueryPattern(
                pattern_type="rating",
                description="Rating-based movie queries",
                example_query="Find highest rated movies",
                sql_template="SELECT * FROM movies WHERE vote_average > {rating} AND vote_count > {min_votes}",
                metadata={"complexity": "medium"}
            ),
            "genre": QueryPattern(
                pattern_type="genre",
                description="Genre-based movie queries",
                example_query="Find action movies",
                sql_template="SELECT * FROM movies WHERE genres LIKE '%{genre}%'",
                metadata={"complexity": "easy"}
            ),
            "language": QueryPattern(
                pattern_type="language",
                description="Language-based queries",
                example_query="Show English language movies",
                sql_template="SELECT * FROM movies WHERE original_language = '{language}'",
                metadata={"complexity": "easy"}
            ),
            "popularity": QueryPattern(
                pattern_type="popularity",
                description="Popularity-based queries",
                example_query="Show most popular movies",
                sql_template="SELECT * FROM movies ORDER BY popularity DESC LIMIT {limit}",
                metadata={"complexity": "easy"}
            ),
            "revenue": QueryPattern(
                pattern_type="revenue",
                description="Revenue and budget queries",
                example_query="Show highest grossing movies",
                sql_template="SELECT *, (revenue - budget) as profit FROM movies WHERE revenue > 0 ORDER BY revenue DESC",
                metadata={"complexity": "medium"}
            ),
            "combined_search": QueryPattern(
                pattern_type="combined_search",
                description="Multi-criteria movie search",
                example_query="Find popular action movies with high ratings",
                sql_template="""
                SELECT * FROM movies 
                WHERE genres LIKE '%{genre}%' 
                AND vote_average > {rating} 
                AND vote_count > {min_votes}
                ORDER BY popularity DESC
                """,
                metadata={"complexity": "hard"}
            )
        }
        
        patterns.update(movie_patterns)
        return patterns

    async def _initialize_default_examples(self):
        """Initialize knowledge base with movie-specific examples"""
        movie_examples = [
            {
                "natural_query": "Show me movies released in 2022",
                "sql_query": """
                SELECT title, release_date, vote_average 
                FROM movies 
                WHERE EXTRACT(YEAR FROM release_date) = 2022 
                ORDER BY vote_average DESC
                """,
                "schema_context": "movies schema",
                "metadata": {
                    "type": "temporal",
                    "difficulty": "easy",
                    "pattern": "temporal"
                }
            },
            {
                "natural_query": "What are the most profitable movies with budget under 50 million",
                "sql_query": """
                SELECT 
                    title, 
                    budget, 
                    revenue, 
                    (revenue - budget) as profit
                FROM movies 
                WHERE budget < 50000000 
                AND revenue > budget
                ORDER BY (revenue - budget) DESC 
                LIMIT 10
                """,
                "schema_context": "movies schema",
                "metadata": {
                    "type": "revenue",
                    "difficulty": "medium",
                    "pattern": "revenue"
                }
            },
            {
                "natural_query": "Find action movies with rating above 8",
                "sql_query": """
                SELECT 
                    title, 
                    vote_average, 
                    popularity,
                    release_date
                FROM movies 
                WHERE vote_average > 8 
                AND genres LIKE '%Action%'
                AND vote_count > 1000
                ORDER BY popularity DESC 
                LIMIT 10
                """,
                "schema_context": "movies schema",
                "metadata": {
                    "type": "combined_search",
                    "difficulty": "hard",
                    "pattern": "combined_search"
                }
            },
            {
                "natural_query": "List highest rated Horror movies from 2020",
                "sql_query": """
                SELECT 
                    title, 
                    vote_average,
                    release_date,
                    popularity
                FROM movies 
                WHERE genres LIKE '%Horror%'
                AND EXTRACT(YEAR FROM release_date) = 2020
                AND vote_count > 100
                ORDER BY vote_average DESC 
                LIMIT 15
                """,
                "schema_context": "movies schema",
                "metadata": {
                    "type": "combined_search",
                    "difficulty": "medium",
                    "pattern": "genre"
                }
            }
        ]

        for example_data in movie_examples:
            await self.add_example(**example_data)

    async def find_movies_by_genre(self, 
                                 query: str,
                                 genre: str,
                                 k: int = 3) -> List[Any]:
        """Find similar examples for specific movie genre"""
        all_similar = await self.find_similar_examples(query, k=k*2)
        genre_examples = []
        
        for ex in all_similar:
            sql_lower = ex.sql_query.lower()
            if f"genres like '%{genre.lower()}%'" in sql_lower:
                genre_examples.append(ex)
        
        return genre_examples[:k]

    async def find_by_rating_range(self, 
                                 query: str,
                                 min_rating: float,
                                 k: int = 3) -> List[Any]:
        """Find similar examples for specific rating range"""
        all_similar = await self.find_similar_examples(query, k=k*2)
        rating_examples = []
        
        for ex in all_similar:
            sql_lower = ex.sql_query.lower()
            if f"vote_average > {min_rating}" in sql_lower:
                rating_examples.append(ex)
        
        return rating_examples[:k]

    def get_movie_schema(self) -> str:
        """Get the movies table schema"""
        return """
        movies(
            id INTEGER,
            title TEXT,
            vote_average FLOAT,
            vote_count INTEGER,
            status TEXT,
            release_date DATE,
            revenue BIGINT,
            runtime INTEGER,
            adult BOOLEAN,
            budget BIGINT,
            imdb_id TEXT,
            original_language TEXT,
            original_title TEXT,
            overview TEXT,
            popularity FLOAT,
            tagline TEXT,
            genres TEXT,
            production_companies TEXT,
            production_countries TEXT,
            spoken_languages TEXT,
            keywords TEXT
        )
        """