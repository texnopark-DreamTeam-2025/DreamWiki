package config

import (
	"fmt"
	"os"
)

type Config struct {
	LogMode         string
	ServerPort      string
	YDBDSN          string
	InferenceAPIURL string
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
	})
	if err != nil {
		return err
	}
	return nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func LoadConfig() (*Config, error) {
	err := validateEnv()
	if err != nil {
		return nil, fmt.Errorf("LoadConfig: %w", err)
	}

	return &Config{
		LogMode:         getEnv("LOG_MODE", "dev"),
		ServerPort:      getEnv("SERVER_PORT", "8080"),
		YDBDSN:          getEnv("YDB_DSN", "grpc://localhost:2136/?database=/local"),
		InferenceAPIURL: getEnv("INFERENCE_API_URL", "http://localhost:8000"),
	}, nil
}
