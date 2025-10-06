# Базовые команды
run:
	SERVER_PORT=7000 LOG_MODE=dev YDB_DSN=grpc://localhost:2136/?database=/local go run cmd/main/main.go

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

dev: generate run

.PHONY: run test run_tests short_test dev generate
