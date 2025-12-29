package internal

import (
	"fmt"

	"github.com/example/block-indexer/pkg/config"
)

type Config struct {
	config.BaseConfig `mapstructure:",squash" yaml:",inline"`
}

func LoadConfig() (Config, error) {
	var cfg Config
	commonPath := config.MustEnv("CONFIG_PATH", "./configs/common.yaml")
	servicePath := config.MustEnv("WS_SERVICE_CONFIG", "")
	if err := config.Load([]string{commonPath, servicePath}, "WS_SERVICE", &cfg); err != nil {
		return cfg, fmt.Errorf("load config: %w", err)
	}
	if cfg.HTTPPort == 0 {
		cfg.HTTPPort = 8090
	}
	return cfg, nil
}
