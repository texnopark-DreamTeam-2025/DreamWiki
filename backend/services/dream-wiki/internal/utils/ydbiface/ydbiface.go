package ydbiface

import (
	"context"
	"database/sql"

	"github.com/ydb-platform/ydb-go-sdk/v3/table"
	"github.com/ydb-platform/ydb-go-sdk/v3/table/result"
)

type YDBIface interface {
	Execute(ctx context.Context, query string, params *table.QueryParameters) error
	ExecuteQuery(ctx context.Context, query string, params *table.QueryParameters) (result.Result, error)
	ExecuteRowQuery(ctx context.Context, query string, params *table.QueryParameters) (*sql.Rows, error)
	BeginTx(ctx context.Context, txSettings table.TransactionSettings) (table.Transaction, error)
	Ping(ctx context.Context) error
	Close(ctx context.Context) error
}

type Transaction interface {
	Execute(ctx context.Context, query string, params *table.QueryParameters) error
	ExecuteQuery(ctx context.Context, query string, params *table.QueryParameters) (result.Result, error)
	ExecuteRowQuery(ctx context.Context, query string, params *table.QueryParameters) (*sql.Rows, error)
	CommitTx(ctx context.Context) error
	Rollback(ctx context.Context) error
}
