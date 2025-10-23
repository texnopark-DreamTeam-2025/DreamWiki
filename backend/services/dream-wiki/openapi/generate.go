package openapi

//go:generate go tool oapi-codegen --config=oapi-codegen.yml ../openapi.yml
//go:generate go tool oapi-codegen --config=oapi-codegen-inference.yml ../../../services/inference/openapi.yml
//go:generate go tool oapi-codegen --config=oapi-codegen-ywiki.yml openapi-ywiki.yml
//go:generate go tool oapi-codegen --config=oapi-codegen-ycloud.yml openapi-ycloud.yml
//go:generate go tool oapi-codegen --config=oapi-codegen-github.yml openapi-github.yml
