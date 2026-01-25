package repository

import (
	"context"
	model "distributed-configuration/pkg/models"
	"distributed-configuration/pkg/utils"
	"errors"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type ConfigRepository interface {
	Create(ctx context.Context, config *model.Configuration) error
	Get(ctx context.Context, config *model.Configuration) error
	Count(ctx context.Context, config *model.Configuration) (int64, error)
}

type configRepository struct {
	db  *gorm.DB
	log *utils.Logger
}

func NewConfigRepository(db *gorm.DB, log *utils.Logger) ConfigRepository {
	return &configRepository{
		db: db, log: log,
	}
}

func (r *configRepository) Create(ctx context.Context, config *model.Configuration) error {
	err := r.db.Create(&config).Error
	if err != nil {
		r.log.Error("failed create new config", zap.Error(err))
		return utils.ErrInternal
	}

	return nil
}
func (r *configRepository) Get(ctx context.Context, config *model.Configuration) error {
	err := r.db.Last(&config).Error
	if err != nil {
		r.log.Error("failed get config data", zap.Error(err))
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utils.ErrNotFound
		}
		return utils.ErrInternal
	}

	return nil
}

func (r *configRepository) Count(ctx context.Context, config *model.Configuration) (int64, error) {
	var count int64
	err := r.db.Model(config).Count(&count).Error
	if err != nil {
		r.log.Error("failed get config data", zap.Error(err))
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, utils.ErrNotFound
		}
		return 0, utils.ErrInternal
	}

	return count, nil
}
