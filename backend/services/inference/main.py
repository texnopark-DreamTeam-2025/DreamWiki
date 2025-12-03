from fastapi import FastAPI, HTTPException
from transformers import AutoTokenizer, AutoModel
import torch
import numpy as np
from models import (
    V1EmbeddingRequest,
    V1EmbeddingResponse,
    V1ErrorResponse,
    V1StemmingRequest,
    V1StemmingResponse,
)
import logging
import re
import nltk
from nltk.stem import SnowballStemmer
import pymorphy3 as pymorphy2
from stop_words import STOP_WORDS
import unicodedata

logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)


russian_chars_regex = re.compile(r"[а-яА-ЯёЁ]")

app = FastAPI(
    title="Inference API",
    description="API Specification for RuBERT Transformer Inference Service",
    version="0.0.1",
)

model = None
tokenizer = None

stemmer = None
morph = None


def load_rubert_model():
    global model, tokenizer
    try:
        logger.info("Loading RuBERT model...")
        tokenizer = AutoTokenizer.from_pretrained(
            "/app/rubert", local_files_only=True, use_fast=False
        )
        model = AutoModel.from_pretrained("/app/rubert", local_files_only=True)
        logger.info("RuBERT model loaded successfully")
    except Exception as e:
        logger.error(f"Failed to load RuBERT model: {e}")
        raise


def load_stemming_tools():
    global stemmer, morph
    try:
        logger.info("Loading stemming tools...")
        stemmer = SnowballStemmer("russian")
        morph = pymorphy2.MorphAnalyzer()
        logger.info("Stemming tools loaded successfully")
    except Exception as e:
        logger.error(f"Failed to load stemming tools: {e}")
        raise


@app.on_event("startup")
async def startup_event():
    load_rubert_model()
    load_stemming_tools()


@app.post(
    "/v1/create-text-embedding",
    response_model=V1EmbeddingResponse,
    responses={422: {"model": V1ErrorResponse}, 500: {"model": V1ErrorResponse}},
    summary="Generate text embeddings using RuBERT transformer",
    operation_id="generateEmbedding",
)
async def generate_embedding(request: V1EmbeddingRequest):
    try:
        if model is None or tokenizer is None:
            raise HTTPException(status_code=500, detail="Model not loaded")

        inputs = tokenizer(
            request.texts,
            return_tensors="pt",
            padding=True,
            truncation=True,
            max_length=512,
        )

        with torch.no_grad():
            outputs = model(**inputs)
            embeddings = outputs.last_hidden_state.mean(dim=1)

        embeddings_list = embeddings.tolist()

        return V1EmbeddingResponse(embeddings=embeddings_list)

    except Exception as e:
        logger.error(f"Error generating embeddings: {e}")
        raise HTTPException(status_code=500, detail=str(e))


@app.get("/health", summary="Health check endpoint")
async def health_check():
    return {"status": "ok"}


def is_russian(word):
    return bool(russian_chars_regex.fullmatch(word))


def is_number(token):
    return bool(re.match(r"^\d+([.,]\d+)?$", token))


def contains_middle_dot(token):
    return bool(re.search(r".\..", token))


def split_and_remove_non_letters(word):
    # Split by any character that is not a Latin or Cyrillic letter
    parts = re.split(r"[^a-zA-Zа-яА-ЯёЁ]", word)
    return [part for part in parts if part]


@app.post(
    "/v1/stemming",
    response_model=V1StemmingResponse,
    responses={422: {"model": V1ErrorResponse}, 500: {"model": V1ErrorResponse}},
    summary="Generate stems for paragraphs",
    operation_id="generateStems",
)
async def generate_stems(request: V1StemmingRequest):
    try:
        if stemmer is None or morph is None:
            raise HTTPException(status_code=500, detail="Stemming tools not loaded")

        result = []
        for paragraph in request.paragraphs:
            words = re.split(r"[\\ ]+", paragraph)
            words = [unicodedata.normalize('NFC', word) for word in words]
            stems = []
            for word in words:
                word = word.lower()
                if word in STOP_WORDS:
                    continue

                if not word:
                    continue

                if is_number(word):
                    stems.append(word)
                    continue

                if contains_middle_dot(word):
                    stems.append(word)
                    continue

                cleaned_words = split_and_remove_non_letters(word)

                for cleaned_word in cleaned_words:
                    if not cleaned_word:
                        continue

                    if russian_chars_regex.search(cleaned_word):
                        parsed = morph.parse(cleaned_word)[0]
                        stems.append(parsed.normal_form)
                    else:
                        stems.append(stemmer.stem(cleaned_word))
            result.append(stems)

        return V1StemmingResponse(stems=result)

    except Exception as e:
        logger.error(f"Error generating stems: {e}")
        raise HTTPException(status_code=500, detail=str(e))


if __name__ == "__main__":
    import uvicorn

    uvicorn.run(app, host="0.0.0.0", port=8000)
