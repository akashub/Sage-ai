# # app/main.py
# from fastapi import FastAPI, HTTPException
# from fastapi.middleware.cors import CORSMiddleware
# from pydantic import BaseModel
# from typing import Optional, Dict, Any
# from langchain_openai import ChatOpenAI

# from app.orchestration.orchestrator import SQLAssistantOrchestrator
# from backend.app.models.schema import QueryRequest, QueryResponse

# app = FastAPI(title="SQL Assistant API")

# # CORS middleware configuration
# app.add_middleware(
#     CORSMiddleware,
#     allow_origins=["*"],  # In production, replace with actual frontend origin
#     allow_credentials=True,
#     allow_methods=["*"],
#     allow_headers=["*"],
# )

# # Initialize orchestrator
# llm = ChatOpenAI(temperature=0)
# orchestrator = SQLAssistantOrchestrator(llm=llm)

# class QueryRequest(BaseModel):
#     query: str
#     schema: Optional[str] = None
#     database_url: Optional[str] = None

# class QueryResponse(BaseModel):
#     sql: str
#     results: Optional[Dict[str, Any]] = None
#     error: Optional[str] = None

# @app.post("/api/query", response_model=QueryResponse)
# async def process_query(request: QueryRequest):
#     try:
#         # Step 1: Analyze the query
#         analysis_result = await orchestrator.analyzer.analyze(
#             question=request.query,
#             schema=request.schema
#         )
#         if not analysis_result["success"]:
#             raise HTTPException(status_code=400, detail=analysis_result["error"])

#         # Step 2: Generate SQL
#         generation_result = await orchestrator.generator.generate(
#             analysis=analysis_result["analysis"],
#             question=request.query,
#             schema=request.schema
#         )
#         if not generation_result["success"]:
#             raise HTTPException(status_code=400, detail=generation_result["error"])

#         # Step 3: Validate SQL
#         validation_result = await orchestrator.validator.validate(
#             sql=generation_result["sql"],
#             schema=request.schema
#         )
#         if not validation_result["success"] or not validation_result["is_valid"]:
#             raise HTTPException(
#                 status_code=400, 
#                 detail=f"SQL Validation failed: {validation_result.get('issues', [])}"
#             )

#         # Step 4: Execute SQL if database_url is provided
#         if request.database_url:
#             execution_result = await orchestrator.executor.execute(
#                 sql=generation_result["sql"],
#                 database_url=request.database_url
#             )
#             if not execution_result["success"]:
#                 raise HTTPException(status_code=400, detail=execution_result["error"])
            
#             return QueryResponse(
#                 sql=generation_result["sql"],
#                 results=execution_result["results"]
#             )
        
#         # If no database_url, return just the generated SQL
#         return QueryResponse(
#             sql=generation_result["sql"],
#             results=None
#         )

#     except Exception as e:
#         raise HTTPException(status_code=500, detail=str(e))

# if __name__ == "__main__":
#     import uvicorn
#     uvicorn.run(app, host="0.0.0.0", port=8000)

# # app/main.py
# from fastapi import FastAPI, HTTPException, Depends
# from fastapi.middleware.cors import CORSMiddleware
# from typing import Optional, Dict, Any
# from langchain_openai import ChatOpenAI, OpenAIEmbeddings
# import os
# from dotenv import load_dotenv

# # Use absolute imports from project root
# from backend.app.models.schema import QueryRequest, QueryResponse
# from backend.app.orchestration.orchestrator import SQLAssistantOrchestrator
# from backend.app.knowledge.movie_knowledge_base import MovieKnowledgeBase
# from backend.app.knowledge.few_shot import FewShotManager
# from backend.app.data.imdb_data_loader import IMDBDataLoader

# # Load environment variables
# load_dotenv()

# app = FastAPI(
#     title="SQL Assistant API",
#     description="Natural Language to SQL Query Assistant",
#     version="1.0.0"
# )

# # CORS middleware configuration
# app.add_middleware(
#     CORSMiddleware,
#     allow_origins=["*"],  # In production, replace with actual frontend origin
#     allow_credentials=True,
#     allow_methods=["*"],
#     allow_headers=["*"],
# )

# # Initialize components
# def get_components():
#     """Initialize and return all required components"""
#     try:
#         # Initialize LLM and embeddings
#         llm = ChatOpenAI(temperature=0)
#         embeddings = OpenAIEmbeddings()

#         # Initialize Few-Shot Manager
#         few_shot_manager = FewShotManager(
#             examples_file='imdb_few_shot_examples.json',
#             embedding_model=embeddings
#         )

#         # Initialize Knowledge Base
#         knowledge_base = MovieKnowledgeBase(
#             embedding_model=embeddings,
#             few_shot_manager=few_shot_manager
#         )

#         # Initialize Orchestrator
#         orchestrator = SQLAssistantOrchestrator(
#             knowledge_base=knowledge_base,
#             few_shot_manager=few_shot_manager
#         )

#         return orchestrator
#     except Exception as e:
#         print(f"Error initializing components: {e}")
#         raise

# # Get components as dependency
# async def get_orchestrator():
#     """Dependency to get orchestrator"""
#     return get_components()

# @app.post("/api/query", response_model=QueryResponse)
# async def process_query(
#     request: QueryRequest,
#     orchestrator: SQLAssistantOrchestrator = Depends(get_orchestrator)
# ):
#     """Process natural language query to SQL"""
#     try:
#         # Process query through orchestrator
#         result = await orchestrator.process_query(
#             query=request.query,
#             schema=request.schema
#         )

#         if not result["success"]:
#             raise HTTPException(
#                 status_code=400,
#                 detail=result.get("error", "Processing failed")
#             )

#         return QueryResponse(
#             sql=result["sql"],
#             results=result.get("results"),
#             metadata=result.get("metadata")
#         )

#     except HTTPException:
#         raise
#     except Exception as e:
#         raise HTTPException(
#             status_code=500,
#             detail=f"Internal server error: {str(e)}"
#         )

# @app.get("/api/examples")
# async def get_examples(
#     orchestrator: SQLAssistantOrchestrator = Depends(get_orchestrator)
# ):
#     """Get few-shot examples for reference"""
#     try:
#         examples = await orchestrator.few_shot_manager.get_relevant_examples(
#             query="",  # Empty query to get default examples
#             k=5
#         )
#         return {"examples": examples}
#     except Exception as e:
#         raise HTTPException(
#             status_code=500,
#             detail=f"Error fetching examples: {str(e)}"
#         )

# @app.get("/api/health")
# async def health_check():
#     """Health check endpoint"""
#     return {"status": "healthy"}

# @app.get("/api/patterns")
# async def get_patterns(
#     orchestrator: SQLAssistantOrchestrator = Depends(get_orchestrator)
# ):
#     """Get available query patterns"""
#     try:
#         patterns = orchestrator.knowledge_base.patterns
#         return {"patterns": patterns}
#     except Exception as e:
#         raise HTTPException(
#             status_code=500,
#             detail=f"Error fetching patterns: {str(e)}"
#         )

# # Error handlers
# @app.exception_handler(HTTPException)
# async def http_exception_handler(request, exc):
#     return {
#         "success": False,
#         "error": exc.detail,
#         "status_code": exc.status_code
#     }

# @app.exception_handler(Exception)
# async def general_exception_handler(request, exc):
#     return {
#         "success": False,
#         "error": str(exc),
#         "status_code": 500
#     }

# if __name__ == "__main__":
#     import uvicorn
    
#     # Get port from environment or default to 8000
#     port = int(os.getenv("PORT", 8000))
    
#     # Run application
#     uvicorn.run(
#         "main:app",
#         host="0.0.0.0",
#         port=port,
#         reload=True  # Enable auto-reload during development
#     )

# backend/app/main.py
import logging
from fastapi import FastAPI, HTTPException, Depends
from fastapi.middleware.cors import CORSMiddleware
from fastapi.responses import JSONResponse
from typing import Optional, Dict, Any
from langchain_openai import ChatOpenAI, OpenAIEmbeddings
import os
from dotenv import load_dotenv

from backend.app.models.schema import QueryRequest, QueryResponse
from backend.app.orchestration.orchestrator import SQLAssistantOrchestrator
from backend.app.knowledge.movie_knowledge_base import MovieKnowledgeBase
from backend.app.knowledge.few_shot import FewShotManager
import json
import os
from pathlib import Path
from ..app.llm.client import llm_client

# Set up logging
logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

# Load environment variables
load_dotenv()

print("Starting server initialization...")

app = FastAPI(
    title="SQL Assistant API",
    description="Natural Language to SQL Query Assistant",
    version="1.0.0"
)

# CORS middleware
app.add_middleware(
    CORSMiddleware,
    allow_origins=json.loads(os.getenv("BACKEND_CORS_ORIGINS", '["http://localhost:3000"]')),
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

print("FastAPI app created...")


# def get_components():
#     """Initialize and return all required components"""
#     print("Initializing components...")
#     try:
#         # Initialize LLM and embeddings
#         llm = ChatOpenAI(
#             temperature=0,
#             model_name=os.getenv("MODEL_NAME", "gpt-4")
#         )
#         embeddings = OpenAIEmbeddings()
        
#         # Initialize Knowledge Base
#         knowledge_base = MovieKnowledgeBase(
#             examples_file='imdb_few_shot_examples.json',
#             embedding_model=embeddings
#         )
        
#         # Initialize Few-Shot Manager with same examples file
#         few_shot_manager = FewShotManager(
#             examples_file='imdb_few_shot_examples.json',
#             embedding_model=embeddings
#         )
        
#         # Initialize Orchestrator
#         orchestrator = SQLAssistantOrchestrator(
#             knowledge_base=knowledge_base,
#             few_shot_manager=few_shot_manager
#         )
        
#         return orchestrator
#     except Exception as e:
#         logger.error(f"Error initializing components: {str(e)}")
#         raise

BASE_DIR = Path(__file__).resolve().parent
DATA_DIR = BASE_DIR / "data"
EXAMPLES_FILE = DATA_DIR / "imdb_few_shot_examples.json"

# def get_components():
#     """Initialize and return all required components"""
#     try:
#         # Initialize LLM and embeddings
#         llm = ChatOpenAI(
#             temperature=0,
#             model_name=os.getenv("GEMINI_MODEL_NAME", "gemini-1.5-flash")
#         )
#         embeddings = OpenAIEmbeddings()
        
#         # Initialize Few-Shot Manager
#         few_shot_manager = FewShotManager(
#             # examples_file='imdb_few_shot_examples.json',
#             examples_file = str(EXAMPLES_FILE),
#             embedding_model=embeddings
#         )
        
#         # Initialize Knowledge Base
#         knowledge_base = MovieKnowledgeBase(
#             embedding_model=embeddings
#         )
        
#         # Initialize Orchestrator
#         orchestrator = SQLAssistantOrchestrator(
#             knowledge_base=knowledge_base,
#             few_shot_manager=few_shot_manager
#         )
        
#         return orchestrator
#     except Exception as e:
#         logger.error(f"Error initializing components: {str(e)}")
#         raise
def get_components():
    """Initialize and return all required components"""
    try:
        # Initialize Few-Shot Manager with hash-based embeddings
        base_path = Path(__file__).parent
        examples_file = base_path / "data" / "imdb_few_shot_examples.json"
        few_shot_manager = FewShotManager(
            examples_file=str(examples_file)
        )
        
        # Initialize Knowledge Base
        knowledge_base = MovieKnowledgeBase()
        
        # Initialize Orchestrator
        orchestrator = SQLAssistantOrchestrator(
            knowledge_base=knowledge_base,
            few_shot_manager=few_shot_manager
        )
        
        return orchestrator
    except Exception as e:
        logger.error(f"Error initializing components: {e}")
        raise

@app.post("/api/query", response_model=QueryResponse)
async def process_query(request: QueryRequest):
    """Process natural language query to SQL"""
    try:
        logger.info(f"Received query: {request.query}")
        
        orchestrator = get_components()
        result = await orchestrator.process_query(
            question=request.query,
            db_schema=request.db_schema
        )
        
        logger.info(f"Query processed successfully: {result}")
        return QueryResponse(**result)
        
    except Exception as e:
        logger.error(f"Error processing query: {str(e)}")
        return JSONResponse(
            status_code=500,
            content={
                "success": False,
                "error": f"Error processing query: {str(e)}"
            }
        )

@app.get("/api/examples")
async def get_examples():
    """Get few-shot examples for reference"""
    try:
        orchestrator = get_components()
        examples = await orchestrator.few_shot_manager.get_relevant_examples(
            query="",
            k=5
        )
        return {"examples": examples}
        
    except Exception as e:
        logger.error(f"Error fetching examples: {str(e)}")
        return JSONResponse(
            status_code=500,
            content={
                "success": False,
                "error": f"Error fetching examples: {str(e)}"
            }
        )

@app.get("/api/health")
async def health_check():
    """Health check endpoint"""
    return {
        "status": "healthy",
        "environment": os.getenv("ENVIRONMENT", "development"),
        "version": "1.0.0"
    }

@app.get("/api/patterns")
async def get_patterns():
    """Get available query patterns"""
    try:
        orchestrator = get_components()
        patterns = orchestrator.knowledge_base.patterns
        return {"patterns": patterns}
    except Exception as e:
        logger.error(f"Error fetching patterns: {str(e)}")
        return JSONResponse(
            status_code=500,
            content={
                "success": False,
                "error": f"Error fetching patterns: {str(e)}"
            }
        )

@app.get("/api/config")
async def get_config():
    """Get API configuration"""
    return {
        "project_name": os.getenv("PROJECT_NAME"),
        "environment": os.getenv("ENVIRONMENT", "development"),
        "model": os.getenv("MODEL_NAME", "gpt-4"),
        "version": "1.0.0"
    }

# Error handlers
@app.exception_handler(HTTPException)
async def http_exception_handler(request, exc):
    return JSONResponse(
        status_code=exc.status_code,
        content={
            "success": False,
            "error": exc.detail
        }
    )

@app.exception_handler(Exception)
async def general_exception_handler(request, exc):
    logger.error(f"Unhandled exception: {str(exc)}")
    return JSONResponse(
        status_code=500,
        content={
            "success": False,
            "error": f"Internal server error: {str(exc)}"
        }
    )

@app.get("/api/test")
async def test_llm():
    """Test LLM functionality"""
    try:
        simple_prompt = """Return this exact JSON:
        {
            "test": "success"
        }"""
        response = await llm_client.get_completion(simple_prompt)
        return {"raw_response": response}
    except Exception as e:
        return {"error": str(e)}

if __name__ == "__main__":
    import uvicorn
    port = int(os.getenv("PORT", 8000))
    uvicorn.run(
        "main:app",
        host="0.0.0.0",
        port=port,
        reload=True,
        log_level="info"
    )