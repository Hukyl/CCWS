package config

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	ClockifyAPIKey string `envconfig:"CLOCKIFY_API_KEY" required:"true"`
}

func Load() (*Config, error) {
	// Load .env file if it exists (ignore error if file doesn't exist)
	godotenv.Load()

	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
