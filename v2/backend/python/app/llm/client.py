# # backend/python/app/llm/client.py
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
#     level=logging.DEBUG,
#     format='%(asctime)s - %(name)s - %(levelname)s - %(message)s',
#     handlers=[
#         logging.FileHandler('llm_service.log'),
#         logging.StreamHandler()
#     ]
# )
# logger = logging.getLogger(__name__)

# class DatasetAnalyzer:
#     """Helper class for dataset-specific analysis"""
#     @staticmethod
#     def infer_column_type(sample_value: str) -> str:
#         """Infer the type of a column from a sample value"""
#         try:
#             int(sample_value)
#             return "integer"
#         except ValueError:
#             try:
#                 float(sample_value)
#                 return "float"
#             except ValueError:
#                 return "string"

#     @staticmethod
#     def identify_special_columns(schema: Dict) -> Dict[str, List[str]]:
#         """Identify special column types"""
#         columns = {
#             "numeric": [],
#             "categorical": [],
#             "temporal": [],
#             "textual": [],
#             "identifier": []
#         }
        
#         for col, info in schema.items():
#             col_type = info.get("inferred_type", "").lower()
#             sample = info.get("sample", "")
            
#             if col_type in ["integer", "float"]:
#                 columns["numeric"].append(col)
#             elif "date" in col.lower() or "time" in col.lower():
#                 columns["temporal"].append(col)
#             elif "id" in col.lower():
#                 columns["identifier"].append(col)
#             elif len(sample.split()) > 3:
#                 columns["textual"].append(col)
#             else:
#                 columns["categorical"].append(col)
                
#         return columns

# class LLMClient:
#     def __init__(self):
#         try:
#             self._initialize_llm()
#             self._initialize_rate_limiting()
#             self._initialize_context_store()
#             self.dataset_analyzer = DatasetAnalyzer()
#             self.request_timeout = int(os.getenv("REQUEST_TIMEOUT", "90"))
#             logger.info("LLM Client initialized successfully")
#         except Exception as e:
#             logger.error(f"Failed to initialize LLM Client: {str(e)}", exc_info=True)
#             raise

#     def _initialize_llm(self):
#         """Initialize LLM configuration"""
#         api_key = os.getenv("GEMINI_API_KEY")
#         if not api_key:
#             raise ValueError("GEMINI_API_KEY not found in environment variables")
            
#         try:
#             genai.configure(api_key=api_key)
#             self.model = genai.GenerativeModel(os.getenv("GEMINI_MODEL_NAME", "gemini-1.5-flash"))
#             self.generation_config = {
#                 'temperature': float(os.getenv("LLM_TEMPERATURE", "0.3")),
#                 'top_p': float(os.getenv("LLM_TOP_P", "0.8")),
#                 'top_k': int(os.getenv("LLM_TOP_K", "40")),
#                 'max_output_tokens': int(os.getenv("LLM_MAX_TOKENS", "2048")),
#             }
#             logger.info("LLM initialized successfully with model: %s", os.getenv("GEMINI_MODEL_NAME", "gemini-1.5-flash"))
#         except Exception as e:
#             logger.error("Failed to initialize LLM: %s", str(e), exc_info=True)
#             raise

#     def _initialize_rate_limiting(self):
#         """Initialize rate limiting settings"""
#         self._last_call_time = time.time()
#         self._rate_limit_delay = float(os.getenv("RATE_LIMIT_DELAY", "1.0"))
#         self.max_retries = int(os.getenv("LLM_MAX_RETRIES", "3"))
#         logger.info("Rate limiting initialized with delay: %s seconds", self._rate_limit_delay)

#     def _initialize_context_store(self):
#         """Initialize context storage"""
#         self._context = {}
#         self._schema_metadata = {}
#         logger.info("Context store initialized")

#     async def _handle_rate_limit(self):
#         """Handle rate limiting between requests"""
#         now = time.time()
#         time_since_last_call = now - self._last_call_time
#         if time_since_last_call < self._rate_limit_delay:
#             delay = self._rate_limit_delay - time_since_last_call
#             logger.debug("Rate limiting: waiting for %.2f seconds", delay)
#             await asyncio.sleep(delay)
#         self._last_call_time = time.time()

#     async def _generate_content(self, prompt: str) -> str:
#         """Generate content with retries and error handling"""
#         for attempt in range(self.max_retries):
#             try:
#                 await self._handle_rate_limit()
#                 logger.debug("Generating content (attempt %d/%d)", attempt + 1, self.max_retries)
                
#                 response = await asyncio.wait_for(
#                     asyncio.to_thread(
#                         self.model.generate_content,
#                         prompt,
#                         generation_config=self.generation_config
#                     ),
#                     timeout=self.request_timeout
#                 )

#                 if not response.text:
#                     raise ValueError("Empty response from LLM")
                
#                 logger.debug("Content generated successfully")
#                 return response.text
                
#             except Exception as e:
#                 logger.error("Content generation failed (attempt %d/%d): %s", 
#                            attempt + 1, self.max_retries, str(e))
#                 if attempt == self.max_retries - 1:
#                     raise
#                 await asyncio.sleep(1)

#     def _clean_response(self, response: str) -> str:
#         """Clean LLM response text"""
#         try:
#             text = response.strip()
#             # Remove any markdown formatting
#             if '```json' in text:
#                 text = text.split('```json')[1].split('```')[0]
#             elif '```' in text:
#                 text = text.split('```')[1].split('```')[0]
#             return text.strip()
#         except Exception as e:
#             logger.error("Failed to clean response: %s", str(e))
#             return response

#     async def analyze_query(self, question: str, schema: Dict[str, Any]) -> Dict[str, Any]:
#         """Analyze natural language query"""
#         try:
#             logger.info("Starting query analysis for: %s", question)
            
#             # Analyze schema structure
#             column_types = self.dataset_analyzer.identify_special_columns(schema)
#             logger.debug("Identified column types: %s", column_types)
            
#             prompt = f"""Analyze this question for the given dataset:

# Question: {question}

# Available Schema:
# {json.dumps(schema, indent=2)}

# Column Categories:
# {json.dumps(column_types, indent=2)}

# Return ONLY a JSON object with this structure (no additional text or explanations):
# {{
#     "query_type": "select",
#     "required_columns": ["list", "of", "columns"],
#     "conditions": ["list", "of", "conditions"],
#     "sort": {{"column": "name", "order": "desc"}},
#     "limit": number,
#     "explanation": "brief explanation"
# }}"""

#             logger.debug("Analysis prompt created")
#             response_text = await self._generate_content(prompt)
#             logger.debug("Raw LLM response: %s", response_text)
            
#             # Clean response
#             cleaned_response = self._clean_response(response_text)
#             logger.debug("Cleaned response: %s", cleaned_response)
            
#             try:
#                 analysis = json.loads(cleaned_response)
#                 logger.info("Analysis completed successfully")
#                 return analysis
#             except json.JSONDecodeError as e:
#                 logger.error("Failed to parse analysis response: %s", str(e))
#                 raise HTTPException(
#                     status_code=500,
#                     detail=f"Invalid response format: {str(e)}"
#                 )
                
#         except Exception as e:
#             logger.error("Query analysis failed: %s", str(e), exc_info=True)
#             raise HTTPException(
#                 status_code=500,
#                 detail=f"Analysis failed: {str(e)}"
#             )
    
#     async def analyze_with_knowledge(self, question: str, schema: Dict[str, Any], knowledge_context: Dict[str, Any]) -> Dict[str, Any]:
#         """Analyze query with knowledge context"""
#         try:
#             logger.info("Starting query analysis with knowledge context")
            
#             # Format knowledge context elements
#             ddl_schemas = knowledge_context.get("ddl_schemas", [])
#             documentation = knowledge_context.get("documentation", [])
#             examples = knowledge_context.get("examples", [])
            
#             # Analyze schema structure
#             column_types = self.dataset_analyzer.identify_special_columns(schema)
            
#             # Construct a prompt that includes knowledge context
#             prompt = f"""Analyze this question for the given dataset with additional context:

# Question: {question}

# Available Schema:
# {json.dumps(schema, indent=2)}

# Column Categories:
# {json.dumps(column_types, indent=2)}

# """
            
#             # Add DDL schemas if available
#             if ddl_schemas:
#                 prompt += f"""
# Database DDL Schemas:
# {json.dumps(ddl_schemas, indent=2)}
# """
            
#             # Add documentation if available
#             if documentation:
#                 prompt += f"""
# Business Documentation:
# {json.dumps(documentation, indent=2)}
# """
            
#             # Add example question-SQL pairs if available
#             if examples:
#                 prompt += f"""
# Similar Question-SQL Examples:
# {json.dumps(examples, indent=2)}
# """
            
#             # Complete the prompt with output instructions
#             prompt += """
# Return ONLY a JSON object with this structure (no additional text or explanations):
# {
#     "query_type": "select",
#     "required_columns": ["list", "of", "columns"],
#     "conditions": ["list", "of", "conditions"],
#     "sort": {"column": "name", "order": "desc"},
#     "limit": number,
#     "explanation": "brief explanation"
# }

# Use the knowledge context to better understand the database structure, terminology, and query patterns.
# """

#             logger.debug("Analysis with knowledge prompt created")
#             response_text = await self._generate_content(prompt)
#             logger.debug("Raw LLM response with knowledge: %s", response_text)
            
#             # Clean response
#             cleaned_response = self._clean_response(response_text)
#             logger.debug("Cleaned response with knowledge: %s", cleaned_response)
            
#             try:
#                 analysis = json.loads(cleaned_response)
#                 logger.info("Analysis with knowledge completed successfully")
#                 return analysis
#             except json.JSONDecodeError as e:
#                 logger.error("Failed to parse analysis with knowledge response: %s", str(e))
#                 raise HTTPException(
#                     status_code=500,
#                     detail=f"Invalid response format: {str(e)}"
#                 )
                
#         except Exception as e:
#             logger.error("Query analysis with knowledge failed: %s", str(e), exc_info=True)
#             raise HTTPException(
#                 status_code=500,
#                 detail=f"Analysis with knowledge failed: {str(e)}"
#             )
#     async def generate_query(self, analysis: Dict[str, Any], schema: Dict[str, Any]) -> str:
#         """Generate SQL query based on analysis"""
#         try:
#             logger.info("Starting query generation")
            
#             prompt = f"""Generate a SQL query based on this analysis:

#     Analysis:
#     {json.dumps(analysis, indent=2)}

#     Schema:
#     {json.dumps(schema, indent=2)}

#     Return ONLY the raw SQL query, no explanations, no additional text, and most importantly NO prefixes.
#     Your response must begin directly with "SELECT" with no prefix.
#     """

#             response_text = await self._generate_content(prompt)
#             cleaned_query = self._clean_response(response_text)
            
#             # Additional cleaning to ensure no "sql" prefix
#             cleaned_query = cleaned_query.strip()
#             for prefix in ["sql", "SQL"]:
#                 if cleaned_query.startswith(prefix):
#                     cleaned_query = cleaned_query[len(prefix):].strip()
            
#             # Ensure it's a SELECT statement
#             if not cleaned_query.upper().startswith("SELECT"):
#                 logger.warning(f"Generated query does not start with SELECT: {cleaned_query}")
#                 # Force it to be a SELECT statement
#                 cleaned_query = f"SELECT * FROM table LIMIT 5"
#                 logger.info(f"Replaced with generic SELECT query: {cleaned_query}")
            
#             logger.info("Query generated successfully")
#             return cleaned_query
                
#         except Exception as e:
#             logger.error("Query generation failed: %s", str(e), exc_info=True)
#             raise HTTPException(
#                 status_code=500,
#                 detail=f"Query generation failed: {str(e)}"
#             )
#     # In app/llm/client.py - modify the generate_query_with_knowledge function

#     async def generate_query_with_knowledge(self, analysis: Dict[str, Any], schema: Dict[str, Any], 
#                                     knowledge_context: Dict[str, Any]) -> str:
#         """Generate SQL query with knowledge context"""
#         try:
#             logger.info("Starting query generation with knowledge context")
            
#             # Format knowledge context elements
#             ddl_schemas = knowledge_context.get("ddl_schemas", [])
#             examples = knowledge_context.get("examples", [])
            
#             prompt = f"""Generate a SQL query based on this analysis:

#     Analysis:
#     {json.dumps(analysis, indent=2)}

#     Schema:
#     {json.dumps(schema, indent=2)}

#     """
            
#             # Add DDL schemas if available
#             if ddl_schemas:
#                 prompt += f"""
#     Database DDL Schemas:
#     {json.dumps(ddl_schemas, indent=2)}
#     """
            
#             # Add example question-SQL pairs if available
#             if examples:
#                 prompt += f"""
#     Similar SQL Examples:
#     {json.dumps(examples, indent=2)}
#     """
            
#             # Complete the prompt with explicit instructions to not use "sql" prefix
#             prompt += """
#     Return ONLY the raw SQL query, no explanations, no additional text, and most importantly NO prefixes.
#     Your response must begin directly with "SELECT" with no prefix.
#     """

#             response_text = await self._generate_content(prompt)
#             cleaned_query = self._clean_response(response_text)
            
#             # Additional cleaning to ensure no "sql" prefix
#             cleaned_query = cleaned_query.strip()
#             for prefix in ["sql", "SQL"]:
#                 if cleaned_query.startswith(prefix):
#                     cleaned_query = cleaned_query[len(prefix):].strip()
            
#             # Ensure it's a SELECT statement
#             if not cleaned_query.upper().startswith("SELECT"):
#                 logger.warning(f"Generated query does not start with SELECT: {cleaned_query}")
#                 # Force it to be a SELECT statement
#                 cleaned_query = f"SELECT * FROM table LIMIT 5"
#                 logger.info(f"Replaced with generic SELECT query: {cleaned_query}")
            
#             logger.info("Query generated successfully with knowledge context")
#             return cleaned_query
                
#         except Exception as e:
#             logger.error("Query generation with knowledge failed: %s", str(e), exc_info=True)
#             raise HTTPException(
#                 status_code=500,
#                 detail=f"Query generation failed: {str(e)}"
#             )
#     async def validate_query(self, query: str, schema: Dict[str, Any]) -> Dict[str, Any]:
#         """Validate generated SQL query"""
#         try:
#             logger.info("Starting query validation")
            
#             prompt = f"""Validate this SQL query:

# Query:
# {query}

# Schema:
# {json.dumps(schema, indent=2)}

# Return ONLY a JSON object with this structure:
# {{
#     "isValid": true/false,
#     "issues": ["list", "of", "issues"],
#     "suggestedFixes": ["list", "of", "fixes"],
#     "explanation": "validation explanation"
# }}"""

#             response_text = await self._generate_content(prompt)
#             cleaned_response = self._clean_response(response_text)
            
#             try:
#                 validation = json.loads(cleaned_response)
#                 logger.info("Validation completed: isValid=%s", validation.get("isValid", False))
#                 return validation
#             except json.JSONDecodeError as e:
#                 logger.error("Failed to parse validation response: %s", str(e))
#                 raise HTTPException(
#                     status_code=500,
#                     detail=f"Invalid validation response: {str(e)}"
#                 )
                
#         except Exception as e:
#             logger.error("Query validation failed: %s", str(e), exc_info=True)
#             raise HTTPException(
#                 status_code=500,
#                 detail=f"Validation failed: {str(e)}"
#             )

#     async def heal_query(self, 
#                         validation_result: Dict[str, Any],
#                         original_query: str,
#                         analysis: Dict[str, Any],
#                         schema: Dict[str, Any]) -> Dict[str, Any]:
#         """Attempt to heal an invalid query"""
#         try:
#             logger.info("Starting query healing process")
#             logger.debug("Original query: %s", original_query)
#             logger.debug("Validation issues: %s", validation_result.get("issues", []))
            
#             prompt = f"""Fix this SQL query based on the validation results:

#     Original Query:
#     {original_query}

#     Validation Issues:
#     {json.dumps(validation_result, indent=2)}

#     Original Analysis:
#     {json.dumps(analysis, indent=2)}

#     Schema:
#     {json.dumps(schema, indent=2)}

#     Return ONLY a JSON object with this structure:
#     {{
#         "healed_query": "fixed SQL query",
#         "changes_made": [
#             {{
#                 "issue": "description of what was wrong",
#                 "fix": "description of how it was fixed",
#                 "reasoning": "explanation of why this fix works"
#             }}
#         ],
#         "requires_reanalysis": false,
#         "confidence": 0.0-1.0,
#         "requires_human_review": false,
#         "notes": "explanation of changes"
#     }}"""

#             response_text = await self._generate_content(prompt)
#             cleaned_response = self._clean_response(response_text)
            
#             try:
#                 healing_result = json.loads(cleaned_response)
#                 logger.info("Healing completed: confidence=%s", healing_result.get("confidence", 0))
#                 return healing_result
#             except json.JSONDecodeError as e:
#                 logger.error("Failed to parse healing response: %s", str(e))
#                 raise HTTPException(
#                     status_code=500,
#                     detail=f"Invalid healing response: {str(e)}"
#                 )
                
#         except Exception as e:
#             logger.error("Query healing failed: %s", str(e), exc_info=True)
#             raise HTTPException(
#                 status_code=500,
#                 detail=f"Healing failed: {str(e)}"
#             )

#     async def process_with_healing(self,
#                                 question: str,
#                                 schema: Dict[str, Any],
#                                 max_healing_attempts: int = 3) -> Dict[str, Any]:
#         """Process a query with automatic healing attempts"""
#         try:
#             healing_attempts = 0
#             current_analysis = None
#             current_query = None

#             while healing_attempts < max_healing_attempts:
#                 try:
#                     logger.info("Processing attempt %d/%d", healing_attempts + 1, max_healing_attempts)

#                     # If we need reanalysis or this is the first attempt
#                     if current_analysis is None:
#                         current_analysis = await self.analyze_query(question, schema)
#                         logger.info("Analysis completed")

#                     # Generate SQL
#                     current_query = await self.generate_query(current_analysis, schema)
#                     logger.info("Query generated: %s", current_query)

#                     # Validate
#                     validation = await self.validate_query(current_query, schema)
#                     logger.info("Validation result: isValid=%s", validation.get("isValid", False))

#                     if validation.get("isValid", False):
#                         return {
#                             "success": True,
#                             "query": current_query,
#                             "analysis": current_analysis,
#                             "validation": validation,
#                             "healing_attempts": healing_attempts
#                         }

#                     # If invalid, attempt healing
#                     healing_attempts += 1
#                     logger.info("Starting healing attempt %d/%d", healing_attempts, max_healing_attempts)

#                     healing_result = await self.heal_query(
#                         validation,
#                         current_query,
#                         current_analysis,
#                         schema
#                     )

#                     if healing_result.get("requires_human_review", False):
#                         logger.warning("Query requires human review")
#                         return {
#                             "success": False,
#                             "error": "Query requires human review",
#                             "validation": validation,
#                             "healing_attempts": healing_attempts,
#                             "notes": healing_result.get("notes", "")
#                         }

#                     if healing_result.get("requires_reanalysis", False):
#                         logger.info("Healing suggests reanalysis")
#                         current_analysis = None  # Force reanalysis
#                         continue

#                     current_query = healing_result.get("healed_query")
#                     logger.info("Applied healed query: %s", current_query)

#                 except Exception as e:
#                     logger.error("Error in healing attempt %d: %s", healing_attempts, str(e))
#                     healing_attempts += 1
#                     if healing_attempts >= max_healing_attempts:
#                         raise

#             logger.warning("Max healing attempts reached")
#             return {
#                 "success": False,
#                 "error": "Max healing attempts reached",
#                 "healing_attempts": healing_attempts
#             }

#         except Exception as e:
#             logger.error("Process with healing failed: %s", str(e), exc_info=True)
#             raise HTTPException(
#                 status_code=500,
#                 detail=f"Processing failed: {str(e)}"
#             )

# try:
#     llm_client = LLMClient()
#     logger.info("LLM client singleton created successfully")
# except Exception as e:
#     logger.error("Failed to create LLM client singleton: %s", str(e), exc_info=True)
#     raise

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
import openai
import anthropic
from mistralai.client import MistralClient
from mistralai.models.chat_completion import ChatMessage

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
    def __init__(self, api_key=None, provider="gemini", model_name = None):
        self.provider = provider
        self.api_key = api_key or self._get_default_api_key(provider)
        self.model_name = model_name
        
        try:
            if self.provider == "gemini":
                self._initialize_gemini()
            elif self.provider == "openai":
                self._initialize_openai()
            elif self.provider == "anthropic":
                self._initialize_anthropic()
            elif self.provider == "mistral":
                self._initialize_mistral()
            else:
                raise ValueError(f"Unsupported provider: {self.provider}")
                
            self._initialize_rate_limiting()
            self._initialize_context_store()
            self.dataset_analyzer = DatasetAnalyzer()
            self.request_timeout = int(os.getenv("REQUEST_TIMEOUT", "90"))
            logger.info("LLM Client initialized successfully for provider: %s", self.provider)
        except Exception as e:
            logger.error(f"Failed to initialize LLM Client: {str(e)}", exc_info=True)
            raise

    def _get_default_api_key(self, provider):
        """Get default API key from environment variables"""
        if provider == "gemini":
            return os.getenv("GEMINI_API_KEY")
        elif provider == "openai":
            return os.getenv("OPENAI_API_KEY")
        elif provider == "anthropic":
            return os.getenv("ANTHROPIC_API_KEY")
        elif provider == "mistral":
            return os.getenv("MISTRAL_API_KEY")
        return None

    def _initialize_gemini(self):
        """Initialize Gemini LLM configuration"""
        if not self.api_key:
            raise ValueError("Gemini API key not provided")
            
        try:
            genai.configure(api_key=self.api_key)
            model_name = self.model_name or os.getenv("GEMINI_MODEL_NAME", "gemini-1.5-flash")
            self.model = genai.GenerativeModel(model_name)
            self.generation_config = {
                'temperature': float(os.getenv("LLM_TEMPERATURE", "0.3")),
                'top_p': float(os.getenv("LLM_TOP_P", "0.8")),
                'top_k': int(os.getenv("LLM_TOP_K", "40")),
                'max_output_tokens': int(os.getenv("LLM_MAX_TOKENS", "2048")),
            }
            logger.info(f"Gemini initialized successfully with model: {model_name}")
        except Exception as e:
            logger.error("Failed to initialize Gemini: %s", str(e), exc_info=True)
            raise

    def _initialize_openai(self):
        """Initialize OpenAI LLM configuration"""
        if not self.api_key:
            raise ValueError("OpenAI API key not provided")
            
        try:
            openai.api_key = self.api_key
            self.model_name = self.model_name or os.getenv("OPENAI_MODEL_NAME", "gpt-4")  # Use self.model_name
            self.generation_config = {
                'temperature': float(os.getenv("LLM_TEMPERATURE", "0.3")),
                'max_tokens': int(os.getenv("LLM_MAX_TOKENS", "2048")),
                'top_p': float(os.getenv("LLM_TOP_P", "0.8")),
            }
            logger.info(f"OpenAI initialized successfully with model: {self.model_name}")
        except Exception as e:
            logger.error("Failed to initialize OpenAI: %s", str(e), exc_info=True)
            raise

    def _initialize_anthropic(self):
        """Initialize Anthropic Claude LLM configuration"""
        if not self.api_key:
            raise ValueError("Anthropic API key not provided")
            
        try:
            self.client = anthropic.Anthropic(api_key=self.api_key)
            self.model_name = self.model_name or os.getenv("ANTHROPIC_MODEL_NAME", "claude-3-opus-20240229")  # Use self.model_name
            self.generation_config = {
                'temperature': float(os.getenv("LLM_TEMPERATURE", "0.3")),
                'max_tokens': int(os.getenv("LLM_MAX_TOKENS", "2048")),
            }
            logger.info(f"Anthropic Claude initialized successfully with model: {self.model_name}")
        except Exception as e:
            logger.error("Failed to initialize Anthropic: %s", str(e), exc_info=True)
            raise

    def _initialize_mistral(self):
        """Initialize Mistral LLM configuration"""
        if not self.api_key:
            raise ValueError("Mistral API key not provided")
            
        try:
            self.client = MistralClient(api_key=self.api_key)
            self.model_name = self.model_name or os.getenv("MISTRAL_MODEL_NAME", "mistral-large-latest")  # Use self.model_name
            self.generation_config = {
                'temperature': float(os.getenv("LLM_TEMPERATURE", "0.3")),
                'max_tokens': int(os.getenv("LLM_MAX_TOKENS", "2048")),
            }
            logger.info(f"Mistral initialized successfully with model: {self.model_name}")
        except Exception as e:
            logger.error("Failed to initialize Mistral: %s", str(e), exc_info=True)
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
        """Generate content using the selected provider"""
        for attempt in range(self.max_retries):
            try:
                await self._handle_rate_limit()
                logger.debug("Generating content (attempt %d/%d) with provider: %s", 
                           attempt + 1, self.max_retries, self.provider)
                
                if self.provider == "gemini":
                    response = await self._generate_gemini(prompt)
                elif self.provider == "openai":
                    response = await self._generate_openai(prompt)
                elif self.provider == "anthropic":
                    response = await self._generate_anthropic(prompt)
                elif self.provider == "mistral":
                    response = await self._generate_mistral(prompt)
                else:
                    raise ValueError(f"Unsupported provider: {self.provider}")

                if not response:
                    raise ValueError("Empty response from LLM")
                
                logger.debug("Content generated successfully")
                return response
                
            except Exception as e:
                logger.error("Content generation failed (attempt %d/%d): %s", 
                           attempt + 1, self.max_retries, str(e))
                if attempt == self.max_retries - 1:
                    raise
                await asyncio.sleep(1)

    async def _generate_gemini(self, prompt: str) -> str:
        """Generate content using Gemini"""
        response = await asyncio.wait_for(
            asyncio.to_thread(
                self.model.generate_content,
                prompt,
                generation_config=self.generation_config
            ),
            timeout=self.request_timeout
        )
        return response.text

    async def _generate_openai(self, prompt: str) -> str:
        """Generate content using OpenAI"""
        response = await asyncio.wait_for(
            asyncio.to_thread(
                openai.ChatCompletion.create,
                model=self.model_name,
                messages=[{"role": "user", "content": prompt}],
                **self.generation_config
            ),
            timeout=self.request_timeout
        )
        return response.choices[0].message.content

    async def _generate_anthropic(self, prompt: str) -> str:
        """Generate content using Anthropic Claude"""
        response = await asyncio.wait_for(
            asyncio.to_thread(
                self.client.messages.create,
                model=self.model_name,
                messages=[{"role": "user", "content": prompt}],
                **self.generation_config
            ),
            timeout=self.request_timeout
        )
        return response.content[0].text

    async def _generate_mistral(self, prompt: str) -> str:
        """Generate content using Mistral"""
        messages = [ChatMessage(role="user", content=prompt)]
        response = await asyncio.wait_for(
            asyncio.to_thread(
                self.client.chat,
                model=self.model_name,
                messages=messages,
                **self.generation_config
            ),
            timeout=self.request_timeout
        )
        return response.choices[0].message.content

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
    
    async def analyze_with_knowledge(self, question: str, schema: Dict[str, Any], knowledge_context: Dict[str, Any]) -> Dict[str, Any]:
        """Analyze query with knowledge context"""
        try:
            logger.info("Starting query analysis with knowledge context")
            
            # Format knowledge context elements
            ddl_schemas = knowledge_context.get("ddl_schemas", [])
            documentation = knowledge_context.get("documentation", [])
            examples = knowledge_context.get("examples", [])
            
            # Analyze schema structure
            column_types = self.dataset_analyzer.identify_special_columns(schema)
            
            # Construct a prompt that includes knowledge context
            prompt = f"""Analyze this question for the given dataset with additional context:

Question: {question}

Available Schema:
{json.dumps(schema, indent=2)}

Column Categories:
{json.dumps(column_types, indent=2)}

"""
            
            # Add DDL schemas if available
            if ddl_schemas:
                prompt += f"""
Database DDL Schemas:
{json.dumps(ddl_schemas, indent=2)}
"""
            
            # Add documentation if available
            if documentation:
                prompt += f"""
Business Documentation:
{json.dumps(documentation, indent=2)}
"""
            
            # Add example question-SQL pairs if available
            if examples:
                prompt += f"""
Similar Question-SQL Examples:
{json.dumps(examples, indent=2)}
"""
            
            # Complete the prompt with output instructions
            prompt += """
Return ONLY a JSON object with this structure (no additional text or explanations):
{
    "query_type": "select",
    "required_columns": ["list", "of", "columns"],
    "conditions": ["list", "of", "conditions"],
    "sort": {"column": "name", "order": "desc"},
    "limit": number,
    "explanation": "brief explanation"
}

Use the knowledge context to better understand the database structure, terminology, and query patterns.
"""

            logger.debug("Analysis with knowledge prompt created")
            response_text = await self._generate_content(prompt)
            logger.debug("Raw LLM response with knowledge: %s", response_text)
            
            # Clean response
            cleaned_response = self._clean_response(response_text)
            logger.debug("Cleaned response with knowledge: %s", cleaned_response)
            
            try:
                analysis = json.loads(cleaned_response)
                logger.info("Analysis with knowledge completed successfully")
                return analysis
            except json.JSONDecodeError as e:
                logger.error("Failed to parse analysis with knowledge response: %s", str(e))
                raise HTTPException(
                    status_code=500,
                    detail=f"Invalid response format: {str(e)}"
                )
                
        except Exception as e:
            logger.error("Query analysis with knowledge failed: %s", str(e), exc_info=True)
            raise HTTPException(
                status_code=500,
                detail=f"Analysis with knowledge failed: {str(e)}"
            )

    async def generate_query(self, analysis: Dict[str, Any], schema: Dict[str, Any]) -> str:
        """Generate SQL query based on analysis"""
        try:
            logger.info("Starting query generation")
            
            prompt = f"""Generate a SQL query based on this analysis:

    Analysis:
    {json.dumps(analysis, indent=2)}

    Schema:
    {json.dumps(schema, indent=2)}

    Return ONLY the raw SQL query, no explanations, no additional text, and most importantly NO prefixes.
    Your response must begin directly with "SELECT" with no prefix.
    """

            response_text = await self._generate_content(prompt)
            cleaned_query = self._clean_response(response_text)
            
            # Additional cleaning to ensure no "sql" prefix
            cleaned_query = cleaned_query.strip()
            for prefix in ["sql", "SQL"]:
                if cleaned_query.startswith(prefix):
                    cleaned_query = cleaned_query[len(prefix):].strip()
            
            # Ensure it's a SELECT statement
            if not cleaned_query.upper().startswith("SELECT"):
                logger.warning(f"Generated query does not start with SELECT: {cleaned_query}")
                # Force it to be a SELECT statement
                cleaned_query = f"SELECT * FROM table LIMIT 5"
                logger.info(f"Replaced with generic SELECT query: {cleaned_query}")
            
            logger.info("Query generated successfully")
            return cleaned_query
                
        except Exception as e:
            logger.error("Query generation failed: %s", str(e), exc_info=True)
            raise HTTPException(
                status_code=500,
                detail=f"Query generation failed: {str(e)}"
            )

    async def generate_query_with_knowledge(self, analysis: Dict[str, Any], schema: Dict[str, Any], 
                                    knowledge_context: Dict[str, Any]) -> str:
        """Generate SQL query with knowledge context"""
        try:
            logger.info("Starting query generation with knowledge context")
            
            # Format knowledge context elements
            ddl_schemas = knowledge_context.get("ddl_schemas", [])
            examples = knowledge_context.get("examples", [])
            
            prompt = f"""Generate a SQL query based on this analysis:

    Analysis:
    {json.dumps(analysis, indent=2)}

    Schema:
    {json.dumps(schema, indent=2)}

    """
            
            # Add DDL schemas if available
            if ddl_schemas:
                prompt += f"""
    Database DDL Schemas:
    {json.dumps(ddl_schemas, indent=2)}
    """
            
            # Add example question-SQL pairs if available
            if examples:
                prompt += f"""
    Similar SQL Examples:
    {json.dumps(examples, indent=2)}
    """
            
            # Complete the prompt with explicit instructions to not use "sql" prefix
            prompt += """
    Return ONLY the raw SQL query, no explanations, no additional text, and most importantly NO prefixes.
    Your response must begin directly with "SELECT" with no prefix.
    """

            response_text = await self._generate_content(prompt)
            cleaned_query = self._clean_response(response_text)
            
            # Additional cleaning to ensure no "sql" prefix
            cleaned_query = cleaned_query.strip()
            for prefix in ["sql", "SQL"]:
                if cleaned_query.startswith(prefix):
                    cleaned_query = cleaned_query[len(prefix):].strip()
            
            # Ensure it's a SELECT statement
            if not cleaned_query.upper().startswith("SELECT"):
                logger.warning(f"Generated query does not start with SELECT: {cleaned_query}")
                # Force it to be a SELECT statement
                cleaned_query = f"SELECT * FROM table LIMIT 5"
                logger.info(f"Replaced with generic SELECT query: {cleaned_query}")
            
            logger.info("Query generated successfully with knowledge context")
            return cleaned_query
                
        except Exception as e:
            logger.error("Query generation with knowledge failed: %s", str(e), exc_info=True)
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
    # Create a fallback client that can be replaced later
    llm_client = None

def get_default_client():
    global llm_client
    if llm_client is None:
        try:
            # Try to create without API key - it can be added later
            llm_client = LLMClient(api_key=None, provider="gemini", model_name=None)
            logger.info("LLM client singleton created successfully")
        except Exception as e:
            logger.warning(f"Failed to create default LLM client: {e}")
            llm_client = None
    return llm_client
