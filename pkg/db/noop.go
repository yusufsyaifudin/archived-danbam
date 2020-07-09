package db

import (
	"context"
)

type noopTx struct{}

func (n noopTx) ExecContext(_ context.Context, _ string, _ ...interface{}) error { return nil }

func (n noopTx) QueryContext(_ context.Context, _ interface{}, _ string, _ ...interface{}) error {
	return nil
}

func (n noopTx) Commit() error { return nil }

func (n noopTx) Rollback() error { return nil }

type noopStmt struct{}

func (n noopStmt) ExecContext(_ context.Context, _ ...interface{}) error { return nil }

func (n noopStmt) QueryContext(_ context.Context, _ interface{}, _ ...interface{}) error { return nil }

func (n noopStmt) Close() error { return nil }

type noopWriter struct{}

func (n noopWriter) ExecContext(_ context.Context, _ string, _ ...interface{}) error { return nil }

func (n noopWriter) QueryContext(_ context.Context, _ interface{}, _ string, _ ...interface{}) error {
	return nil
}

func (n noopWriter) BeginTx() (SQLTx, error) {
	return &noopTx{}, nil
}

func (n noopWriter) Prepare(query string) (SQLPreparedStatement, error) {
	return &noopStmt{}, nil
}

func (n noopSql) Writer() SQLWriter {
	return &noopWriter{}
}

func (n noopSql) Reader() SQLReader {
	return &noopWriter{}
}

type noopSql struct{}

func (n noopSql) Close() error { return nil }

func NewNoop() SQL {
	return &noopSql{}
}
