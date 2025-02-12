# # # # backend/python/app/llm/client.py
# # # from datetime import time
# # # import google.generativeai as genai
# # # from typing import Dict, Any
# # # import os
# # # import json
# # # import asyncio

# # # class LLMClient:
# # #     def __init__(self):
# # #         genai.configure(api_key=os.getenv("GEMINI_API_KEY"))
# # #         self.model = genai.GenerativeModel(os.getenv("GEMINI_MODEL_NAME", "gemini-1.5-pro"))
# # #         self._last_call_time = 0
# # #         self._rate_limit_delay = float(os.getenv("RATE_LIMIT_DELAY", "1.0"))

# # #     async def _handle_rate_limit(self):
# # #         """Handle rate limiting"""
# # #         now = time.time()
# # #         time_since_last_call = now - self._last_call_time
# # #         if time_since_last_call < self._rate_limit_delay:
# # #             await asyncio.sleep(self._rate_limit_delay - time_since_last_call)
# # #         self._last_call_time = time.time()

# # #     async def analyze_query(self, question: str, schema: Dict[str, Any]) -> Dict[str, Any]:
# # #         await self._handle_rate_limit()
        
# # #         prompt = f"""Analyze this natural language query:
# # #         Question: {question}
        
# # #         Available CSV Schema:
# # #         {json.dumps(schema, indent=2)}
        
# # #         Return JSON with:
# # #         1. Required columns
# # #         2. Operations needed
# # #         3. Conditions/filters
# # #         4. Any aggregations
        
# # #         Format:
# # #         {{
# # #             "columns": ["col1", "col2"],
# # #             "operations": ["sum", "group_by"],
# # #             "conditions": ["date > '2024-01-01'"],
# # #             "aggregations": [
# # #                 {{"type": "sum", "column": "amount"}}
# # #             ]
# # #         }}
# # #         """
        
# # #         response = await self.model.generate_content(prompt)
# # #         return json.loads(response.text)

# # #     async def generate_query(self, analysis: Dict[str, Any], schema: Dict[str, Any]) -> str:
# # #         await self._handle_rate_limit()
        
# # #         prompt = f"""Generate a pandas query based on:
        
# # #         Analysis:
# # #         {json.dumps(analysis, indent=2)}
        
# # #         Schema:
# # #         {json.dumps(schema, indent=2)}
        
# # #         Return only the pandas code to execute this query.
# # #         """
        
# # #         response = await self.model.generate_content(prompt)
# # #         return response.text

# # #     async def validate_query(self, query: str, schema: Dict[str, Any]) -> Dict[str, Any]:
# # #         await self._handle_rate_limit()
        
# # #         prompt = f"""Validate this pandas query:
        
# # #         Query:
# # #         {query}
        
# # #         Schema:
# # #         {json.dumps(schema, indent=2)}
        
# # #         Return JSON:
# # #         {{
# # #             "isValid": true/false,
# # #             "issues": ["issue1", "issue2"]
# # #         }}
# # #         """
        
# # #         response = await self.model.generate_content(prompt)
# # #         return json.loads(response.text)

# # # llm_client = LLMClient()


# # # backend/python/app/llm/client.py
# # import google.generativeai as genai
# # from typing import Dict, Any, List
# # from dotenv import load_dotenv
# # import os
# # import json
# # import asyncio
# # import time
# # from enum import Enum

# # load_dotenv()

# # class QueryType(Enum):
# #     SELECT = "select"
# #     AGGREGATE = "aggregate"
# #     JOIN = "join"
# #     COMPLEX = "complex"

# # class LLMClient:
# #     def __init__(self):
# #         # genai.configure(api_key=os.getenv("GEMINI_API_KEY"))
# #         # self.model = genai.GenerativeModel(os.getenv("GEMINI_MODEL_NAME", "gemini-1.5-pro"))
# #         # self._last_call_time = 0
# #         # self._rate_limit_delay = float(os.getenv("RATE_LIMIT_DELAY", "1.0"))
# #         def __init__(self):
# #             api_key = os.getenv("GEMINI_API_KEY")
# #             if not api_key:
# #                 raise ValueError("GEMINI_API_KEY not found in environment variables")
                
# #             genai.configure(api_key=api_key)
# #             self.model = genai.GenerativeModel(os.getenv("GEMINI_MODEL_NAME", "gemini-1.5-pro"))
# #             self._last_call_time = 0
# #             self._rate_limit_delay = float(os.getenv("RATE_LIMIT_DELAY", "1.0"))

# #     async def _handle_rate_limit(self):
# #         now = time.time()
# #         time_since_last_call = now - self._last_call_time
# #         if time_since_last_call < self._rate_limit_delay:
# #             await asyncio.sleep(self._rate_limit_delay - time_since_last_call)
# #         self._last_call_time = time.time()

# #     async def analyze_query(self, question: str, schema: Dict[str, Any]) -> Dict[str, Any]:
# #         await self._handle_rate_limit()
        
# #         prompt = f"""Let's analyze this natural language query step by step.

# # Question: {question}

# # CSV Schema:
# # {json.dumps(schema, indent=2)}

# # Follow these steps:

# # 1. Identify key components:
# #    - What information is being requested?
# #    - Are there any temporal conditions (dates, time periods)?
# #    - Are there any specific filters or conditions?
# #    - Is aggregation needed (sum, average, count)?

# # 2. Map requirements to available columns:
# #    - Which columns are needed to answer this query?
# #    - Are there any required transformations?
# #    - What data types are we working with?

# # 3. Determine query complexity:
# #    - Is this a simple SELECT?
# #    - Do we need aggregations?
# #    - Are multiple conditions required?
# #    - Do we need any joins or complex operations?

# # Return a detailed analysis in JSON format:
# # {{
# #     "query_type": "select/aggregate/join/complex",
# #     "required_columns": ["col1", "col2"],
# #     "transformations": [
# #         {{"column": "col1", "type": "date_format/number/text", "details": "..."}},
# #     ],
# #     "conditions": [
# #         {{"column": "col1", "operator": "=/>/<", "value": "x", "logic": "AND/OR"}}
# #     ],
# #     "aggregations": [
# #         {{"function": "sum/avg/count", "column": "col1", "alias": "result"}}
# #     ],
# #     "group_by": ["col1"],
# #     "order_by": [
# #         {{"column": "col1", "direction": "asc/desc"}}
# #     ],
# #     "reasoning": "Detailed explanation of the analysis...",
# #     "suggested_approach": "Step-by-step approach to generate the query..."
# # }}

# # Think through this carefully and provide a comprehensive analysis."""
        
# #         response = await self.model.generate_content(prompt)
# #         return json.loads(response.text)

# #     async def generate_query(self, analysis: Dict[str, Any], schema: Dict[str, Any]) -> str:
# #         await self._handle_rate_limit()
        
# #         prompt = f"""Generate an SQL query based on this analysis.

# # Analysis:
# # {json.dumps(analysis, indent=2)}

# # Schema:
# # {json.dumps(schema, indent=2)}

# # Follow these steps:

# # 1. Review the analysis and schema:
# #    - Confirm all required columns exist
# #    - Verify data types match operations
# #    - Check for any potential issues

# # 2. Build the query components:
# #    - Start with SELECT clause
# #    - Add FROM and any necessary JOINs
# #    - Include WHERE conditions
# #    - Add GROUP BY if aggregating
# #    - Include HAVING if needed
# #    - Add ORDER BY and LIMIT

# # 3. Validate query structure:
# #    - Check syntax
# #    - Verify column references
# #    - Confirm logical operations
# #    - Ensure proper ordering of clauses

# # 4. Explain your reasoning:
# #    - Why this query structure was chosen
# #    - Any assumptions made
# #    - Potential alternatives considered

# # Return both the SQL query and your reasoning in JSON format:
# # {{
# #     "sql": "The complete SQL query",
# #     "explanation": "Step-by-step explanation of the query construction",
# #     "alternatives": "Any alternative approaches considered",
# #     "assumptions": "List of assumptions made",
# #     "potential_issues": "Any potential issues to be aware of"
# # }}

# # Generate a clear, efficient, and well-structured query."""
        
# #         response = await self.model.generate_content(prompt)
# #         return json.loads(response.text)

# #     async def validate_query(self, query: str, schema: Dict[str, Any]) -> Dict[str, Any]:
# #         await self._handle_rate_limit()
        
# #         prompt = f"""Validate this SQL query thoroughly.

# # Query:
# # {query}

# # Schema:
# # {json.dumps(schema, indent=2)}

# # Perform a comprehensive validation:

# # 1. Syntax Analysis:
# #    - Check SQL syntax correctness
# #    - Verify clause ordering
# #    - Validate parentheses and operators
# #    - Check for SQL injection risks

# # 2. Schema Validation:
# #    - Verify all referenced columns exist
# #    - Check data type compatibility
# #    - Validate join conditions
# #    - Verify aggregation usage

# # 3. Logical Validation:
# #    - Check condition logic
# #    - Verify group by completeness
# #    - Validate aggregate functions
# #    - Check for logical contradictions

# # 4. Performance Considerations:
# #    - Identify potential bottlenecks
# #    - Check for inefficient patterns
# #    - Suggest optimizations

# # Return detailed validation results in JSON:
# # {{
# #     "isValid": true/false,
# #     "syntaxIssues": ["list of syntax issues"],
# #     "schemaIssues": ["list of schema-related issues"],
# #     "logicalIssues": ["list of logical problems"],
# #     "performanceIssues": ["list of performance concerns"],
# #     "suggestedFixes": [
# #         {{
# #             "issue": "Description of the issue",
# #             "fix": "Suggested fix",
# #             "priority": "high/medium/low"
# #         }}
# #     ],
# #     "requiresHealing": true/false,
# #     "healingType": "generator/analyzer",
# #     "healingDetails": "Explanation of why healing is needed and what approach to take"
# # }}

# # Provide thorough validation and clear suggestions for improvement."""
        
# #         response = await self.model.generate_content(prompt)
# #         return json.loads(response.text)

# #     async def heal_query(self, validation_result: Dict[str, Any], original_query: str, 
# #                         analysis: Dict[str, Any], schema: Dict[str, Any]) -> Dict[str, Any]:
# #         await self._handle_rate_limit()
        
# #         prompt = f"""Help heal this SQL query based on validation issues.

# # Original Query:
# # {original_query}

# # Validation Issues:
# # {json.dumps(validation_result, indent=2)}

# # Original Analysis:
# # {json.dumps(analysis, indent=2)}

# # Schema:
# # {json.dumps(schema, indent=2)}

# # Follow these healing steps:

# # 1. Analyze Issues:
# #    - Review each validation issue
# #    - Identify root causes
# #    - Determine if analysis or generation needs revision

# # 2. Propose Fixes:
# #    - Address each issue systematically
# #    - Consider multiple approaches
# #    - Evaluate trade-offs

# # 3. Generate New Query:
# #    - Apply fixes
# #    - Verify fixes don't introduce new issues
# #    - Ensure query maintains original intent

# # Return healing results in JSON:
# # {{
# #     "healed_query": "The fixed SQL query",
# #     "changes_made": [
# #         {{
# #             "issue": "Original issue",
# #             "fix": "Applied fix",
# #             "reasoning": "Why this fix was chosen"
# #         }}
# #     ],
# #     "validation_needed": true/false,
# #     "confidence": 0.0-1.0,
# #     "requires_human_review": true/false,
# #     "notes": "Additional observations or warnings"
# # }}

# # Focus on maintaining query intent while fixing issues."""
        
# #         response = await self.model.generate_content(prompt)
# #         return json.loads(response.text)

# # llm_client = LLMClient()

# # backend/python/app/llm/client.py
# import os
# import logging
# from dotenv import load_dotenv
# from fastapi import HTTPException
# import google.generativeai as genai
# from typing import Dict, Any
# import json
# import asyncio
# import time
# from datetime import datetime

# # Load environment variables
# load_dotenv()

# # Configure logging
# logging.basicConfig(
#     level=logging.INFO,
#     format='%(asctime)s - %(name)s - %(levelname)s - %(message)s',
#     handlers=[
#         logging.FileHandler('llm_service.log'),
#         logging.StreamHandler()
#     ]
# )
# logger = logging.getLogger(__name__)

# class LLMClient:
#     def __init__(self):
#         self._last_call_time = time.time()
#         self._rate_limit_delay = float(os.getenv("RATE_LIMIT_DELAY", "1.0"))
#         self._context = {}  # Store context per session

#         api_key = os.getenv("GEMINI_API_KEY")
#         if not api_key:
#             logger.error("GEMINI_API_KEY not found in environment variables")
#             raise ValueError("GEMINI_API_KEY not found in environment variables")
            
#         genai.configure(api_key=api_key)
#         self.model = genai.GenerativeModel(os.getenv("GEMINI_MODEL_NAME", "gemini-1.5-pro"))
#         logger.info("LLM Client initialized successfully")
#         self.safety_settings = {
#             "HARM_CATEGORY_HARASSMENT": "BLOCK_NONE",
#             "HARM_CATEGORY_HATE_SPEECH": "BLOCK_NONE",
#             "HARM_CATEGORY_SEXUALLY_EXPLICIT": "BLOCK_NONE",
#             "HARM_CATEGORY_DANGEROUS_CONTENT": "BLOCK_NONE",
#         }
#         self._last_call_time = time.time()
#         self._rate_limit_delay = float(os.getenv("RATE_LIMIT_DELAY", "1.0"))
#         self.timeout = 45

#     async def _handle_rate_limit(self):
#         now = time.time()
#         time_since_last_call = now - self._last_call_time
#         if time_since_last_call < self._rate_limit_delay:
#             delay = self._rate_limit_delay - time_since_last_call
#             logger.debug(f"Rate limiting: waiting for {delay:.2f} seconds")
#             await asyncio.sleep(delay)
#         self._last_call_time = time.time()

#     # def _generate_content(self, prompt: str) -> str:
#     #     """Synchronous wrapper for content generation with retries and error handling"""
#     #     max_retries = 3
#     #     for attempt in range(max_retries):
#     #         try:
#     #             response = self.model.generate_content(
#     #                 prompt,
#     #                 generation_config={
#     #                     'temperature': 0.3,
#     #                     'top_p': 0.8,
#     #                     'top_k': 40,
#     #                     'max_output_tokens': 2048,
#     #                 },
#     #                 safety_settings=self.safety_settings
#     #             )
                
#     #             if not response.text:
#     #                 raise ValueError("Empty response from Gemini")
                
#     #             return response.text
                
#     #         except Exception as e:
#     #             logger.error(f"Attempt {attempt + 1} failed: {str(e)}")
#     #             if attempt == max_retries - 1:
#     #                 raise
#     #             time.sleep(1)  # Wait before retry

#     async def generate_query(self, analysis: Dict[str, Any], schema: Dict[str, Any]) -> str:
#         try:
#             async with asyncio.timeout(self.timeout):
#                 prompt = self._create_generation_prompt(analysis, schema)
#                 response = await self._generate_with_retry(prompt)
#                 return response
#         except asyncio.TimeoutError:
#             logger.error("LLM request timed out")
#             raise HTTPException(status_code=504, detail="LLM request timed out")
#         except Exception as e:
#             logger.error(f"Error in generate_query: {str(e)}", exc_info=True)
#             raise

#     async def _generate_with_retry(self, prompt: str, max_retries: int = 3) -> str:
#         for attempt in range(max_retries):
#             try:
#                 response = await asyncio.wait_for(
#                     asyncio.to_thread(self.model.generate_content, prompt),
#                     timeout=self.timeout
#                 )
#                 return response.text
#             except Exception as e:
#                 logger.error(f"Attempt {attempt + 1} failed: {str(e)}")
#                 if attempt == max_retries - 1:
#                     raise
#                 await asyncio.sleep(1)  # Wait before retry

#     async def analyze_query(self, question: str, schema: Dict[str, Any]) -> Dict[str, Any]:
#         logger.info(f"Analyzing query: {question}")
#         await self._handle_rate_limit()

#         session_id = "default"  
#         self._context[session_id] = {
#             "last_question": question,
#             "schema": schema
#         }
        
#         try:
#             prompt = f"""You are a SQL query analyzer for a movie database. Analyze this question:

# Question: {question}

# Available columns in the database:
# {json.dumps(schema, indent=2)}

# Follow these steps:
# 1. First, identify the core requirements:
#    - What type of movies are we looking for? (Horror, Action, etc.)
#    - What's the main metric for filtering/sorting? (revenue, rating, etc.)
#    - How many results are needed?

# 2. Then, determine the necessary columns:
#    - We definitely need 'title' for movie names
#    - Include 'revenue' for financial queries
#    - Include 'genres' for genre filtering
#    - Add any other relevant columns

# 3. Create the appropriate conditions:
#    - Genre filtering using LIKE
#    - Any numerical thresholds
#    - Date ranges if needed

# 4. Specify sorting and limits:
#    - Order by the main metric
#    - Apply any LIMIT clause

# Return ONLY a JSON object in this exact format:
# {{
#     "query_type": "select",
#     "required_columns": ["title", "revenue", "genres"],
#     "conditions": ["genres LIKE '%Horror%'"],
#     "sort": {{"column": "revenue", "order": "desc"}},
#     "limit": 10,
#     "explanation": "Finding top 10 horror movies by revenue"
# }}

# Do not include any other text or explanation outside the JSON object."""

#             logger.debug(f"Generated prompt:\n{prompt}")
            
#             response_text = await asyncio.get_event_loop().run_in_executor(
#                 None, self._generate_content, prompt
#             )
            
#             logger.debug(f"Raw LLM response:\n{response_text}")
            
#             # Clean the response
#             response_text = response_text.strip()
#             if response_text.startswith("```json"):
#                 response_text = response_text[7:]
#             if response_text.endswith("```"):
#                 response_text = response_text[:-3]
#             response_text = response_text.strip()
            
#             try:
#                 analysis = json.loads(response_text)
#             except json.JSONDecodeError as e:
#                 logger.error(f"Failed to parse JSON: {response_text}")
#                 raise ValueError(f"Invalid JSON response from LLM: {str(e)}")
                
#             logger.info("Query analysis completed successfully")
#             return analysis

#         except Exception as e:
#             logger.error(f"Error in analyze_query: {str(e)}", exc_info=True)
#             raise

# #     async def generate_query(self, analysis: Dict[str, Any], schema: Dict[str, Any]) -> str:
# #         logger.info("Generating SQL query")
# #         await self._handle_rate_limit()

# #         try:
# #             prompt = f"""Generate a SQL query based on this analysis:

# # Analysis:
# # {json.dumps(analysis, indent=2)}

# # Schema:
# # {json.dumps(schema, indent=2)}

# # Consider:
# # 1. Proper column names and quotes
# # 2. Join conditions if needed
# # 3. Proper SQL syntax for movie-related queries
# # 4. Efficient query structure

# # Return only the SQL query without any explanations."""

# #             logger.debug(f"Generated prompt:\n{prompt}")
# #             response_text = await asyncio.get_event_loop().run_in_executor(
# #                 None, self._generate_content, prompt
# #             )
# #             query = response_text.strip()
# #             logger.debug(f"Generated SQL query:\n{query}")
            
# #             return query

# #         except Exception as e:
# #             logger.error(f"Error in generate_query: {str(e)}", exc_info=True)
# #             raise

#     async def generate_query(self, analysis: Dict[str, Any], schema: Dict[str, Any]) -> str:
#             logger.info("Generating SQL query")
#             await self._handle_rate_limit()

#             try:
#                 prompt = f"""Generate a SQL query based on this analysis for a movie database:

#     Analysis:
#     {json.dumps(analysis, indent=2)}

#     Schema:
#     {json.dumps(schema, indent=2)}

#     Follow these steps:
#     1. Use exact column names from the schema
#     2. Format genre conditions correctly (using LIKE)
#     3. Include proper SQL syntax for movies
#     4. Return response in this exact JSON format:
#     {{
#         "query": "SELECT [your SQL query here]",
#         "explanation": "Brief explanation of the query"
#     }}

#     Do not include any other text or explanation outside the JSON object."""

#                 logger.debug(f"Generated prompt:\n{prompt}")
#                 response_text = await asyncio.get_event_loop().run_in_executor(
#                     None, self._generate_content, prompt
#                 )
                
#                 # Clean the response
#                 response_text = response_text.strip()
#                 if response_text.startswith("```json"):
#                     response_text = response_text[7:]
#                 if response_text.endswith("```"):
#                     response_text = response_text[:-3]
#                 response_text = response_text.strip()
                
#                 try:
#                     response = json.loads(response_text)
#                     logger.debug(f"Generated SQL query:\n{response['query']}")
#                     return response
#                 except json.JSONDecodeError as e:
#                     logger.error(f"Failed to parse JSON: {response_text}")
#                     raise ValueError(f"Invalid JSON response from LLM: {str(e)}")

#             except Exception as e:
#                 logger.error(f"Error in generate_query: {str(e)}", exc_info=True)
#                 raise
# #     async def validate_query(self, query: str, schema: Dict[str, Any]) -> Dict[str, Any]:
# #         logger.info("Validating SQL query")
# #         await self._handle_rate_limit()

# #         try:
# #             prompt = f"""Validate this SQL query for a movie database:

# # Query:
# # {query}

# # Schema:
# # {json.dumps(schema, indent=2)}

# # Check for:
# # 1. SQL syntax correctness
# # 2. Column name validity
# # 3. Proper conditions and joins
# # 4. Performance considerations

# # Return JSON:
# # {{
# #     "isValid": true/false,
# #     "issues": ["list", "of", "issues"],
# #     "suggestedFixes": {{
# #         "issue": "description",
# #         "fix": "suggested solution"
# #     }},
# #     "explanation": "Detailed validation explanation"
# # }}"""

# #             logger.debug(f"Generated prompt:\n{prompt}")
# #             response_text = await asyncio.get_event_loop().run_in_executor(
# #                 None, self._generate_content, prompt
# #             )
# #             logger.debug(f"Raw validation response:\n{response_text}")
            
# #             validation = json.loads(response_text)
# #             logger.info(f"Query validation completed. Valid: {validation.get('isValid', False)}")
# #             return validation

# #         except Exception as e:
# #             logger.error(f"Error in validate_query: {str(e)}", exc_info=True)
# #             raise

#     async def validate_query(self, query: str, schema: Dict[str, Any]) -> Dict[str, Any]:
#         logger.info("Validating SQL query")
#         await self._handle_rate_limit()

#         try:
#             prompt = f"""Validate this SQL query for a movie database:

# SQL Query:
# {query}

# Available Schema:
# {json.dumps(schema, indent=2)}

# Validate the following:
# 1. SQL syntax correctness
# 2. Column names match the schema exactly
# 3. Proper use of conditions for movie queries
# 4. Genre filtering syntax
# 5. Performance considerations

# Return ONLY a JSON object in this exact format:
# {{
#     "isValid": true or false,
#     "issues": ["issue1", "issue2"],
#     "suggestedFixes": ["fix1", "fix2"],
#     "explanation": "Detailed explanation"
# }}

# For a valid query example:
# {{
#     "isValid": true,
#     "issues": [],
#     "suggestedFixes": [],
#     "explanation": "Query is valid and uses correct schema"
# }}

# For an invalid query example:
# {{
#     "isValid": false,
#     "issues": ["Column name 'movie_name' doesn't exist", "Missing ORDER BY for LIMIT"],
#     "suggestedFixes": ["Use 'title' instead of 'movie_name'", "Add ORDER BY revenue DESC"],
#     "explanation": "Query needs column name correction and proper ordering"
# }}"""

#             logger.debug(f"Generated prompt:\n{prompt}")
#             response_text = await asyncio.get_event_loop().run_in_executor(
#                 None, self._generate_content, prompt
#             )
            
#             # Clean and validate response
#             response_text = response_text.strip()
#             if response_text.startswith("```json"):
#                 response_text = response_text[7:]
#             if response_text.endswith("```"):
#                 response_text = response_text[:-3]
#             response_text = response_text.strip()

#             try:
#                 validation = json.loads(response_text)
#                 # Ensure correct types
#                 validation["isValid"] = bool(validation.get("isValid", False))
#                 validation["issues"] = list(validation.get("issues", []))
#                 validation["suggestedFixes"] = list(validation.get("suggestedFixes", []))
#                 validation["explanation"] = str(validation.get("explanation", ""))
                
#                 logger.info(f"Validation completed: isValid={validation['isValid']}")
#                 return validation
#             except json.JSONDecodeError as e:
#                 logger.error(f"Failed to parse JSON: {response_text}")
#                 raise ValueError(f"Invalid JSON response from LLM: {str(e)}")

#         except Exception as e:
#             logger.error(f"Error in validate_query: {str(e)}", exc_info=True)
#             raise

#     async def heal_query(self, 
#                         validation_result: Dict[str, Any],
#                         original_query: str,
#                         analysis: Dict[str, Any],
#                         schema: Dict[str, Any]) -> Dict[str, Any]:
#         logger.info("Attempting to heal invalid query")
#         await self._handle_rate_limit()

#         try:
#             prompt = f"""Help fix this SQL query for a movie database:

# Original Query:
# {original_query}

# Validation Issues:
# {json.dumps(validation_result, indent=2)}

# Original Analysis:
# {json.dumps(analysis, indent=2)}

# Schema:
# {json.dumps(schema, indent=2)}

# Provide a solution in JSON:
# {{
#     "healed_query": "fixed SQL query",
#     "changes_made": [
#         {{
#             "issue": "what was wrong",
#             "fix": "how it was fixed",
#             "reasoning": "why this fix works"
#         }}
#     ],
#     "requires_reanalysis": false,
#     "confidence": 0.95
# }}"""

#             logger.debug(f"Generated prompt:\n{prompt}")
#             response_text = await asyncio.get_event_loop().run_in_executor(
#                 None, self._generate_content, prompt
#             )
#             logger.debug(f"Raw healing response:\n{response_text}")
            
#             healing_result = json.loads(response_text)
#             logger.info("Query healing completed")
#             return healing_result

#         except Exception as e:
#             logger.error(f"Error in heal_query: {str(e)}", exc_info=True)
#             raise

# # Create singleton instance
# llm_client = LLMClient()
# logger.info("LLM client singleton instance created")

# backend/python/app/llm/client.py
# import os
# import logging
# from dotenv import load_dotenv
# import google.generativeai as genai
# from typing import Dict, Any, List
# import json
# import asyncio
# import time
# from fastapi import HTTPException

# # Load environment variables
# load_dotenv()

# # Configure logging
# logging.basicConfig(
#     level=logging.INFO,
#     format='%(asctime)s - %(name)s - %(levelname)s - %(message)s',
#     handlers=[
#         logging.FileHandler('llm_service.log'),
#         logging.StreamHandler()
#     ]
# )
# logger = logging.getLogger(__name__)

# class LLMClient:
#     def __init__(self):
#         self._initialize_llm()
#         self._initialize_rate_limiting()
#         self._initialize_context_store()

#     def _initialize_llm(self):
#         """Initialize LLM configuration"""
#         api_key = os.getenv("GEMINI_API_KEY")
#         if not api_key:
#             raise ValueError("GEMINI_API_KEY not found in environment variables")
            
#         genai.configure(api_key=api_key)
#         self.model = genai.GenerativeModel(os.getenv("GEMINI_MODEL_NAME", "gemini-1.5-pro"))
#         self.generation_config = {
#             'temperature': float(os.getenv("LLM_TEMPERATURE", "0.3")),
#             'top_p': float(os.getenv("LLM_TOP_P", "0.8")),
#             'top_k': int(os.getenv("LLM_TOP_K", "40")),
#             'max_output_tokens': int(os.getenv("LLM_MAX_TOKENS", "2048")),
#         }
#         logger.info("LLM initialized successfully")

#     def _initialize_rate_limiting(self):
#         """Initialize rate limiting settings"""
#         self._last_call_time = time.time()
#         self._rate_limit_delay = float(os.getenv("RATE_LIMIT_DELAY", "1.0"))
#         self.timeout = int(os.getenv("LLM_TIMEOUT", "45"))
#         self.max_retries = int(os.getenv("LLM_MAX_RETRIES", "3"))

#     def _initialize_context_store(self):
#         """Initialize context storage"""
#         self._context = {}
#         self._schema_insights = {}

#     async def _handle_rate_limit(self):
#         """Handle rate limiting between requests"""
#         now = time.time()
#         time_since_last_call = now - self._last_call_time
#         if time_since_last_call < self._rate_limit_delay:
#             await asyncio.sleep(self._rate_limit_delay - time_since_last_call)
#         self._last_call_time = time.time()

#     async def _generate_with_retry(self, prompt: str) -> str:
#         """Generate content with retry logic"""
#         for attempt in range(self.max_retries):
#             try:
#                 await self._handle_rate_limit()
#                 logger.debug(f"Sending prompt to LLM (attempt {attempt + 1}/{self.max_retries})")
                
#                 response = await asyncio.to_thread(
#                     self.model.generate_content,
#                     prompt,
#                     generation_config=self.generation_config
#                 )
                
#                 if not response.text:
#                     raise ValueError("Empty response from LLM")
                    
#                 logger.debug(f"Received response from LLM: {response.text[:100]}...")
#                 return response.text
                
#             except Exception as e:
#                 logger.error(f"Attempt {attempt + 1} failed: {str(e)}")
#                 if attempt == self.max_retries - 1:
#                     raise
#                 await asyncio.sleep(1)

#     async def analyze_schema(self, schema: Dict[str, Any]) -> Dict[str, Any]:
#         """Analyze dataset schema for insights"""
#         try:
#             prompt = f"""Analyze this dataset schema and provide insights:

# Schema:
# {json.dumps(schema, indent=2)}

# Think through:
# 1. Data types and their implications
# 2. Potential relationships between columns
# 3. Key metrics and dimensions
# 4. Possible aggregations
# 5. Common query patterns

# Return a JSON with analysis:
# {{
#     "metrics": ["list of numerical columns suitable for aggregation"],
#     "dimensions": ["list of categorical/grouping columns"],
#     "date_columns": ["list of temporal columns"],
#     "relationships": [
#         {{"column": "col_name", "related_to": ["related_columns"], "type": "relationship_type"}}
#     ],
#     "suggested_queries": [
#         {{"description": "query description", "columns_needed": ["cols"]}}
#     ]
# }}"""

#             response_text = await self._generate_with_retry(prompt)
#             insights = json.loads(response_text)
#             self._schema_insights[schema.get("name", "default")] = insights
#             return insights
            
#         except Exception as e:
#             logger.error(f"Error analyzing schema: {str(e)}", exc_info=True)
#             raise

#     async def analyze_query(self, question: str, schema: Dict[str, Any]) -> Dict[str, Any]:
#         """Analyze natural language query"""
#         try:
#             # Get schema insights if available
#             schema_name = schema.get("name", "default")
#             insights = self._schema_insights.get(schema_name, {})
            
#             prompt = f"""Analyze this question for the given dataset:

# Question: {question}

# Schema:
# {json.dumps(schema, indent=2)}

# Schema Insights:
# {json.dumps(insights, indent=2)}

# Think through these steps:
# 1. Understand the main requirement (what data is being requested)
# 2. Identify required columns and their types
# 3. Determine necessary filters and conditions
# 4. Consider appropriate sorting and limits
# 5. Identify any needed aggregations
# 6. Consider data quality implications

# Return a JSON analysis:
# {{
#     "query_type": "select/aggregate/compare",
#     "required_columns": ["list", "of", "columns"],
#     "conditions": ["list", "of", "conditions"],
#     "sort": {{"column": "col_name", "order": "asc/desc"}},
#     "limit": number_of_results,
#     "aggregations": [
#         {{"function": "sum/avg/count", "column": "col_name"}}
#     ],
#     "grouping": ["group_by_columns"],
#     "explanation": "detailed explanation of analysis"
# }}"""

#             response_text = await self._generate_with_retry(prompt)
#             analysis = json.loads(response_text)
            
#             # Store in context
#             self._context[question] = {
#                 "schema": schema,
#                 "analysis": analysis,
#                 "timestamp": time.time()
#             }
            
#             return analysis
            
#         except Exception as e:
#             logger.error(f"Error in analyze_query: {str(e)}", exc_info=True)
#             raise

#     async def generate_query(self, analysis: Dict[str, Any], schema: Dict[str, Any]) -> str:
#         """Generate SQL query from analysis"""
#         try:
#             prompt = f"""Generate a SQL query based on this analysis:

# Analysis:
# {json.dumps(analysis, indent=2)}

# Schema:
# {json.dumps(schema, indent=2)}

# Follow these guidelines:
# 1. Use exact column names from schema
# 2. Include appropriate table aliases
# 3. Handle NULL values appropriately
# 4. Use proper SQL syntax and formatting
# 5. Consider performance implications
# 6. Add comments for complex parts

# Generate only the SQL query, no other text."""

#             return await self._generate_with_retry(prompt)
            
#         except Exception as e:
#             logger.error(f"Error in generate_query: {str(e)}", exc_info=True)
#             raise

#     async def validate_query(self, query: str, schema: Dict[str, Any]) -> Dict[str, Any]:
#         """Validate generated SQL query"""
#         try:
#             prompt = f"""Validate this SQL query:

# Query:
# {query}

# Schema:
# {json.dumps(schema, indent=2)}

# Perform these checks:
# 1. SQL syntax correctness
# 2. Column name validity
# 3. Data type compatibility
# 4. NULL handling
# 5. Performance considerations
# 6. Potential edge cases

# Return a JSON validation result:
# {{
#     "isValid": true/false,
#     "issues": ["list", "of", "issues"],
#     "suggestedFixes": ["list", "of", "fixes"],
#     "explanation": "detailed explanation",
#     "performance_notes": "performance considerations",
#     "edge_cases": ["potential", "edge", "cases"]
# }}"""

#             response_text = await self._generate_with_retry(prompt)
#             return json.loads(response_text)
            
#         except Exception as e:
#             logger.error(f"Error in validate_query: {str(e)}", exc_info=True)
#             raise

#     async def format_results(self, results: List[Dict], analysis: Dict[str, Any]) -> str:
#         """Format query results into natural language"""
#         try:
#             prompt = f"""Format these query results into a natural language response:

# Results:
# {json.dumps(results, indent=2)}

# Analysis Context:
# {json.dumps(analysis, indent=2)}

# Format the response to:
# 1. Summarize the key findings
# 2. Highlight important patterns
# 3. Provide relevant context
# 4. Include numerical summaries
# 5. Suggest follow-up insights

# Return only the formatted response."""

#             return await self._generate_with_retry(prompt)
            
#         except Exception as e:
#             logger.error(f"Error in format_results: {str(e)}", exc_info=True)
#             raise

#     async def heal_query(self, 
#                             validation_result: Dict[str, Any],
#                             original_query: str,
#                             analysis: Dict[str, Any],
#                             schema: Dict[str, Any]) -> Dict[str, Any]:
#             """Attempt to heal an invalid query"""
#             try:
#                 prompt = f"""Fix this SQL query based on validation results:

#     Original Query: {original_query}

#     Validation Issues:
#     {json.dumps(validation_result, indent=2)}

#     Original Analysis:
#     {json.dumps(analysis, indent=2)}

#     Schema:
#     {json.dumps(schema, indent=2)}

#     Return a JSON with this exact structure:
#     {{
#         "healed_query": "fixed SQL query",
#         "changes_made": [
#             {{
#                 "issue": "description of issue",
#                 "fix": "how it was fixed",
#                 "reasoning": "why this fix works"
#             }}
#         ],
#         "requires_reanalysis": false,
#         "confidence": 0.95,
#         "requires_human_review": false,
#         "notes": "explanation of changes"
#     }}

#     If the query needs complete reanalysis, set requires_reanalysis to true.
#     If the issues are too complex, set requires_human_review to true."""

#                 logger.debug(f"Healing prompt: {prompt}")
#                 response_text = await self._generate_with_retry(prompt)
#                 logger.debug(f"Raw healing response: {response_text}")
                
#                 healing_result = json.loads(response_text)
                
#                 # Add metadata about healing attempt
#                 healing_result["timestamp"] = time.time()
#                 healing_result["original_query"] = original_query
                
#                 # Store healing attempt in context
#                 query_hash = hash(original_query)
#                 if query_hash not in self._context:
#                     self._context[query_hash] = {
#                         "healing_attempts": []
#                     }
#                 self._context[query_hash]["healing_attempts"].append(healing_result)
                
#                 logger.info(f"Query healing completed: {healing_result.get('healed_query', '')[:100]}")
#                 return healing_result
                
#             except json.JSONDecodeError as e:
#                 logger.error(f"JSON parsing error in healing: {str(e)}, Response: {response_text}")
#                 raise HTTPException(status_code=500, detail="Failed to parse healing response")
#             except Exception as e:
#                 logger.error(f"Error in heal_query: {str(e)}", exc_info=True)
#                 raise

#     async def process_with_healing(self,
#                                  question: str,
#                                  schema: Dict[str, Any],
#                                  max_healing_attempts: int = 3) -> Dict[str, Any]:
#         """Process a query with automatic healing attempts"""
#         healing_attempts = 0
#         current_analysis = None
#         current_query = None

#         while healing_attempts < max_healing_attempts:
#             try:
#                 # If we need reanalysis or this is the first attempt
#                 if current_analysis is None:
#                     current_analysis = await self.analyze_query(question, schema)
#                     logger.info(f"Analysis completed (attempt {healing_attempts + 1})")

#                 # Generate SQL
#                 current_query = await self.generate_query(current_analysis, schema)
#                 logger.info(f"Query generated: {current_query}")

#                 # Validate
#                 validation = await self.validate_query(current_query, schema)
#                 logger.info(f"Validation result: {validation.get('isValid')}")

#                 if validation.get("isValid", False):
#                     return {
#                         "success": True,
#                         "query": current_query,
#                         "analysis": current_analysis,
#                         "validation": validation,
#                         "healing_attempts": healing_attempts
#                     }

#                 # If invalid, attempt healing
#                 healing_attempts += 1
#                 logger.info(f"Starting healing attempt {healing_attempts}/{max_healing_attempts}")

#                 healing_result = await self.heal_query(
#                     validation,
#                     current_query,
#                     current_analysis,
#                     schema
#                 )

#                 if healing_result.get("requires_human_review", False):
#                     return {
#                         "success": False,
#                         "error": "Query requires human review",
#                         "validation": validation,
#                         "healing_attempts": healing_attempts,
#                         "notes": healing_result.get("notes", "")
#                     }

#                 if healing_result.get("requires_reanalysis", False):
#                     current_analysis = None  # Force reanalysis
#                     continue

#                 current_query = healing_result.get("healed_query")

#             except Exception as e:
#                 logger.error(f"Error in healing attempt {healing_attempts}: {str(e)}")
#                 healing_attempts += 1
#                 if healing_attempts >= max_healing_attempts:
#                     raise

#         return {
#             "success": False,
#             "error": "Max healing attempts reached",
#             "healing_attempts": healing_attempts
#         }

# # Create singleton instance
# llm_client = LLMClient()

# backend/python/app/llm/client.py
import os
import logging
from dotenv import load_dotenv
import google.generativeai as genai
from typing import Dict, Any, List
import json
import asyncio
import time
from fastapi import HTTPException

# Load environment variables
load_dotenv()

# Configure logging
logging.basicConfig(
    level=logging.DEBUG,
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s',
    handlers=[
        logging.FileHandler('llm_service.log'),
        logging.StreamHandler()
    ]
)
logger = logging.getLogger(__name__)

class DatasetAnalyzer:
    """Helper class for dataset-specific analysis"""
    @staticmethod
    def infer_column_type(sample_value: str) -> str:
        """Infer the type of a column from a sample value"""
        try:
            int(sample_value)
            return "integer"
        except ValueError:
            try:
                float(sample_value)
                return "float"
            except ValueError:
                return "string"

    @staticmethod
    def identify_special_columns(schema: Dict) -> Dict[str, List[str]]:
        """Identify special column types"""
        columns = {
            "numeric": [],
            "categorical": [],
            "temporal": [],
            "textual": [],
            "identifier": []
        }
        
        for col, info in schema.items():
            col_type = info.get("inferred_type", "").lower()
            sample = info.get("sample", "")
            
            if col_type in ["integer", "float"]:
                columns["numeric"].append(col)
            elif "date" in col.lower() or "time" in col.lower():
                columns["temporal"].append(col)
            elif "id" in col.lower():
                columns["identifier"].append(col)
            elif len(sample.split()) > 3:
                columns["textual"].append(col)
            else:
                columns["categorical"].append(col)
                
        return columns

class LLMClient:
    def __init__(self):
        try:
            self._initialize_llm()
            self._initialize_rate_limiting()
            self._initialize_context_store()
            self.dataset_analyzer = DatasetAnalyzer()
            self.request_timeout = int(os.getenv("REQUEST_TIMEOUT", "90"))
            logger.info("LLM Client initialized successfully")
        except Exception as e:
            logger.error(f"Failed to initialize LLM Client: {str(e)}", exc_info=True)
            raise

    def _initialize_llm(self):
        """Initialize LLM configuration"""
        api_key = os.getenv("GEMINI_API_KEY")
        if not api_key:
            raise ValueError("GEMINI_API_KEY not found in environment variables")
            
        try:
            genai.configure(api_key=api_key)
            self.model = genai.GenerativeModel(os.getenv("GEMINI_MODEL_NAME", "gemini-1.5-flash"))
            self.generation_config = {
                'temperature': float(os.getenv("LLM_TEMPERATURE", "0.3")),
                'top_p': float(os.getenv("LLM_TOP_P", "0.8")),
                'top_k': int(os.getenv("LLM_TOP_K", "40")),
                'max_output_tokens': int(os.getenv("LLM_MAX_TOKENS", "2048")),
            }
            logger.info("LLM initialized successfully with model: %s", os.getenv("GEMINI_MODEL_NAME", "gemini-1.5-flash"))
        except Exception as e:
            logger.error("Failed to initialize LLM: %s", str(e), exc_info=True)
            raise

    def _initialize_rate_limiting(self):
        """Initialize rate limiting settings"""
        self._last_call_time = time.time()
        self._rate_limit_delay = float(os.getenv("RATE_LIMIT_DELAY", "1.0"))
        self.max_retries = int(os.getenv("LLM_MAX_RETRIES", "3"))
        logger.info("Rate limiting initialized with delay: %s seconds", self._rate_limit_delay)

    def _initialize_context_store(self):
        """Initialize context storage"""
        self._context = {}
        self._schema_metadata = {}
        logger.info("Context store initialized")

    async def _handle_rate_limit(self):
        """Handle rate limiting between requests"""
        now = time.time()
        time_since_last_call = now - self._last_call_time
        if time_since_last_call < self._rate_limit_delay:
            delay = self._rate_limit_delay - time_since_last_call
            logger.debug("Rate limiting: waiting for %.2f seconds", delay)
            await asyncio.sleep(delay)
        self._last_call_time = time.time()

    async def _generate_content(self, prompt: str) -> str:
        """Generate content with retries and error handling"""
        for attempt in range(self.max_retries):
            try:
                await self._handle_rate_limit()
                logger.debug("Generating content (attempt %d/%d)", attempt + 1, self.max_retries)
                
                # response = await asyncio.to_thread(
                #     self.model.generate_content,
                #     prompt,
                #     generation_config=self.generation_config
                # ),
                response = await asyncio.wait_for(
                    asyncio.to_thread(
                        self.model.generate_content,
                        prompt,
                        generation_config=self.generation_config
                    ),
                    timeout=self.request_timeout
                )

                if not response.text:
                    raise ValueError("Empty response from LLM")
                
                logger.debug("Content generated successfully")
                return response.text
                
            except Exception as e:
                logger.error("Content generation failed (attempt %d/%d): %s", 
                           attempt + 1, self.max_retries, str(e))
                if attempt == self.max_retries - 1:
                    raise
                await asyncio.sleep(1)

    def _clean_response(self, response: str) -> str:
        """Clean LLM response text"""
        try:
            text = response.strip()
            # Remove any markdown formatting
            if '```json' in text:
                text = text.split('```json')[1].split('```')[0]
            elif '```' in text:
                text = text.split('```')[1].split('```')[0]
            return text.strip()
        except Exception as e:
            logger.error("Failed to clean response: %s", str(e))
            return response

    async def analyze_query(self, question: str, schema: Dict[str, Any]) -> Dict[str, Any]:
        """Analyze natural language query"""
        try:
            logger.info("Starting query analysis for: %s", question)
            
            # Analyze schema structure
            column_types = self.dataset_analyzer.identify_special_columns(schema)
            logger.debug("Identified column types: %s", column_types)
            
            prompt = f"""Analyze this question for the given dataset:

Question: {question}

Available Schema:
{json.dumps(schema, indent=2)}

Column Categories:
{json.dumps(column_types, indent=2)}

Return ONLY a JSON object with this structure (no additional text or explanations):
{{
    "query_type": "select",
    "required_columns": ["list", "of", "columns"],
    "conditions": ["list", "of", "conditions"],
    "sort": {{"column": "name", "order": "desc"}},
    "limit": number,
    "explanation": "brief explanation"
}}"""

            logger.debug("Analysis prompt created")
            response_text = await self._generate_content(prompt)
            logger.debug("Raw LLM response: %s", response_text)
            
            # Clean response
            cleaned_response = self._clean_response(response_text)
            logger.debug("Cleaned response: %s", cleaned_response)
            
            try:
                analysis = json.loads(cleaned_response)
                logger.info("Analysis completed successfully")
                return analysis
            except json.JSONDecodeError as e:
                logger.error("Failed to parse analysis response: %s", str(e))
                raise HTTPException(
                    status_code=500,
                    detail=f"Invalid response format: {str(e)}"
                )
                
        except Exception as e:
            logger.error("Query analysis failed: %s", str(e), exc_info=True)
            raise HTTPException(
                status_code=500,
                detail=f"Analysis failed: {str(e)}"
            )

    async def generate_query(self, analysis: Dict[str, Any], schema: Dict[str, Any]) -> str:
        """Generate SQL query from analysis"""
        try:
            logger.info("Starting query generation")
            
            prompt = f"""Generate a SQL query based on this analysis:

Analysis:
{json.dumps(analysis, indent=2)}

Schema:
{json.dumps(schema, indent=2)}

Return ONLY the SQL query, no explanations or additional text."""

            response_text = await self._generate_content(prompt)
            cleaned_query = self._clean_response(response_text)
            logger.info("Query generated successfully")
            return cleaned_query
            
        except Exception as e:
            logger.error("Query generation failed: %s", str(e), exc_info=True)
            raise HTTPException(
                status_code=500,
                detail=f"Query generation failed: {str(e)}"
            )

    async def validate_query(self, query: str, schema: Dict[str, Any]) -> Dict[str, Any]:
        """Validate generated SQL query"""
        try:
            logger.info("Starting query validation")
            
            prompt = f"""Validate this SQL query:

Query:
{query}

Schema:
{json.dumps(schema, indent=2)}

Return ONLY a JSON object with this structure:
{{
    "isValid": true/false,
    "issues": ["list", "of", "issues"],
    "suggestedFixes": ["list", "of", "fixes"],
    "explanation": "validation explanation"
}}"""

            response_text = await self._generate_content(prompt)
            cleaned_response = self._clean_response(response_text)
            
            try:
                validation = json.loads(cleaned_response)
                logger.info("Validation completed: isValid=%s", validation.get("isValid", False))
                return validation
            except json.JSONDecodeError as e:
                logger.error("Failed to parse validation response: %s", str(e))
                raise HTTPException(
                    status_code=500,
                    detail=f"Invalid validation response: {str(e)}"
                )
                
        except Exception as e:
            logger.error("Query validation failed: %s", str(e), exc_info=True)
            raise HTTPException(
                status_code=500,
                detail=f"Validation failed: {str(e)}"
            )

    async def heal_query(self, 
                        validation_result: Dict[str, Any],
                        original_query: str,
                        analysis: Dict[str, Any],
                        schema: Dict[str, Any]) -> Dict[str, Any]:
        """Attempt to heal an invalid query"""
        try:
            logger.info("Starting query healing process")
            logger.debug("Original query: %s", original_query)
            logger.debug("Validation issues: %s", validation_result.get("issues", []))
            
            prompt = f"""Fix this SQL query based on the validation results:

    Original Query:
    {original_query}

    Validation Issues:
    {json.dumps(validation_result, indent=2)}

    Original Analysis:
    {json.dumps(analysis, indent=2)}

    Schema:
    {json.dumps(schema, indent=2)}

    Return ONLY a JSON object with this structure:
    {{
        "healed_query": "fixed SQL query",
        "changes_made": [
            {{
                "issue": "description of what was wrong",
                "fix": "description of how it was fixed",
                "reasoning": "explanation of why this fix works"
            }}
        ],
        "requires_reanalysis": false,
        "confidence": 0.0-1.0,
        "requires_human_review": false,
        "notes": "explanation of changes"
    }}"""

            response_text = await self._generate_content(prompt)
            cleaned_response = self._clean_response(response_text)
            
            try:
                healing_result = json.loads(cleaned_response)
                logger.info("Healing completed: confidence=%s", healing_result.get("confidence", 0))
                return healing_result
            except json.JSONDecodeError as e:
                logger.error("Failed to parse healing response: %s", str(e))
                raise HTTPException(
                    status_code=500,
                    detail=f"Invalid healing response: {str(e)}"
                )
                
        except Exception as e:
            logger.error("Query healing failed: %s", str(e), exc_info=True)
            raise HTTPException(
                status_code=500,
                detail=f"Healing failed: {str(e)}"
            )

    async def process_with_healing(self,
                                question: str,
                                schema: Dict[str, Any],
                                max_healing_attempts: int = 3) -> Dict[str, Any]:
        """Process a query with automatic healing attempts"""
        try:
            healing_attempts = 0
            current_analysis = None
            current_query = None

            while healing_attempts < max_healing_attempts:
                try:
                    logger.info("Processing attempt %d/%d", healing_attempts + 1, max_healing_attempts)

                    # If we need reanalysis or this is the first attempt
                    if current_analysis is None:
                        current_analysis = await self.analyze_query(question, schema)
                        logger.info("Analysis completed")

                    # Generate SQL
                    current_query = await self.generate_query(current_analysis, schema)
                    logger.info("Query generated: %s", current_query)

                    # Validate
                    validation = await self.validate_query(current_query, schema)
                    logger.info("Validation result: isValid=%s", validation.get("isValid", False))

                    if validation.get("isValid", False):
                        return {
                            "success": True,
                            "query": current_query,
                            "analysis": current_analysis,
                            "validation": validation,
                            "healing_attempts": healing_attempts
                        }

                    # If invalid, attempt healing
                    healing_attempts += 1
                    logger.info("Starting healing attempt %d/%d", healing_attempts, max_healing_attempts)

                    healing_result = await self.heal_query(
                        validation,
                        current_query,
                        current_analysis,
                        schema
                    )

                    if healing_result.get("requires_human_review", False):
                        logger.warning("Query requires human review")
                        return {
                            "success": False,
                            "error": "Query requires human review",
                            "validation": validation,
                            "healing_attempts": healing_attempts,
                            "notes": healing_result.get("notes", "")
                        }

                    if healing_result.get("requires_reanalysis", False):
                        logger.info("Healing suggests reanalysis")
                        current_analysis = None  # Force reanalysis
                        continue

                    current_query = healing_result.get("healed_query")
                    logger.info("Applied healed query: %s", current_query)

                except Exception as e:
                    logger.error("Error in healing attempt %d: %s", healing_attempts, str(e))
                    healing_attempts += 1
                    if healing_attempts >= max_healing_attempts:
                        raise

            logger.warning("Max healing attempts reached")
            return {
                "success": False,
                "error": "Max healing attempts reached",
                "healing_attempts": healing_attempts
            }

        except Exception as e:
            logger.error("Process with healing failed: %s", str(e), exc_info=True)
            raise HTTPException(
                status_code=500,
                detail=f"Processing failed: {str(e)}"
            )

try:
    llm_client = LLMClient()
    logger.info("LLM client singleton created successfully")
except Exception as e:
    logger.error("Failed to create LLM client singleton: %s", str(e), exc_info=True)
    raise