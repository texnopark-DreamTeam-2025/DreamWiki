from fastapi import FastAPI, HTTPException
from transformers import AutoTokenizer, AutoModel
import torch
import numpy as np
from models import V1EmbeddingRequest, V1EmbeddingResponse, V1ErrorResponse
import logging

# Configure logging
logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

# Initialize FastAPI app
app = FastAPI(
    title="Inference API",
    description="API Specification for RuBERT Transformer Inference Service",
    version="0.0.1"
)

# Global variables for model and tokenizer
model = None
tokenizer = None


def load_rubert_model():
    """Load the RuBERT model and tokenizer"""
    global model, tokenizer
    try:
        logger.info("Loading RuBERT model...")
        # Load the pre-downloaded tokenizer and model
        tokenizer = AutoTokenizer.from_pretrained("/app/rubert", local_files_only=True, use_fast=False)
        model = AutoModel.from_pretrained("/app/rubert", local_files_only=True)
        logger.info("RuBERT model loaded successfully")
    except Exception as e:
        logger.error(f"Failed to load RuBERT model: {e}")
        raise


@app.on_event("startup")
async def startup_event():
    """Load the model on startup"""
    load_rubert_model()


@app.post("/v1/create-text-embedding",
          response_model=V1EmbeddingResponse,
          responses={422: {"model": V1ErrorResponse}, 500: {"model": V1ErrorResponse}},
          summary="Generate text embeddings using RuBERT transformer",
          operation_id="generateEmbedding")
async def generate_embedding(request: V1EmbeddingRequest):
    """
    Generate text embeddings using RuBERT transformer model.

    - **texts**: List of input texts for embedding generation
    """
    try:
        if model is None or tokenizer is None:
            raise HTTPException(status_code=500, detail="Model not loaded")

        # Tokenize all input texts
        inputs = tokenizer(request.texts, return_tensors="pt", padding=True, truncation=True, max_length=512)

        # Generate embeddings
        with torch.no_grad():
            outputs = model(**inputs)
            # Use the mean of the last hidden states as the embedding
            embeddings = outputs.last_hidden_state.mean(dim=1)

        # Convert to list of lists of floats
        embeddings_list = embeddings.tolist()

        return V1EmbeddingResponse(embeddings=embeddings_list)

    except Exception as e:
        logger.error(f"Error generating embeddings: {e}")
        raise HTTPException(status_code=500, detail=str(e))


@app.get("/health", summary="Health check endpoint")
async def health_check():
    """Health check endpoint"""
    return {"status": "ok"}


if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, host="0.0.0.0", port=8000)
