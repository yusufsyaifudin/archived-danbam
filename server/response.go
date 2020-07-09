package server

import (
	"net/http"
)

// Response represents an api response
type Response interface {
	StatusCode() int
	Body() ([]byte, error)
	Header() http.Header
	ContentType() string
}
