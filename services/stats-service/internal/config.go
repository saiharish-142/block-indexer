package internal

import (
	"fmt"

	"github.com/example/block-indexer/pkg/config"
)

type AggregationConfig struct {
	IntervalSeconds int `mapstructure:"intervalSeconds" yaml:"intervalSeconds"`
}

type Config struct {
	config.BaseConfig `mapstructure:",squash" yaml:",inline"`
	Aggregation       AggregationConfig `mapstructure:"aggregation" yaml:"aggregation"`
}

func LoadConfig() (Config, error) {
	var cfg Config
	commonPath := config.MustEnv("CONFIG_PATH", "./configs/common.yaml")
	servicePath := config.MustEnv("STATS_SERVICE_CONFIG", "")
	if err := config.Load([]string{commonPath, servicePath}, "STATS_SERVICE", &cfg); err != nil {
		return cfg, fmt.Errorf("load config: %w", err)
	}
	if cfg.HTTPPort == 0 {
		cfg.HTTPPort = 9003
	}
	return cfg, nil
}
