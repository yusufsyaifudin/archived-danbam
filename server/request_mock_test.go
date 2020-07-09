package server

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"

	"github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/mock"
)

func TestNewRequestMock(t *testing.T) {
	convey.Convey("Request Mock", t, func() {
		convey.Convey("Should not be nil", func() {
			req := NewRequestMock()
			convey.So(req, convey.ShouldNotBeNil)
		})
	})
}

func TestRequestMock_ContentType(t *testing.T) {
	convey.Convey("Request Mock", t, func() {
		convey.Convey("Get content type", func() {
			req := NewRequestMock()
			convey.So(req, convey.ShouldNotBeNil)

			req.On("ContentType").Return(ContentTypeJSON).Once()
			convey.So(req.ContentType(), convey.ShouldResemble, ContentTypeJSON)
		})
	})
}

func TestRequestMock_Bind(t *testing.T) {
	convey.Convey("Request Mock", t, func() {
		convey.Convey("Get bind data", func() {
			req := NewRequestMock()
			convey.So(req, convey.ShouldNotBeNil)

			type data struct {
				Foo string
				Bar int64
			}

			want := &data{
				Foo: "string",
				Bar: 1,
			}

			req.On("Bind", mock.Anything).Return(want, nil).Once()

			var actual = &data{}
			err := req.Bind(actual)
			convey.So(actual, convey.ShouldResemble, want)
			convey.So(err, convey.ShouldBeNil)
		})

		convey.Convey("On error", func() {
			req := NewRequestMock()
			convey.So(req, convey.ShouldNotBeNil)

			type data struct {
				Foo string
				Bar int64
			}

			want := &data{}
			wantErr := fmt.Errorf("error")

			req.On("Bind", mock.Anything).Return(want, wantErr).Once()

			var actual = &data{}
			err := req.Bind(actual)
			convey.So(actual, convey.ShouldResemble, want)
			convey.So(err, convey.ShouldResemble, wantErr)
		})
	})
}

func TestRequestMock_RawRequest(t *testing.T) {
	convey.Convey("Request Mock", t, func() {
		convey.Convey("Get raw request", func() {
			req := NewRequestMock()
			convey.So(req, convey.ShouldNotBeNil)

			want := &http.Request{
				Method: "",
				URL: &url.URL{
					Scheme:     "https",
					Opaque:     "",
					User:       &url.Userinfo{},
					Host:       "example.com",
					Path:       "/path",
					RawPath:    "",
					ForceQuery: false,
					RawQuery:   "",
					Fragment:   "",
				},
				Proto:      "",
				ProtoMajor: 0,
				ProtoMinor: 0,
				Header:     http.Header{},
				Body:       ioutil.NopCloser(bytes.NewReader([]byte("hello world"))),
			}

			req.On("RawRequest").Return(want).Once()

			raw := req.RawRequest()
			convey.So(raw, convey.ShouldResemble, want)
		})

	})
}

func TestRequestMock_GetParam(t *testing.T) {
	convey.Convey("Get Param test", t, func() {
		convey.Convey("Success", func() {
			req := NewRequestMock()
			convey.So(req, convey.ShouldNotBeNil)

			arg := "foo"
			ret := "bar"
			req.On("GetParam", arg).Return(ret).Once()

			convey.So(req.GetParam(arg), convey.ShouldResemble, ret)
		})
	})
}

func TestRequestMock_GetQueryParam(t *testing.T) {
	convey.Convey("Get QueryParam test", t, func() {
		convey.Convey("Success", func() {
			req := NewRequestMock()
			convey.So(req, convey.ShouldNotBeNil)

			arg := "foo"
			ret := "bar"
			req.On("GetQueryParam", arg).Return(ret).Once()

			convey.So(req.GetQueryParam(arg), convey.ShouldResemble, ret)
		})
	})
}
