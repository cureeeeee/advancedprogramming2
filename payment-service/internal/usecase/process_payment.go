package usecase

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/cureeeeee/payment-service/internal/domain"
	"github.com/cureeeeee/payment-service/internal/messaging"
	"github.com/google/uuid"
)

type PaymentRepository interface {
	Save(transactionID string, payment domain.Payment) error
}

type EventPublisher interface {
	PublishPaymentCompleted(ctx context.Context, event messaging.PaymentCompletedEvent) error
}

type PaymentUseCase struct {
	repo      PaymentRepository
	publisher EventPublisher
}

func NewPaymentUseCase(repo PaymentRepository, publisher EventPublisher) *PaymentUseCase {
	return &PaymentUseCase{repo: repo, publisher: publisher}
}

func (u *PaymentUseCase) ProcessPayment(ctx context.Context, payment domain.Payment) (domain.PaymentResult, error) {
	transactionID := uuid.NewString()
	if err := u.repo.Save(transactionID, payment); err != nil {
		return domain.PaymentResult{}, fmt.Errorf("save transaction: %w", err)
	}

	event := messaging.PaymentCompletedEvent{
		OrderID:       payment.OrderID,
		Amount:        payment.Amount,
		CustomerEmail: "user@example.com",
		Status:        "PAID",
	}

	if u.publisher != nil {
		log.Printf("[Payment] Publishing event: order_id=%s, amount=%.2f, customer_email=%s", event.OrderID, event.Amount, event.CustomerEmail)
		if err := u.publisher.PublishPaymentCompleted(ctx, event); err != nil {
			return domain.PaymentResult{}, fmt.Errorf("publish payment event: %w", err)
		}
		log.Printf("[Payment] Event published successfully to RabbitMQ")
	}

	return domain.PaymentResult{
		Success:       true,
		TransactionID: transactionID,
		Message:       "payment processed",
		ProcessedAt:   time.Now().UTC(),
	}, nil
}
