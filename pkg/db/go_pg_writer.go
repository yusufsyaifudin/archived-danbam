package db

import (
	"context"

	"github.com/go-pg/pg/v9"
	"github.com/opentracing/opentracing-go"
)

type goPgWriter struct {
	conf Conf
	db   *pg.DB
}

func (g goPgWriter) ExecContext(ctx context.Context, query string, params ...interface{}) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "goPgWriter.ExecContext")
	defer func() {
		ctx.Done()
		span.Finish()
	}()

	_, err := g.db.ExecContext(context.WithValue(ctx, "context", ctx), query, params...)
	return err
}

func (g goPgWriter) QueryContext(ctx context.Context, model interface{}, query string, params ...interface{}) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "goPgWriter.QueryContext")
	defer func() {
		ctx.Done()
		span.Finish()
	}()

	_, err := g.db.QueryContext(context.WithValue(ctx, "context", ctx), model, query, params...)
	return err
}

func (g goPgWriter) BeginTx() (SQLTx, error) {
	tx, err := g.db.Begin()
	if err != nil {
		return nil, err
	}

	return newGoPgTx(g.conf, tx)
}

func (g goPgWriter) Prepare(query string) (SQLPreparedStatement, error) {
	stmt, err := g.db.Prepare(query)
	if err != nil {
		return nil, err
	}

	return newGoPgPreparedStmt(g.conf, stmt)
}

func newGoPgWriter(conf Conf, db *pg.DB) (SQLWriter, error) {
	return &goPgWriter{
		conf: conf,
		db:   db,
	}, nil
}
