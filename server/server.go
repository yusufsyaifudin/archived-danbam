package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"
	"strings"
	"time"

	"github.com/labstack/echo-contrib/jaegertracing"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/bytes"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
)

// server is an server for http
type server struct {
	e       *echo.Echo
	stopped bool
	routes  []*Route

	enableProfiling bool
	listenAddress   string
	writeTimeout    time.Duration
	readTimeout     time.Duration

	zapLogger   *zap.Logger
	openTracing opentracing.Tracer
}

func (s *server) init() {
	echo.NotFoundHandler = func(c echo.Context) error {
		// render your 404 page
		return c.JSON(http.StatusNotFound, ReplyStructure{
			Error: &ReplyErrorStructure{
				Code:    fmt.Sprintf("HTTP_%d", http.StatusNotFound),
				Title:   "Path Not Found",
				Message: fmt.Sprintf("path %s not found", c.Path()),
			},
			Type: ReplyError,
			Data: nil,
		})
	}

	// init echo server
	s.e = echo.New()
	s.e.HideBanner = true
	s.e.HidePort = true
	s.e.Use(jaegertracing.Trace(s.openTracing))
	s.e.Use(stoppingRequest(s.stopped))
	s.e.Pre(middleware.RemoveTrailingSlash())

	s.e.Use(middleware.BodyLimitWithConfig(middleware.BodyLimitConfig{
		Limit: "500KB",
		Skipper: func(eCtx echo.Context) bool {
			// Based on content length
			limit, _ := bytes.Parse("500KB")
			if eCtx.Request().ContentLength > limit {
				_ = eCtx.JSON(http.StatusRequestEntityTooLarge, ReplyStructure{
					Error: &ReplyErrorStructure{
						Code:    fmt.Sprintf("HTTP_%d", http.StatusRequestEntityTooLarge),
						Title:   "Request Entity Too Large",
						Message: "entity too large",
					},
					Type: ReplyError,
					Data: nil,
				})

				return false
			}

			return true
		},
	}))

	if s.enableProfiling {
		s.e.GET("/debug/pprof/*", echo.WrapHandler(http.DefaultServeMux))
	}

	s.e.Use(middleware.BodyDump(func(eCtx echo.Context, reqBody, resBody []byte) {
		span, ctx := opentracing.StartSpanFromContext(eCtx.Request().Context(), "BodyDump")
		defer func() {
			span.Finish()
			ctx.Done()
		}()

		// for ping path, don't do logging
		if eCtx.Path() == "/ping" {
			return
		}

		st := eCtx.Get(startTimeKey)
		startTime := time.Now()
		if v, ok := st.(time.Time); ok {
			startTime = v
		}

		latency := float64(time.Now().Sub(startTime).Nanoseconds()) / float64(time.Millisecond)

		var reqBodyObj interface{}
		_ = json.Unmarshal(reqBody, &reqBodyObj)

		var respBodyObj interface{}
		_ = json.Unmarshal(resBody, &respBodyObj)

		s.zapLogger.Info(
			"requested",
			zap.String("method", eCtx.Request().Method),
			zap.String("path", eCtx.Path()),
			zap.Float64("latency", latency),
			zap.Any("req_body", reqBodyObj),
			zap.Int("resp_status", eCtx.Response().Status),
			zap.Any("resp_body", respBodyObj),
		)
	}))
}

// RegisterRoutes will register all routes
func (s *server) RegisterRoutes(routes []*Route) {
	if s.e == nil {
		return
	}

	for _, r := range routes {
		if r == nil {
			continue
		}

		s.routes = append(s.routes, r)
	}
}

// Start will start the server
func (s *server) Start() error {
	for _, r := range s.routes {

		// register middleware then define routes
		m := ChainMiddleware(r.Middleware...)

		switch strings.ToUpper(strings.TrimSpace(r.Method)) {
		case http.MethodConnect:
			s.e.Add(http.MethodConnect, r.Path, wrapEcho(m(r.Handler)))

		case http.MethodDelete:
			s.e.Add(http.MethodDelete, r.Path, wrapEcho(m(r.Handler)))

		case http.MethodGet:

			s.e.GET(r.Path, wrapEcho(m(r.Handler)))

		case http.MethodHead:
			s.e.Add(http.MethodHead, r.Path, wrapEcho(m(r.Handler)))

		case http.MethodOptions:
			s.e.Add(http.MethodOptions, r.Path, wrapEcho(m(r.Handler)))

		case http.MethodPatch:
			s.e.Add(http.MethodPatch, r.Path, wrapEcho(m(r.Handler)))

		case http.MethodPost:
			s.e.Add(http.MethodPost, r.Path, wrapEcho(m(r.Handler)))

		case echo.PROPFIND:
			s.e.Add(echo.PROPFIND, r.Path, wrapEcho(m(r.Handler)))

		case http.MethodPut:
			s.e.Add(http.MethodPut, r.Path, wrapEcho(m(r.Handler)))

		case http.MethodTrace:
			s.e.Add(http.MethodTrace, r.Path, wrapEcho(m(r.Handler)))

		case "ANY":
			s.e.Any(r.Path, wrapEcho(m(r.Handler)))
		}
	}

	for _, route := range s.e.Routes() {
		s.zapLogger.Info(route.Path, zap.String("method", route.Method))
	}

	_, _ = fmt.Fprintf(os.Stdout, "Starting server at %s\n", s.listenAddress)
	return s.e.StartServer(&http.Server{
		Addr:         s.listenAddress,
		ReadTimeout:  s.readTimeout,
		WriteTimeout: s.writeTimeout,
	})
}

// Shutdown will inform the server to gracefully shutdown.
func (s *server) Shutdown() {
	s.stopped = true
}

// NewServer return new server instance
func NewServer(conf Config) *server {
	tracer := conf.OpenTracing
	if tracer == nil {
		tracer = new(opentracing.NoopTracer)
	}

	s := &server{
		e:       echo.New(),
		stopped: false,
		routes:  make([]*Route, 0),

		enableProfiling: conf.EnableProfiling,
		listenAddress:   conf.ListenAddress,
		writeTimeout:    conf.WriteTimeout,
		readTimeout:     conf.ReadTimeout,

		zapLogger:   conf.ZapLogger,
		openTracing: tracer,
	}

	s.init()
	return s
}
