package server

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/smartystreets/goconvey/convey"
)

func Test_wrapEcho(t *testing.T) {
	convey.Convey("Test wrapEcho", t, func() {
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(""))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		convey.Convey("Success with no error (body is not empty)", func() {
			var userHandler = func(ctx context.Context, req Request) Response {
				return mockResp(http.StatusOK, http.Header{}, []byte(`{}`), nil)
			}

			h := wrapEcho(userHandler)
			err := h(c)
			convey.So(err, convey.ShouldBeNil)
		})

		convey.Convey("Return success when body error", func() {
			var userHandler = func(ctx context.Context, req Request) Response {
				return mockResp(http.StatusOK, http.Header{}, []byte(`{}`), fmt.Errorf("error"))
			}

			h := wrapEcho(userHandler)
			err := h(c)
			convey.So(err, convey.ShouldBeNil)
		})

		convey.Convey("Return success when body is nil", func() {
			var userHandler = func(ctx context.Context, req Request) Response {
				return mockResp(http.StatusOK, http.Header{}, nil, nil)
			}

			h := wrapEcho(userHandler)
			err := h(c)
			convey.So(err, convey.ShouldBeNil)
		})

		convey.Convey("Return success (using header)", func() {
			var userHandler = func(ctx context.Context, req Request) Response {
				return mockResp(http.StatusOK, http.Header{
					"Additional-Header": []string{"A", "B"},
				}, []byte(`{}`), nil)
			}

			h := wrapEcho(userHandler)
			err := h(c)
			convey.So(err, convey.ShouldBeNil)
		})
	})
}
