package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	RabbitMQURL       string
	QueueName         string
	RedisURL          string
	ProviderMode      string
	ProviderDelayMin  time.Duration
	ProviderDelayMax  time.Duration
	ProviderFailRate  int
	ProcessedTTL      time.Duration
	RetryMaxAttempts  int
	RetryBaseDelay    time.Duration
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

	redisURL := os.Getenv("NOTIFICATION_REDIS_URL")
	if redisURL == "" {
		redisURL = "redis:6379"
	}

	providerMode := os.Getenv("PROVIDER_MODE")
	if providerMode == "" {
		providerMode = "SIMULATED"
	}

	failureRate := 20
	if rate := os.Getenv("PROVIDER_FAILURE_RATE"); rate != "" {
		if parsed, err := strconv.Atoi(rate); err == nil {
			failureRate = parsed
		}
	}

	minDelay := 500 * time.Millisecond
	if d := os.Getenv("PROVIDER_DELAY_MIN"); d != "" {
		if parsed, err := time.ParseDuration(d); err == nil {
			minDelay = parsed
		}
	}

	maxDelay := 1500 * time.Millisecond
	if d := os.Getenv("PROVIDER_DELAY_MAX"); d != "" {
		if parsed, err := time.ParseDuration(d); err == nil {
			maxDelay = parsed
		}
	}

	processedTTL := 24 * time.Hour
	if ttl := os.Getenv("PROCESSED_TTL"); ttl != "" {
		if parsed, err := time.ParseDuration(ttl); err == nil {
			processedTTL = parsed
		}
	}

	retryMax := 3
	if retry := os.Getenv("RETRY_MAX_ATTEMPTS"); retry != "" {
		if parsed, err := strconv.Atoi(retry); err == nil {
			retryMax = parsed
		}
	}

	retryBaseDelay := 2 * time.Second
	if d := os.Getenv("RETRY_BASE_DELAY"); d != "" {
		if parsed, err := time.ParseDuration(d); err == nil {
			retryBaseDelay = parsed
		}
	}

	return Config{
		RabbitMQURL:      rabbitURL,
		QueueName:        queue,
		RedisURL:         redisURL,
		ProviderMode:     providerMode,
		ProviderDelayMin: minDelay,
		ProviderDelayMax: maxDelay,
		ProviderFailRate: failureRate,
		ProcessedTTL:     processedTTL,
		RetryMaxAttempts: retryMax,
		RetryBaseDelay:   retryBaseDelay,
	}
}
