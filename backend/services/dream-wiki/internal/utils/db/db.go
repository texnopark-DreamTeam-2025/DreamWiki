package db

import (
	"context"
	"time"

	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/config"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/utils/logger"
	"github.com/ydb-platform/ydb-go-sdk/v3"
	"github.com/ydb-platform/ydb-go-sdk/v3/table"
)

// ConnectToYDB establishes a connection to YDB using the provided configuration and logger
func ConnectToYDB(config *config.Config, logger logger.Logger) (*ydb.Driver, error) {
	ydbDSN := config.YDBDSN

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	logger.Info("Connecting to YDB with DSN: ", ydbDSN)

	driver, err := ydb.Open(ctx, ydbDSN)
	if err != nil {
		logger.Error("Failed to connect to YDB: ", err)
		return nil, err
	}

	// Test the connection
	err = driver.Table().Do(ctx, func(ctx context.Context, s table.Session) error {
		return nil
	})
	if err != nil {
		logger.Error("Failed to ping YDB: ", err)
		driver.Close(context.Background())
		return nil, err
	}

	logger.Info("Successfully connected to YDB")
	return driver, nil
}
