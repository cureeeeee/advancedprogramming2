package provider

import (
	context"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/cureeeeee/notification-service/internal/domain"
)

type SimulatedEmailProvider struct {
	failureRate int
	minDelay    time.Duration
	maxDelay    time.Duration
	rand        *rand.Rand
}

func NewSimulatedEmailProvider(failureRate int, minDelay, maxDelay time.Duration) *SimulatedEmailProvider {
	return &SimulatedEmailProvider{
		failureRate: failureRate,
		minDelay:    minDelay,
		maxDelay:    maxDelay,
		rand:        rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (p *SimulatedEmailProvider) Send(ctx context.Context, event domain.PaymentCompletedEvent) error {
	delay := p.minDelay
	if p.maxDelay > p.minDelay {
		delta := p.maxDelay - p.minDelay
		delay += time.Duration(p.rand.Int63n(int64(delta)))
	}

	select {
	case <-time.After(delay):
		// continue
	case <-ctx.Done():
		return ctx.Err()
	}

	if p.rand.Intn(100) < p.failureRate {
		return errors.New("simulated email provider failure")
	}

	fmt.Printf("[Notification Provider] Sent email to %s for Order #%s. Amount: $%.2f\n", event.CustomerEmail, event.OrderID, event.Amount)
	return nil
}
