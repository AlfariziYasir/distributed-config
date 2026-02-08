package config

import (
	"fmt"
	"time"

	"github.com/caarlos0/env/v11"
	_ "github.com/joho/godotenv/autoload"
)

type Config struct {
	AdminSecret      string        `env:"ADMIN_SECRET"`
	ControllerSecret string        `env:"CONTROLLER_SECRET"`
	DBDSN            string        `env:"DBDSN"`
	RedisAddr        string        `env:"REDIS_ADDR"`
	RedisPass        string        `env:"REDIS_PASSWORD"`
	HTTPPort         int           `env:"CONTROLLER_PORT"`
	PollUrl          string        `env:"POLL_URL"`
	PollInterval     time.Duration `env:"POLL_INTERVAL"`
	ChannelKey       string        `env:"CHANNEL_KEY"`
}

func NewConfig() (*Config, error) {
	cfg := Config{}

	err := env.Parse(&cfg)
	if err != nil {
		return nil, fmt.Errorf("read env error: %w", err)
	}

	return &cfg, nil
}
