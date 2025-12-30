package usecase

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"time"

	"ec-sample/services/product-service/domain"

	"go.uber.org/zap"
)

// ProductUseCase ユースケース
type ProductUseCase struct {
	repo      domain.ProductRepository // MySQL (主データソース)
	esRepo    domain.ProductRepository // Elasticsearch (検索インデックス)
	cacheRepo domain.CacheRepository   // Redis (キャッシュ)
	logger    *zap.Logger
}

// NewProductUseCase コンストラクタ
func NewProductUseCase(repo domain.ProductRepository, esRepo domain.ProductRepository, cacheRepo domain.CacheRepository, logger *zap.Logger) *ProductUseCase {
	return &ProductUseCase{
		repo:      repo,
		esRepo:    esRepo,
		cacheRepo: cacheRepo,
		logger:    logger,
	}
}

// GetProductByID 商品詳細取得（キャッシュ優先）
func (uc *ProductUseCase) GetProductByID(ctx context.Context, id string) (*domain.Product, error) {
	if id == "" {
		return nil, fmt.Errorf("product id is required")
	}

	// キャッシュキー生成
	cacheKey := fmt.Sprintf("product:%s", id)

	// キャッシュから取得試行
	var cachedProduct domain.Product
	if uc.cacheRepo != nil {
		if err := uc.cacheRepo.Get(ctx, cacheKey, &cachedProduct); err == nil {
			uc.logger.Debug("Product cache hit", zap.String("product_id", id))
			return &cachedProduct, nil
		}
	}

	// キャッシュミス時、MySQLから取得
	product, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// キャッシュに保存（5分間）
	if uc.cacheRepo != nil {
		if err := uc.cacheRepo.Set(ctx, cacheKey, product, 5*time.Minute); err != nil {
			uc.logger.Warn("Failed to cache product", zap.Error(err), zap.String("product_id", id))
		}
	}

	return product, nil
}

// SearchProducts 商品検索（キャッシュ優先）
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

	// キャッシュキー生成（クエリパラメータのハッシュ）
	cacheKey := uc.generateSearchCacheKey(query)

	// キャッシュから取得試行
	var cachedResult domain.SearchResult
	if uc.cacheRepo != nil {
		if err := uc.cacheRepo.Get(ctx, cacheKey, &cachedResult); err == nil {
			uc.logger.Debug("Search cache hit", zap.String("cache_key", cacheKey))
			return &cachedResult, nil
		}
	}

	// キャッシュミス時、DBから検索
	result, err := uc.repo.Search(ctx, query)
	if err != nil {
		return nil, err
	}

	// キャッシュに保存（1分間）
	if uc.cacheRepo != nil {
		if err := uc.cacheRepo.Set(ctx, cacheKey, result, 1*time.Minute); err != nil {
			uc.logger.Warn("Failed to cache search result", zap.Error(err))
		}
	}

	return result, nil
}

// generateSearchCacheKey 検索クエリからキャッシュキーを生成
func (uc *ProductUseCase) generateSearchCacheKey(query *domain.SearchQuery) string {
	// クエリパラメータを結合してMD5ハッシュ化
	queryStr := fmt.Sprintf("search:%s:%s:%d:%d:%s:%f:%f",
		query.Keyword, query.Category, query.Page, query.PageSize,
		query.SortBy, query.PriceMin, query.PriceMax,
	)
	hash := md5.Sum([]byte(queryStr))
	return fmt.Sprintf("search:%s", hex.EncodeToString(hash[:]))
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

// CreateProduct 商品作成 (MySQL + ES同期 + キャッシュ削除)
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

	// 検索キャッシュとカテゴリキャッシュを無効化
	uc.invalidateSearchCache(ctx)

	uc.logger.Info("Product created successfully", zap.String("product_id", product.ID))
	return nil
}

// UpdateProduct 商品更新 (MySQL + ES同期 + キャッシュ無効化)
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

	// 商品キャッシュと検索キャッシュを無効化
	if uc.cacheRepo != nil {
		// 個別商品キャッシュ削除
		cacheKey := fmt.Sprintf("product:%s", product.ID)
		if err := uc.cacheRepo.Delete(ctx, cacheKey); err != nil {
			uc.logger.Warn("Failed to delete product cache", zap.Error(err), zap.String("product_id", product.ID))
		}
	}
	uc.invalidateSearchCache(ctx)

	uc.logger.Info("Product updated successfully", zap.String("product_id", product.ID))
	return nil
}

// invalidateSearchCache 検索関連のキャッシュを無効化
func (uc *ProductUseCase) invalidateSearchCache(ctx context.Context) {
	if uc.cacheRepo == nil {
		return
	}

	// 検索キャッシュ削除（パターンマッチ）
	if err := uc.cacheRepo.DeletePattern(ctx, "search:*"); err != nil {
		uc.logger.Warn("Failed to invalidate search cache", zap.Error(err))
	}

	// カテゴリキャッシュ削除
	if err := uc.cacheRepo.Delete(ctx, "categories"); err != nil {
		uc.logger.Warn("Failed to invalidate categories cache", zap.Error(err))
	}
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
