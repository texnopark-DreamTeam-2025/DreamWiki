from pydantic import BaseModel
from typing import List, Union


class V1EmbeddingRequest(BaseModel):
    texts: List[str]


class V1EmbeddingResponse(BaseModel):
    embeddings: List[List[Union[float, int]]]


class V1ErrorResponse(BaseModel):
    error: str
