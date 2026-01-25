package client

import (
	"bytes"
	"context"
	"distributed-configuration/internal/agent/config"
	model "distributed-configuration/pkg/models"
	"distributed-configuration/pkg/utils"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"go.uber.org/zap"
)

type ControllerClient interface {
	Register(ctx context.Context, agentName, hostname string) (model.AgentResponse, error)
	FetchConfig(ctx context.Context, agentID, etag, pollUrl string) (model.ConfigResponse, error)
}

type controllerClient struct {
	log        *utils.Logger
	cfg        *config.Config
	httpClient *http.Client
}

func NewControllerClient(log *utils.Logger, cfg *config.Config) ControllerClient {
	return &controllerClient{
		log: log,
		cfg: cfg,
		httpClient: &http.Client{
			Timeout: cfg.Timeout,
		},
	}
}

func (c *controllerClient) Register(ctx context.Context, agentName, hostname string) (model.AgentResponse, error) {
	var res model.AgentResponse

	payload := map[string]any{
		"name": agentName,
		"host": hostname,
	}

	body, _ := json.Marshal(payload)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.cfg.ControllerUrl+"/register", bytes.NewBuffer(body))
	if err != nil {
		c.log.Error("failed create new request", zap.Error(err))
		return model.AgentResponse{}, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		c.log.Error("network error, failed to register", zap.Error(err))
		return model.AgentResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		errBody, _ := io.ReadAll(resp.Body)
		c.log.Error(string(errBody))
		return model.AgentResponse{}, fmt.Errorf("register failed (status %d): %s", resp.StatusCode, string(errBody))
	}

	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		c.log.Error("invalid json response", zap.Error(err))
		return model.AgentResponse{}, err
	}

	return res, nil
}
func (c *controllerClient) FetchConfig(ctx context.Context, agentID, etag, pollUrl string) (model.ConfigResponse, error) {
	var res model.ConfigResponse

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.cfg.ControllerUrl+pollUrl, nil)
	if err != nil {
		c.log.Error("failed create new request", zap.Error(err))
		return model.ConfigResponse{}, err
	}

	req.Header.Set("Authorization", "Bearer "+c.cfg.ControllerSecret)
	req.Header.Set("X-Agent-ID", agentID)
	req.Header.Set("Accept", "application/json")
	if etag != "" {
		req.Header.Set("If-None-Match", etag)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		c.log.Error("network error, failed to get config", zap.Error(err))
		return model.ConfigResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotModified {
		c.log.Info("data not modified", zap.Int("code", resp.StatusCode))
		return model.ConfigResponse{}, nil
	}

	if resp.StatusCode != http.StatusOK {
		errBody, _ := io.ReadAll(resp.Body)
		c.log.Error(string(errBody))
		return model.ConfigResponse{}, fmt.Errorf("register failed (status %d): %s", resp.StatusCode, string(errBody))
	}

	res.ETag = resp.Header.Get("ETag")
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		c.log.Error("invalid json response", zap.Error(err))
		return model.ConfigResponse{}, err
	}

	return res, nil
}
