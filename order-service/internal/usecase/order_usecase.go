package usecase

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/youruser/order-service/internal/domain"
)

var ErrValidation = errors.New("validation error")
var ErrNotFound = errors.New("not found")

type OrderRepository interface {
	Create(order domain.Order) error
	UpdateStatus(orderID, status string, updatedAt time.Time) (domain.Order, error)
	GetByID(orderID string) (domain.Order, error)
}

type PaymentGateway interface {
	ProcessPayment(ctx context.Context, orderID string, amount float64, currency string) (string, error)
}

type Notifier interface {
	Publish(update domain.StatusUpdate)
}

type OrderUseCase struct {
	repo     OrderRepository
	payment  PaymentGateway
	notifier Notifier
}

func NewOrderUseCase(repo OrderRepository, payment PaymentGateway, notifier Notifier) *OrderUseCase {
	return &OrderUseCase{repo: repo, payment: payment, notifier: notifier}
}

func (u *OrderUseCase) CreateOrder(ctx context.Context, amount float64, currency string) (domain.Order, string, error) {
	if amount <= 0 {
		return domain.Order{}, "", fmt.Errorf("%w: amount must be positive", ErrValidation)
	}
	if strings.TrimSpace(currency) == "" {
		return domain.Order{}, "", fmt.Errorf("%w: currency is required", ErrValidation)
	}

	now := time.Now().UTC()
	order := domain.Order{
		ID:        uuid.NewString(),
		Amount:    amount,
		Currency:  currency,
		Status:    "CREATED",
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := u.repo.Create(order); err != nil {
		return domain.Order{}, "", fmt.Errorf("create order: %w", err)
	}

	txnID, err := u.payment.ProcessPayment(ctx, order.ID, order.Amount, order.Currency)
	if err != nil {
		return domain.Order{}, "", fmt.Errorf("process payment: %w", err)
	}

	order, err = u.repo.UpdateStatus(order.ID, "PAID", time.Now().UTC())
	if err == nil {
		u.notifier.Publish(domain.StatusUpdate{OrderID: order.ID, Status: order.Status, UpdatedAt: order.UpdatedAt})
	}

	return order, txnID, nil
}

func (u *OrderUseCase) UpdateStatus(ctx context.Context, orderID, status string) (domain.Order, error) {
	_ = ctx
	if strings.TrimSpace(orderID) == "" || strings.TrimSpace(status) == "" {
		return domain.Order{}, fmt.Errorf("%w: order_id and status are required", ErrValidation)
	}

	updated, err := u.repo.UpdateStatus(orderID, status, time.Now().UTC())
	if err != nil {
		return domain.Order{}, fmt.Errorf("%w: %v", ErrNotFound, err)
	}

	u.notifier.Publish(domain.StatusUpdate{OrderID: updated.ID, Status: updated.Status, UpdatedAt: updated.UpdatedAt})
	return updated, nil
}

func (u *OrderUseCase) GetOrder(ctx context.Context, orderID string) (domain.Order, error) {
	_ = ctx
	if strings.TrimSpace(orderID) == "" {
		return domain.Order{}, fmt.Errorf("%w: order_id is required", ErrValidation)
	}

	order, err := u.repo.GetByID(orderID)
	if err != nil {
		return domain.Order{}, fmt.Errorf("%w: %v", ErrNotFound, err)
	}
	return order, nil
}
