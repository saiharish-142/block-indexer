package internal

import (
	"fmt"

	"github.com/example/block-indexer/pkg/config"
)

type RPCConfig struct {
	DAGEndpoint string `mapstructure:"dagEndpoint" yaml:"dagEndpoint"`
	WSEndpoint  string `mapstructure:"wsEndpoint" yaml:"wsEndpoint"`
}

type Config struct {
	config.BaseConfig `mapstructure:",squash" yaml:",inline"`
	RPC               RPCConfig `mapstructure:"rpc" yaml:"rpc"`
}

func LoadConfig() (Config, error) {
	var cfg Config
	commonPath := config.MustEnv("CONFIG_PATH", "./configs/common.yaml")
	servicePath := config.MustEnv("DAG_INDEXER_CONFIG", "")
	if err := config.Load([]string{commonPath, servicePath}, "DAG_INDEXER", &cfg); err != nil {
		return cfg, fmt.Errorf("load config: %w", err)
	}
	if cfg.HTTPPort == 0 {
		cfg.HTTPPort = 9002
	}
	return cfg, nil
}
