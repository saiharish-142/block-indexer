package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
)

// BaseConfig holds common service settings.
type BaseConfig struct {
	ServiceName string      `mapstructure:"serviceName" yaml:"serviceName"`
	HTTPPort    int         `mapstructure:"httpPort" yaml:"httpPort"`
	LogLevel    string      `mapstructure:"logLevel" yaml:"logLevel"`
	DB          DBConfig    `mapstructure:"db" yaml:"db"`
	Redis       RedisConfig `mapstructure:"redis" yaml:"redis"`
	NATS        NATSConfig  `mapstructure:"nats" yaml:"nats"`
}

type DBConfig struct {
	Host     string `mapstructure:"host" yaml:"host"`
	Port     int    `mapstructure:"port" yaml:"port"`
	User     string `mapstructure:"user" yaml:"user"`
	Password string `mapstructure:"password" yaml:"password"`
	Name     string `mapstructure:"name" yaml:"name"`
	SSLMode  string `mapstructure:"sslMode" yaml:"sslMode"`
	MaxConns int32  `mapstructure:"maxConns" yaml:"maxConns"`
	MinConns int32  `mapstructure:"minConns" yaml:"minConns"`
}

type RedisConfig struct {
	Addr     string `mapstructure:"addr" yaml:"addr"`
	DB       int    `mapstructure:"db" yaml:"db"`
	Password string `mapstructure:"password" yaml:"password"`
}

type NATSConfig struct {
	URL           string `mapstructure:"url" yaml:"url"`
	SubjectPrefix string `mapstructure:"subjectPrefix" yaml:"subjectPrefix"`
}

// Load reads YAML configs in order and applies environment overrides.
func Load(paths []string, envPrefix string, cfg any) error {
	v := viper.New()
	for _, p := range paths {
		if p == "" {
			continue
		}
		v.SetConfigFile(p)
		if err := v.MergeInConfig(); err != nil {
			return fmt.Errorf("read config %s: %w", p, err)
		}
	}
	if envPrefix != "" {
		v.SetEnvPrefix(envPrefix)
	}
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()
	if err := v.Unmarshal(cfg); err != nil {
		return fmt.Errorf("unmarshal config: %w", err)
	}
	return nil
}

func MustEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return fallback
}
