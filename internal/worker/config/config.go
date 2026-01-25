package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
	_ "github.com/joho/godotenv/autoload"
)

type Config struct {
	AgentSecret  string `env:"AGENT_SECRET"`
	ClientSecret string `env:"CLIENT_SECRET"`
	HTTPPort     int    `env:"HTTP_PORT"`
}

func NewConfig() (*Config, error) {
	cfg := Config{}

	err := env.Parse(&cfg)
	if err != nil {
		return nil, fmt.Errorf("read env error: %w", err)
	}

	return &cfg, nil
}
