package config

import (
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	ClockifyAPIKey string `envconfig:"CLOCKIFY_API_KEY" required:"true"`
}

func Load() (*Config, error) {
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
