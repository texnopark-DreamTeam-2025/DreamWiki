package db

import (
	"context"

	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/config"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/utils/logger"
	"github.com/ydb-platform/ydb-go-sdk/v3"
	"github.com/ydb-platform/ydb-go-sdk/v3/table"
)

func CreateYDBDriver(config *config.Config) (*ydb.Driver, error) {
	ydbDSN := config.YDBDSN

	ctx := context.Background()

	driver, err := ydb.Open(ctx, ydbDSN)
	if err != nil {
		logger.Fatalf("failed to connect to YDB with DSN: %s, error: %v", ydbDSN, err)
		return nil, err
	}

	err = driver.Table().Do(ctx, func(ctx context.Context, s table.Session) error {
		return nil
	})
	if err != nil {
		logger.Fatalf("failed to ping YDB: %v", err)
		return nil, err
	}

	logger.Info("successful connection to YDB", nil)

	return driver, nil
}
