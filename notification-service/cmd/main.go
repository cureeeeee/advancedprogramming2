package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"

	"github.com/cureeeeee/notification-service/internal/config"
	"github.com/cureeeeee/notification-service/internal/consumer"
	"github.com/cureeeeee/notification-service/internal/domain"
	"github.com/cureeeeee/notification-service/internal/store"
)

func main() {
	_ = godotenv.Load()

	cfg := config.Load()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	consumerClient, err := consumer.NewRabbitMQConsumer(cfg.RabbitMQURL, cfg.QueueName)
	if err != nil {
		log.Fatalf("create rabbitmq consumer: %v", err)
	}
	defer func() {
		if err := consumerClient.Close(); err != nil {
			log.Printf("close rabbitmq consumer: %v", err)
		}
	}()

	messages, err := consumerClient.Consume(ctx)
	if err != nil {
		log.Fatalf("start consuming messages: %v", err)
	}

	processed := store.NewProcessedStore()
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

			var event domain.PaymentCompletedEvent
			if err := json.Unmarshal(msg.Body, &event); err != nil {
				log.Printf("invalid message payload: %v", err)
				_ = msg.Nack(false, false)
				continue
			}

			if event.OrderID == "" {
				log.Println("skipping payment event with missing order_id")
				_ = msg.Ack(false)
				continue
			}

			if processed.Seen(event.OrderID) {
				log.Printf("duplicate payment event ignored for order %s", event.OrderID)
				_ = msg.Ack(false)
				continue
			}

			if event.CustomerEmail == "" {
				event.CustomerEmail = "user@example.com"
			}

			log.Printf("[Notification] Sent email to %s for Order #%s. Amount: $%.2f", event.CustomerEmail, event.OrderID, event.Amount)
			processed.Mark(event.OrderID)
			if err := msg.Ack(false); err != nil {
				log.Printf("ack message failed: %v", err)
			}
		}
	}
}
