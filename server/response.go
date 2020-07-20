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

// mockResponse for testing purpose
type mockResponse struct {
	statusCode int
	header     http.Header
	body       []byte
	err        error
}

// StatusCode return HTTP status code
func (m mockResponse) StatusCode() int {
	return m.statusCode
}

// Body always return byte body without error
func (m mockResponse) Body() ([]byte, error) {
	return m.body, m.err
}

// Header always return http.Header
func (m mockResponse) Header() http.Header {
	return m.header
}

// ContentType always return application/json
func (m mockResponse) ContentType() string {
	return ContentTypeJSON
}

// mockResp response mocking
func mockResp(statusCode int, header http.Header, body []byte, err error) Response {
	return &mockResponse{
		statusCode: statusCode,
		header:     header,
		body:       body,
		err:        err,
	}
}
