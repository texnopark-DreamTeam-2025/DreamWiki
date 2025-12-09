package db_adapter

import (
	"context"
	"fmt"
	"sync/atomic"

	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/app/models"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/config"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/utils/logger"
	"github.com/ydb-platform/ydb-go-sdk/v3"
	"github.com/ydb-platform/ydb-go-sdk/v3/query"
	"github.com/ydb-platform/ydb-go-sdk/v3/table"
	"github.com/ydb-platform/ydb-go-sdk/v3/topic"
	"github.com/ydb-platform/ydb-go-sdk/v3/topic/topicoptions"
)

type TransactionMode = int

const (
	SnapshotReadOnly TransactionMode = iota
	SerializableReadWrite
)

type (
	DBAdapter interface {
		NewTransaction(ctx context.Context, mode TransactionMode) Transaction
		NewTopicReader(topicName string, messageCallback TopicReaderCallback, options ...topicoptions.TopicOption) TopicReader
		Close()
	}

	Transaction interface {
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

	TopicReader interface {
		ReadMessages(ctx context.Context) error
	}

	TopicReaderCallback = func(message []byte)
)

type (
	dbAdapterImpl struct {
		config *config.Config
		log    logger.Logger
	}

	transactionImpl struct {
		log logger.Logger
		ctx context.Context
		tx  *query.TxActor
		db  *ydb.Driver

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

	topicReaderImpl struct {
		db              *ydb.Driver
		topicName       string
		messageCallback TopicReaderCallback
		options         []topicoptions.TopicOption
		log             logger.Logger
	}
)

var (
	_ DBAdapter   = &dbAdapterImpl{}
	_ TopicReader = &topicReaderImpl{}
	_ Transaction = &transactionImpl{}
	_ Actor       = &actorImpl{}
	_ ResultSet   = &resultImpl{}
)

func NewDBAdapter(config *config.Config, log logger.Logger) DBAdapter {
	return &dbAdapterImpl{
		config: config,
		log:    log,
	}
}

func (d *dbAdapterImpl) NewTransaction(ctx context.Context, mode TransactionMode) Transaction {
	driver, err := ydb.Open(ctx, d.config.YDBDSN)
	if err != nil {
		panic(err.Error())
	}
	transactionWrapper := &transactionImpl{
		db:          driver,
		ctx:         ctx,
		log:         d.log,
		success:     make(chan bool),
		commitError: make(chan error),
	}
	txRetrievingChannel := make(chan query.TxActor)

	action := func(ctx context.Context, tx query.TxActor) error {
		txRetrievingChannel <- tx
		d.log.Infof("YDB transaction %s started", tx.ID())

		shouldCommit := <-transactionWrapper.success
		close(transactionWrapper.success)
		if shouldCommit {
			d.log.Infof("YDB transaction %s committed", tx.ID())
			return nil
		}
		d.log.Infof("YDB transaction %s rolled back", tx.ID())
		return fmt.Errorf("transaction rolled back")
	}

	go func() {
		transactionWrapper.commitError <- driver.Query().DoTx(ctx, action)
	}()
	tx := <-txRetrievingChannel
	transactionWrapper.tx = &tx
	close(txRetrievingChannel)

	return transactionWrapper
}

func (d *dbAdapterImpl) Close() {
}

func (d *dbAdapterImpl) NewTopicReader(topicName string, messageCallback TopicReaderCallback, options ...topicoptions.TopicOption) TopicReader {
	driver, err := ydb.Open(context.Background(), d.config.YDBDSN)
	if err != nil {
		panic(err.Error())
	}
	return &topicReaderImpl{
		db:              driver,
		topicName:       topicName,
		messageCallback: messageCallback,
		options:         options,
		log:             d.log,
	}
}

func (y *transactionImpl) TopicClient() topic.Client {
	return y.db.Topic()
}

func (t *topicReaderImpl) ReadMessages(ctx context.Context) error {
	topicClient := t.db.Topic()
	reader, err := topicClient.StartReader("dream_wiki", topicoptions.ReadTopic(t.topicName))
	if err != nil {
		return err
	}
	defer func() {
		err := reader.Close(context.Background())
		if err != nil {
			t.log.Error("failed to close reader", "error", err)
		} else {
			t.log.Info("reader closed")
		}
	}()

	for {
		mess, err := reader.ReadMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return ctx.Err()
			}
			t.log.Error("failed to read message", err)
			continue
		}

		// Commit the message immediately to avoid reprocessing
		commitError := reader.Commit(mess.Context(), mess)
		if commitError != nil {
			t.log.Error("failed to commit message", commitError)
			continue
		}

		messageContent := make([]byte, 1024)
		n, readError := mess.Read(messageContent)
		if readError != nil {
			t.log.Error("failed to get message data", "error", readError)
			continue
		}
		if n == 0 {
			t.log.Error("zero message length")
			continue
		}

		t.log.Debugf("Got message on queue %s: %s", t.topicName, string(messageContent))
		messageContent = messageContent[:n]

		t.messageCallback(messageContent)
	}
}

func (y *transactionImpl) Commit() error {
	if n := atomic.AddInt32(&(y.closed), 1); n == 1 {
		y.success <- true
		return <-y.commitError
	}
	return nil
}

func (y *transactionImpl) Rollback() {
	if n := atomic.AddInt32(&(y.closed), 1); n == 1 {
		y.success <- false
		<-y.commitError
	}
}

func (y *transactionImpl) GetTX() table.TransactionIdentifier {
	return *y.tx
}

func (y *transactionImpl) InTX() Actor {
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

func (y *transactionImpl) OutsideTX() Actor {
	txRetriever := make(chan query.TxActor)
	actorInstance := actorImpl{
		ctx:            y.ctx,
		closingChannel: make(chan any),
		tx:             nil,
	}

	action := func(ctx context.Context, tx query.Session) error {
		return nil
	}

	go func() {
		_ = y.db.Query().Do(y.ctx, action)
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
