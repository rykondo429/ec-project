package domain

import (
	"context"
)

// CouponRepository クーポンリポジトリインターフェース
type CouponRepository interface {
	GetByCode(ctx context.Context, code string) (*Coupon, error)
	GetByID(ctx context.Context, id string) (*Coupon, error)
	GetActiveCoupons(ctx context.Context) ([]*Coupon, error)
	Create(ctx context.Context, coupon *Coupon) error
	Update(ctx context.Context, coupon *Coupon) error
	IncrementUsageCount(ctx context.Context, couponID string) error
}

// SaleRepository セールリポジトリインターフェース
type SaleRepository interface {
	GetByID(ctx context.Context, id string) (*Sale, error)
	GetActiveSales(ctx context.Context) ([]*Sale, error)
	GetSalesByCategory(ctx context.Context, category string) ([]*Sale, error)
	GetSalesByProduct(ctx context.Context, productID string) ([]*Sale, error)
	Create(ctx context.Context, sale *Sale) error
	Update(ctx context.Context, sale *Sale) error
}
