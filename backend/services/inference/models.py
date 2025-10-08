from pydantic import BaseModel
from typing import List, Union


class V1EmbeddingRequest(BaseModel):
    text: str


class V1EmbeddingResponse(BaseModel):
    embedding: List[Union[float, int]]


class V1ErrorResponse(BaseModel):
    error: str
