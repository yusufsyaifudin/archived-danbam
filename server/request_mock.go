package server

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/stretchr/testify/mock"
)

type requestMock struct {
	mock.Mock

	startTime time.Time
	traceID   string
}

func (r requestMock) ContentType() string {
	args := r.Called()
	return args.String(0)
}

func (r requestMock) Bind(out interface{}) error {
	args := r.Called(out)

	resultByte, _ := json.Marshal(args.Get(0))
	_ = json.Unmarshal(resultByte, out)

	return args.Error(1)
}

func (r requestMock) RawRequest() *http.Request {
	args := r.Called()
	return args.Get(0).(*http.Request)
}

func (r requestMock) GetParam(param string) string {
	args := r.Called(param)
	return args.String(0)
}

func (r requestMock) GetQueryParam(param string) string {
	args := r.Called(param)
	return args.String(0)
}

func NewRequestMock() *requestMock {
	return &requestMock{}
}
