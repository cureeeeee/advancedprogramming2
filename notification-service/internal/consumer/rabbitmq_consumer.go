package consumer

import (
	"context"
	"fmt"

	"github.com/streadway/amqp"
)

type EventConsumer struct {
	conn      *amqp.Connection
	channel   *amqp.Channel
	queueName string
}

func NewRabbitMQConsumer(url, queueName string) (*EventConsumer, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("dial rabbitmq: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		_ = conn.Close()
		return nil, fmt.Errorf("open channel: %w", err)
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

	return &EventConsumer{conn: conn, channel: ch, queueName: queueName}, nil
}

func (c *EventConsumer) Consume(ctx context.Context) (<-chan amqp.Delivery, error) {
	if err := c.channel.Qos(1, 0, false); err != nil {
		return nil, fmt.Errorf("set qos: %w", err)
	}

	msgs, err := c.channel.Consume(
		c.queueName,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("consume queue: %w", err)
	}

	return msgs, nil
}

func (c *EventConsumer) Close() error {
	var err error
	if c.channel != nil {
		if closeErr := c.channel.Close(); closeErr != nil {
			err = closeErr
		}
	}
	if c.conn != nil {
		if closeErr := c.conn.Close(); closeErr != nil && err == nil {
			err = closeErr
		}
	}
	return err
}
