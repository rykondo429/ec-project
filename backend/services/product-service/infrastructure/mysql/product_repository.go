package mysql

import (
	"context"
	"fmt"

	"ec-sample/services/product-service/domain"

	"gorm.io/gorm"
)

// MySQLRepository MySQL実装
type MySQLRepository struct {
	db *gorm.DB
}

// NewMySQLRepository コンストラクタ
func NewMySQLRepository(db *gorm.DB) domain.ProductRepository {
	return &MySQLRepository{db: db}
}

// GetByID IDで商品を取得
func (r *MySQLRepository) GetByID(ctx context.Context, id string) (*domain.Product, error) {
	var product domain.Product
	if err := r.db.WithContext(ctx).Where("id = ? AND is_active = ?", id, true).First(&product).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("product not found")
		}
		return nil, fmt.Errorf("failed to get product: %w", err)
	}
	return &product, nil
}

// Search 検索（MySQLのFULLTEXTインデックス使用）
func (r *MySQLRepository) Search(ctx context.Context, query *domain.SearchQuery) (*domain.SearchResult, error) {
	var products []*domain.Product
	var total int64

	db := r.db.WithContext(ctx).Model(&domain.Product{}).Where("is_active = ?", true)

	// キーワード検索（FULLTEXT）
	if query.Keyword != "" {
		db = db.Where("MATCH(name, description) AGAINST(? IN NATURAL LANGUAGE MODE)", query.Keyword)
	}

	// カテゴリフィルタ
	if query.Category != "" {
		db = db.Where("category = ?", query.Category)
	}

	// 価格範囲フィルタ
	if query.PriceMin > 0 {
		db = db.Where("price >= ?", query.PriceMin)
	}
	if query.PriceMax > 0 {
		db = db.Where("price <= ?", query.PriceMax)
	}

	// 総件数取得
	if err := db.Count(&total).Error; err != nil {
		return nil, fmt.Errorf("failed to count products: %w", err)
	}

	// ソート
	orderClause := "created_at DESC"
	if query.SortBy != "" {
		order := "ASC"
		if query.SortOrder == "desc" {
			order = "DESC"
		}
		orderClause = fmt.Sprintf("%s %s", query.SortBy, order)
	}
	db = db.Order(orderClause)

	// ページネーション
	offset := (query.Page - 1) * query.PageSize
	if err := db.Offset(offset).Limit(query.PageSize).Find(&products).Error; err != nil {
		return nil, fmt.Errorf("failed to search products: %w", err)
	}

	totalPages := int(total) / query.PageSize
	if int(total)%query.PageSize > 0 {
		totalPages++
	}

	return &domain.SearchResult{
		Products:    products,
		Total:       total,
		CurrentPage: query.Page,
		TotalPages:  totalPages,
		PageSize:    query.PageSize,
	}, nil
}

// GetByCategory カテゴリで取得
func (r *MySQLRepository) GetByCategory(ctx context.Context, category string, page, pageSize int) (*domain.SearchResult, error) {
	query := &domain.SearchQuery{
		Category: category,
		Page:     page,
		PageSize: pageSize,
	}
	return r.Search(ctx, query)
}

// GetCategories ユニークなカテゴリー一覧を取得
func (r *MySQLRepository) GetCategories(ctx context.Context) ([]string, error) {
	var categories []string
	if err := r.db.WithContext(ctx).
		Model(&domain.Product{}).
		Where("is_active = ? AND category != ''", true).
		Distinct("category").
		Pluck("category", &categories).Error; err != nil {
		return nil, fmt.Errorf("failed to get categories: %w", err)
	}
	return categories, nil
}

// GetFeatured おすすめ商品取得
func (r *MySQLRepository) GetFeatured(ctx context.Context) ([]*domain.Product, error) {
	var products []*domain.Product
	if err := r.db.WithContext(ctx).
		Where("is_active = ?", true).
		Order("rating DESC, review_count DESC").
		Limit(10).
		Find(&products).Error; err != nil {
		return nil, fmt.Errorf("failed to get featured products: %w", err)
	}
	return products, nil
}

// Create 商品を作成
func (r *MySQLRepository) Create(ctx context.Context, product *domain.Product) error {
	if err := r.db.WithContext(ctx).Create(product).Error; err != nil {
		return fmt.Errorf("failed to create product: %w", err)
	}
	return nil
}

// Update 商品を更新
func (r *MySQLRepository) Update(ctx context.Context, product *domain.Product) error {
	if err := r.db.WithContext(ctx).Save(product).Error; err != nil {
		return fmt.Errorf("failed to update product: %w", err)
	}
	return nil
}

// GetAll 全商品取得（ES同期用）
func (r *MySQLRepository) GetAll(ctx context.Context) ([]*domain.Product, error) {
	var products []*domain.Product
	if err := r.db.WithContext(ctx).Find(&products).Error; err != nil {
		return nil, fmt.Errorf("failed to get all products: %w", err)
	}
	return products, nil
}
