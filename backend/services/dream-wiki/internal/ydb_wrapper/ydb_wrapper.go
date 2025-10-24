package ydb_wrapper

import (
	"context"
	"fmt"
	"sync/atomic"

	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/app/models"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/deps"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/utils/logger"
	"github.com/ydb-platform/ydb-go-sdk/v3/table"
	"github.com/ydb-platform/ydb-go-sdk/v3/table/result"
	"github.com/ydb-platform/ydb-go-sdk/v3/table/result/indexed"
)

type (
	YDBWrapper interface {
		Commit()
		Rollback()

		InTX() Actor

		// OutsideTX returns actor that creates independent transaction for single expression.
		// May be useful for log writing.
		OutsideTX() Actor
	}

	Actor interface {
		// Execute executes YQL statement. Must be called exactly once.
		// If no error was returned, ResultSet must be closed.
		Execute(yql string, opts ...table.ParameterOption) (ResultSet, error)
	}

	ResultSet interface {
		RowCount() int
		NextRow() bool
		FetchRow(...indexed.RequiredOrOptional) error
		FetchExactlyOne(...indexed.RequiredOrOptional) error

		Close()
	}

	ydbWrapperImpl struct {
		log         logger.Logger
		ctx         context.Context
		tableClient table.Client
		tx          *table.TransactionActor

		// success and closed used in transaction commit/rollback mechanics.
		success chan bool
		closed  int32
	}

	actorImpl struct {
		ctx            context.Context
		tx             table.TransactionActor
		closingChannel chan any
	}

	resultImpl struct {
		ctx           context.Context
		result        result.Result
		closeCallback func()
	}
)

var (
	_ YDBWrapper = &ydbWrapperImpl{}
	_ Actor      = &actorImpl{}
	_ ResultSet  = &resultImpl{}
)

func NewYDBWrapper(ctx context.Context, deps *deps.Deps, withTransaction bool) YDBWrapper {
	result := &ydbWrapperImpl{
		ctx:         ctx,
		log:         deps.Logger,
		tableClient: deps.YDBDriver.Table(),
	}
	if withTransaction {
		result.beginTX()
	}
	return result
}

func (y *ydbWrapperImpl) beginTX() {
	y.success = make(chan bool)
	txRetriever := make(chan table.TransactionActor)

	action := func(ctx context.Context, tx table.TransactionActor) error {
		txRetriever <- tx
		y.log.Infof("YDB transaction %s started", tx.ID())

		shouldCommit := <-y.success
		close(y.success)
		if shouldCommit {
			y.log.Infof("YDB transaction %s committed", tx.ID())
			return nil
		}
		y.log.Infof("YDB transaction %s rolled back", tx.ID())
		return fmt.Errorf("transaction rolled back")
	}

	go func() {
		_ = y.tableClient.DoTx(y.ctx, action)
	}()
	tx := <-txRetriever
	y.tx = &tx
	close(txRetriever)
}

func (y *ydbWrapperImpl) Commit() {
	if n := atomic.AddInt32(&(y.closed), 1); n == 1 {
		y.success <- true
	}
}

func (y *ydbWrapperImpl) InTX() Actor {
	if y.tx == nil {
		panic("YDB wrapper has no TX. Maybe you forgot set flag?")
	}
	actor := actorImpl{
		ctx:            y.ctx,
		closingChannel: make(chan any),
		tx:             *y.tx,
	}

	go func() {
		<-actor.closingChannel
		close(actor.closingChannel)
	}()

	return &actor
}

func (y *ydbWrapperImpl) Rollback() {
	if n := atomic.AddInt32(&(y.closed), 1); n == 1 {
		y.success <- false
	}
}

func (y *ydbWrapperImpl) OutsideTX() Actor {
	txRetriever := make(chan table.TransactionActor)
	actorInstance := actorImpl{
		ctx:            y.ctx,
		closingChannel: make(chan any),
		tx:             nil,
	}

	action := func(ctx context.Context, tx table.TransactionActor) error {
		txRetriever <- tx
		<-actorInstance.closingChannel
		close(actorInstance.closingChannel)
		return nil
	}

	go func() {
		_ = y.tableClient.DoTx(y.ctx, action)
	}()

	tx := <-txRetriever
	close(txRetriever)
	actorInstance.tx = tx
	return &actorInstance
}

func (a *actorImpl) Execute(yql string, opts ...table.ParameterOption) (ResultSet, error) {
	result, err := a.tx.Execute(a.ctx, yql, table.NewQueryParameters(opts...))
	if err != nil {
		a.closingChannel <- nil
		return nil, err
	}
	if result.Err() != nil {
		a.closingChannel <- nil
		return nil, result.Err()
	}
	if result.ResultSetCount() < 1 {
		a.closingChannel <- nil
		return nil, fmt.Errorf("expected exactly one result set")
	}
	result.NextResultSet(a.ctx)

	return &resultImpl{
		ctx:    a.ctx,
		result: result,
		closeCallback: func() {
			a.closingChannel <- nil
		},
	}, nil
}

func (r *resultImpl) NextRow() bool {
	return r.result.NextRow()
}

func (r *resultImpl) RowCount() int {
	return r.result.CurrentResultSet().RowCount()
}

func (r *resultImpl) FetchExactlyOne(values ...indexed.RequiredOrOptional) error {
	rowCount := r.RowCount()
	if rowCount == 0 {
		return models.ErrNoRows
	}
	if rowCount > 1 {
		return fmt.Errorf("expected exactly one row")
	}
	r.result.NextRow()
	err := r.result.Scan(values...)
	return err
}

func (r *resultImpl) FetchRow(values ...indexed.RequiredOrOptional) error {
	return r.result.Scan(values...)
}

func (r *resultImpl) Close() {
	_ = r.result.Close()
	r.closeCallback()
}
