# app/orchestration/factories.py
from typing import Optional
from ..knowledge.base import KnowledgeBase
from ..knowledge.few_shot import FewShotManager
from .nodes.analyzer import AnalyzerNode
from .nodes.generator import GeneratorNode
from .nodes.validator import ValidatorNode
from .nodes.executor import ExecutorNode

class NodeFactory:
    """Factory for creating and managing processing nodes"""
    def __init__(self, 
                 knowledge_base: Optional[KnowledgeBase] = None,
                 few_shot_manager: Optional[FewShotManager] = None):
        """Initialize factory with knowledge components"""
        self._instances = {}
        self.knowledge_base = knowledge_base
        self.few_shot_manager = few_shot_manager

    def get_analyzer(self) -> AnalyzerNode:
        """Get or create AnalyzerNode instance"""
        if 'analyzer' not in self._instances:
            self._instances['analyzer'] = AnalyzerNode()
        return self._instances['analyzer']

    def get_generator(self) -> GeneratorNode:
        """Get or create GeneratorNode instance"""
        if 'generator' not in self._instances:
            self._instances['generator'] = GeneratorNode(
                knowledge_base=self.knowledge_base
            )
        return self._instances['generator']

    def get_validator(self) -> ValidatorNode:
        """Get or create ValidatorNode instance"""
        if 'validator' not in self._instances:
            self._instances['validator'] = ValidatorNode()
        return self._instances['validator']

    def get_executor(self) -> ExecutorNode:
        """Get or create ExecutorNode instance"""
        if 'executor' not in self._instances:
            self._instances['executor'] = ExecutorNode()
        return self._instances['executor']