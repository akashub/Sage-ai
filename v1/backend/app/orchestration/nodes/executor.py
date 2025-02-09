# app/orchestration/nodes/executor.py
from sqlalchemy import create_engine, text
from sqlalchemy.exc import SQLAlchemyError
from typing import Dict, Any, Optional
from langchain_openai import ChatOpenAI
from pydantic import BaseModel


class ExecutorNode:
    def __init__(self):
        self.engines = {}
    
    def _get_engine(self, database_url: str):
        """Get or create SQLAlchemy engine"""
        if database_url not in self.engines:
            self.engines[database_url] = create_engine(database_url)
        return self.engines[database_url]

    async def execute(self, 
                     sql: str, 
                     database_url: str,
                     params: Optional[Dict] = None) -> Dict[str, Any]:
        """Execute SQL query"""
        try:
            engine = self._get_engine(database_url)
            
            with engine.connect() as connection:
                # Start transaction
                with connection.begin():
                    # Execute query
                    result = connection.execute(
                        text(sql),
                        parameters=params or {}
                    )
                    
                    # Fetch results
                    if result.returns_rows:
                        rows = result.fetchall()
                        results = [dict(row._mapping) for row in rows]
                    else:
                        results = {"affected_rows": result.rowcount}
                    
                    return {
                        "success": True,
                        "results": results,
                        "error": None
                    }
                    
        except SQLAlchemyError as e:
            return {
                "success": False,
                "results": None,
                "error": f"Database error: {str(e)}"
            }
        except Exception as e:
            return {
                "success": False,
                "results": None,
                "error": f"Execution failed: {str(e)}"
            }