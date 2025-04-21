# # backend/python/app/main.py
# import asyncio
# import logging
# from fastapi import FastAPI, HTTPException, Request
# from pydantic import BaseModel
# from typing import Dict, Any, List, Optional
# from .llm.client import llm_client
# from fastapi.middleware.cors import CORSMiddleware
# import logging

# try:
#     from .llm.client import LLMClient
# except ImportError as e:
#     logging.Logger.error(f"Failed to import LLM client: {str(e)}")
#     raise

# # Configure logging
# logging.basicConfig(
#     level=logging.DEBUG,
#     format='%(asctime)s - %(name)s - %(levelname)s - %(message)s'
# )
# logger = logging.getLogger(__name__)

# app = FastAPI()

# # Add CORS middleware
# app.add_middleware(
#     CORSMiddleware,
#     allow_origins=["*"],
#     allow_credentials=True,
#     allow_methods=["*"],
#     allow_headers=["*"],
# )

# class SessionRequest(BaseModel):
#     session_id: str
#     data: Dict[str, Any]

# @app.middleware("http")
# async def timeout_middleware(request: Request, call_next):
#     try:
#         return await asyncio.wait_for(call_next(request), timeout=120.0)
#     except asyncio.TimeoutError:
#         raise HTTPException(status_code=504, detail="Request timeout")
    
# def create_llm_client(provider: str, api_key: str) -> LLMClient:
#     """Create an LLM client for the specified provider"""
#     try:
#         return LLMClient(api_key=api_key, provider=provider)
#     except Exception as e:
#         logger.error(f"Failed to create LLM client for provider {provider}: {str(e)}")
#         raise HTTPException(status_code=400, detail=f"Invalid API key or configuration for {provider}")


# @app.post("/analyze")
# async def analyze_query(request: SessionRequest):
#     logger.info(f"Received analyze request for session: {request.session_id}")
#     logger.debug(f"Request data: {request.data}")

#     try:
#         question = request.data.get("question")
#         schema = request.data.get("schema")
        
#         logger.debug(f"Question: {question}")
#         logger.debug(f"Schema: {schema}")

#         if not question:
#             raise HTTPException(status_code=400, detail="Question is required")
#         if not schema:
#             raise HTTPException(status_code=400, detail="Schema is required")

#         analysis = await llm_client.analyze_query(question, schema)
#         logger.info(f"Analysis completed: {analysis}")
        
#         return {"analysis": analysis}
#     except Exception as e:
#         logger.error(f"Error in analyze_query: {str(e)}", exc_info=True)
#         raise HTTPException(status_code=500, detail=str(e))

# @app.post("/analyze_with_knowledge")
# async def analyze_with_knowledge(request: SessionRequest):
#     logger.info(f"Received analyze with knowledge request for session: {request.session_id}")
#     logger.debug(f"Request data: {request.data}")

#     try:
#         # Extract required fields
#         query = request.data.get("query")
#         schema = request.data.get("schema")
#         knowledge_context = request.data.get("knowledge_context", {})
        
#         if not query:
#             raise HTTPException(status_code=400, detail="Query is required")
#         if not schema:
#             raise HTTPException(status_code=400, detail="Schema is required")

#         # Call LLM client with knowledge context
#         analysis = await llm_client.analyze_with_knowledge(query, schema, knowledge_context)
#         logger.info(f"Analysis with knowledge completed")
        
#         return {"analysis": analysis}
#     except Exception as e:
#         logger.error(f"Error in analyze_with_knowledge: {str(e)}", exc_info=True)
#         raise HTTPException(status_code=500, detail=str(e))

# @app.post("/generate")
# async def generate_query(request: SessionRequest):
#     logger.info(f"Received generate request for session: {request.session_id}")
#     try:
#         analysis = request.data.get("analysis")
#         schema = request.data.get("schema")
        
#         if not analysis:
#             raise HTTPException(status_code=400, detail="Analysis is required")
#         if not schema:
#             raise HTTPException(status_code=400, detail="Schema is required")
            
#         query = await llm_client.generate_query(analysis, schema)
#         return {"query": query}
#     except Exception as e:
#         logger.error(f"Error in generate_query: {str(e)}", exc_info=True)
#         raise HTTPException(status_code=500, detail=str(e))

# @app.post("/generate_with_knowledge")
# async def generate_with_knowledge(request: SessionRequest):
#     logger.info(f"Received generate with knowledge request for session: {request.session_id}")
    
#     try:
#         # Extract required fields
#         analysis = request.data.get("analysis")
#         schema = request.data.get("schema")
#         knowledge_context = request.data.get("knowledge_context", {})
        
#         if not analysis:
#             raise HTTPException(status_code=400, detail="Analysis is required")
#         if not schema:
#             raise HTTPException(status_code=400, detail="Schema is required")
            
#         # Call LLM client with knowledge context
#         query = await llm_client.generate_query_with_knowledge(analysis, schema, knowledge_context)
        
#         return {"query": query}
#     except Exception as e:
#         logger.error(f"Error in generate_with_knowledge: {str(e)}", exc_info=True)
#         raise HTTPException(status_code=500, detail=str(e))

# @app.post("/validate")
# async def validate_query(request: SessionRequest):
#     try:
#         query = request.data.get("query")
#         schema = request.data.get("schema")
        
#         if not query:
#             raise HTTPException(status_code=400, detail="Query is required")
#         if not schema:
#             raise HTTPException(status_code=400, detail="Schema is required")
            
#         validation = await llm_client.validate_query(query, schema)
#         return validation
#     except Exception as e:
#         logger.error(f"Error in validate_query: {str(e)}", exc_info=True)
#         raise HTTPException(status_code=500, detail=str(e))

# @app.post("/heal")
# async def heal_query(request: SessionRequest):
#     logger.info("Received healing request")
#     try:
#         validation_result = request.data.get("validation_result")
#         original_query = request.data.get("original_query")
#         analysis = request.data.get("analysis")
#         schema = request.data.get("schema")
        
#         if not all([validation_result, original_query, analysis, schema]):
#             raise HTTPException(status_code=400, detail="Missing required healing parameters")
            
#         healing_result = await llm_client.heal_query(
#             validation_result,
#             original_query,
#             analysis,
#             schema
#         )
#         return healing_result
#     except Exception as e:
#         logger.error(f"Error in heal_query: {str(e)}", exc_info=True)
#         raise HTTPException(status_code=500, detail=str(e))

# @app.post("/process")
# async def process_query(request: Dict[str, Any]):
#     logger.info("Received process request")
#     try:
#         question = request.get("question")
#         schema = request.get("schema")
        
#         if not question or not schema:
#             raise HTTPException(status_code=400, detail="Missing question or schema")
            
#         # Step 1: Analyze query
#         analysis = await llm_client.analyze_query(question, schema)
        
#         # Step 2: Generate SQL
#         query = await llm_client.generate_query(analysis, schema)
        
#         # Step 3: Validate query
#         validation = await llm_client.validate_query(query, schema)
        
#         # Step 4: Handle validation issues
#         if not validation.get("isValid", False):
#             # Attempt to heal the query
#             healing_result = await llm_client.heal_query(
#                 validation,
#                 query,
#                 analysis,
#                 schema
#             )
            
#             if healing_result.get("healed_query"):
#                 query = healing_result["healed_query"]
#                 # Re-validate the healed query
#                 validation = await llm_client.validate_query(query, schema)
                
#                 if not validation.get("isValid", False):
#                     return {
#                         "success": False,
#                         "error": "Failed to generate a valid query after healing",
#                         "attempted_query": query
#                     }
#             else:
#                 return {
#                     "success": False,
#                     "error": "Failed to heal query",
#                     "validation": validation
#                 }
        
#         return {
#             "success": True,
#             "query": query,
#             "analysis": analysis,
#             "validation": validation
#         }
#     except Exception as e:
#         logger.error(f"Error in process_query: {str(e)}", exc_info=True)
#         raise HTTPException(status_code=500, detail=str(e))

# @app.post("/process_with_knowledge")
# async def process_with_knowledge(request: Dict[str, Any]):
#     logger.info("Received process with knowledge request")
#     try:
#         question = request.get("question")
#         schema = request.get("schema")
#         knowledge_context = request.get("knowledge_context", {})
        
#         if not question or not schema:
#             raise HTTPException(status_code=400, detail="Missing question or schema")
            
#         # Step 1: Analyze query with knowledge context
#         analysis = await llm_client.analyze_with_knowledge(question, schema, knowledge_context)
        
#         # Step 2: Generate SQL with knowledge context
#         query = await llm_client.generate_query_with_knowledge(analysis, schema, knowledge_context)
        
#         # Step 3: Validate query
#         validation = await llm_client.validate_query(query, schema)
        
#         # Step 4: Handle validation issues
#         if not validation.get("isValid", False):
#             # Attempt to heal the query
#             healing_result = await llm_client.heal_query(
#                 validation,
#                 query,
#                 analysis,
#                 schema
#             )
            
#             if healing_result.get("healed_query"):
#                 query = healing_result["healed_query"]
#                 # Re-validate the healed query
#                 validation = await llm_client.validate_query(query, schema)
                
#                 if not validation.get("isValid", False):
#                     return {
#                         "success": False,
#                         "error": "Failed to generate a valid query after healing",
#                         "attempted_query": query
#                     }
#             else:
#                 return {
#                     "success": False,
#                     "error": "Failed to heal query",
#                     "validation": validation
#                 }
        
#         return {
#             "success": True,
#             "query": query,
#             "analysis": analysis,
#             "validation": validation,
#             "knowledge_used": True
#         }
#     except Exception as e:
#         logger.error(f"Error in process_with_knowledge: {str(e)}", exc_info=True)
#         raise HTTPException(status_code=500, detail=str(e))

# @app.get("/health")
# async def health_check():
#     return {"status": "ok"}

# if __name__ == "__main__":
#     import uvicorn
#     uvicorn.run(app, host="0.0.0.0", port=8000)

# backend/python/app/main.py
import asyncio
import logging
from fastapi import FastAPI, HTTPException, Request
from pydantic import BaseModel
from typing import Dict, Any, List, Optional
from fastapi.middleware.cors import CORSMiddleware

# Import with error handling for optional dependencies
try:
    from .llm.client import LLMClient
except ImportError as e:
    logging.Logger.error(f"Failed to import LLM client: {str(e)}")
    raise

# Configure logging
logging.basicConfig(
    level=logging.DEBUG,
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s'
)
logger = logging.getLogger(__name__)

app = FastAPI()

# Add CORS middleware
app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

# Store session-specific clients
session_clients = {}

class SessionRequest(BaseModel):
    session_id: str
    data: Dict[str, Any]
    llm_config: Optional[Dict[str, str]] = None

def create_llm_client(provider: str, api_key: str) -> LLMClient:
    """Create an LLM client for the specified provider"""
    try:
        return LLMClient(api_key=api_key, provider=provider)
    except Exception as e:
        logger.error(f"Failed to create LLM client for provider {provider}: {str(e)}")
        raise HTTPException(status_code=400, detail=f"Invalid API key or configuration for {provider}")

@app.middleware("http")
async def timeout_middleware(request: Request, call_next):
    try:
        return await asyncio.wait_for(call_next(request), timeout=120.0)
    except asyncio.TimeoutError:
        raise HTTPException(status_code=504, detail="Request timeout")

@app.post("/analyze")
async def analyze_query(request: SessionRequest):
    logger.info(f"Received analyze request for session: {request.session_id}")
    logger.debug(f"Request data: {request.data}")
    
    try:
        # Get or create client for this session
        if request.llm_config and request.session_id not in session_clients:
            provider = request.llm_config.get("provider", "gemini")
            api_key = request.llm_config.get("api_key")
            session_clients[request.session_id] = create_llm_client(provider, api_key)
        
        client = session_clients.get(request.session_id, None)
        if client is None:
            # Use default client if none exists
            client = LLMClient()
            session_clients[request.session_id] = client
        
        question = request.data.get("question")
        schema = request.data.get("schema")
        
        if not question:
            raise HTTPException(status_code=400, detail="Question is required")
        if not schema:
            raise HTTPException(status_code=400, detail="Schema is required")

        analysis = await client.analyze_query(question, schema)
        logger.info(f"Analysis completed: {analysis}")
        
        return {"analysis": analysis}
    except Exception as e:
        logger.error(f"Error in analyze_query: {str(e)}", exc_info=True)
        raise HTTPException(status_code=500, detail=str(e))

@app.post("/analyze_with_knowledge")
async def analyze_with_knowledge(request: SessionRequest):
    logger.info(f"Received analyze with knowledge request for session: {request.session_id}")
    logger.debug(f"Request data: {request.data}")

    try:
        # Get or create client for this session
        if request.llm_config and request.session_id not in session_clients:
            provider = request.llm_config.get("provider", "gemini")
            api_key = request.llm_config.get("api_key")
            session_clients[request.session_id] = create_llm_client(provider, api_key)
        
        client = session_clients.get(request.session_id, None)
        if client is None:
            client = LLMClient()
            session_clients[request.session_id] = client
        
        # Extract required fields
        query = request.data.get("query")
        schema = request.data.get("schema")
        knowledge_context = request.data.get("knowledge_context", {})
        
        if not query:
            raise HTTPException(status_code=400, detail="Query is required")
        if not schema:
            raise HTTPException(status_code=400, detail="Schema is required")

        # Call LLM client with knowledge context
        analysis = await client.analyze_with_knowledge(query, schema, knowledge_context)
        logger.info(f"Analysis with knowledge completed")
        
        return {"analysis": analysis}
    except Exception as e:
        logger.error(f"Error in analyze_with_knowledge: {str(e)}", exc_info=True)
        raise HTTPException(status_code=500, detail=str(e))

@app.post("/generate")
async def generate_query(request: SessionRequest):
    logger.info(f"Received generate request for session: {request.session_id}")
    try:
        # Get or create client for this session
        if request.llm_config and request.session_id not in session_clients:
            provider = request.llm_config.get("provider", "gemini")
            api_key = request.llm_config.get("api_key")
            session_clients[request.session_id] = create_llm_client(provider, api_key)
        
        client = session_clients.get(request.session_id, None)
        if client is None:
            client = LLMClient()
            session_clients[request.session_id] = client
        
        analysis = request.data.get("analysis")
        schema = request.data.get("schema")
        
        if not analysis:
            raise HTTPException(status_code=400, detail="Analysis is required")
        if not schema:
            raise HTTPException(status_code=400, detail="Schema is required")
            
        query = await client.generate_query(analysis, schema)
        return {"query": query}
    except Exception as e:
        logger.error(f"Error in generate_query: {str(e)}", exc_info=True)
        raise HTTPException(status_code=500, detail=str(e))

@app.post("/generate_with_knowledge")
async def generate_with_knowledge(request: SessionRequest):
    logger.info(f"Received generate with knowledge request for session: {request.session_id}")
    
    try:
        # Get or create client for this session
        if request.llm_config and request.session_id not in session_clients:
            provider = request.llm_config.get("provider", "gemini")
            api_key = request.llm_config.get("api_key")
            session_clients[request.session_id] = create_llm_client(provider, api_key)
        
        client = session_clients.get(request.session_id, None)
        if client is None:
            client = LLMClient()
            session_clients[request.session_id] = client
        
        # Extract required fields
        analysis = request.data.get("analysis")
        schema = request.data.get("schema")
        knowledge_context = request.data.get("knowledge_context", {})
        
        if not analysis:
            raise HTTPException(status_code=400, detail="Analysis is required")
        if not schema:
            raise HTTPException(status_code=400, detail="Schema is required")
            
        # Call LLM client with knowledge context
        query = await client.generate_query_with_knowledge(analysis, schema, knowledge_context)
        
        return {"query": query}
    except Exception as e:
        logger.error(f"Error in generate_with_knowledge: {str(e)}", exc_info=True)
        raise HTTPException(status_code=500, detail=str(e))

@app.post("/validate")
async def validate_query(request: SessionRequest):
    try:
        # Get or create client for this session
        if request.llm_config and request.session_id not in session_clients:
            provider = request.llm_config.get("provider", "gemini")
            api_key = request.llm_config.get("api_key")
            session_clients[request.session_id] = create_llm_client(provider, api_key)
        
        client = session_clients.get(request.session_id, None)
        if client is None:
            client = LLMClient()
            session_clients[request.session_id] = client
        
        query = request.data.get("query")
        schema = request.data.get("schema")
        
        if not query:
            raise HTTPException(status_code=400, detail="Query is required")
        if not schema:
            raise HTTPException(status_code=400, detail="Schema is required")
            
        validation = await client.validate_query(query, schema)
        return validation
    except Exception as e:
        logger.error(f"Error in validate_query: {str(e)}", exc_info=True)
        raise HTTPException(status_code=500, detail=str(e))

@app.post("/heal")
async def heal_query(request: SessionRequest):
    logger.info("Received healing request")
    try:
        # Get or create client for this session
        if request.llm_config and request.session_id not in session_clients:
            provider = request.llm_config.get("provider", "gemini")
            api_key = request.llm_config.get("api_key")
            session_clients[request.session_id] = create_llm_client(provider, api_key)
        
        client = session_clients.get(request.session_id, None)
        if client is None:
            client = LLMClient()
            session_clients[request.session_id] = client
        
        validation_result = request.data.get("validation_result")
        original_query = request.data.get("original_query")
        analysis = request.data.get("analysis")
        schema = request.data.get("schema")
        
        if not all([validation_result, original_query, analysis, schema]):
            raise HTTPException(status_code=400, detail="Missing required healing parameters")
            
        healing_result = await client.heal_query(
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
async def process_query(request: Dict[str, Any]):
    logger.info("Received process request")
    try:
        session_id = request.get("session_id", "default")
        llm_config = request.get("llm_config")
        
        # Get or create client for this session
        if llm_config and session_id not in session_clients:
            provider = llm_config.get("provider", "gemini")
            api_key = llm_config.get("api_key")
            session_clients[session_id] = create_llm_client(provider, api_key)
        
        client = session_clients.get(session_id, None)
        if client is None:
            client = LLMClient()
            session_clients[session_id] = client
        
        question = request.get("question")
        schema = request.get("schema")
        
        if not question or not schema:
            raise HTTPException(status_code=400, detail="Missing question or schema")
            
        # Step 1: Analyze query
        analysis = await client.analyze_query(question, schema)
        
        # Step 2: Generate SQL
        query = await client.generate_query(analysis, schema)
        
        # Step 3: Validate query
        validation = await client.validate_query(query, schema)
        
        # Step 4: Handle validation issues
        if not validation.get("isValid", False):
            # Attempt to heal the query
            healing_result = await client.heal_query(
                validation,
                query,
                analysis,
                schema
            )
            
            if healing_result.get("healed_query"):
                query = healing_result["healed_query"]
                # Re-validate the healed query
                validation = await client.validate_query(query, schema)
                
                if not validation.get("isValid", False):
                    return {
                        "success": False,
                        "error": "Failed to generate a valid query after healing",
                        "attempted_query": query
                    }
            else:
                return {
                    "success": False,
                    "error": "Failed to heal query",
                    "validation": validation
                }
        
        return {
            "success": True,
            "query": query,
            "analysis": analysis,
            "validation": validation
        }
    except Exception as e:
        logger.error(f"Error in process_query: {str(e)}", exc_info=True)
        raise HTTPException(status_code=500, detail=str(e))

@app.post("/process_with_knowledge")
async def process_with_knowledge(request: Dict[str, Any]):
    logger.info("Received process with knowledge request")
    try:
        session_id = request.get("session_id", "default")
        llm_config = request.get("llm_config")
        
        # Get or create client for this session
        if llm_config and session_id not in session_clients:
            provider = llm_config.get("provider", "gemini")
            api_key = llm_config.get("api_key")
            session_clients[session_id] = create_llm_client(provider, api_key)
        
        client = session_clients.get(session_id, None)
        if client is None:
            client = LLMClient()
            session_clients[session_id] = client
        
        question = request.get("question")
        schema = request.get("schema")
        knowledge_context = request.get("knowledge_context", {})
        
        if not question or not schema:
            raise HTTPException(status_code=400, detail="Missing question or schema")
            
        # Step 1: Analyze query with knowledge context
        analysis = await client.analyze_with_knowledge(question, schema, knowledge_context)
        
        # Step 2: Generate SQL with knowledge context
        query = await client.generate_query_with_knowledge(analysis, schema, knowledge_context)
        
        # Step 3: Validate query
        validation = await client.validate_query(query, schema)
        
        # Step 4: Handle validation issues
        if not validation.get("isValid", False):
            # Attempt to heal the query
            healing_result = await client.heal_query(
                validation,
                query,
                analysis,
                schema
            )
            
            if healing_result.get("healed_query"):
                query = healing_result["healed_query"]
                # Re-validate the healed query
                validation = await client.validate_query(query, schema)
                
                if not validation.get("isValid", False):
                    return {
                        "success": False,
                        "error": "Failed to generate a valid query after healing",
                        "attempted_query": query
                    }
            else:
                return {
                    "success": False,
                    "error": "Failed to heal query",
                    "validation": validation
                }
        
        return {
            "success": True,
            "query": query,
            "analysis": analysis,
            "validation": validation,
            "knowledge_used": True
        }
    except Exception as e:
        logger.error(f"Error in process_with_knowledge: {str(e)}", exc_info=True)
        raise HTTPException(status_code=500, detail=str(e))

@app.on_event("shutdown")
async def shutdown_event():
    """Clean up resources on shutdown"""
    logger.info("Cleaning up LLM clients...")
    # Clean up any resources if needed
    session_clients.clear()

@app.get("/health")
async def health_check():
    return {"status": "ok"}

if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, host="0.0.0.0", port=8000)