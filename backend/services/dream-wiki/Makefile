# Базовые команды
run:
	go run cmd/main/main.go

# Генерация кода из OpenAPI спецификации
generate:
	@echo "==> Generating API code..."
	oapi-codegen --config=oapi-codegen.yml openapi.yml

# Тестирование
run_tests:
	@echo "==> Running tests..."
	go test -coverprofile coverage_raw.out -v ./...

test: run_tests
	@echo "==> Calculating coverage..."
	grep -v "mock" coverage_raw.out > coverage.out
	go tool cover -func=coverage.out
	go tool cover -html=coverage.out -o=coverage.html
	@echo "==> Done! Check dream-wiki/coverage.html file!"

short_test:
	go test -cover ./...

# Development команды
dev: generate run

.PHONY: run test run_tests short_test docker_build docker_up docker_down docker_restart docker_logs docker_rebuild run_docker dev generate
