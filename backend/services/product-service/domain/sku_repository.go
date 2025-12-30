package domain

import "context"

// SKURepository SKUリポジトリインターフェース
type SKURepository interface {
	// 基本CRUD
	Create(ctx context.Context, sku *SKU) error
	Update(ctx context.Context, sku *SKU) error
	Delete(ctx context.Context, skuID string) error
	GetByID(ctx context.Context, skuID string) (*SKU, error)
	GetBySKUCode(ctx context.Context, skuCode string) (*SKU, error)

	// 商品関連
	GetByProductID(ctx context.Context, productID string) ([]*SKU, error)
	GetActiveByProductID(ctx context.Context, productID string) ([]*SKU, error)

	// 在庫管理
	AdjustStock(ctx context.Context, skuID string, quantity int) error
	GetLowStockSKUs(ctx context.Context, threshold int) ([]*SKU, error)

	// バルク操作
	CreateBulk(ctx context.Context, skus []*SKU) error
	UpdateBulk(ctx context.Context, skus []*SKU) error
}
