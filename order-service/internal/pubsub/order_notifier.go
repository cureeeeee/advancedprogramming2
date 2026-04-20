package pubsub

import (
	"sync"

	"github.com/cureeeeee/order-service/internal/domain"
)

type OrderNotifier struct {
	mu   sync.RWMutex
	subs map[string]map[chan domain.StatusUpdate]struct{}
}

func NewOrderNotifier() *OrderNotifier {
	return &OrderNotifier{subs: make(map[string]map[chan domain.StatusUpdate]struct{})}
}

func (n *OrderNotifier) Subscribe(orderID string) chan domain.StatusUpdate {
	n.mu.Lock()
	defer n.mu.Unlock()

	ch := make(chan domain.StatusUpdate, 8)
	if _, ok := n.subs[orderID]; !ok {
		n.subs[orderID] = make(map[chan domain.StatusUpdate]struct{})
	}
	n.subs[orderID][ch] = struct{}{}
	return ch
}

func (n *OrderNotifier) Unsubscribe(orderID string, ch chan domain.StatusUpdate) {
	n.mu.Lock()
	defer n.mu.Unlock()

	if m, ok := n.subs[orderID]; ok {
		if _, exists := m[ch]; exists {
			delete(m, ch)
			close(ch)
		}
		if len(m) == 0 {
			delete(n.subs, orderID)
		}
	}
}

func (n *OrderNotifier) Publish(update domain.StatusUpdate) {
	n.mu.RLock()
	defer n.mu.RUnlock()

	for ch := range n.subs[update.OrderID] {
		select {
		case ch <- update:
		default:
		}
	}
}
