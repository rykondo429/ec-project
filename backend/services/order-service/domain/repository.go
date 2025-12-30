package domain

import (
	"context"
)

// OrderRepository 注文リポジトリインターフェース
type OrderRepository interface {
	GetByID(ctx context.Context, id string) (*Order, error)
	GetByUserID(ctx context.Context, userID string, page, pageSize int) ([]*Order, int64, error)
	Create(ctx context.Context, order *Order) error
	Update(ctx context.Context, order *Order) error
	UpdateStatus(ctx context.Context, id string, status OrderStatus) error
}
