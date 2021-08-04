package config

import (
	"fmt"

	"github.com/caarlos0/env/v6"
)

var cfg = &Config{}

type Config struct {
	ListenPort  int `env:"port" envDefault:"8082"`
	MetricsPort int `env:"metrics_port" envDefault:"9090"`
}

// GetConfig - return config.
func GetConfig() (*Config, error) {
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("failed to ENV param: %w", err)
	}

	return cfg, nil
}
