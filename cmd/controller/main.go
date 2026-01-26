package main

import (
	"context"
	_ "distributed-configuration/docs/controller"
	"distributed-configuration/internal/controller/config"
	"distributed-configuration/internal/controller/handler"
	"distributed-configuration/internal/controller/repository"
	"distributed-configuration/internal/controller/service"
	model "distributed-configuration/pkg/models"
	"distributed-configuration/pkg/utils"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/redis/go-redis/v9"
	httpSwagger "github.com/swaggo/http-swagger"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// @title           Distributed Config System API
// @version         1.0
// @description     API for managing distributed configurations.

// @securityDefinitions.apiKey BearerAuth
// @in                         header
// @name                       Authorization
// @description                Type "Bearer" followed by a space and your token.

// @host            localhost:8080
// @BasePath        /

func main() {
	log, err := utils.New(zapcore.DebugLevel, "controller-service", "v1", false)
	if err != nil {
		log.Fatal(err.Error())
	}

	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal(err.Error())
		return
	}

	db, err := gorm.Open(sqlite.Open(cfg.DBDSN), &gorm.Config{
		Logger: utils.NewZapGormLogger(log.Logger, logger.Error, time.Duration(10*time.Second)),
	})
	if err != nil {
		log.Fatal(err.Error())
		return
	}

	err = db.AutoMigrate(&model.Configuration{}, &model.Agent{})
	if err != nil {
		log.Fatal(err.Error())
		return
	}

	rds := redis.NewClient(&redis.Options{Addr: cfg.RedisAddr, Password: "64mUrO4eXR3D"})

	agentRepo := repository.NewAgentRepository(db, &log)
	configRepo := repository.NewConfigRepository(db, &log)

	agentSvc := service.NewAgentService(&log, agentRepo, cfg)
	configSvc := service.NewConfigService(&log, configRepo)
	notif := service.NewRedisNotifier(rds, cfg.ChannelKey, &log)

	handler := handler.NewHandler(configSvc, agentSvc, &log, cfg, notif)

	mux := http.NewServeMux()

	mux.Handle(
		"/admin/config",
		handler.Authentication(
			handler.RoleBase(utils.RoleAdmin)(
				http.HandlerFunc(handler.Save),
			),
		),
	)
	mux.Handle(
		"/register",
		handler.Authentication(
			handler.RoleBase(utils.RoleAgent)(
				http.HandlerFunc(handler.Register),
			),
		),
	)
	mux.Handle(
		"/config",
		handler.Authentication(
			handler.RoleBase(utils.RoleAgent)(
				http.HandlerFunc(handler.Config),
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
