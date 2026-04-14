package config

import "os"

type Config struct {
	HTTPPort        string
	GRPCPort        string
	PaymentGRPCAddr string
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

	return Config{
		HTTPPort:        httpPort,
		GRPCPort:        grpcPort,
		PaymentGRPCAddr: paymentAddr,
	}
}
