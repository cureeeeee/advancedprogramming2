package memory

import (
	"sync"

	"github.com/youruser/payment-service/internal/domain"
)

type PaymentRepository struct {
	mu           sync.RWMutex
	transactions map[string]domain.Payment
}

func NewPaymentRepository() *PaymentRepository {
	return &PaymentRepository{transactions: make(map[string]domain.Payment)}
}

func (r *PaymentRepository) Save(transactionID string, payment domain.Payment) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.transactions[transactionID] = payment
	return nil
}
