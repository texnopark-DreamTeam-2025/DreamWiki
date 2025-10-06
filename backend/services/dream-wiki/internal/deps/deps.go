package deps

import (
	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/utils/logger"
	"github.com/ydb-platform/ydb-go-sdk/v3/table"
)

type Deps struct {
	DB     table.Client
	Logger *logger.Logger
}
