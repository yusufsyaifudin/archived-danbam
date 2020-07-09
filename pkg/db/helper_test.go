package db

import (
	"context"
	"encoding/json"
	"io"

	"github.com/go-pg/pg/v9/orm"
	"github.com/stretchr/testify/mock"
)

type goPGORMNoop struct {
	mock.Mock
}

type formatter struct {
}

func (f formatter) FormatQuery(b []byte, query string, params ...interface{}) []byte {
	formatter := new(orm.Formatter)
	got := formatter.FormatQuery(b, query, params...)
	return got
}

func (g goPGORMNoop) Model(model ...interface{}) *orm.Query {
	return nil
}

func (g goPGORMNoop) ModelContext(c context.Context, model ...interface{}) *orm.Query {
	return nil
}

func (g goPGORMNoop) Select(model interface{}) error {
	args := g.Called(model)

	resultByte, _ := json.Marshal(args.Get(0))
	_ = json.Unmarshal(resultByte, model)
	return args.Error(1)
}

func (g goPGORMNoop) Insert(model ...interface{}) error {
	args := g.Called(model...)
	return args.Error(0)
}

func (g goPGORMNoop) Update(model interface{}) error {
	args := g.Called(model)

	resultByte, _ := json.Marshal(args.Get(0))
	_ = json.Unmarshal(resultByte, model)
	return args.Error(1)
}

func (g goPGORMNoop) Delete(model interface{}) error {
	args := g.Called(model)

	resultByte, _ := json.Marshal(args.Get(0))
	_ = json.Unmarshal(resultByte, model)
	return args.Error(1)
}

func (g goPGORMNoop) ForceDelete(model interface{}) error {
	args := g.Called(model)

	resultByte, _ := json.Marshal(args.Get(0))
	_ = json.Unmarshal(resultByte, model)
	return args.Error(1)
}

func (g goPGORMNoop) Exec(query interface{}, params ...interface{}) (orm.Result, error) {
	var called = make([]interface{}, 0)
	called = append(called, query)
	for _, p := range params {
		called = append(called, p)
	}

	args := g.Called(called...)
	return nil, args.Error(1)
}

func (g goPGORMNoop) ExecContext(c context.Context, query interface{}, params ...interface{}) (orm.Result, error) {
	var called = make([]interface{}, 0)
	called = append(called, c)
	called = append(called, query)
	for _, p := range params {
		called = append(called, p)
	}

	args := g.Called(called...)
	return nil, args.Error(1)
}

func (g goPGORMNoop) ExecOne(query interface{}, params ...interface{}) (orm.Result, error) {
	var called = make([]interface{}, 0)
	called = append(called, query)
	for _, p := range params {
		called = append(called, p)
	}

	args := g.Called(called...)
	return nil, args.Error(1)
}

func (g goPGORMNoop) ExecOneContext(c context.Context, query interface{}, params ...interface{}) (orm.Result, error) {
	args := g.Called(c, query, params)
	return nil, args.Error(1)
}

func (g goPGORMNoop) Query(model, query interface{}, params ...interface{}) (orm.Result, error) {
	var called = make([]interface{}, 0)
	called = append(called, model)
	called = append(called, query)
	for _, p := range params {
		called = append(called, p)
	}

	args := g.Called(called...)

	resultByte, _ := json.Marshal(args.Get(0))
	_ = json.Unmarshal(resultByte, model)

	return nil, args.Error(1)
}

func (g goPGORMNoop) QueryContext(c context.Context, model, query interface{}, params ...interface{}) (orm.Result, error) {
	var called = make([]interface{}, 0)
	called = append(called, c)
	called = append(called, model)
	called = append(called, query)
	for _, p := range params {
		called = append(called, p)
	}

	args := g.Called(called...)

	resultByte, _ := json.Marshal(args.Get(0))
	_ = json.Unmarshal(resultByte, model)

	return nil, args.Error(1)
}

func (g goPGORMNoop) QueryOne(model, query interface{}, params ...interface{}) (orm.Result, error) {
	var called = make([]interface{}, 0)
	called = append(called, model)
	called = append(called, query)
	for _, p := range params {
		called = append(called, p)
	}

	args := g.Called(called...)

	resultByte, _ := json.Marshal(args.Get(0))
	_ = json.Unmarshal(resultByte, model)

	return nil, args.Error(1)
}

func (g goPGORMNoop) QueryOneContext(c context.Context, model, query interface{}, params ...interface{}) (orm.Result, error) {
	var called = make([]interface{}, 0)
	called = append(called, c)
	called = append(called, model)
	called = append(called, query)
	for _, p := range params {
		called = append(called, p)
	}

	args := g.Called(called...)

	resultByte, _ := json.Marshal(args.Get(0))
	_ = json.Unmarshal(resultByte, model)

	return nil, args.Error(1)
}

func (g goPGORMNoop) CopyFrom(r io.Reader, query interface{}, params ...interface{}) (orm.Result, error) {
	args := g.Called(r, query, params)
	return nil, args.Error(1)
}

func (g goPGORMNoop) CopyTo(w io.Writer, query interface{}, params ...interface{}) (orm.Result, error) {
	args := g.Called(w, query, params)
	return nil, args.Error(1)
}

func (g goPGORMNoop) Context() context.Context {
	return context.Background()
}

func (g goPGORMNoop) Formatter() orm.QueryFormatter {
	return &formatter{}
}

func newGoPGDBTest() *goPGORMNoop {
	return &goPGORMNoop{}
}
