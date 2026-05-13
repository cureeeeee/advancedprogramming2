package provider

import (
	"context"

	"github.com/cureeeeee/notification-service/internal/domain"
)

type EmailProvider interface {
	Send(ctx context.Context, event domain.PaymentCompletedEvent) error
}
