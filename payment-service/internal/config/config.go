package config

import (
	"os"
)

type Config struct {
	GRPCPort     string
	RabbitMQURL  string
	PostgresURL  string
	PaymentQueue string
}

func Load() Config {
	port := os.Getenv("PAYMENT_GRPC_PORT")
	if port == "" {
		port = "50051"
	}

	rabbitURL := os.Getenv("RABBITMQ_URL")
	if rabbitURL == "" {
		rabbitURL = "amqp://guest:guest@rabbitmq:5672/"
	}

	postgresURL := os.Getenv("POSTGRES_URL")
	if postgresURL == "" {
		postgresURL = "postgres://postgres:password@db:5432/payments?sslmode=disable"
	}

	queueName := os.Getenv("PAYMENT_EVENT_QUEUE")
	if queueName == "" {
		queueName = "payment.completed"
	}

	return Config{
		GRPCPort:     port,
		RabbitMQURL:  rabbitURL,
		PostgresURL:  postgresURL,
		PaymentQueue: queueName,
	}
}
