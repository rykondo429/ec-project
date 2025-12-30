package mysql

import (
	"context"
	"fmt"

	"ec-sample/services/product-service/domain"

	"gorm.io/gorm"
)

// SKURepository SKUのMySQLリポジトリ実装
type SKURepository struct {
	db *gorm.DB
}

// NewSKURepository コンストラクタ
func NewSKURepository(db *gorm.DB) domain.SKURepository {
	return &SKURepository{db: db}
}

// Create SKUを作成
func (r *SKURepository) Create(ctx context.Context, sku *domain.SKU) error {
	if err := r.db.WithContext(ctx).Create(sku).Error; err != nil {
		return fmt.Errorf("failed to create SKU: %w", err)
	}
	return nil
}

// Update SKUを更新
func (r *SKURepository) Update(ctx context.Context, sku *domain.SKU) error {
	if err := r.db.WithContext(ctx).Save(sku).Error; err != nil {
		return fmt.Errorf("failed to update SKU: %w", err)
	}
	return nil
}

// Delete SKUを削除（論理削除）
func (r *SKURepository) Delete(ctx context.Context, skuID string) error {
	if err := r.db.WithContext(ctx).
		Model(&domain.SKU{}).
		Where("id = ?", skuID).
		Update("is_active", false).Error; err != nil {
		return fmt.Errorf("failed to delete SKU: %w", err)
	}
	return nil
}

// GetByID IDでSKUを取得
func (r *SKURepository) GetByID(ctx context.Context, skuID string) (*domain.SKU, error) {
	var sku domain.SKU
	if err := r.db.WithContext(ctx).
		Where("id = ?", skuID).
		First(&sku).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("SKU not found: %s", skuID)
		}
		return nil, fmt.Errorf("failed to get SKU: %w", err)
	}
	return &sku, nil
}

// GetBySKUCode SKUコードでSKUを取得
func (r *SKURepository) GetBySKUCode(ctx context.Context, skuCode string) (*domain.SKU, error) {
	var sku domain.SKU
	if err := r.db.WithContext(ctx).
		Where("sku_code = ?", skuCode).
		First(&sku).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("SKU not found: %s", skuCode)
		}
		return nil, fmt.Errorf("failed to get SKU by code: %w", err)
	}
	return &sku, nil
}

// GetByProductID 商品IDに紐づくSKU一覧を取得
func (r *SKURepository) GetByProductID(ctx context.Context, productID string) ([]*domain.SKU, error) {
	var skus []*domain.SKU
	if err := r.db.WithContext(ctx).
		Where("product_id = ?", productID).
		Order("sku_code ASC").
		Find(&skus).Error; err != nil {
		return nil, fmt.Errorf("failed to get SKUs by product ID: %w", err)
	}
	return skus, nil
}

// GetActiveByProductID 商品IDに紐づくアクティブなSKU一覧を取得
func (r *SKURepository) GetActiveByProductID(ctx context.Context, productID string) ([]*domain.SKU, error) {
	var skus []*domain.SKU
	if err := r.db.WithContext(ctx).
		Where("product_id = ? AND is_active = ?", productID, true).
		Order("sku_code ASC").
		Find(&skus).Error; err != nil {
		return nil, fmt.Errorf("failed to get active SKUs: %w", err)
	}
	return skus, nil
}

// AdjustStock 在庫を調整（増減）
func (r *SKURepository) AdjustStock(ctx context.Context, skuID string, quantity int) error {
	// トランザクション内で在庫を更新
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 現在の在庫を取得
		var sku domain.SKU
		if err := tx.Where("id = ?", skuID).First(&sku).Error; err != nil {
			return fmt.Errorf("SKU not found: %w", err)
		}

		// 在庫更新
		newStock := sku.Stock + quantity
		if newStock < 0 {
			return fmt.Errorf("insufficient stock: current=%d, requested=%d", sku.Stock, quantity)
		}

		if err := tx.Model(&domain.SKU{}).
			Where("id = ?", skuID).
			Update("stock", newStock).Error; err != nil {
			return fmt.Errorf("failed to update stock: %w", err)
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to adjust stock: %w", err)
	}
	return nil
}

// GetLowStockSKUs 在庫が閾値以下のSKUを取得
func (r *SKURepository) GetLowStockSKUs(ctx context.Context, threshold int) ([]*domain.SKU, error) {
	var skus []*domain.SKU
	if err := r.db.WithContext(ctx).
		Where("stock <= ? AND is_active = ?", threshold, true).
		Order("stock ASC").
		Find(&skus).Error; err != nil {
		return nil, fmt.Errorf("failed to get low stock SKUs: %w", err)
	}
	return skus, nil
}

// CreateBulk 複数のSKUを一括作成
func (r *SKURepository) CreateBulk(ctx context.Context, skus []*domain.SKU) error {
	if len(skus) == 0 {
		return nil
	}

	if err := r.db.WithContext(ctx).Create(&skus).Error; err != nil {
		return fmt.Errorf("failed to create SKUs in bulk: %w", err)
	}
	return nil
}

// UpdateBulk 複数のSKUを一括更新
func (r *SKURepository) UpdateBulk(ctx context.Context, skus []*domain.SKU) error {
	if len(skus) == 0 {
		return nil
	}

	// トランザクション内で1件ずつ更新
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, sku := range skus {
			if err := tx.Save(sku).Error; err != nil {
				return fmt.Errorf("failed to update SKU %s: %w", sku.ID, err)
			}
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to update SKUs in bulk: %w", err)
	}
	return nil
}
