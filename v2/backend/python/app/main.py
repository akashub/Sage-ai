# # backend/python/app/llm/client.py
# from datetime import time
# import google.generativeai as genai
# from typing import Dict, Any
# import os
# import json
# import asyncio

# class LLMClient:
#     def __init__(self):
#         genai.configure(api_key=os.getenv("GEMINI_API_KEY"))
#         self.model = genai.GenerativeModel(os.getenv("GEMINI_MODEL_NAME", "gemini-1.5-pro"))
#         self._last_call_time = 0
#         self._rate_limit_delay = float(os.getenv("RATE_LIMIT_DELAY", "1.0"))

#     async def _handle_rate_limit(self):
#         """Handle rate limiting"""
#         now = time.time()
#         time_since_last_call = now - self._last_call_time
#         if time_since_last_call < self._rate_limit_delay:
#             await asyncio.sleep(self._rate_limit_delay - time_since_last_call)
#         self._last_call_time = time.time()

#     async def analyze_query(self, question: str, schema: Dict[str, Any]) -> Dict[str, Any]:
#         await self._handle_rate_limit()
        
#         prompt = f"""Analyze this natural language query:
#         Question: {question}
        
#         Available CSV Schema:
#         {json.dumps(schema, indent=2)}
        
#         Return JSON with:
#         1. Required columns
#         2. Operations needed
#         3. Conditions/filters
#         4. Any aggregations
        
#         Format:
#         {{
#             "columns": ["col1", "col2"],
#             "operations": ["sum", "group_by"],
#             "conditions": ["date > '2024-01-01'"],
#             "aggregations": [
#                 {{"type": "sum", "column": "amount"}}
#             ]
#         }}
#         """
        
#         response = await self.model.generate_content(prompt)
#         return json.loads(response.text)

#     async def generate_query(self, analysis: Dict[str, Any], schema: Dict[str, Any]) -> str:
#         await self._handle_rate_limit()
        
#         prompt = f"""Generate a pandas query based on:
        
#         Analysis:
#         {json.dumps(analysis, indent=2)}
        
#         Schema:
#         {json.dumps(schema, indent=2)}
        
#         Return only the pandas code to execute this query.
#         """
        
#         response = await self.model.generate_content(prompt)
#         return response.text

#     async def validate_query(self, query: str, schema: Dict[str, Any]) -> Dict[str, Any]:
#         await self._handle_rate_limit()
        
#         prompt = f"""Validate this pandas query:
        
#         Query:
#         {query}
        
#         Schema:
#         {json.dumps(schema, indent=2)}
        
#         Return JSON:
#         {{
#             "isValid": true/false,
#             "issues": ["issue1", "issue2"]
#         }}
#         """
        
#         response = await self.model.generate_content(prompt)
#         return json.loads(response.text)

# llm_client = LLMClient()

#TODO: Implement Self-healing loop after validation

import asyncio
import logging
from venv import logger
from fastapi import FastAPI, HTTPException, Request
from pydantic import BaseModel
from typing import Dict, Any, Optional
from .llm.client import llm_client
# from fastapi.middleware.timeout import TimeoutMiddleware
from fastapi.middleware.cors import CORSMiddleware

logging.basicConfig(
    level=logging.DEBUG,
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s'
)
logger = logging.getLogger(__name__)


app = FastAPI()

# class AnalyzeRequest(BaseModel):
#     question: str
#     schema: Dict[str, Any]
class AnalyzeRequest(BaseModel):
    session_id: str
    data: Dict[str, Any]

class GenerateRequest(BaseModel):
    analysis: Dict[str, Any]
    schema: Dict[str, Any]
    session_id: str

class ValidateRequest(BaseModel):
    query: str
    schema: Dict[str, Any]

class HealingRequest(BaseModel):
    query: str
    validation_result: Dict[str, Any]
    analysis: Dict[str, Any]
    schema: Dict[str, Any]

class QueryRequest(BaseModel):
    question: str
    schema: Dict[str, Any]

# @app.post("/analyze")
# async def analyze_query(request: AnalyzeRequest):
#     try:
#         analysis = await llm_client.analyze_query(request.question, request.schema)
#         return {"analysis": analysis}
#     except ValueError as e:
#         raise HTTPException(status_code=400, detail=str(e))
#     except Exception as e:
#         logger.error(f"Error processing request: {str(e)}", exc_info=True)
#         raise HTTPException(status_code=500, detail=str(e))
@app.post("/analyze")
async def analyze_query(request: AnalyzeRequest):
    logger.info(f"Received analyze request for session: {request.session_id}")
    logger.debug(f"Request data: {request.data}")

    try:
        question = request.data.get("question")
        schema = request.data.get("schema")
        
        logger.debug(f"Question: {question}")
        logger.debug(f"Schema: {schema}")

        if not question:
            raise HTTPException(status_code=400, detail="Question is required")
        if not schema:
            raise HTTPException(status_code=400, detail="Schema is required")

        analysis = await llm_client.analyze_query(question, schema)
        logger.info(f"Analysis completed: {analysis}")
        
        return {"analysis": analysis}
    except Exception as e:
        logger.error(f"Error in analyze_query: {str(e)}", exc_info=True)
        raise HTTPException(status_code=500, detail=str(e))

# Add CORS middleware
app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

@app.middleware("http")
async def timeout_middleware(request: Request, call_next):
    try:
        return await asyncio.wait_for(call_next(request), timeout=120.0)
    except asyncio.TimeoutError:
        raise HTTPException(status_code=504, detail="Request timeout")

# @app.post("/generate")
# async def generate_query(request: GenerateRequest):
#     try:
#         result = await llm_client.generate_query(request.analysis, request.schema)
#         return result
#     except ValueError as e:
#         raise HTTPException(status_code=400, detail=str(e))
#     except Exception as e:
#         logger.error(f"Error processing request: {str(e)}", exc_info=True)
#         raise HTTPException(status_code=500, detail=str(e))
@app.post("/generate")
async def generate_query(request: GenerateRequest):
    logger.info(f"Received generate request for session: {request.session_id}")
    try:
        query = await llm_client.generate_query(request.analysis, request.schema)
        return {"query": query}
    except Exception as e:
        logger.error(f"Error in generate_query: {str(e)}", exc_info=True)
        raise HTTPException(status_code=500, detail=str(e))

# @app.post("/validate")
# async def validate_query(request: ValidateRequest):
#     try:
#         validation = await llm_client.validate_query(
#             request.query, 
#             request.schema
#         )
        
#         # If validation fails, attempt healing
#         if not validation["isValid"] and validation.get("requiresHealing", False):
#             healing_result = await llm_client.heal_query(
#                 validation,
#                 request.query,
#                 request.analysis,
#                 request.schema
#             )
#             return {
#                 "validation": validation,
#                 "healing": healing_result
#             }
        
#         return {"validation": validation}
#     except Exception as e:
#         raise HTTPException(status_code=500, detail=str(e))

@app.post("/validate")
async def validate_query(request: ValidateRequest):
    try:
        validation = await llm_client.validate_query(request.query, request.schema)
        return validation
    except Exception as e:
        logger.error(f"Error in validate_query: {str(e)}", exc_info=True)
        raise HTTPException(status_code=500, detail=str(e))

# backend/python/app/main.py
@app.post("/heal")
async def heal_query(request: dict):
    logger.info("Received healing request")
    try:
        validation_result = request.get("validation_result")
        original_query = request.get("original_query")
        analysis = request.get("analysis")
        schema = request.get("schema")
        
        if not all([validation_result, original_query, analysis, schema]):
            raise HTTPException(status_code=400, detail="Missing required healing parameters")
            
        healing_result = await llm_client.heal_query(
            validation_result,
            original_query,
            analysis,
            schema
        )
        return healing_result
    except Exception as e:
        logger.error(f"Error in heal_query: {str(e)}", exc_info=True)
        raise HTTPException(status_code=500, detail=str(e))

@app.post("/process")
async def process_query(request: dict):
    logger.info("Received process request")
    try:
        question = request.get("question")
        schema = request.get("schema")
        
        if not question or not schema:
            raise HTTPException(status_code=400, detail="Missing question or schema")
            
        result = await llm_client.process_with_healing(
            question,
            schema
        )
        return result
    except Exception as e:
        logger.error(f"Error in process_query: {str(e)}", exc_info=True)
        raise HTTPException(status_code=500, detail=str(e))