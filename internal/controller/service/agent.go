package service

import (
	"context"
	"distributed-configuration/internal/controller/config"
	"distributed-configuration/internal/controller/repository"
	model "distributed-configuration/pkg/models"
	"distributed-configuration/pkg/utils"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type AgentService interface {
	Register(ctx context.Context, req *model.AgentRequest) (string, error)
	Verify(ctx context.Context, agentID string) error
}

type agentService struct {
	log  *utils.Logger
	repo repository.AgentRepository
	cfg  *config.Config
}

func NewAgentService(log *utils.Logger, repo repository.AgentRepository, cfg *config.Config) AgentService {
	return &agentService{
		log:  log,
		repo: repo,
		cfg:  cfg,
	}
}

func (s *agentService) Register(ctx context.Context, req *model.AgentRequest) (string, error) {
	agentID := uuid.New().String()

	agent := model.Agent{
		Id:                  agentID,
		Name:                req.Name,
		Host:                req.Host,
		PollIntervalSeconds: int(s.cfg.PollInterval.Seconds()),
		CreatedAt:           time.Now(),
		LastSeen:            time.Now(),
	}
	err := s.repo.Create(ctx, &agent)
	if err != nil {
		s.log.Error("failed create new agent", zap.Error(err))
		return "", err
	}

	return agentID, nil
}

func (s *agentService) Verify(ctx context.Context, agentID string) error {
	agent := model.Agent{Id: agentID}
	err := s.repo.Get(ctx, &agent)
	if err != nil {
		s.log.Error("failed get agent data", zap.Error(err))
		return err
	}

	threshold := time.Duration(agent.PollIntervalSeconds*2) * time.Second
	if time.Since(agent.LastSeen) > threshold {
		return utils.ErrInActive
	}

	agent.LastSeen = time.Now()
	err = s.repo.Update(ctx, &agent)
	if err != nil {
		s.log.Error("failed get agent data", zap.Error(err))
		return err
	}

	return nil
}
