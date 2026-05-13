package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"github.com/streadway/amqp"

	"github.com/cureeeeee/notification-service/internal/config"
	"github.com/cureeeeee/notification-service/internal/consumer"
	"github.com/cureeeeee/notification-service/internal/domain"
	"github.com/cureeeeee/notification-service/internal/provider"
	"github.com/cureeeeee/notification-service/internal/store"
)

func main() {
	_ = godotenv.Load()

	cfg := config.Load()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	redisClient := redis.NewClient(&redis.Options{Addr: cfg.RedisURL})
	defer func() {
		if err := redisClient.Close(); err != nil {
			log.Printf("close redis client: %v", err)
		}
	}()

	consumerClient, err := consumer.NewRabbitMQConsumer(cfg.RabbitMQURL, cfg.QueueName)
	if err != nil {
		log.Fatalf("create rabbitmq consumer: %v", err)
	}
	defer func() {
		if err := consumerClient.Close(); err != nil {
			log.Printf("close rabbitmq consumer: %v", err)
		}
	}()

	processedStore := store.NewRedisProcessedStore(redisClient, cfg.ProcessedTTL)

	var emailProvider provider.EmailProvider
	switch cfg.ProviderMode {
	case "REAL":
		log.Println("PROVIDER_MODE=REAL is not implemented, falling back to SIMULATED")
		fallthrough
	default:
		emailProvider = provider.NewSimulatedEmailProvider(cfg.ProviderFailRate, cfg.ProviderDelayMin, cfg.ProviderDelayMax)
	}

	messages, err := consumerClient.Consume(ctx)
	if err != nil {
		log.Fatalf("start consuming messages: %v", err)
	}

	log.Printf("notification service is listening for messages on queue %s", cfg.QueueName)

	for {
		select {
		case <-ctx.Done():
			log.Println("notification service shutdown requested")
			return
		case msg, ok := <-messages:
			if !ok {
				log.Println("rabbitmq messages channel closed")
				return
			}

			if err := handleMessage(ctx, msg, emailProvider, processedStore, cfg); err != nil {
				log.Printf("message processing error: %v", err)
			}
		}
	}
}

func handleMessage(ctx context.Context, msg amqp.Delivery, provider provider.EmailProvider, processed *store.RedisProcessedStore, cfg config.Config) error {
	var event domain.PaymentCompletedEvent
	if err := json.Unmarshal(msg.Body, &event); err != nil {
		log.Printf("invalid message payload: %v", err)
		return msg.Nack(false, false)
	}

	if event.OrderID == "" {
		log.Println("skipping payment event with missing order_id")
		return msg.Ack(false)
	}

	processedBefore, err := processed.Seen(ctx, event.OrderID)
	if err != nil {
		return fmt.Errorf("check processed store: %w", err)
	}
	if processedBefore {
		log.Printf("duplicate payment event ignored for order %s", event.OrderID)
		return msg.Ack(false)
	}

	if event.CustomerEmail == "" {
		event.CustomerEmail = "user@example.com"
	}

	if err := processWithRetry(ctx, event, provider, cfg.RetryMaxAttempts, cfg.RetryBaseDelay); err != nil {
		log.Printf("failed to send notification for order %s after retries: %v", event.OrderID, err)
		return msg.Nack(false, true)
	}

	if err := processed.Mark(ctx, event.OrderID); err != nil {
		return fmt.Errorf("mark processed: %w", err)
	}

	log.Printf("[Notification] Sent email to %s for Order #%s. Amount: $%.2f", event.CustomerEmail, event.OrderID, event.Amount)
	return msg.Ack(false)
}

func processWithRetry(ctx context.Context, event domain.PaymentCompletedEvent, provider provider.EmailProvider, maxAttempts int, baseDelay time.Duration) error {
	var err error
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		err = provider.Send(ctx, event)
		if err == nil {
			return nil
		}

		if attempt < maxAttempts {
			delay := baseDelay * time.Duration(1<<uint(attempt-1))
			log.Printf("notification send failed, retrying in %s (%d/%d): %v", delay, attempt, maxAttempts, err)
			select {
			case <-time.After(delay):
			case <-ctx.Done():
				return ctx.Err()
			}
		}
	}
	return err
}
