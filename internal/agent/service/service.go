package service

import (
	"context"
	"distributed-configuration/internal/agent/client"
	"distributed-configuration/internal/agent/config"
	"distributed-configuration/internal/agent/repository"
	model "distributed-configuration/pkg/models"
	"distributed-configuration/pkg/utils"
	"os"
	"time"

	"go.uber.org/zap"
)

type AgentService struct {
	state      *model.AgentState
	repo       *repository.FileStore
	log        *utils.Logger
	controller client.ControllerClient
	worker     client.WorkerClient
	cfg        *config.Config
}

func NewAgentService(
	controller client.ControllerClient,
	worker client.WorkerClient,
	repo *repository.FileStore,
	log *utils.Logger,
	cfg *config.Config,
) *AgentService {
	return &AgentService{
		state:      &model.AgentState{},
		repo:       repo,
		controller: controller,
		log:        log,
		worker:     worker,
		cfg:        cfg,
	}
}

func (s *AgentService) Start(ctx context.Context) {
	var state model.AgentState

	err := s.repo.Load(&state)
	if err == nil {
		s.log.Info("restore state value")
		s.state = &state

		if s.state.Config != nil {
			err = s.worker.PushConfig(ctx, s.state.Config)
			if err != nil {
				s.log.Error("failed push update to worker", zap.Error(err))
			}
		}
	}

	if s.state.AgentID == "" {
		s.register(ctx)
	}

	s.polling(ctx)
}

func (s *AgentService) register(ctx context.Context) {
	backoff := 1 * time.Second
	for s.state.AgentID == "" {
		s.log.Info("attempting to register")

		hostname, err := os.Hostname()
		if err != nil {
			hostname = "unknows"
		}

		res, err := s.controller.Register(ctx, s.cfg.AgentName, hostname)
		if err == nil {
			s.state.RegistraionData(res.AgentId, res.PollUrl, res.PollIntervalSeconds)
			s.repo.Save(s.state.Snapshot())
			s.log.Info(
				"registered agent",
				zap.String("agent_id", res.AgentId),
				zap.String("poll_url", res.PollUrl),
				zap.Int("poll_interval", res.PollIntervalSeconds),
			)
			return
		}

		s.log.Warn("failed register agent", zap.Error(err))
		time.Sleep(backoff)

		backoff *= 2
		if backoff > 30*time.Second {
			backoff = 30 * time.Second
		}
	}
}

func (s *AgentService) polling(ctx context.Context) {
	interval := time.Duration(s.state.GetInterval()) * time.Second
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	backoff := 1 * time.Second

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			agenID, etag, pollUrl := s.state.Get()
			res, err := s.controller.FetchConfig(ctx, agenID, etag, pollUrl)
			if err != nil {
				s.log.Error("poll failed", zap.Error(err), zap.Int("backoff", int(backoff)))
				time.Sleep(backoff)
				if backoff < 1*time.Minute {
					backoff *= 2
				} else {
					backoff = 1 * time.Minute
				}
				continue
			}

			if res.Data == nil {
				s.log.Info("data not modified")
				continue
			}

			s.state.UpdateConfig(res.ETag, s.state.Config)
			s.repo.Save(s.state.Snapshot())

			err = s.worker.PushConfig(ctx, s.state.Config)
			if err != nil {
				s.log.Error("failed push update to worker", zap.Error(err))
			}
		}
	}
}
