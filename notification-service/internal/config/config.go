package config

import "os"

type Config struct {
	RabbitMQURL string
	QueueName   string
}

func Load() Config {
	rabbitURL := os.Getenv("RABBITMQ_URL")
	if rabbitURL == "" {
		rabbitURL = "amqp://guest:guest@rabbitmq:5672/"
	}

	queue := os.Getenv("RABBITMQ_QUEUE")
	if queue == "" {
		queue = "payment.completed"
	}

	return Config{
		RabbitMQURL: rabbitURL,
		QueueName:   queue,
	}
}
