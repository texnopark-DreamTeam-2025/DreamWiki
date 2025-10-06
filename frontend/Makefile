.PHONY: generate-client build

generate-client:
	npx @openapitools/openapi-generator-cli generate -i ../backend/services/dream-wiki/openapi.yml -g typescript -o codegen

build: generate-client
	npm run build
