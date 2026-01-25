package repository

import (
	model "distributed-configuration/pkg/models"
	"distributed-configuration/pkg/utils"
	"encoding/json"
	"os"

	"go.uber.org/zap"
)

const Filename = "config.json"

type FileStore struct {
	log *utils.Logger
}

func NewFileStore(log *utils.Logger) *FileStore {
	return &FileStore{log: log}
}

func (r *FileStore) Save(state *model.AgentState) error {
	data, _ := json.MarshalIndent(state, "", " ")
	return os.WriteFile(Filename, data, 0644)
}

func (r *FileStore) Load(state *model.AgentState) error {
	data, err := os.ReadFile(Filename)
	if err != nil {
		r.log.Error("failed read file store", zap.Error(err))
		return err
	}

	err = json.Unmarshal(data, &state)
	if err != nil {
		r.log.Error("failed parse json encoded", zap.Error(err))
		return err
	}

	return nil
}
