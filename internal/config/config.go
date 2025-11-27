package config

import (
	"os"
	"strconv"
	"time"
)

// Config holds shared configuration for all services.
type Config struct {
	Env               string
	APIAddr           string
	WSAddr            string
	GRPCAddr          string
	MetricsAddr       string
	PostgresURL       string
	RedisAddr         string
	RedisPassword     string
	ChainRPCURL       string
	ConfirmationDepth int
	PollInterval      time.Duration
	BatchSize         int
	GrpcTarget        string
}

// Load builds configuration from environment variables with sensible defaults.
func Load() Config {
	return Config{
		Env:               getEnv("APP_ENV", "development"),
		APIAddr:           getEnv("API_ADDR", ":8080"),
		WSAddr:            getEnv("WS_ADDR", ":8090"),
		GRPCAddr:          getEnv("GRPC_ADDR", ":9100"),
		MetricsAddr:       getEnv("METRICS_ADDR", ":9101"),
		PostgresURL:       getEnv("POSTGRES_URL", "postgres://user:pass@localhost:5432/block_indexer?sslmode=disable"),
		RedisAddr:         getEnv("REDIS_ADDR", "localhost:6379"),
		RedisPassword:     getEnv("REDIS_PASSWORD", ""),
		ChainRPCURL:       getEnv("CHAIN_RPC_URL", "ws://localhost:8546"),
		ConfirmationDepth: getEnvInt("CONFIRM_DEPTH", 50),
		PollInterval:      getEnvDuration("POLL_INTERVAL", 2*time.Second),
		BatchSize:         getEnvInt("BATCH_SIZE", 200),
		GrpcTarget:        getEnv("GRPC_TARGET", "dns:///localhost:9100"),
	}
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func getEnvInt(key string, def int) int {
	if v := os.Getenv(key); v != "" {
		parsed, err := strconv.Atoi(v)
		if err == nil {
			return parsed
		}
	}
	return def
}

func getEnvDuration(key string, def time.Duration) time.Duration {
	if v := os.Getenv(key); v != "" {
		parsed, err := time.ParseDuration(v)
		if err == nil {
			return parsed
		}
	}
	return def
}
