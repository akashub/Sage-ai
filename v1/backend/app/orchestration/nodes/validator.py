import json
import sqlparse
import re
from typing import List, Tuple
from typing import Dict, Any, Optional
from langchain_openai import ChatOpenAI
from pydantic import BaseModel
from ....app.llm.client import  llm_client

class ValidatorNode:
    # def __init__(self, llm: Optional[ChatOpenAI] = None):
    #     self.llm = llm or ChatOpenAI(temperature=0)
    #     self.risk_patterns = [
    #         (r';\s*DROP', 'Potential DROP command'),
    #         (r';\s*DELETE', 'Potential DELETE command'),
    #         (r';\s*UPDATE', 'Potential UPDATE command'),
    #         (r'UNION\s+SELECT', 'UNION SELECT pattern'),
    #         (r'--', 'SQL comment'),
    #         (r'/\*.*?\*/', 'Block comment')
    #     ]
    def __init__(self):
        self.risk_patterns = [
            (r';\s*DROP', 'Potential DROP command'),
            (r';\s*DELETE', 'Potential DELETE command'),
            (r';\s*UPDATE', 'Potential UPDATE command'),
            (r'UNION\s+SELECT', 'UNION SELECT pattern'),
            (r'--', 'SQL comment'),
            (r'/\*.*?\*/', 'Block comment')
        ]

    def _check_syntax(self, sql: str) -> Tuple[bool, List[str]]:
        """Validate SQL syntax"""
        try:
            parsed = sqlparse.parse(sql)
            if not parsed or not parsed[0].tokens:
                return False, ["Invalid SQL syntax"]
            return True, []
        except Exception as e:
            return False, [f"Syntax error: {str(e)}"]

    def _check_security(self, sql: str) -> Tuple[bool, List[str]]:
        """Check for security issues"""
        issues = []
        for pattern, message in self.risk_patterns:
            if re.search(pattern, sql, re.IGNORECASE):
                issues.append(f"Security risk: {message}")
        return len(issues) == 0, issues

    async def _check_semantics(self, sql: str, schema: Optional[str]) -> Tuple[bool, List[str]]:
        """Validate SQL semantics using LLM"""
        try:
            prompt = f"""Validate this SQL query semantics:
            SQL: {sql}
            
            Schema (if available):
            {schema or 'No schema provided'}
            
            Check for:
            1. Table existence
            2. Column existence
            3. Join conditions
            4. Type compatibility
            5. Logical errors
            
            Return JSON:
            {{
                "is_valid": boolean,
                "issues": ["list_of_issues"]
            }}
            """
            
            response = await llm_client.get_completion(prompt)
            # validation = response.content
            validation = json.loads(response)

            return validation["is_valid"], validation.get("issues", [])
        except Exception as e:
            return False, [f"Semantic validation failed: {str(e)}"]

    async def validate(self, 
                      sql: str, 
                      schema: Optional[str] = None) -> Dict[str, Any]:
        """Complete SQL validation"""
        try:
            # Syntax check
            syntax_valid, syntax_issues = self._check_syntax(sql)
            if not syntax_valid:
                return {
                    "success": False,
                    "is_valid": False,
                    "issues": syntax_issues,
                    "error": None
                }

            # Security check
            security_valid, security_issues = self._check_security(sql)
            if not security_valid:
                return {
                    "success": True,
                    "is_valid": False,
                    "issues": security_issues,
                    "error": None
                }

            # Semantic check
            semantic_valid, semantic_issues = await self._check_semantics(sql, schema)
            
            all_issues = syntax_issues + security_issues + semantic_issues
            is_valid = syntax_valid and security_valid and semantic_valid
            
            return {
                "success": True,
                "is_valid": is_valid,
                "issues": all_issues,
                "error": None
            }
        except Exception as e:
            return {
                "success": False,
                "is_valid": False,
                "issues": [],
                "error": f"Validation failed: {str(e)}"
            }