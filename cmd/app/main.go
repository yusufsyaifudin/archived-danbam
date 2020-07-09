package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	"github.com/opentracing/opentracing-go"
	App "github.com/yusufsyaifudin/danbam/internal/domain/app"
	"github.com/yusufsyaifudin/danbam/pkg/tracer"
	"github.com/yusufsyaifudin/danbam/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

const (
	tracerURL   = "localhost:5775"
	serviceName = "DanBam App Service"
	port        = ":3009"
)

func main() {
	zapLogger, err := zap.NewProduction(
		zap.AddCaller(),
	)
	if err != nil {
		log.Fatal(err)
		return
	}

	tracerService, closer := tracer.New(true, serviceName, tracerURL, 1)
	defer func() {
		if closer == nil {
			_, _ = fmt.Fprintf(os.Stderr, "tracer closer is nil\n")
			return
		}

		if err := closer.Close(); err != nil {
			_, _ = fmt.Fprintf(os.Stdout, "closing tracer error: %s\n", err.Error())
			return
		}
	}()

	// set global tracer of this application
	opentracing.SetGlobalTracer(tracerService)

	appStore, err := App.DB(zapLogger)
	if err != nil {
		log.Fatal(err)
		return
	}

	defer func() {
		if err := appStore.Close(); err != nil {
			log.Fatal(err)
			return
		}
	}()

	appRpc, err := App.NewGRpc(appStore)
	if err != nil {
		log.Fatal(err)
		return
	}

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
		return
	}

	opts := []grpc.ServerOption{
		grpc.UnaryInterceptor(
			otgrpc.OpenTracingServerInterceptor(tracerService, otgrpc.LogPayloads()),
		),
		grpc.StreamInterceptor(
			otgrpc.OpenTracingStreamServerInterceptor(tracerService, otgrpc.LogPayloads()),
		),
	}

	s := grpc.NewServer(opts...)
	proto.RegisterAppServiceServer(s, appRpc)

	zapLogger.Info(fmt.Sprintf("Listening gRPC server at: %s\n", port))
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
		return
	}
}
