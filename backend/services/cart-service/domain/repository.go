package domain

import (
	"context"
)

// CartRepository カートリポジトリインターフェース
type CartRepository interface {
	GetCart(ctx context.Context, userID string) (*Cart, error)
	SaveCart(ctx context.Context, cart *Cart) error
	AddItem(ctx context.Context, userID string, item *CartItem) error
	UpdateItem(ctx context.Context, userID, productID string, quantity int) error
	RemoveItem(ctx context.Context, userID, productID string) error
	ClearCart(ctx context.Context, userID string) error
}
