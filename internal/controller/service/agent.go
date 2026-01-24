package service

import (
	"bytes"
	"context"
	"distributed-configuration/internal/controller/repository"
	model "distributed-configuration/pkg/models"
	"distributed-configuration/pkg/utils"
	"regexp"
	"strconv"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type AgentConfig interface {
	Register(ctx context.Context, req *model.AgentRequest) (string, error)
	Verify(ctx context.Context, agentID string) error
}

type agentService struct {
	log  *utils.Logger
	repo repository.AgentRepository
}

func NewAgentService(log *utils.Logger, repo repository.AgentRepository) AgentConfig {
	return &agentService{
		log:  log,
		repo: repo,
	}
}

func (s *agentService) Register(ctx context.Context, req *model.AgentRequest) (string, error) {
	agentID := uuid.New().String()

	agent := model.Agent{
		Id:        agentID,
		Name:      req.Name,
		Host:      req.Host,
		CreatedAt: time.Now(),
	}
	err := s.repo.Create(ctx, &agent)
	if err != nil {
		s.log.Error("failed create new agent", zap.Error(err))
		return "", err
	}

	return agentID, nil
}

func (s *agentService) Verify(ctx context.Context, agentID string) error {
	return  nil
}
