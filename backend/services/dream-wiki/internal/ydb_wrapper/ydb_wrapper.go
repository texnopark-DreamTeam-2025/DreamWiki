package ydb_wrapper

import (
	"context"
	"fmt"
	"sync/atomic"

	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/app/models"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/deps"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/utils/logger"
	"github.com/ydb-platform/ydb-go-sdk/v3"
	"github.com/ydb-platform/ydb-go-sdk/v3/query"
	"github.com/ydb-platform/ydb-go-sdk/v3/table"
	"github.com/ydb-platform/ydb-go-sdk/v3/topic"
)

type (
	YDBWrapper interface {
		Commit() error
		Rollback()
		GetTX() table.TransactionIdentifier

		InTX() Actor

		// OutsideTX returns actor that creates independent transaction for single expression.
		// May be useful for log writing.
		OutsideTX() Actor

		TopicClient() topic.Client
	}

	Actor interface {
		// Execute executes YQL statement. Must be called exactly once.
		// If no error was returned, ResultSet must be closed.
		Execute(yql string, opts ...table.ParameterOption) (ResultSet, error)
	}

	ResultSet interface {
		RowCount() int
		NextRow() bool
		FetchRow(...any) error
		FetchExactlyOne(...any) error

		Close()
	}

	ydbWrapperImpl struct {
		db  *ydb.Driver
		log logger.Logger
		ctx context.Context
		tx  *query.TxActor

		// success, commitError and closed used in transaction commit/rollback mechanics.
		success     chan bool
		commitError chan error
		closed      int32
	}

	actorImpl struct {
		ctx            context.Context
		tx             query.TxActor
		log            logger.Logger
		closingChannel chan any
	}

	resultImpl struct {
		ctx           context.Context
		log           logger.Logger
		rows          []query.Row
		rowIdx        int
		closeCallback func()
	}
)

var (
	_ YDBWrapper = &ydbWrapperImpl{}
	_ Actor      = &actorImpl{}
	_ ResultSet  = &resultImpl{}
)

func NewYDBWrapper(ctx context.Context, deps *deps.Deps, withTransaction bool) YDBWrapper {
	driver, _ := ydb.Open(ctx, deps.Config.YDBDSN)

	result := &ydbWrapperImpl{
		ctx: ctx,
		db:  driver,
		log: deps.Logger,
	}
	if withTransaction {
		result.beginTX()
	}
	return result
}

func (y *ydbWrapperImpl) beginTX() {
	y.success = make(chan bool)
	y.commitError = make(chan error)
	txRetriever := make(chan query.TxActor)

	action := func(ctx context.Context, tx query.TxActor) error {
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
		y.commitError <- y.db.Query().DoTx(y.ctx, action)
		y.db.Close(context.Background())
	}()
	tx := <-txRetriever
	y.tx = &tx
	close(txRetriever)
}

func (y *ydbWrapperImpl) TopicClient() topic.Client {
	return y.db.Topic()
}

func (y *ydbWrapperImpl) Commit() error {
	if n := atomic.AddInt32(&(y.closed), 1); n == 1 {
		y.success <- true
		return <-y.commitError
	}
	return nil
}

func (y *ydbWrapperImpl) Rollback() {
	if n := atomic.AddInt32(&(y.closed), 1); n == 1 {
		y.success <- false
		<-y.commitError
	}
}

func (y *ydbWrapperImpl) GetTX() table.TransactionIdentifier {
	return *y.tx
}

func (y *ydbWrapperImpl) InTX() Actor {
	if y.tx == nil {
		panic("YDB wrapper has no TX. Maybe you forgot set flag?")
	}
	actor := actorImpl{
		ctx:            y.ctx,
		closingChannel: make(chan any),
		tx:             *y.tx,
		log:            y.log,
	}

	go func() {
		<-actor.closingChannel
		close(actor.closingChannel)
	}()

	return &actor
}

func (y *ydbWrapperImpl) OutsideTX() Actor {
	txRetriever := make(chan query.TxActor)
	actorInstance := actorImpl{
		ctx:            y.ctx,
		closingChannel: make(chan any),
		tx:             nil,
	}

	action := func(ctx context.Context, tx query.TxActor) error {
		txRetriever <- tx
		<-actorInstance.closingChannel
		close(actorInstance.closingChannel)
		return nil
	}

	go func() {
		_ = y.db.Query().DoTx(y.ctx, action)
	}()

	tx := <-txRetriever
	close(txRetriever)
	actorInstance.tx = tx
	return &actorInstance
}

func (a *actorImpl) Execute(yql string, opts ...table.ParameterOption) (ResultSet, error) {

	paramsBuilder := ydb.ParamsBuilder()
	for _, opt := range opts {
		paramsBuilder = paramsBuilder.Param(opt.Name()).Any(opt.Value())
	}

	resultSet, err := a.tx.QueryResultSet(a.ctx, yql, query.WithParameters(paramsBuilder.Build()))
	if err != nil {
		a.log.Error(err)
		a.closingChannel <- nil
		return nil, err
	}
	defer resultSet.Close(a.ctx)

	rows := make([]query.Row, 0)
	for row, err := range resultSet.Rows(a.ctx) {
		if err != nil {
			return nil, err
		}
		rows = append(rows, row)
	}

	return &resultImpl{
		ctx:    a.ctx,
		rows:   rows,
		log:    a.log,
		rowIdx: -1,
		closeCallback: func() {
			a.closingChannel <- nil
		},
	}, nil
}

func (r *resultImpl) NextRow() bool {
	r.rowIdx++
	return r.rowIdx < int(r.RowCount())
}

func (r *resultImpl) RowCount() int {
	return len(r.rows)
}

func (r *resultImpl) FetchExactlyOne(values ...any) error {
	rowCount := r.RowCount()
	if rowCount <= 0 {
		r.log.Error("no rows")
		return models.ErrNoRows
	}
	if rowCount > 1 {
		r.log.Error(rowCount, "rows instead of 1")
		return fmt.Errorf("expected exactly one row")
	}
	r.NextRow()
	err := r.rows[r.rowIdx].Scan(values...)
	if err != nil {
		r.log.Error(err)
	}
	return err
}

func (r *resultImpl) FetchRow(values ...any) error {
	err := r.rows[r.rowIdx].Scan(values...)
	if err != nil {
		r.log.Error(err)
		return err
	}
	return nil
}

func (r *resultImpl) Close() {
	r.closeCallback()
}
