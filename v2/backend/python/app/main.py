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

# class GenerateRequest(BaseModel):
#     analysis: Dict[str, Any]
#     schema: Dict[str, Any]
#     session_id: str
class GenerateRequest(BaseModel):
    session_id: str
    data: Dict[str, Any]

class ValidateRequest(BaseModel):
    session_id: str
    data: Dict[str, Any]

class HealingRequest(BaseModel):
    query: str
    validation_result: Dict[str, Any]
    analysis: Dict[str, Any]
    schema: Dict[str, Any]

class QueryRequest(BaseModel):
    question: str
    schema: Dict[str, Any]

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


@app.post("/generate")
async def generate_query(request: GenerateRequest):
    logger.info(f"Received generate request for session: {request.session_id}")
    try:
        analysis = request.data["analysis"]
        schema = request.data["schema"]
        query = await llm_client.generate_query(analysis, schema)
        # Return the query string directly in the response
        return {"query": query} # Making sure the query is a string and not an object
    except Exception as e:
        logger.error(f"Error in generate_query: {str(e)}", exc_info=True)
        raise HTTPException(status_code=500, detail=str(e))


@app.post("/validate")
async def validate_query(request: ValidateRequest):
    try:
        query = request.data["query"]
        schema = request.data["schema"]
        validation = await llm_client.validate_query(query, schema)
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