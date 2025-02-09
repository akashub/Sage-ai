from typing import Dict, Any, List, Optional
from langchain_openai import ChatOpenAI
from pydantic import BaseModel
from ....app.knowledge.base import KnowledgeBase
from ....app.llm.client import llm_client

class GeneratorNode:
    # def __init__(self, 
    #              llm: Optional[ChatOpenAI] = None,
    #              knowledge_base: Optional[KnowledgeBase] = None):
    #     self.llm = llm or ChatOpenAI(temperature=0)
    #     self.knowledge_base = knowledge_base or KnowledgeBase()
    
    def __init__(self, 
                 knowledge_base: Optional[KnowledgeBase] = None):
        self.knowledge_base = knowledge_base or KnowledgeBase()
        
    async def generate(self, 
                      analysis: Dict[str, Any], 
                      question: str,
                      schema: Optional[str] = None,
                      examples: Optional[List[Any]] = None) -> Dict[str, Any]:
        """Generate SQL from analyzed components"""
        try:
            # Get similar examples
            similar_examples = await self.knowledge_base.find_similar_examples(
                query=question,
                k=3
            )
            
            prompt = self._create_generation_prompt(
                analysis=analysis,
                question=question,
                schema=schema,
                similar_examples=similar_examples
            )
            
            response = await llm_client.get_completion(prompt)
            return {
                "success": True,
                # "sql": response.content.strip(),
                "sql": response.strip(),
                "error": None
            }
        except Exception as e:
            return {
                "success": False,
                "sql": None,
                "error": f"SQL generation failed: {str(e)}"
            }

    def _create_generation_prompt(self,
                                analysis: Dict[str, Any],
                                question: str,
                                schema: Optional[str],
                                similar_examples: List[Any]) -> str:
        prompt = f"""Generate a SQL query based on this analysis and similar examples.

Analysis: {analysis}
Original Question: {question}
Schema Context:
{schema or 'No schema provided'}

Similar examples from our knowledge base:
"""

        for i, example in enumerate(similar_examples, 1):
            prompt += f"""
Example {i}:
Natural Query: {example.natural_query}
SQL Query: {example.sql_query}
"""

        prompt += """
Requirements:
1. Use proper SQL syntax
2. Include necessary JOINs
3. Handle NULL values
4. Use appropriate quoting
5. Include any required type casting

Return only the SQL query without explanations.
"""
        return prompt