# Inference Service

This service provides text embedding generation using the RuBERT transformer model.

## Overview

The service exposes a single endpoint for generating text embeddings using the RuBERT transformer model. It's designed to be used by other services in the system that need text embeddings for natural language processing tasks.

## API Specification

The API specification is defined in the [openapi.yml](openapi.yml) file.

## Endpoints

- `POST /v1/create-text-embedding` - Generate text embedding using RuBERT transformer

## Running the Service

### Using Docker (Recommended)

```bash
# Build the Docker image
docker build -t inference-service .

# Run the service
docker run -p 8000:8000 inference-service
```

### Using Python Directly

```bash
# Install dependencies
pip install -r requirements.txt

# Run the service
uvicorn main:app --host 0.0.0.0 --port 8000
```

## Usage Example

```bash
curl -X POST "http://localhost:8000/v1/create-text-embedding" \
     -H "Content-Type: application/json" \
     -d '{"text": "Пример текста для генерации эмбеддинга"}'
```

## Health Check

```bash
curl "http://localhost:8000/health"
