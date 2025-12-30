package domain

import (
	"context"
	"time"
)

// ProductRepository リポジトリインターフェース
type ProductRepository interface {
	GetByID(ctx context.Context, id string) (*Product, error)
	Search(ctx context.Context, query *SearchQuery) (*SearchResult, error)
	GetByCategory(ctx context.Context, category string, page, pageSize int) (*SearchResult, error)
	GetCategories(ctx context.Context) ([]string, error)
	GetFeatured(ctx context.Context) ([]*Product, error)
	GetAll(ctx context.Context) ([]*Product, error)
	Create(ctx context.Context, product *Product) error
	Update(ctx context.Context, product *Product) error
}

// CacheRepository Redisキャッシュリポジトリインターフェース
type CacheRepository interface {
	Get(ctx context.Context, key string, dest interface{}) error
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Delete(ctx context.Context, key string) error
	DeletePattern(ctx context.Context, pattern string) error
	Ping(ctx context.Context) error
}
