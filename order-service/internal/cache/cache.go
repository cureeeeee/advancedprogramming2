package cache

import (
	"context"

	"github.com/cureeeeee/order-service/internal/domain"
)

type Cache interface {
	Get(ctx context.Context, orderID string) (domain.Order, error)
	Set(ctx context.Context, order domain.Order) error
	Delete(ctx context.Context, orderID string) error
}
