package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/yusufsyaifudin/danbam/internal/controller/app"
	App "github.com/yusufsyaifudin/danbam/internal/domain/app"
	"github.com/yusufsyaifudin/danbam/server"
	"go.uber.org/zap"
)

func main() {
	zapLogger, err := zap.NewDevelopment(
		zap.AddCaller(),
		zap.AddCallerSkip(3),
	)

	if err != nil {
		log.Fatal(err)
		return
	}

	// migration must per each scope,
	//m := []migration.Migrate{
	//	new(migrate.CreateUsersTable1592291861),
	//}

	appStore, err := App.DB()
	if err != nil {
		log.Fatal(err)
		return
	}

	defer func() {
		if appStore != nil {
			if err := appStore.Close(); err != nil {
				log.Fatal(err)
				return
			}
		}
	}()

	// ========= Start server with graceful shutdown
	srv := server.NewServer(server.Config{
		EnableProfiling: true,
		ListenAddress:   ":2222",
		WriteTimeout:    3 * time.Second,
		ReadTimeout:     3 * time.Second,
		ZapLogger:       zapLogger,
		OpenTracing:     nil,
	})

	srv.RegisterRoutes(app.Routes(appStore))

	var apiErrChan = make(chan error, 1)
	go func() {
		apiErrChan <- srv.Start()
	}()

	var signalChan = make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
	select {
	case <-signalChan:
		_, _ = fmt.Fprintf(os.Stdout, "exiting...\n")
		srv.Shutdown()

	case err := <-apiErrChan:
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "error API: %s\n", err.Error())
		}
	}
}
