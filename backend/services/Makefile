# Базовые команды
run:
	cd dream-wiki && go run cmd/main/main.go

# Генерация кода из OpenAPI спецификации
generate:
	@echo "==> Generating API code..."
	cd dream-wiki && oapi-codegen --config=oapi-codegen.yml openapi.yml

# Тестирование
run_tests:
	@echo "==> Running tests..."
	cd dream-wiki && go test -coverprofile coverage_raw.out -v ./...

test: run_tests
	@echo "==> Calculating coverage..."
	cd dream-wiki && grep -v "mock" coverage_raw.out > coverage.out
	cd dream-wiki && go tool cover -func=coverage.out
	cd dream-wiki && go tool cover -html=coverage.out -o=coverage.html
	@echo "==> Done! Check dream-wiki/coverage.html file!"

short_test:
	cd dream-wiki && go test -cover ./...

# Docker команды
docker_build:
	@echo "==> Building Docker images..."
	docker-compose build

docker_up:
	@echo "==> Starting containers..."
	docker-compose up

docker_down:
	@echo "==> Stopping containers..."
	docker-compose down

docker_restart: docker_down docker_up

docker_logs:
	@echo "==> Showing logs..."
	docker-compose logs -f

docker_rebuild: docker_down docker_build docker_up

run_docker: docker_build docker_up

# Development команды
dev: generate run

.PHONY: run test run_tests short_test docker_build docker_up docker_down docker_restart docker_logs docker_rebuild run_docker dev generate