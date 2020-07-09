package db

import (
	"context"

	"github.com/opentracing/opentracing-go"

	"github.com/go-pg/pg/v9"
)

type goPgPreparedStmt struct {
	conf Conf
	db   *pg.Stmt
}

func (g goPgPreparedStmt) ExecContext(ctx context.Context, params ...interface{}) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "goPgPreparedStmt.QueryContext")
	defer func() {
		ctx.Done()
		span.Finish()
	}()

	_, err := g.db.ExecContext(ctx, params...)
	return err
}

func (g goPgPreparedStmt) QueryContext(ctx context.Context, model interface{}, params ...interface{}) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "goPgPreparedStmt.QueryContext")
	defer func() {
		ctx.Done()
		span.Finish()
	}()

	_, err := g.db.QueryContext(ctx, model, params...)
	return err
}

func (g goPgPreparedStmt) Close() error {
	return g.db.Close()
}

func newGoPgPreparedStmt(conf Conf, db *pg.Stmt) (SQLPreparedStatement, error) {
	return &goPgPreparedStmt{
		conf: conf,
		db:   db,
	}, nil
}
