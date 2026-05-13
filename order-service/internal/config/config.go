package config

import (
	"os"
	"time"
)

type Config struct {
	HTTPPort        string
	GRPCPort        string
	PaymentGRPCAddr string
	RedisURL        string
	CacheTTL        time.Duration
}

func Load() Config {
	httpPort := os.Getenv("ORDER_HTTP_PORT")
	if httpPort == "" {
		httpPort = "8080"
	}

	grpcPort := os.Getenv("ORDER_GRPC_PORT")
	if grpcPort == "" {
		grpcPort = "50052"
	}

	paymentAddr := os.Getenv("PAYMENT_GRPC_ADDRESS")
	if paymentAddr == "" {
		paymentAddr = "localhost:50051"
	}

	redisURL := os.Getenv("ORDER_REDIS_URL")
	if redisURL == "" {
		redisURL = "redis:6379"
	}

	cacheTTL := 5 * time.Minute
	if ttl := os.Getenv("ORDER_CACHE_TTL"); ttl != "" {
		if parsed, err := time.ParseDuration(ttl); err == nil {
			cacheTTL = parsed
		}
	}

	return Config{
		HTTPPort:        httpPort,
		GRPCPort:        grpcPort,
		PaymentGRPCAddr: paymentAddr,
		RedisURL:        redisURL,
		CacheTTL:        cacheTTL,
	}
}
