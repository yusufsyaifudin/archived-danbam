package server

import (
	"time"

	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
)

type Config struct {
	EnableProfiling bool
	ListenAddress   string
	WriteTimeout    time.Duration
	ReadTimeout     time.Duration

	ZapLogger   *zap.Logger
	OpenTracing opentracing.Tracer
}
