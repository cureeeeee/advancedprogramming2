package memory

import (
	"errors"
	"sync"
	"time"

	"github.com/youruser/order-service/internal/domain"
)

var ErrOrderNotFound = errors.New("order not found")

type OrderRepository struct {
	mu     sync.RWMutex
	orders map[string]domain.Order
}

func NewOrderRepository() *OrderRepository {
	return &OrderRepository{orders: make(map[string]domain.Order)}
}

func (r *OrderRepository) Create(order domain.Order) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.orders[order.ID] = order
	return nil
}

func (r *OrderRepository) UpdateStatus(orderID, status string, updatedAt time.Time) (domain.Order, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	order, ok := r.orders[orderID]
	if !ok {
		return domain.Order{}, ErrOrderNotFound
	}

	order.Status = status
	order.UpdatedAt = updatedAt
	r.orders[orderID] = order
	return order, nil
}

func (r *OrderRepository) GetByID(orderID string) (domain.Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	order, ok := r.orders[orderID]
	if !ok {
		return domain.Order{}, ErrOrderNotFound
	}

	return order, nil
}
