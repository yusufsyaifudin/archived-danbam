package server

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

type echoRequest struct {
	echoCtx echo.Context
}

func newEchoRequest(echoCtx echo.Context) Request {
	return &echoRequest{
		echoCtx: echoCtx,
	}
}

func (r *echoRequest) ContentType() string {
	return r.echoCtx.Request().Header.Get("Content-Type")
}

func (r *echoRequest) Bind(out interface{}) error {
	err := r.echoCtx.Bind(out)
	if err == nil {
		return nil
	}

	return fmt.Errorf(err.Error())
}

func (r *echoRequest) RawRequest() *http.Request {
	return r.echoCtx.Request()
}

func (r *echoRequest) GetParam(key string) string {
	return r.echoCtx.Param(key)
}

func (r *echoRequest) GetQueryParam(key string) string {
	return r.echoCtx.QueryParam(key)
}
