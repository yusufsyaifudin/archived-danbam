package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/uber/jaeger-client-go"
)

// wrapEcho wraps a Handler and turns it into gin compatible handler
// This method should be called with a context.Context
func wrapEcho(h Handler) echo.HandlerFunc {
	return func(eCtx echo.Context) error {
		// Get context from request context, see here:
		// https://github.com/labstack/echo-contrib/blob/v0.9.0/jaegertracing/jaegertracing.go#L159
		span, ctx := opentracing.StartSpanFromContext(eCtx.Request().Context(), "wrapEcho")
		defer func() {
			span.Finish()
			ctx.Done()
		}()

		traceID := "no-trace-id"
		if sc, ok := span.Context().(jaeger.SpanContext); ok {
			traceID = sc.String()
		}

		// create request and run the handler
		var req = newEchoRequest(eCtx)
		resp := h(opentracing.ContextWithSpan(ctx, span), req)

		// get the body first
		body, err := resp.Body()
		if err != nil {
			span.SetTag("error", fmt.Sprintf("error while writing response body %s", err.Error()))

			err = eCtx.JSON(http.StatusTeapot, ReplyStructure{
				Error: &ReplyErrorStructure{
					Code:    fmt.Sprintf("HTTP_%d", http.StatusTeapot),
					Title:   "Response writing error",
					Message: fmt.Sprintf("error while writing response body %s", err.Error()),
				},
				Type: ReplyError,
				Data: nil,
			})
			return err
		}

		if body == nil || len(body) <= 0 {
			return eCtx.JSON(http.StatusUnprocessableEntity, ReplyStructure{
				Error: &ReplyErrorStructure{
					Code:    fmt.Sprintf("HTTP_%d", http.StatusUnprocessableEntity),
					Title:   "Response body is nil",
					Message: "Sorry, we're about writing response body, but this error come to rescue.",
				},
				Type: ReplyError,
				Data: nil,
			})
		}

		span.LogFields(log.String("response", string(body)))

		// inject to response header
		_ = opentracing.GlobalTracer().Inject(
			span.Context(),
			opentracing.HTTPHeaders,
			opentracing.HTTPHeadersCarrier(resp.Header()),
		)

		// then write header
		for k, v := range resp.Header() {
			for _, h := range v {
				eCtx.Response().Writer.Header().Add(k, h)
			}
		}

		eCtx.Response().Writer.Header().Add("Content-Type", resp.ContentType())
		eCtx.Response().Writer.Header().Add("X-Trace-ID", traceID)
		eCtx.Response().Writer.Header().Add("X-Uber-ID", traceID)

		// the last is writing the body
		eCtx.Response().Writer.WriteHeader(resp.StatusCode())
		eCtx.Response().Status = resp.StatusCode() // for echoLogger to log real value, we need pass this
		eCtx.Response().Size = int64(len(body))
		_, err = eCtx.Response().Writer.Write(body)
		eCtx.Response().Committed = true // to ensure that there is no error "header already written"

		return err
	}
}

func stoppingRequest(stopped bool) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(eCtx echo.Context) error {
			span, ctx := opentracing.StartSpanFromContext(eCtx.Request().Context(), "stoppingRequest")
			defer func() {
				span.Finish()
				ctx.Done()
			}()

			eCtx.Set(startTimeKey, time.Now())

			traceID := "no-trace-id"
			if sc, ok := span.Context().(jaeger.SpanContext); ok {
				traceID = sc.String()
			}

			// check if server is shutting down
			// if it's the case then don't receive anymore requests
			if stopped {
				err := eCtx.JSON(http.StatusLocked, ReplyStructure{
					Error: &ReplyErrorStructure{
						Code:    fmt.Sprintf("HTTP_%d", http.StatusLocked),
						Title:   "Server is shutting down",
						Message: "server is on command to gracefully shutdown",
					},
					Type: ReplyError,
					Data: nil,
				})

				eCtx.Response().Writer.Header().Add("Content-Type", ContentTypeJSON)
				eCtx.Response().Writer.Header().Add("X-Trace-Id", traceID)
				eCtx.Response().Writer.Header().Add("Uber-Trace-ID", traceID)
				eCtx.Response().Committed = true
				return err
			}

			return next(eCtx)
		}
	}
}
