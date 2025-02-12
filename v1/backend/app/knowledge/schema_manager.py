# v1/backend/app/knowledge/schema_manager.py
from typing import Dict, Optional
import json

class SchemaManager:
    def __init__(self):
        self.schemas: Dict[str, dict] = {}
    
    def add_schema(self, name: str, schema_dict: dict):
        """Add or update a schema"""
        self.schemas[name] = schema_dict
    
    def get_schema(self, name: str) -> Optional[dict]:
        """Retrieve a schema by name"""
        return self.schemas.get(name)
    
    def format_schema_for_prompt(self, name: str) -> str:
        """Format schema for inclusion in prompts"""
        schema = self.get_schema(name)
        if not schema:
            return "No schema available"
            
        formatted = "Database Schema:\n"
        for table, columns in schema.items():
            formatted += f"Table: {table}\n"
            formatted += "Columns:\n"
            for col, details in columns.items():
                formatted += f"- {col}: {details['type']}"
                if details.get('description'):
                    formatted += f" ({details['description']})"
                formatted += "\n"
            formatted += "\n"
        return formatted