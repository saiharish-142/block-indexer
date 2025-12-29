package internal

import (
	"fmt"

	"github.com/example/block-indexer/pkg/config"
)

type TraceConfig struct {
	Concurrency int `mapstructure:"concurrency" yaml:"concurrency"`
}

type Config struct {
	config.BaseConfig `mapstructure:",squash" yaml:",inline"`
	Trace             TraceConfig `mapstructure:"trace" yaml:"trace"`
	RPC               struct {
		EVMEndpoint string `mapstructure:"evmEndpoint" yaml:"evmEndpoint"`
	} `mapstructure:"rpc" yaml:"rpc"`
}

func LoadConfig() (Config, error) {
	var cfg Config
	commonPath := config.MustEnv("CONFIG_PATH", "./configs/common.yaml")
	servicePath := config.MustEnv("TRACE_SERVICE_CONFIG", "")
	if err := config.Load([]string{commonPath, servicePath}, "TRACE_SERVICE", &cfg); err != nil {
		return cfg, fmt.Errorf("load config: %w", err)
	}
	if cfg.HTTPPort == 0 {
		cfg.HTTPPort = 9004
	}
	if cfg.Trace.Concurrency == 0 {
		cfg.Trace.Concurrency = 2
	}
	return cfg, nil
}
