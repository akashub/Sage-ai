# # backend/app/llm/client.py
# from langchain_openai import ChatOpenAI, OpenAIEmbeddings
# import google.generativeai as genai
# from typing import Optional, Literal
# import os
# import time
# import asyncio
# import json
# from dotenv import load_dotenv

# load_dotenv()

# class LLMClient:
#     _instance = None
#     _initialized = False

#     def __new__(cls):
#         if cls._instance is None:
#             cls._instance = super().__new__(cls)
#         return cls._instance

#     def __init__(self):
#         if not LLMClient._initialized:
#             self._last_call_time = 0
#             self._rate_limit_delay = float(os.getenv("RATE_LIMIT_DELAY", "1.0"))
#             self._provider = os.getenv("LLM_PROVIDER", "gemini")  # 'openai' or 'gemini'
#             self._initialize_clients()
#             LLMClient._initialized = True

#     def _initialize_clients(self):
#         """Initialize LLM clients based on provider"""
#         if self._provider == "openai":
#             self.llm = ChatOpenAI(
#                 temperature=float(os.getenv("LLM_TEMPERATURE", "0")),
#                 model_name=os.getenv("MODEL_NAME", "gpt-4"),
#                 api_key=os.getenv("OPENAI_API_KEY"),
#             )
#             self.embeddings = OpenAIEmbeddings(
#                 api_key=os.getenv("OPENAI_API_KEY")
#             )
#         elif self._provider == "gemini":
#             genai.configure(api_key=os.getenv("GOOGLE_API_KEY"))
#             self.llm = genai.GenerativeModel(os.getenv("GEMINI_MODEL", "gemini-1.5-pro"))
#             # Note: You might need a different embedding model for Gemini
#             self.embeddings = OpenAIEmbeddings(
#                 api_key=os.getenv("OPENAI_API_KEY")
#             )
#         else:
#             raise ValueError(f"Unsupported LLM provider: {self._provider}")

#     async def _handle_rate_limit(self):
#         """Handle rate limiting"""
#         now = time.time()
#         time_since_last_call = now - self._last_call_time
#         if time_since_last_call < self._rate_limit_delay:
#             await asyncio.sleep(self._rate_limit_delay - time_since_last_call)
#         self._last_call_time = time.time()

#     async def get_completion(self, prompt: str) -> str:
#         """Get completion with rate limiting"""
#         await self._handle_rate_limit()
#         try:
#             if self._provider == "openai":
#                 response = await self.llm.ainvoke(prompt)
#                 return response.content
#             else:  # gemini
#                 response = await asyncio.to_thread(
#                     self.llm.generate_content, prompt
#                 )
#                 return response.text
#         except Exception as e:
#             print(f"Error getting completion: {e}")
#             raise

#     async def get_embedding(self, text: str) -> list:
#         """Get embedding with rate limiting"""
#         await self._handle_rate_limit()
#         try:
#             # Currently using OpenAI embeddings for both providers
#             embeddings = await self.embeddings.aembed_documents([text])
#             return embeddings[0]
#         except Exception as e:
#             print(f"Error getting embedding: {e}")
#             raise

#     def switch_provider(self, provider: Literal["openai", "gemini"]):
#         """Switch between LLM providers"""
#         if provider not in ["openai", "gemini"]:
#             raise ValueError("Provider must be either 'openai' or 'gemini'")
#         self._provider = provider
#         self._initialize_clients()

# # Global instance
# llm_client = LLMClient()

# backend/app/llm/client.py
import hashlib
from typing import Optional
import google.generativeai as genai
import os
import time
import asyncio
import json
from dotenv import load_dotenv

load_dotenv()

class LLMClient:
    _instance = None
    _initialized = False

    def __new__(cls):
        if cls._instance is None:
            cls._instance = super().__new__(cls)
        return cls._instance

    def __init__(self):
        if not LLMClient._initialized:
            self._last_call_time = 0
            self._rate_limit_delay = float(os.getenv("RATE_LIMIT_DELAY", "1.0"))
            self._initialize_gemini()
            LLMClient._initialized = True

    def _initialize_gemini(self):
        """Initialize Gemini LLM client"""
        try:
            genai.configure(api_key=os.getenv("GEMINI_API_KEY"))
            self.model = genai.GenerativeModel(os.getenv("GEMINI_MODEL_NAME", "gemini-1.5-pro"))
            print("Gemini initialized successfully")
        except Exception as e:
            print(f"Error initializing Gemini: {e}")
            raise

    async def _handle_rate_limit(self):
        """Handle rate limiting"""
        now = time.time()
        time_since_last_call = now - self._last_call_time
        if time_since_last_call < self._rate_limit_delay:
            await asyncio.sleep(self._rate_limit_delay - time_since_last_call)
        self._last_call_time = time.time()

    async def get_completion(self, prompt: str) -> str:
        """Get completion with rate limiting"""
        await self._handle_rate_limit()
        try:
            response = await asyncio.to_thread(
                self.model.generate_content, 
                prompt
            )
            return response.text
        except Exception as e:
            print(f"Error getting completion from Gemini: {e}")
            raise

    async def get_embedding(self, text: str) -> list:
        """Get embedding for text using a simple hash-based approach"""
        await self._handle_rate_limit()
        try:
            # Using hash for now as Gemini doesn't provide embeddings
            hash_value = hashlib.sha256(text.encode()).digest()
            # Convert to 1536 dimensional vector (same as OpenAI's embeddings)
            vector = []
            for byte in hash_value:
                # Convert each byte to 6 float values between -1 and 1
                for i in range(6):
                    vector.append(((byte >> i) & 1) * 2 - 1)
            return vector * 8  # Repeat to get 1536 dimensions
        except Exception as e:
            print(f"Error generating embedding: {e}")
            raise

    async def get_completion(self, prompt: str) -> str:
        await self._handle_rate_limit()
        try:
            response = await asyncio.to_thread(
                self.model.generate_content,
                prompt,
                generation_config={
                    'temperature': 0.1,
                    'top_p': 1,
                    'top_k': 1,
                    'max_output_tokens': 2048,
                }
            )
            
            # Add debug logging
            print(f"Raw Gemini response: {response}")
            
            if hasattr(response, 'text'):
                # Clean the response to ensure it's valid JSON
                text = response.text.strip()
                # If response starts after JSON, find the JSON part
                if '{' in text:
                    text = text[text.find('{'):]
                if '}' in text:
                    text = text[:text.rfind('}')+1]
                return text
            return str(response)
            
        except Exception as e:
            print(f"Error getting completion from Gemini: {e}")
            raise

# Global instance
llm_client = LLMClient()