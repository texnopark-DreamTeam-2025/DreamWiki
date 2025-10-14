package config

import (
	"fmt"
	"os"

	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/utils/logger"
	"go.uber.org/zap"
)

type Config struct {
	LogMode          string
	ServerPort       string
	YDBDSN           string
	InferenceAPIURL  string
	JWTSecretKey     string
	YWikiToken       string
	YandexCloudOrgID string
}

func checkEnv(envVars []string) error {
	var missingVars []string

	for _, envVar := range envVars {
		if value, exists := os.LookupEnv(envVar); !exists || value == "" {
			missingVars = append(missingVars, envVar)
		}
	}

	if len(missingVars) > 0 {
		return fmt.Errorf("error: this env vars are missing: %v", missingVars)
	}
	return nil
}

func validateEnv() error {
	err := checkEnv([]string{
		"LOG_MODE",
		"YDB_DSN",
		"SERVER_PORT",
		"INFERENCE_API_URL",
		"JWT_SECRET_KEY",
		"YWIKI_TOKEN",
		"YANDEX_CLOUD_ORG_ID",
	})
	if err != nil {
		return err
	}
	return nil
}

func getEnv(key string) string {
	result := os.Getenv(key)
	if result == "" {
		panic("Invalid env key")
	}
	return result
}

func LoadConfig() (*Config, error) {
	err := validateEnv()
	if err != nil {
		return nil, fmt.Errorf("LoadConfig: %w", err)
	}

	return &Config{
		LogMode:          getEnv("LOG_MODE"),
		ServerPort:       getEnv("SERVER_PORT"),
		YDBDSN:           getEnv("YDB_DSN"),
		InferenceAPIURL:  getEnv("INFERENCE_API_URL"),
		JWTSecretKey:     getEnv("JWT_SECRET_KEY"),
		YWikiToken:       getEnv("YWIKI_TOKEN"),
		YandexCloudOrgID: getEnv("YANDEX_CLOUD_ORG_ID"),
	}, nil
}

func LogConfig(config *Config, log logger.Logger) {
	// DO NOT include secrets in this slice, because logs can leak
	loggedFields := []string{
		"LOG_MODE",
		"YDB_DSN",
		"SERVER_PORT",
		"INFERENCE_API_URL",
		"YANDEX_CLOUD_ORG_ID",
	}
	fields := make([]any, 0, len(loggedFields)+1)
	fields = append(fields, "config loaded")
	for _, field := range loggedFields {
		fields = append(fields, zap.String(field, os.Getenv(field)))
	}
	log.Info(fields...)
}
