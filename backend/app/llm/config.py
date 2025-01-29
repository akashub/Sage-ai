# backend/app/llm/config.py
from pydantic import BaseModel
from typing import Optional, Literal

class LLMConfig(BaseModel):
    """LLM Configuration Settings"""
    provider: Literal["openai", "gemini"] = "gemini"
    model_name: str = "gemini-1.5-flash"  # or "gemini-1.5-pro"
    temperature: float = 0.0
    rate_limit_delay: float = 1.0
    max_retries: int = 3
    timeout: int = 30

class EmbeddingConfig(BaseModel):
    """Embedding Configuration Settings"""
    model_name: str = "text-embedding-ada-002"
    rate_limit_delay: float = 1.0
    dimension: int = 1536

# Default configurations
DEFAULT_LLM_CONFIG = LLMConfig()
DEFAULT_EMBEDDING_CONFIG = EmbeddingConfig()