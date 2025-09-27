package database

import (
	"context"
	"time"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Transaction interface {
	Executor

	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error

	// Begin(ctx context.Context) (pgx.Tx, error)
	// BeginFunc(ctx context.Context, f func(pgx.Tx) error) (err error)
	// CopyFrom(ctx context.Context, tableName pgx.Identifier, columnNames []string, rowSrc pgx.CopyFromSource) (int64, error)
	// SendBatch(ctx context.Context, b *pgx.Batch) pgx.BatchResults
	// LargeObjects() pgx.LargeObjects
	// Prepare(ctx context.Context, name, sql string) (*pgconn.StatementDescription, error)
	// QueryFunc(ctx context.Context, sql string, args []interface{}, scans []interface{}, f func(pgx.QueryFuncRow) error) (pgconn.CommandTag, error)
	// Conn() *pgx.Conn

	Get(ctx context.Context, returnValue interface{}, sql string, args ...interface{}) error
	GetOne(ctx context.Context, returnValue interface{}, sql string, args ...interface{}) error
	Insert(ctx context.Context, returnValue interface{}, sql string, args ...interface{}) error
	Update(ctx context.Context, returnValue interface{}, sql string, args ...interface{}) error
	Delete(ctx context.Context, returnValue interface{}, sql string, args ...interface{}) error
}

type Executor interface {
	Exec(ctx context.Context, sql string, arguments ...interface{}) (commandTag pgconn.CommandTag, err error)
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
}

type IDatabase interface {
	GetResult(tx pgx.Tx, result interface{}, q string, args ...interface{}) error
	GetTransaction() (*Tx, error)
	Get(ctx context.Context, returnValue interface{}, sql string, args ...interface{}) error
	GetOne(ctx context.Context, returnValue interface{}, sql string, args ...interface{}) error
	GetRows(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	Count(ctx context.Context, returnValue *int64, sql string, args ...interface{}) error
	Insert(ctx context.Context, returnValue interface{}, sql string, args ...interface{}) error
	Update(ctx context.Context, returnValue interface{}, sql string, args ...interface{}) error
	Delete(ctx context.Context, returnValue interface{}, sql string, args ...interface{}) error
	GetContext() context.Context
	GetContextWithCache(ctx context.Context, key string, ttl ...time.Duration) context.Context
	GetReadConn(ctx context.Context) (*pgxpool.Conn, error)
	GetWriteConn(ctx context.Context) (*pgxpool.Conn, error)
	GetWithReplace(ctx context.Context, returnValue interface{}, sql string, args ...interface{}) error
	GetOneWithReplace(ctx context.Context, returnValue interface{}, sql string, args ...interface{}) error
	InsertWithReplace(ctx context.Context, returnValue interface{}, sql string, args ...interface{}) error
	UpdateWithReplace(ctx context.Context, returnValue interface{}, sql string, args ...interface{}) error
	UpdateManyWithReplace(ctx context.Context, returnValue interface{}, sql string, args ...interface{}) error
}
