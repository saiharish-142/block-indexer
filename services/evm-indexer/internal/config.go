package internal

import (
	"fmt"

	"github.com/example/block-indexer/pkg/config"
)

type RPCConfig struct {
	EVMEndpoint string `mapstructure:"evmEndpoint" yaml:"evmEndpoint"`
	WSEndpoint  string `mapstructure:"wsEndpoint" yaml:"wsEndpoint"`
}

type BackfillConfig struct {
	Enabled         bool `mapstructure:"enabled" yaml:"enabled"`
	IntervalSeconds int  `mapstructure:"intervalSeconds" yaml:"intervalSeconds"`
}

type Config struct {
	config.BaseConfig `mapstructure:",squash" yaml:",inline"`
	RPC               RPCConfig      `mapstructure:"rpc" yaml:"rpc"`
	Backfill          BackfillConfig `mapstructure:"backfill" yaml:"backfill"`
}

func LoadConfig() (Config, error) {
	var cfg Config
	commonPath := config.MustEnv("CONFIG_PATH", "./configs/common.yaml")
	servicePath := config.MustEnv("EVM_INDEXER_CONFIG", "")
	if err := config.Load([]string{commonPath, servicePath}, "EVM_INDEXER", &cfg); err != nil {
		return cfg, fmt.Errorf("load config: %w", err)
	}
	if cfg.HTTPPort == 0 {
		cfg.HTTPPort = 9001
	}
	return cfg, nil
}
