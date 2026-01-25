package main

import (
	"context"
	"distributed-configuration/internal/agent/client"
	"distributed-configuration/internal/agent/config"
	"distributed-configuration/internal/agent/repository"
	"distributed-configuration/internal/agent/service"
	"distributed-configuration/pkg/utils"

	"go.uber.org/zap/zapcore"
)

func main() {
	log, err := utils.New(zapcore.DebugLevel, "agent-service", "v1", false)
	if err != nil {
		log.Fatal(err.Error())
	}

	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal(err.Error())
		return
	}

	repo := repository.NewFileStore(&log)
	controller := client.NewControllerClient(&log, cfg)
	worker := client.NewWorkerClient(&log, cfg)

	service := service.NewAgentService(controller, worker, repo, &log, cfg)

	ctx := context.Background()
	service.Start(ctx)
}
