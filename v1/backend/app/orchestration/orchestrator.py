
# app/orchestration/orchestrator.py
from typing import Dict, Any, Optional
from .factories import NodeFactory
from ..knowledge.base import KnowledgeBase
from ..knowledge.few_shot import FewShotManager
from langchain_openai import ChatOpenAI
from pydantic import BaseModel

class QueryContext(BaseModel):
    question: str
    db_schema: Optional[str] = None
    database_url: Optional[str] = None  # Made optional to match process_query
    current_step: str = "INIT"
    analyzed_intent: Optional[Dict] = None
    generated_sql: Optional[str] = None
    validation_result: Optional[Dict] = None
    execution_result: Optional[Dict] = None
    error: Optional[str] = None

class SQLAssistantOrchestrator:
    # def __init__(self, 
    #              llm: Optional[ChatOpenAI] = None,
    #              few_shot_manager: Optional[FewShotManager] = None,
    #              knowledge_base: Optional[KnowledgeBase] = None):
    #     """Initialize orchestrator with optional components"""
    #     self.llm = llm or ChatOpenAI(temperature=0)
    #     self.few_shot_manager = few_shot_manager
    #     self.knowledge_base = knowledge_base
    #     self.factory = NodeFactory()
        
    #     # Initialize nodes from NodeFactory   
    #     self.analyzer = self.factory.get_analyzer()
    #     self.generator = self.factory.get_generator()
    #     self.validator = self.factory.get_validator()
    #     self.executor = self.factory.get_executor()
    def __init__(self, 
                 few_shot_manager: Optional[FewShotManager] = None,
                 knowledge_base: Optional[KnowledgeBase] = None):
        """Initialize orchestrator with optional components"""
        self.few_shot_manager = few_shot_manager
        self.knowledge_base = knowledge_base
        self.factory = NodeFactory()
        
        # Initialize nodes from NodeFactory   
        self.analyzer = self.factory.get_analyzer()
        self.generator = self.factory.get_generator()
        self.validator = self.factory.get_validator()
        self.executor = self.factory.get_executor()
    async def process_query(
        self,
        question: str,
        db_schema: Optional[str] = None,
        database_url: Optional[str] = None,
    ) -> Dict[str, Any]:
        """Process a natural language query through the entire pipeline"""
        try:
            # Create query context
            context = QueryContext(
                question=question,
                db_schema=db_schema,
                database_url=database_url,
            )

            # 1. Get examples and patterns from Knowledge Layer
            if self.few_shot_manager and self.knowledge_base:
                relevant_examples = await self.few_shot_manager.get_relevant_examples(question)
                patterns = self.knowledge_base.identify_patterns(question)
            else:
                relevant_examples = []
                patterns = []

            # 2. Analyze query
            context.current_step = "ANALYSIS"
            analysis_result = await self.analyzer.analyze(
                question=question,
                schema=db_schema,
                examples=relevant_examples,
                patterns=patterns
            )
            if not analysis_result["success"]:
                return {
                    "success": False,
                    "error": analysis_result["error"]
                }
            context.analyzed_intent = analysis_result["analysis"]

            # 3. Generate SQL
            context.current_step = "GENERATION"
            generation_result = await self.generator.generate(
                analysis=analysis_result["analysis"],
                question=question,
                schema=db_schema,
                examples=relevant_examples
            )
            if not generation_result["success"]:
                return {
                    "success": False,
                    "error": generation_result["error"]
                }
            context.generated_sql = generation_result["sql"]

            # 4. Validate SQL
            context.current_step = "VALIDATION"
            validation_result = await self.validator.validate(
                sql=generation_result["sql"],
                schema=db_schema
            )
            if not validation_result["success"] or not validation_result["is_valid"]:
                return {
                    "success": False,
                    "error": f"Validation failed: {validation_result.get('issues', [])}"
                }
            context.validation_result = validation_result

            # 5. Execute SQL if database_url provided
            if database_url:
                context.current_step = "EXECUTION"
                execution_result = await self.executor.execute(
                    sql=generation_result["sql"],
                    database_url=database_url
                )
                if not execution_result["success"]:
                    return {
                        "success": False,
                        "error": execution_result["error"]
                    }
                context.execution_result = execution_result

                return {
                    "success": True,
                    "sql": generation_result["sql"],
                    "results": execution_result["results"],
                    "context": context.dict(exclude_none=True),
                    "error": None
                }

            # Return just SQL if no database_url
            return {
                "success": True,
                "sql": generation_result["sql"],
                "results": None,
                "context": context.dict(exclude_none=True),
                "error": None
            }

        except Exception as e:
            return {
                "success": False,
                "error": f"Processing failed: {str(e)}"
            }