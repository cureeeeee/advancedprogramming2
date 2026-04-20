package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/cureeeeee/payment-service/internal/domain"
)

type PaymentRepository interface {
	Save(transactionID string, payment domain.Payment) error
}

type PaymentUseCase struct {
	repo PaymentRepository
}

func NewPaymentUseCase(repo PaymentRepository) *PaymentUseCase {
	return &PaymentUseCase{repo: repo}
}

func (u *PaymentUseCase) ProcessPayment(ctx context.Context, payment domain.Payment) (domain.PaymentResult, error) {
	_ = ctx

	transactionID := uuid.NewString()
	if err := u.repo.Save(transactionID, payment); err != nil {
		return domain.PaymentResult{}, fmt.Errorf("save transaction: %w", err)
	}

	return domain.PaymentResult{
		Success:       true,
		TransactionID: transactionID,
		Message:       "payment processed",
		ProcessedAt:   time.Now().UTC(),
	}, nil
}
