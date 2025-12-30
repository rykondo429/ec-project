package usecase

import (
	"context"
	"fmt"

	"ec-sample/services/product-service/domain"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// SKUUseCase SKUユースケース
type SKUUseCase struct {
	skuRepo     domain.SKURepository
	productRepo domain.ProductRepository
	logger      *zap.Logger
}

// NewSKUUseCase コンストラクタ
func NewSKUUseCase(skuRepo domain.SKURepository, productRepo domain.ProductRepository, logger *zap.Logger) *SKUUseCase {
	return &SKUUseCase{
		skuRepo:     skuRepo,
		productRepo: productRepo,
		logger:      logger,
	}
}

// CreateSKU SKU作成
func (uc *SKUUseCase) CreateSKU(ctx context.Context, req *domain.SKUCreateRequest) (*domain.SKU, error) {
	// 商品の存在確認
	if _, err := uc.productRepo.GetByID(ctx, req.ProductID); err != nil {
		return nil, fmt.Errorf("product not found: %w", err)
	}

	// SKUコードの重複確認
	if _, err := uc.skuRepo.GetBySKUCode(ctx, req.SKUCode); err == nil {
		return nil, fmt.Errorf("SKU code already exists: %s", req.SKUCode)
	}

	sku := &domain.SKU{
		ID:         uuid.New().String(),
		ProductID:  req.ProductID,
		SKUCode:    req.SKUCode,
		Name:       req.Name,
		Attributes: req.Attributes,
		Price:      req.Price,
		Stock:      req.Stock,
		IsActive:   true,
	}

	if err := uc.skuRepo.Create(ctx, sku); err != nil {
		uc.logger.Error("Failed to create SKU", zap.Error(err), zap.String("sku_code", req.SKUCode))
		return nil, fmt.Errorf("failed to create SKU: %w", err)
	}

	uc.logger.Info("SKU created successfully", zap.String("sku_id", sku.ID), zap.String("sku_code", sku.SKUCode))
	return sku, nil
}

// GetSKUByID SKU取得
func (uc *SKUUseCase) GetSKUByID(ctx context.Context, skuID string) (*domain.SKU, error) {
	if skuID == "" {
		return nil, fmt.Errorf("SKU ID is required")
	}
	return uc.skuRepo.GetByID(ctx, skuID)
}

// GetSKUsBySKUCode SKUコードでSKU取得
func (uc *SKUUseCase) GetSKUBySKUCode(ctx context.Context, skuCode string) (*domain.SKU, error) {
	if skuCode == "" {
		return nil, fmt.Errorf("SKU code is required")
	}
	return uc.skuRepo.GetBySKUCode(ctx, skuCode)
}

// GetSKUsByProductID 商品のSKU一覧取得
func (uc *SKUUseCase) GetSKUsByProductID(ctx context.Context, productID string, activeOnly bool) ([]*domain.SKU, error) {
	if productID == "" {
		return nil, fmt.Errorf("product ID is required")
	}

	if activeOnly {
		return uc.skuRepo.GetActiveByProductID(ctx, productID)
	}
	return uc.skuRepo.GetByProductID(ctx, productID)
}

// UpdateSKU SKU更新
func (uc *SKUUseCase) UpdateSKU(ctx context.Context, skuID string, req *domain.SKUUpdateRequest) (*domain.SKU, error) {
	// 既存SKUを取得
	sku, err := uc.skuRepo.GetByID(ctx, skuID)
	if err != nil {
		return nil, err
	}

	// 更新項目を反映
	if req.Name != nil {
		sku.Name = *req.Name
	}
	if req.Attributes != nil {
		sku.Attributes = *req.Attributes
	}
	if req.Price != nil {
		sku.Price = *req.Price
	}
	if req.Stock != nil {
		sku.Stock = *req.Stock
	}
	if req.IsActive != nil {
		sku.IsActive = *req.IsActive
	}

	if err := uc.skuRepo.Update(ctx, sku); err != nil {
		uc.logger.Error("Failed to update SKU", zap.Error(err), zap.String("sku_id", skuID))
		return nil, fmt.Errorf("failed to update SKU: %w", err)
	}

	uc.logger.Info("SKU updated successfully", zap.String("sku_id", skuID))
	return sku, nil
}

// DeleteSKU SKU削除（論理削除）
func (uc *SKUUseCase) DeleteSKU(ctx context.Context, skuID string) error {
	if skuID == "" {
		return fmt.Errorf("SKU ID is required")
	}

	if err := uc.skuRepo.Delete(ctx, skuID); err != nil {
		uc.logger.Error("Failed to delete SKU", zap.Error(err), zap.String("sku_id", skuID))
		return fmt.Errorf("failed to delete SKU: %w", err)
	}

	uc.logger.Info("SKU deleted successfully", zap.String("sku_id", skuID))
	return nil
}

// AdjustStock 在庫調整
func (uc *SKUUseCase) AdjustStock(ctx context.Context, adjustment *domain.SKUStockAdjustment) (*domain.SKU, error) {
	if adjustment.SKUID == "" {
		return nil, fmt.Errorf("SKU ID is required")
	}
	if adjustment.Quantity == 0 {
		return nil, fmt.Errorf("quantity must not be zero")
	}

	// 在庫調整実行
	if err := uc.skuRepo.AdjustStock(ctx, adjustment.SKUID, adjustment.Quantity); err != nil {
		uc.logger.Error("Failed to adjust stock",
			zap.Error(err),
			zap.String("sku_id", adjustment.SKUID),
			zap.Int("quantity", adjustment.Quantity),
		)
		return nil, fmt.Errorf("failed to adjust stock: %w", err)
	}

	// 更新後のSKUを取得
	sku, err := uc.skuRepo.GetByID(ctx, adjustment.SKUID)
	if err != nil {
		return nil, err
	}

	uc.logger.Info("Stock adjusted successfully",
		zap.String("sku_id", adjustment.SKUID),
		zap.Int("quantity_change", adjustment.Quantity),
		zap.Int("new_stock", sku.Stock),
		zap.String("reason", adjustment.Reason),
	)

	return sku, nil
}

// GetLowStockSKUs 低在庫SKU一覧取得
func (uc *SKUUseCase) GetLowStockSKUs(ctx context.Context, threshold int) ([]*domain.SKU, error) {
	if threshold < 0 {
		threshold = 10 // デフォルト閾値
	}

	return uc.skuRepo.GetLowStockSKUs(ctx, threshold)
}

// CreateBulkSKUs 複数SKUを一括作成
func (uc *SKUUseCase) CreateBulkSKUs(ctx context.Context, productID string, requests []*domain.SKUCreateRequest) ([]*domain.SKU, error) {
	if productID == "" {
		return nil, fmt.Errorf("product ID is required")
	}
	if len(requests) == 0 {
		return nil, fmt.Errorf("no SKUs to create")
	}

	// 商品の存在確認
	if _, err := uc.productRepo.GetByID(ctx, productID); err != nil {
		return nil, fmt.Errorf("product not found: %w", err)
	}

	skus := make([]*domain.SKU, len(requests))
	for i, req := range requests {
		skus[i] = &domain.SKU{
			ID:         uuid.New().String(),
			ProductID:  productID,
			SKUCode:    req.SKUCode,
			Name:       req.Name,
			Attributes: req.Attributes,
			Price:      req.Price,
			Stock:      req.Stock,
			IsActive:   true,
		}
	}

	if err := uc.skuRepo.CreateBulk(ctx, skus); err != nil {
		uc.logger.Error("Failed to create SKUs in bulk", zap.Error(err), zap.String("product_id", productID))
		return nil, fmt.Errorf("failed to create SKUs in bulk: %w", err)
	}

	uc.logger.Info("SKUs created in bulk successfully",
		zap.String("product_id", productID),
		zap.Int("count", len(skus)),
	)

	return skus, nil
}

// GetProductWithSKUs 商品とそのSKU一覧を取得
func (uc *SKUUseCase) GetProductWithSKUs(ctx context.Context, productID string) (*domain.ProductWithSKUs, error) {
	if productID == "" {
		return nil, fmt.Errorf("product ID is required")
	}

	// 商品取得
	product, err := uc.productRepo.GetByID(ctx, productID)
	if err != nil {
		return nil, err
	}

	// SKU一覧取得
	skus, err := uc.skuRepo.GetActiveByProductID(ctx, productID)
	if err != nil {
		return nil, err
	}

	return &domain.ProductWithSKUs{
		Product: product,
		SKUs:    skus,
	}, nil
}
