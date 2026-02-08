package service

import (
	"context"
	model "distributed-configuration/pkg/models"
	"distributed-configuration/pkg/utils"
	"encoding/json"

	"go.uber.org/zap"
)

type WorkerService interface {
	Save(ctx context.Context, config json.RawMessage) error
	Get(ctx context.Context) (map[string]any, error)
}

type workerService struct {
	log  *utils.Logger
	data *model.DataConfig
}

func NewWorkerService(log *utils.Logger) WorkerService {
	return &workerService{
		log:  log,
		data: &model.DataConfig{},
	}
}

func (s *workerService) Save(ctx context.Context, config json.RawMessage) error {
	if len(config) == 0 {
		s.log.Error("config content cannot be empty")
		return utils.ErrNotFound
	}

	if !json.Valid(config) {
		s.log.Error("invalid json format")
		return utils.ErrConflict
	}

	s.data.UpdateData(config)
	s.log.Info("configuration successfully updated", zap.String("data", string(config)))

	return nil
}
func (s *workerService) Get(ctx context.Context) (map[string]any, error) {
	config := s.data.GetConfig()
	if config == nil {
		s.log.Error("no configuration active")
		return nil, utils.ErrNotFound
	}

	res := map[string]any{}
	json.Unmarshal(config, &res)
	s.log.Info("get config successfully", zap.Any("data", res))

	return res, nil
}
