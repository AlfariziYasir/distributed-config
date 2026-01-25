package config

import (
	"fmt"
	"time"

	"github.com/caarlos0/env/v11"
	_ "github.com/joho/godotenv/autoload"
)

type Config struct {
	AdminSecret  string        `env:"ADMIN_SECRET"`
	AgentSecret  string        `env:"AGENT_SECRET"`
	DBDSN        string        `env:"DBDSN"`
	RedisAddr    string        `env:"REDIS_ADDR"`
	HTTPPort     int           `env:"HTTP_PORT"`
	PollUrl      string        `env:"POLL_URL"`
	PollInterval time.Duration `env:"POLL_INTERVAL"`
}

func NewConfig() (*Config, error) {
	cfg := Config{}

	err := env.Parse(&cfg)
	if err != nil {
		return nil, fmt.Errorf("read env error: %w", err)
	}

	return &cfg, nil
}
