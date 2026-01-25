package handler

import (
	"distributed-configuration/internal/worker/config"
	"distributed-configuration/internal/worker/service"
	"distributed-configuration/pkg/utils"
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
)

type handler struct {
	log    *utils.Logger
	cfg    *config.Config
	worker service.WorkerService
}

func NewHandler(log *utils.Logger, cfg *config.Config, worker service.WorkerService) *handler {
	return &handler{
		log:    log,
		cfg:    cfg,
		worker: worker,
	}
}

// SaveConfig godoc
// @Summary      Receive config
// @Description  Receive config sent by an agent to store at internal storage
// @Tags         agent
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  body      map[string]interface{}  true  "config data"
// @Success      200      {object}  map[string]interface{}
// @Failure      400      {object}  map[string]string "Invalid request body"
// @Router       /config [post]
func (h *handler) Save(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var payload json.RawMessage
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		h.log.Error("invalid request body", zap.Error(err))
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	err = h.worker.Save(r.Context(), payload)
	resp := map[string]any{
		"status":  "success",
		"message": "configuration saved successfully",
	}
	utils.WriteJSON(w, http.StatusOK, resp)
}

// GetConfig godoc
// @Summary      Fetch config data
// @Description  Fetch config data for client
// @Tags         client
// @Produce      json
// @Security     BearerAuth
// @Success      200      		{object}  map[string]interface{}
// @Failure      401            {object}  map[string]string "Unauthorized"
// @Router       /hit [get]
func (h handler) Get(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	res, err := h.worker.Get(r.Context())
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

	utils.WriteJSON(w, http.StatusOK, res)
}
