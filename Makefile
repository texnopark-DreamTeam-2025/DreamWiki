.PHONY: build-frontend build-backend start-all down-all

build-frontend:
	cd frontend && make build

start-all: build-frontend
	docker compose -f infra/docker-compose/docker-compose.yaml --env-file=.env up -d

down-all:
	docker compose -f infra/docker-compose/docker-compose.yaml down
