package internal

import (
	"fmt"

	"github.com/example/block-indexer/pkg/config"
)

type CORSConfig struct {
	AllowedOrigins []string `mapstructure:"allowedOrigins" yaml:"allowedOrigins"`
}

type Config struct {
	config.BaseConfig `mapstructure:",squash" yaml:",inline"`
	CORS              CORSConfig `mapstructure:"cors" yaml:"cors"`
}

func LoadConfig() (Config, error) {
	var cfg Config
	commonPath := config.MustEnv("CONFIG_PATH", "./configs/common.yaml")
	servicePath := config.MustEnv("API_GATEWAY_CONFIG", "")
	if err := config.Load([]string{commonPath, servicePath}, "API_GATEWAY", &cfg); err != nil {
		return cfg, fmt.Errorf("load config: %w", err)
	}
	if cfg.HTTPPort == 0 {
		cfg.HTTPPort = 8080
	}
	return cfg, nil
}
