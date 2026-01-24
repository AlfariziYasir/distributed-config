package handler

import (
	"distributed-configuration/internal/controller/config"
	"distributed-configuration/internal/controller/service"
	model "distributed-configuration/pkg/models"
	"distributed-configuration/pkg/utils"
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
)

type handler struct {
	config service.ConfigService
	cfg    config.Config
	log    *utils.Logger
}

func NewHandler(config service.ConfigService, log *utils.Logger) *handler {
	return &handler{
		config: config,
		log:    log,
	}
}

func (h handler) Save(w http.ResponseWriter, r *http.Request) {
	var payload model.Configuration

	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		h.log.Error("invalid request body", zap.Error(err))
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	err = h.config.Save(r.Context(), &payload)
	if err != nil {
		h.log.Error("failed to save config", zap.Error(err))
		status, msg := utils.MapError(err)
		http.Error(w, msg, status)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]any{
		"status":  "success",
		"message": "configuration saved successfully",
	})
}
