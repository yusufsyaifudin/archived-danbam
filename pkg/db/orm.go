package db

import (
	"context"
)

// SQL can be use writer or read replica only.
type SQL interface {
	Writer() SQLWriter
	Reader() SQLReader
	Close() error
}

// SQLWriter is a common interface for pg.DB and pg.Tx types.
// this copies from go-pg DB interface
type SQLWriter interface {
	QueryExecutor
	BeginTx() (SQLTx, error)
	Prepare(query string) (SQLPreparedStatement, error)
}

// SQLTx is a interface to handle go-pg orm transaction
type SQLTx interface {
	QueryExecutor
	Commit() error
	Rollback() error
}

// SQLReader use slave instance
type SQLReader interface {
	QueryExecutor
}

// SQLPreparedStatement preparing statement
type SQLPreparedStatement interface {
	ExecContext(ctx context.Context, params ...interface{}) error
	QueryContext(ctx context.Context, model interface{}, params ...interface{}) error

	// Close closing the prepared statement
	Close() error
}

// QueryExecutor is interface for doing the query
type QueryExecutor interface {
	ExecContext(ctx context.Context, query string, params ...interface{}) error
	QueryContext(ctx context.Context, model interface{}, query string, params ...interface{}) error
}
