.PHONY: start-backend

start-backend:
	docker compose -f backend/docker-compose.yaml --env-file .env up -d
