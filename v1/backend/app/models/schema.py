# app/models/schema.py
from pydantic import BaseModel
from typing import Optional, Dict, Any, List

class DatabaseConfig(BaseModel):
    url: str
    db_schema: Optional[str] = None

class QueryAnalysis(BaseModel):
    query_type: str
    tables: List[str]
    columns: List[str]
    conditions: Optional[List[str]] = None
    joins: Optional[List[str]] = None
    aggregations: Optional[List[str]] = None
    group_by: Optional[List[str]] = None
    order_by: Optional[List[str]] = None
    limit: Optional[int] = None

class QueryRequest(BaseModel):
    query: str
    db_schema: Optional[str] = None
    database_url: Optional[str] = None

class QueryResponse(BaseModel):
    success: bool
    sql: Optional[str] = None
    results: Optional[Dict[str, Any]] = None
    error: Optional[str] = None

class ValidationResponse(BaseModel):
    success: bool
    is_valid: bool
    issues: List[str]
    error: Optional[str] = None