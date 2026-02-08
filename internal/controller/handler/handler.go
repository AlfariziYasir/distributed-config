package handler

import (
	"context"
	"distributed-configuration/internal/controller/config"
	"distributed-configuration/internal/controller/service"
	model "distributed-configuration/pkg/models"
	"distributed-configuration/pkg/utils"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"go.uber.org/zap"
)

type handler struct {
	config service.ConfigService
	agent  service.AgentService
	cfg    *config.Config
	log    *utils.Logger
	notif  *service.RedisNotifier
}

func NewHandler(
	config service.ConfigService,
	agent service.AgentService,
	log *utils.Logger,
	cfg *config.Config,
	notif *service.RedisNotifier,
) *handler {
	return &handler{
		config: config,
		agent:  agent,
		log:    log,
		cfg:    cfg,
		notif:  notif,
	}
}

// UpdateConfig godoc
// @Summary      Update global configuration
// @Description  Admin endpoint to update the configuration that will be pushed to all agents
// @Tags         admin
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        config  body      model.Configuration  true  "New Configuration"
// @Success      200      {object}  map[string]interface{} "message: config updated"
// @Router       /admin/config [post]
func (h handler) Save(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var payload model.Configuration
	err := json.NewDecoder(r.Body).Decode(&payload.Data)
	if err != nil {
		h.log.Error("invalid request body", zap.Error(err))
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	err = h.config.Save(r.Context(), &payload)
	if err != nil {
		status, msg := utils.MapError(err)
		if status == http.StatusNotModified {
			w.WriteHeader(http.StatusNotModified)
			return
		}
		h.log.Error("failed to get config", zap.Error(err))
		http.Error(w, msg, status)
		return
	}

	err = h.notif.PublishUpdate(context.Background())
	if err != nil {
		h.log.Error("failed to publish update", zap.Error(err))
	}

	resp := map[string]any{
		"status":  "success",
		"message": "configuration saved successfully",
	}
	utils.WriteJSON(w, http.StatusCreated, resp)
}

// RegisterAgent godoc
// @Summary      Register a new agent
// @Description  Register an agent to get a unique ID and polling configuration
// @Tags         agent
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  body      model.AgentRequest  true  "Agent Registration Info"
// @Success      200      {object}  map[string]interface{}
// @Failure      400      {object}  map[string]string "Invalid request body"
// @Router       /register [post]
func (h handler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var payload model.AgentRequest
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		h.log.Error("invalid request body", zap.Error(err))
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	agentID, err := h.agent.Register(r.Context(), &payload)
	if err != nil {
		h.log.Error("failed to register new agent", zap.Error(err))
		status, msg := utils.MapError(err)
		http.Error(w, msg, status)
		return
	}

	resp := map[string]any{
		"agent_id":              agentID,
		"poll_url":              h.cfg.PollUrl,
		"poll_interval_seconds": int(h.cfg.PollInterval.Seconds()),
	}
	utils.WriteJSON(w, http.StatusCreated, resp)
}

// GetConfig godoc
// @Summary      Poll for latest configuration
// @Description  Get the latest config if version has changed. Returns 304 if version matches.
// @Tags         agent
// @Produce      json
// @Security     BearerAuth
// @Param        X-Agent-ID     header    string  true   "Unique Agent ID"
// @Param        If-None-Match  header    string  false  "Current config version (ETag)"
// @Success      200      		{object}  map[string]interface{}
// @Success      304            {string}  string "Not Modified"
// @Failure      401            {object}  map[string]string "Unauthorized"
// @Router       /config [get]
func (h handler) Config(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ctx := r.Context()
	versionx := r.Header.Get("If-None-Match")

	sendLatestConfig := func() bool {
		res, err := h.config.Get(ctx, versionx)
		if err != nil {
			status, msg := utils.MapError(err)
			if status != http.StatusNotModified {
				h.log.Error("failed to get config", zap.Error(err))
				http.Error(w, msg, status)
				return true
			}
			return false
		}

		if fmt.Sprintf("v%d", res.Version) != versionx {
			resp := map[string]any{}
			json.Unmarshal(res.Data, &resp)
			w.Header().Set("ETag", fmt.Sprintf("v%d", res.Version))
			utils.WriteJSON(w, http.StatusOK, resp)
			return true
		}

		return false
	}

	if sent := sendLatestConfig(); sent {
		return
	}

	updateCh := h.notif.Subscribe()

	select {
	case <-time.After(60 * time.Second):
		w.WriteHeader(http.StatusNotModified)
		return
	case <-updateCh:
		sendLatestConfig()
		return
	case <-ctx.Done():
		return
	}
}
