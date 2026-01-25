package config

import (
	"fmt"
	"time"

	"github.com/caarlos0/env/v11"
	_ "github.com/joho/godotenv/autoload"
)

type Config struct {
	AgentName        string        `env:"AGENT_NAME"`
	ControllerSecret string        `env:"CONTROLLER_SECRET"`
	WorkerSecret     string        `env:"WORKER_SECRET"`
	RedisAddr        string        `env:"REDIS_ADDR"`
	ControllerUrl    string        `env:"CONTROLLER_URL"`
	WorkerUrl        string        `env:"WORKER_URL"`
	Timeout          time.Duration `env:"TIMEOUT"`
}

func NewConfig() (*Config, error) {
	cfg := Config{}

	err := env.Parse(&cfg)
	if err != nil {
		return nil, fmt.Errorf("read env error: %w", err)
	}

	return &cfg, nil
}
