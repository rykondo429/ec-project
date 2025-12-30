package mysql

import (
	"context"
	"fmt"
	"time"

	"ec-sample/services/promotion-service/domain"

	"gorm.io/gorm"
)

// MySQLCouponRepository MySQL Coupon実装
type MySQLCouponRepository struct {
	db *gorm.DB
}

// NewMySQLCouponRepository コンストラクタ
func NewMySQLCouponRepository(db *gorm.DB) domain.CouponRepository {
	return &MySQLCouponRepository{db: db}
}

// GetByCode コード で取得
func (r *MySQLCouponRepository) GetByCode(ctx context.Context, code string) (*domain.Coupon, error) {
	var coupon domain.Coupon
	if err := r.db.WithContext(ctx).
		Where("code = ?", code).
		First(&coupon).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get coupon: %w", err)
	}
	return &coupon, nil
}

// GetByID IDで取得
func (r *MySQLCouponRepository) GetByID(ctx context.Context, id string) (*domain.Coupon, error) {
	var coupon domain.Coupon
	if err := r.db.WithContext(ctx).
		Where("id = ?", id).
		First(&coupon).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("coupon not found")
		}
		return nil, fmt.Errorf("failed to get coupon: %w", err)
	}
	return &coupon, nil
}

// GetActiveCoupons アクティブなクーポン取得
func (r *MySQLCouponRepository) GetActiveCoupons(ctx context.Context) ([]*domain.Coupon, error) {
	var coupons []*domain.Coupon
	if err := r.db.WithContext(ctx).
		Where("status = ? AND start_date <= NOW() AND end_date >= NOW()", domain.CouponStatusActive).
		Find(&coupons).Error; err != nil {
		return nil, fmt.Errorf("failed to get active coupons: %w", err)
	}
	return coupons, nil
}

// Create 作成
func (r *MySQLCouponRepository) Create(ctx context.Context, coupon *domain.Coupon) error {
	coupon.CreatedAt = time.Now()
	if err := r.db.WithContext(ctx).Create(coupon).Error; err != nil {
		return fmt.Errorf("failed to create coupon: %w", err)
	}
	return nil
}

// Update 更新
func (r *MySQLCouponRepository) Update(ctx context.Context, coupon *domain.Coupon) error {
	if err := r.db.WithContext(ctx).Save(coupon).Error; err != nil {
		return fmt.Errorf("failed to update coupon: %w", err)
	}
	return nil
}

// IncrementUsageCount 使用回数をインクリメント
func (r *MySQLCouponRepository) IncrementUsageCount(ctx context.Context, couponID string) error {
	if err := r.db.WithContext(ctx).
		Model(&domain.Coupon{}).
		Where("id = ?", couponID).
		Update("current_usage_count", gorm.Expr("current_usage_count + ?", 1)).Error; err != nil {
		return fmt.Errorf("failed to increment usage count: %w", err)
	}
	return nil
}

// MySQLSaleRepository MySQL Sale実装
type MySQLSaleRepository struct {
	db *gorm.DB
}

// NewMySQLSaleRepository コンストラクタ
func NewMySQLSaleRepository(db *gorm.DB) domain.SaleRepository {
	return &MySQLSaleRepository{db: db}
}

// GetByID IDで取得
func (r *MySQLSaleRepository) GetByID(ctx context.Context, id string) (*domain.Sale, error) {
	var sale domain.Sale
	if err := r.db.WithContext(ctx).
		Where("id = ?", id).
		First(&sale).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("sale not found")
		}
		return nil, fmt.Errorf("failed to get sale: %w", err)
	}
	return &sale, nil
}

// GetActiveSales アクティブなセール取得
func (r *MySQLSaleRepository) GetActiveSales(ctx context.Context) ([]*domain.Sale, error) {
	var sales []*domain.Sale
	if err := r.db.WithContext(ctx).
		Where("start_date <= NOW() AND end_date >= NOW()").
		Find(&sales).Error; err != nil {
		return nil, fmt.Errorf("failed to get active sales: %w", err)
	}
	return sales, nil
}

// GetSalesByCategory カテゴリ別取得
func (r *MySQLSaleRepository) GetSalesByCategory(ctx context.Context, category string) ([]*domain.Sale, error) {
	var sales []*domain.Sale
	if err := r.db.WithContext(ctx).
		Where("start_date <= NOW() AND end_date >= NOW() AND JSON_CONTAINS(applicable_categories, ?)", "\""+category+"\"").
		Find(&sales).Error; err != nil {
		return nil, fmt.Errorf("failed to get sales by category: %w", err)
	}
	return sales, nil
}

// GetSalesByProduct 商品別取得
func (r *MySQLSaleRepository) GetSalesByProduct(ctx context.Context, productID string) ([]*domain.Sale, error) {
	var sales []*domain.Sale
	if err := r.db.WithContext(ctx).
		Where("start_date <= NOW() AND end_date >= NOW() AND JSON_CONTAINS(applicable_products, ?)", "\""+productID+"\"").
		Find(&sales).Error; err != nil {
		return nil, fmt.Errorf("failed to get sales by product: %w", err)
	}
	return sales, nil
}

// Create 作成
func (r *MySQLSaleRepository) Create(ctx context.Context, sale *domain.Sale) error {
	sale.CreatedAt = time.Now()
	if err := r.db.WithContext(ctx).Create(sale).Error; err != nil {
		return fmt.Errorf("failed to create sale: %w", err)
	}
	return nil
}

// Update 更新
func (r *MySQLSaleRepository) Update(ctx context.Context, sale *domain.Sale) error {
	if err := r.db.WithContext(ctx).Save(sale).Error; err != nil {
		return fmt.Errorf("failed to update sale: %w", err)
	}
	return nil
}
