package client

import (
	"bytes"
	"context"
	"distributed-configuration/internal/agent/config"
	"distributed-configuration/pkg/utils"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"go.uber.org/zap"
)

type WorkerClient interface {
	PushConfig(ctx context.Context, config json.RawMessage) error
}

type workerClient struct {
	log        *utils.Logger
	cfg        *config.Config
	httpClient *http.Client
}

func NewWorkerClient(log *utils.Logger, cfg *config.Config) WorkerClient {
	return &workerClient{
		log: log,
		cfg: cfg,
		httpClient: &http.Client{
			Timeout: cfg.Timeout,
		},
	}
}

func (c *workerClient) PushConfig(ctx context.Context, config json.RawMessage) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.cfg.WorkerUrl+"/register", bytes.NewBuffer(config))
	if err != nil {
		c.log.Error("failed create new request", zap.Error(err))
		return err
	}

	req.Header.Set("Authorization", "Bearer "+c.cfg.WorkerSecret)
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		c.log.Error("network error, failed to register", zap.Error(err))
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		errBody, _ := io.ReadAll(resp.Body)
		c.log.Error(string(errBody))
		return fmt.Errorf("register failed (status %d): %s", resp.StatusCode, string(errBody))
	}

	return nil
}
