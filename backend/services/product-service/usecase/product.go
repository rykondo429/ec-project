package usecase

import (
	"context"
	"fmt"

	"ec-sample/services/product-service/domain"

	"go.uber.org/zap"
)

// ProductUseCase ユースケース
type ProductUseCase struct {
	repo   domain.ProductRepository // MySQL (主データソース)
	esRepo domain.ProductRepository // Elasticsearch (検索インデックス)
	logger *zap.Logger
}

// NewProductUseCase コンストラクタ
func NewProductUseCase(repo domain.ProductRepository, esRepo domain.ProductRepository, logger *zap.Logger) *ProductUseCase {
	return &ProductUseCase{
		repo:   repo,
		esRepo: esRepo,
		logger: logger,
	}
}

// GetProductByID 商品詳細取得
func (uc *ProductUseCase) GetProductByID(ctx context.Context, id string) (*domain.Product, error) {
	if id == "" {
		return nil, fmt.Errorf("product id is required")
	}
	return uc.repo.GetByID(ctx, id)
}

// SearchProducts 商品検索
func (uc *ProductUseCase) SearchProducts(ctx context.Context, query *domain.SearchQuery) (*domain.SearchResult, error) {
	if query == nil {
		return nil, fmt.Errorf("search query is required")
	}

	// デフォルト値設定
	if query.Page <= 0 {
		query.Page = 1
	}
	if query.PageSize <= 0 {
		query.PageSize = 20
	}
	if query.PageSize > 100 {
		query.PageSize = 100
	}

	return uc.repo.Search(ctx, query)
}

// GetProductsByCategory カテゴリ別商品取得
func (uc *ProductUseCase) GetProductsByCategory(ctx context.Context, category string, page, pageSize int) (*domain.SearchResult, error) {
	if category == "" {
		return nil, fmt.Errorf("category is required")
	}

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}

	return uc.repo.GetByCategory(ctx, category, page, pageSize)
}

// GetFeaturedProducts おすすめ商品取得
func (uc *ProductUseCase) GetFeaturedProducts(ctx context.Context) ([]*domain.Product, error) {
	return uc.repo.GetFeatured(ctx)
}

// GetCategories カテゴリー一覧取得
func (uc *ProductUseCase) GetCategories(ctx context.Context) ([]string, error) {
	return uc.repo.GetCategories(ctx)
}

// CreateProduct 商品作成 (MySQL + ES同期)
func (uc *ProductUseCase) CreateProduct(ctx context.Context, product *domain.Product) error {
	// バリデーション
	if product == nil {
		return fmt.Errorf("product is required")
	}
	if product.Name == "" {
		return fmt.Errorf("product name is required")
	}
	if product.Price <= 0 {
		return fmt.Errorf("product price must be greater than 0")
	}

	// MySQLに保存
	if err := uc.repo.Create(ctx, product); err != nil {
		uc.logger.Error("MySQL product create failed", zap.Error(err), zap.String("product_id", product.ID))
		return fmt.Errorf("failed to create product in MySQL: %w", err)
	}

	// Elasticsearchに同期（非同期的に、エラーはログのみ）
	if uc.esRepo != nil {
		if err := uc.esRepo.Create(ctx, product); err != nil {
			uc.logger.Warn("Elasticsearch sync failed on create",
				zap.Error(err),
				zap.String("product_id", product.ID),
			)
			// ESの同期失敗は致命的エラーとしない
		}
	}

	uc.logger.Info("Product created successfully", zap.String("product_id", product.ID))
	return nil
}

// UpdateProduct 商品更新 (MySQL + ES同期)
func (uc *ProductUseCase) UpdateProduct(ctx context.Context, product *domain.Product) error {
	// バリデーション
	if product == nil {
		return fmt.Errorf("product is required")
	}
	if product.ID == "" {
		return fmt.Errorf("product id is required")
	}

	// MySQLで更新
	if err := uc.repo.Update(ctx, product); err != nil {
		uc.logger.Error("MySQL product update failed", zap.Error(err), zap.String("product_id", product.ID))
		return fmt.Errorf("failed to update product in MySQL: %w", err)
	}

	// Elasticsearchに同期（非同期的に、エラーはログのみ）
	if uc.esRepo != nil {
		if err := uc.esRepo.Update(ctx, product); err != nil {
			uc.logger.Warn("Elasticsearch sync failed on update",
				zap.Error(err),
				zap.String("product_id", product.ID),
			)
			// ESの同期失敗は致命的エラーとしない
		}
	}

	uc.logger.Info("Product updated successfully", zap.String("product_id", product.ID))
	return nil
}

// SyncAllProductsToElasticsearch MySQL→ESの全商品同期
func (uc *ProductUseCase) SyncAllProductsToElasticsearch(ctx context.Context) (int, error) {
	if uc.esRepo == nil {
		return 0, fmt.Errorf("elasticsearch repository is not configured")
	}

	// MySQLから全商品取得
	products, err := uc.repo.GetAll(ctx)
	if err != nil {
		uc.logger.Error("Failed to get all products from MySQL", zap.Error(err))
		return 0, fmt.Errorf("failed to get products from MySQL: %w", err)
	}

	if len(products) == 0 {
		uc.logger.Info("No products to sync")
		return 0, nil
	}

	// Elasticsearchに1件ずつ同期
	successCount := 0
	failedCount := 0

	for _, product := range products {
		if err := uc.esRepo.Update(ctx, product); err != nil {
			uc.logger.Warn("Failed to sync product to Elasticsearch",
				zap.Error(err),
				zap.String("product_id", product.ID),
			)
			failedCount++
			continue
		}
		successCount++
	}

	uc.logger.Info("Elasticsearch sync completed",
		zap.Int("total", len(products)),
		zap.Int("success", successCount),
		zap.Int("failed", failedCount),
	)

	if failedCount > 0 {
		return successCount, fmt.Errorf("sync completed with %d failures out of %d products", failedCount, len(products))
	}

	return successCount, nil
}
