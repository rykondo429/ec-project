package mysql

import (
	"context"
	"fmt"

	"ec-sample/services/order-service/domain"

	"gorm.io/gorm"
)

// MySQLRepository MySQL実装
type MySQLRepository struct {
	db *gorm.DB
}

// NewMySQLRepository コンストラクタ
func NewMySQLRepository(db *gorm.DB) domain.OrderRepository {
	return &MySQLRepository{db: db}
}

// GetByID IDで取得
func (r *MySQLRepository) GetByID(ctx context.Context, id string) (*domain.Order, error) {
	var order domain.Order
	if err := r.db.WithContext(ctx).
		Preload("Items").
		Where("id = ?", id).
		First(&order).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("order not found")
		}
		return nil, fmt.Errorf("failed to get order: %w", err)
	}
	return &order, nil
}

// GetByUserID ユーザーIDで取得
func (r *MySQLRepository) GetByUserID(ctx context.Context, userID string, page, pageSize int) ([]*domain.Order, int64, error) {
	var orders []*domain.Order
	var total int64

	// 総件数取得
	if err := r.db.WithContext(ctx).
		Model(&domain.Order{}).
		Where("user_id = ?", userID).
		Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to get count: %w", err)
	}

	offset := (page - 1) * pageSize

	// オーダー取得
	if err := r.db.WithContext(ctx).
		Preload("Items").
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&orders).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to get orders: %w", err)
	}

	return orders, total, nil
}

// Create 作成
func (r *MySQLRepository) Create(ctx context.Context, order *domain.Order) error {
	if err := r.db.WithContext(ctx).Create(order).Error; err != nil {
		return fmt.Errorf("failed to create order: %w", err)
	}
	return nil
}

// Update 更新
func (r *MySQLRepository) Update(ctx context.Context, order *domain.Order) error {
	if err := r.db.WithContext(ctx).Save(order).Error; err != nil {
		return fmt.Errorf("failed to update order: %w", err)
	}
	return nil
}

// UpdateStatus ステータス更新
func (r *MySQLRepository) UpdateStatus(ctx context.Context, id string, status domain.OrderStatus) error {
	if err := r.db.WithContext(ctx).
		Model(&domain.Order{}).
		Where("id = ?", id).
		Update("status", status).Error; err != nil {
		return fmt.Errorf("failed to update order status: %w", err)
	}
	return nil
}
