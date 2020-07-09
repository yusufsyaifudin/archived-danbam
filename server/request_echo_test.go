package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/smartystreets/goconvey/convey"
)

func Test_newEchoRequest(t *testing.T) {
	convey.Convey("newEchoRequest", t, func() {
		convey.Convey("Should return not nil", func() {
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			server := newEchoRequest(c)
			convey.So(server, convey.ShouldNotBeNil)
		})
	})
}

func TestEchoRequest_ContentType(t *testing.T) {
	convey.Convey("Get Content Type", t, func() {
		convey.Convey("Success", func() {
			e := echo.New()
			req := httptest.NewRequest(http.MethodPost, "/", nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			server := newEchoRequest(c)
			convey.So(server, convey.ShouldNotBeNil)

			convey.So(server.ContentType(), convey.ShouldResemble, echo.MIMEApplicationJSON)
		})
	})
}

func TestEchoRequest_Bind(t *testing.T) {
	convey.Convey("Test binding", t, func() {
		convey.Convey("Should return bind data to struct", func() {
			type data struct {
				ID int64
			}

			want := &data{
				ID: 123,
			}

			wantBytes, _ := json.Marshal(want)

			e := echo.New()
			req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(wantBytes))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			server := newEchoRequest(c)
			convey.So(server, convey.ShouldNotBeNil)

			var form = &data{}
			err := server.Bind(form)
			convey.So(err, convey.ShouldBeNil)
			convey.So(form, convey.ShouldResemble, want)
		})

		convey.Convey("Should return error because body is empty", func() {
			type data struct {
				ID int64
			}

			want := &data{
				ID: 123,
			}

			e := echo.New()
			req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"id": "abc"}`))
			req.Header.Set("X-LANGUAGE", "en")

			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			server := newEchoRequest(c)
			convey.So(server, convey.ShouldNotBeNil)

			var form = &data{}
			err := server.Bind(form)
			convey.So(err, convey.ShouldResemble, fmt.Errorf("code=415, message=Unsupported Media Type"))
			convey.So(form, convey.ShouldNotResemble, want)
		})

	})
}

func TestEchoRequest_RawRequest(t *testing.T) {
	convey.Convey("Get RawRequest", t, func() {
		convey.Convey("Success", func() {
			e := echo.New()
			req := httptest.NewRequest(http.MethodPost, "/", nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			server := newEchoRequest(c)
			convey.So(server, convey.ShouldNotBeNil)

			raw := server.RawRequest()
			convey.So(raw, convey.ShouldResemble, req)
		})
	})
}

func TestEchoRequest_GetParam(t *testing.T) {
	convey.Convey("Get Path Parameters", t, func() {
		convey.Convey("Success", func() {
			e := echo.New()

			f := make(url.Values)
			req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(f.Encode()))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)

			rec := httptest.NewRecorder()
			e.POST("/:id/:email", func(context echo.Context) error {
				return nil
			})
			c := e.NewContext(req, rec)

			c.SetParamNames("id", "email")
			c.SetParamValues("1", "jon@labstack.com")

			server := newEchoRequest(c)
			convey.So(server, convey.ShouldNotBeNil)

			id := server.GetParam("id")
			convey.So(id, convey.ShouldResemble, "1")

			email := server.GetParam("email")
			convey.So(email, convey.ShouldResemble, "jon@labstack.com")
		})
	})
}

func TestEchoRequest_GetQueryParam(t *testing.T) {
	convey.Convey("Get Query Parameters", t, func() {
		convey.Convey("Success", func() {
			e := echo.New()

			q := make(url.Values)
			q.Set("email", "jon@labstack.com")

			req := httptest.NewRequest(http.MethodPost, "/?"+q.Encode(), nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)

			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			server := newEchoRequest(c)
			convey.So(server, convey.ShouldNotBeNil)

			email := server.GetQueryParam("email")
			convey.So(email, convey.ShouldResemble, "jon@labstack.com")
		})
	})
}
