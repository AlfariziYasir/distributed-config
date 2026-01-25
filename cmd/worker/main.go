package main

import (
	"context"
	_ "distributed-configuration/docs"
	"distributed-configuration/internal/worker/config"
	"distributed-configuration/internal/worker/handler"
	"distributed-configuration/internal/worker/service"
	"distributed-configuration/pkg/utils"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	httpSwagger "github.com/swaggo/http-swagger"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// @title           Distributed Config System API
// @version         1.0
// @description     API for managing distributed configurations.

// @securityDefinitions.apiKey BearerAuth
// @in                         header
// @name                       Authorization
// @description                Type "Bearer" followed by a space and your token.

// @host            localhost:8181
// @BasePath        /
func main() {
	log, err := utils.New(zapcore.DebugLevel, "worker-service", "v1", false)
	if err != nil {
		log.Fatal(err.Error())
	}

	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal(err.Error())
		return
	}

	svc := service.NewWorkerService(&log)
	handler := handler.NewHandler(&log, cfg, svc)

	mux := http.NewServeMux()

	mux.Handle(
		"/config",
		handler.Authentication(
			handler.RoleBase(utils.RoleAgent)(
				http.HandlerFunc(handler.Save),
			),
		),
	)
	mux.Handle(
		"/hit",
		handler.Authentication(
			handler.RoleBase(utils.RoleClient)(
				http.HandlerFunc(handler.Get),
			),
		),
	)

	mux.Handle("/swagger/", httpSwagger.WrapHandler)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.HTTPPort),
		Handler: mux,
	}

	go func() {
		log.Info("http server started", zap.Int("addr", cfg.HTTPPort))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err.Error())
		}
	}()

	shutdown(server, &log)
}

func shutdown(srv *http.Server, log *utils.Logger) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	log.Info("shutting down server")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error(err.Error())
	}
}
