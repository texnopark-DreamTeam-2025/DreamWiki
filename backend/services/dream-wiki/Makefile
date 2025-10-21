define setup_env
	$(eval ENV_FILE := ../../../.env)
	@echo "==> Load $(ENV_FILE)"
	$(eval export $(shell sed 's/=.*//' $(ENV_FILE)))
	$(eval include $(ENV_FILE))
endef

load-env:
	$(call setup_env)

run: load-env
	SERVER_PORT=7000 LOG_MODE=dev YDB_DSN=grpc://127.0.0.1:2136/?database=/local \
	INFERENCE_API_URL=http://localhost:8082 JWT_SECRET_KEY=secret \
	go run cmd/main/main.go

generate:
	@echo "==> Generating..."
	go generate ./...

dev: generate run

.PHONY: run dev generate
