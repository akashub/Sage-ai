# app/core/config.py
from pydantic_settings import BaseSettings

class Settings(BaseSettings):
    PROJECT_NAME: str = "Sage.ai"
    API_V1_STR: str = "/api/v1"
    OPENAI_API_KEY: str
    MODEL_NAME: str = "gpt-4"

    class Config:
        env_file = ".env"

settings = Settings()