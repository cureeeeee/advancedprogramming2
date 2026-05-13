package messaging

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/streadway/amqp"
)

type PaymentCompletedEvent struct {
	OrderID       string  `json:"order_id"`
	Amount        float64 `json:"amount"`
	CustomerEmail string  `json:"customer_email"`
	Status        string  `json:"status"`
}

type EventPublisher interface {
	PublishPaymentCompleted(ctx context.Context, event PaymentCompletedEvent) error
	Close() error
}

type RabbitMQPublisher struct {
	conn      *amqp.Connection
	channel   *amqp.Channel
	queueName string
	confirm   chan amqp.Confirmation
}

func NewRabbitMQPublisher(url, queueName string) (*RabbitMQPublisher, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("dial rabbitmq: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		_ = conn.Close()
		return nil, fmt.Errorf("open channel: %w", err)
	}

	if err := ch.Confirm(false); err != nil {
		_ = ch.Close()
		_ = conn.Close()
		return nil, fmt.Errorf("confirm mode: %w", err)
	}

	if _, err := ch.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		_ = ch.Close()
		_ = conn.Close()
		return nil, fmt.Errorf("declare queue: %w", err)
	}

	return &RabbitMQPublisher{
		conn:      conn,
		channel:   ch,
		queueName: queueName,
		confirm:   ch.NotifyPublish(make(chan amqp.Confirmation, 1)),
	}, nil
}

func (p *RabbitMQPublisher) PublishPaymentCompleted(ctx context.Context, event PaymentCompletedEvent) error {
	body, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("marshal payment event: %w", err)
	}

	log.Printf("[RabbitMQ] Publishing to queue '%s': %s", p.queueName, string(body))

	if err := p.channel.Publish(
		"",
		p.queueName,
		false,
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "application/json",
			Body:         body,
		},
	); err != nil {
		return fmt.Errorf("publish payment event: %w", err)
	}

	select {
	case confirmed := <-p.confirm:
		if !confirmed.Ack {
			log.Printf("[RabbitMQ] Message NOT acknowledged by broker")
			return errors.New("message was not acknowledged by rabbitmq")
		}
		log.Printf("[RabbitMQ] Message acknowledged by broker")
		return nil
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(5 * time.Second):
		return errors.New("rabbitmq publish confirmation timeout")
	}
}

func (p *RabbitMQPublisher) Close() error {
	var err error
	if p.channel != nil {
		if closeErr := p.channel.Close(); closeErr != nil {
			err = closeErr
		}
	}
	if p.conn != nil {
		if closeErr := p.conn.Close(); closeErr != nil && err == nil {
			err = closeErr
		}
	}
	return err
}
