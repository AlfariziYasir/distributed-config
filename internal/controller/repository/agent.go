package repository

import (
	"context"
	model "distributed-configuration/pkg/models"
	"distributed-configuration/pkg/utils"
	"errors"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type AgentRepository interface {
	Create(ctx context.Context, agent *model.Agent) error
	Get(ctx context.Context, agent *model.Agent) error
	Update(ctx context.Context, agent *model.Agent) error
}

type agentRepository struct {
	db  *gorm.DB
	log *utils.Logger
}

func NewAgentRepository(db *gorm.DB, log *utils.Logger) AgentRepository {
	return &agentRepository{
		db: db, log: log,
	}
}

func (r *agentRepository) Create(ctx context.Context, agent *model.Agent) error {
	err := r.db.Create(&agent).Error
	if err != nil {
		r.log.Error("failed create new agent", zap.Error(err))
		return utils.ErrInternal
	}

	return nil
}
func (r *agentRepository) Get(ctx context.Context, agent *model.Agent) error {
	err := r.db.First(&agent).Error
	if err != nil {
		r.log.Error("failed get agent data", zap.Error(err))
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utils.ErrNotFound
		}
		return utils.ErrInternal
	}

	return nil
}
func (r *agentRepository) Update(ctx context.Context, agent *model.Agent) error {
	err := r.db.Save(agent).Error
	if err != nil {
		r.log.Error("failed update agent data", zap.Error(err))
		return utils.ErrInternal
	}

	return nil
}
