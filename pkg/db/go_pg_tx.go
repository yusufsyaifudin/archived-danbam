package db

import (
	"context"

	"github.com/go-pg/pg/v9"
	"github.com/opentracing/opentracing-go"
)

// goPgTx handling go pg transaction
type goPgTx struct {
	conf Conf
	tx   *pg.Tx
}

func (g goPgTx) ExecContext(ctx context.Context, query string, params ...interface{}) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "goPgTx.ExecContext")
	defer func() {
		ctx.Done()
		span.Finish()
	}()

	_, err := g.tx.ExecContext(context.WithValue(ctx, "context", ctx), query, params...)
	return err
}

func (g goPgTx) QueryContext(ctx context.Context, model interface{}, query string, params ...interface{}) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "goPgTx.QueryContext")
	defer func() {
		ctx.Done()
		span.Finish()
	}()

	_, err := g.tx.QueryContext(context.WithValue(ctx, "context", ctx), model, query, params...)
	return err
}

func (g goPgTx) Commit() error {
	return g.tx.Commit()
}

func (g goPgTx) Rollback() error {
	return g.tx.Rollback()
}

func newGoPgTx(conf Conf, tx *pg.Tx) (SQLTx, error) {
	return &goPgTx{
		conf: conf,
		tx:   tx,
	}, nil
}
