.PHONY: build-frontend build-backend start-all

build-frontend:
	cd frontend && make build

start-all: build-frontend
	docker compose -f infra/docker-compose/docker-compose.yaml --env-file .env up --build -d
