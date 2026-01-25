package service

import (
	"context"
	"distributed-configuration/internal/controller/repository"
	model "distributed-configuration/pkg/models"
	"distributed-configuration/pkg/utils"
	"encoding/json"
	"reflect"
	"regexp"
	"strconv"
	"time"

	"go.uber.org/zap"
)

type ConfigService interface {
	Save(ctx context.Context, req *model.Configuration) error
	Get(ctx context.Context, version string) (model.Configuration, error)
}

type configService struct {
	log  *utils.Logger
	repo repository.ConfigRepository
}

func NewConfigService(log *utils.Logger, repo repository.ConfigRepository) ConfigService {
	return &configService{
		log:  log,
		repo: repo,
	}
}

func (s *configService) Save(ctx context.Context, req *model.Configuration) error {
	var (
		config           model.Configuration
		newData, oldData map[string]any
	)

	count, err := s.repo.Count(ctx, &config)
	if err != nil {
		s.log.Error("failed get latest config", zap.Error(err))
		return err
	} else if count == 0 {
		newConfig := model.Configuration{
			Version:   1,
			Data:      req.Data,
			CreatedAt: time.Now(),
		}
		err = s.repo.Create(ctx, &newConfig)
		if err != nil {
			s.log.Error("failed create new config", zap.Error(err))
			return err
		}

		return nil
	}

	s.log.Info("total data", zap.Int("count", int(count)))

	err = s.repo.Get(ctx, &config)
	if err != nil {
		s.log.Error("failed get latest config", zap.Error(err))
		return err
	}

	json.Unmarshal(req.Data, &newData)
	json.Unmarshal(config.Data, &oldData)

	ok := reflect.DeepEqual(newData, oldData)
	if ok {
		s.log.Info("data not modified")
		return utils.ErrNotModified
	}

	newConfig := model.Configuration{
		Version:   config.Version + 1,
		Data:      req.Data,
		CreatedAt: time.Now(),
	}
	err = s.repo.Create(ctx, &newConfig)
	if err != nil {
		s.log.Error("failed create new config", zap.Error(err))
		return err
	}

	return nil
}

func (s *configService) Get(ctx context.Context, version string) (model.Configuration, error) {
	var config model.Configuration

	err := s.repo.Get(ctx, &config)
	if err != nil {
		s.log.Error("failed get latest config", zap.Error(err))
		return model.Configuration{}, err
	}

	versionStr := regexp.MustCompile(`\d+`).FindString(version)
	versionInt, _ := strconv.Atoi(versionStr)
	if versionInt == config.Version {
		s.log.Info("data not modified")
		return model.Configuration{}, utils.ErrNotModified
	}

	return config, nil
}
