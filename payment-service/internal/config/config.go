package config

import (
	"os"
)

type Config struct {
	GRPCPort string
}

func Load() Config {
	port := os.Getenv("PAYMENT_GRPC_PORT")
	if port == "" {
		port = "50051"
	}

	return Config{GRPCPort: port}
}
