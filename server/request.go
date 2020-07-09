package server

import (
	"net/http"
)

// Request represents an api request
type Request interface {
	ContentType() string
	Bind(out interface{}) error
	RawRequest() *http.Request
	GetParam(string) string
	GetQueryParam(string) string
}
