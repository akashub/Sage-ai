# # app/orchestration/nodes/analyzer.py
# from typing import Dict, Any, Optional
# from langchain_openai import ChatOpenAI
# from pydantic import BaseModel

# class AnalyzerNode:
#     def __init__(self, llm: Optional[ChatOpenAI] = None):
#         self.llm = llm or ChatOpenAI(temperature=0)

#     async def analyze(self, question: str, schema: Optional[str] = None) -> Dict[str, Any]:
#         """Analyze natural language query to understand intent"""
#         try:
#             prompt = f"""Analyze this natural language query for SQL conversion.
#             Question: {question}
            
#             Schema Context (if available):
#             {schema or 'No schema provided'}
            
#             Provide analysis in JSON format:
#             {{
#                 "query_type": "select/insert/update/delete",
#                 "tables": ["required_tables"],
#                 "columns": ["required_columns"],
#                 "conditions": ["where_conditions"],
#                 "joins": ["needed_joins"],
#                 "aggregations": ["any_aggregations"],
#                 "group_by": ["group_columns"],
#                 "order_by": ["order_columns"],
#                 "limit": optional_limit_value
#             }}
#             """
            
#             response = await self.llm.ainvoke(prompt)
#             return {
#                 "success": True,
#                 "analysis": response.content,
#                 "error": None
#             }
#         except Exception as e:
#             return {
#                 "success": False,
#                 "analysis": None,
#                 "error": f"Analysis failed: {str(e)}"
#             }

from typing import Dict, Any, Optional, List
from langchain_openai import ChatOpenAI
from pydantic import BaseModel, Field
import json
import re
from ....app.llm.client import llm_client


class QueryPattern(BaseModel):
    pattern_type: str = Field(..., description="Type of query pattern (e.g., temporal, aggregation)")
    example_query: str = Field(..., description="Example of this pattern")
    sql_template: str = Field(..., description="SQL template for this pattern")

class QueryIntent(BaseModel):
    query_type: str
    tables: List[str]
    columns: List[str]
    conditions: Optional[List[str]] = None
    joins: Optional[List[str]] = None
    aggregations: Optional[List[str]] = None
    group_by: Optional[List[str]] = None
    order_by: Optional[List[str]] = None
    limit: Optional[int] = None
    patterns: Optional[List[str]] = None

class AnalyzerNode:
    # def __init__(self, llm: Optional[ChatOpenAI] = None):
    #     self.llm = llm or ChatOpenAI(temperature=0)
    #     self.common_patterns = self._initialize_patterns()
    
    def __init__(self):
        self.common_patterns = self._initialize_patterns()

    def _initialize_patterns(self) -> Dict[str, QueryPattern]:
        """Initialize common query patterns"""
        return {
            "temporal": QueryPattern(
                pattern_type="temporal",
                example_query="Show sales from last month",
                sql_template="SELECT * FROM {table} WHERE DATE_TRUNC('month', {date_column}) = DATE_TRUNC('month', CURRENT_DATE - INTERVAL '1 month')"
            ),
            "aggregation": QueryPattern(
                pattern_type="aggregation",
                example_query="Show total sales by category",
                sql_template="SELECT {group_by}, SUM({value_column}) as total FROM {table} GROUP BY {group_by}"
            ),
            "join": QueryPattern(
                pattern_type="join",
                example_query="Show customer orders with their details",
                sql_template="SELECT * FROM {table1} JOIN {table2} ON {table1}.{key1} = {table2}.{key2}"
            ),
            "ranking": QueryPattern(
                pattern_type="ranking",
                example_query="Show top 5 customers by order value",
                sql_template="SELECT * FROM {table} ORDER BY {order_column} DESC LIMIT {limit}"
            )
        }

    def _extract_table_references(self, question: str, schema: Optional[str]) -> List[str]:
        """Extract potential table references from the question"""
        tables = []
        if schema:
            # Extract table names from schema
            schema_tables = re.findall(r'CREATE TABLE (\w+)|Table: (\w+)', schema)
            schema_tables = [t[0] or t[1] for t in schema_tables]
            
            # Look for these tables in the question
            for table in schema_tables:
                if table.lower() in question.lower():
                    tables.append(table)
        return tables

    def _identify_patterns(self, question: str) -> List[str]:
        """Identify common query patterns in the question"""
        patterns = []
        pattern_indicators = {
            "temporal": r"last|next|previous|recent|today|yesterday|month|year|week|date",
            "aggregation": r"total|sum|average|count|minimum|maximum|avg|min|max",
            "join": r"with|including|related|connected|associated",
            "ranking": r"top|bottom|highest|lowest|first|last|best|worst"
        }
        
        for pattern, indicator in pattern_indicators.items():
            if re.search(indicator, question, re.IGNORECASE):
                patterns.append(pattern)
        
        return patterns

    async def analyze(self, 
                 question: str, 
                 schema: Optional[str] = None,
                 examples: Optional[List[Any]] = None,
                 patterns: Optional[List[str]] = None) -> Dict[str, Any]:
        """Analyze natural language query to understand intent"""
        try:
            # Pre-process and extract initial information
            identified_patterns = self._identify_patterns(question)
            table_references = self._extract_table_references(question, schema)
            
            # Create the analysis prompt with examples from knowledge layer
            prompt = self._create_analysis_prompt(
            question=question,
            schema=schema,
            patterns=identified_patterns,
            tables=table_references,
            examples=examples
        )
            print("Prompt sent to Gemini")
            # Get LLM analysis
            response = await llm_client.get_completion(prompt)
            print("Response from Gemini:\n{response}")
            
            try:
                # Parse the JSON response
                analysis = json.loads(response)
                
                # Validate the analysis
                QueryIntent(**analysis)
                
                # Add identified patterns if not present
                if not analysis.get("patterns"):
                    analysis["patterns"] = identified_patterns
                
                return {
                    "success": True,
                    "analysis": analysis,
                    "error": None
                }
            except json.JSONDecodeError:
                return {
                    "success": False,
                    "analysis": None,
                    "error": "Failed to parse analysis response"
                }
            except Exception as e:
                return {
                    "success": False,
                    "analysis": None,
                    "error": f"Invalid analysis format: {str(e)}"
                }
                
        except Exception as e:
            return {
                "success": False,
                "analysis": None,
                "error": f"Analysis failed: {str(e)}"
            }

    def _create_analysis_prompt(self, 
                              question: str, 
                              schema: Optional[str],
                              patterns: List[str],
                              tables: List[str],
                              examples: Optional[List[Any]] = None) -> str:
        """Create the analysis prompt for the LLM"""
        
        # Start with base prompt
        prompt = f"""Analyze this natural language query for SQL conversion.
Question: {question}

Schema Context:
{schema or 'No schema provided'}

Identified Patterns: {', '.join(patterns) if patterns else 'None'}
Potential Tables: {', '.join(tables) if tables else 'None'}

Requirements:
1. Identify the correct query type (SELECT, INSERT, UPDATE, DELETE)
2. List all required tables
3. Specify required columns
4. Determine any conditions, joins, or aggregations needed
5. Identify any grouping, ordering, or limits

Provide analysis in JSON format:
{{
    "query_type": "select/insert/update/delete",
    "tables": ["required_tables"],
    "columns": ["required_columns"],
    "conditions": ["where_conditions"],
    "joins": ["needed_joins"],
    "aggregations": ["any_aggregations"],
    "group_by": ["group_columns"],
    "order_by": ["order_columns"],
    "limit": optional_limit_value,
    "patterns": ["identified_patterns"]
}}
"""

        # Add pattern-specific guidance if patterns were identified
        if patterns:
            prompt += "\nRelevant Patterns:\n"
            for pattern in patterns:
                if pattern in self.common_patterns:
                    pattern_info = self.common_patterns[pattern]
                    prompt += f"""
Pattern: {pattern_info.pattern_type}
Example: {pattern_info.example_query}
Template: {pattern_info.sql_template}
"""
        if examples:
            prompt += "\nSimilar Examples:\n"
            for i, example in enumerate(examples, 1):
                prompt += f"""
    Example {i}:
    Natural Query: {example.natural_query}
    SQL: {example.sql_query}
"""

        return prompt

    def get_pattern_template(self, pattern_type: str) -> Optional[str]:
        """Get SQL template for a specific pattern"""
        pattern = self.common_patterns.get(pattern_type)
        return pattern.sql_template if pattern else None