package domain

import "time"

// Product ドメインモデル
type Product struct {
	ID          string    `json:"id" gorm:"primaryKey;type:varchar(36)"`
	Name        string    `json:"name" gorm:"type:varchar(255);not null"`
	Description string    `json:"description" gorm:"type:longtext"`
	Price       float64   `json:"price" gorm:"type:decimal(10,2);not null"`
	Stock       int       `json:"stock" gorm:"not null;default:0"`
	Category    string    `json:"category" gorm:"type:varchar(100);index"`
	Images      []string  `json:"images" gorm:"-"` // JSONで管理、GORMでは保存しない
	Rating      float64   `json:"rating" gorm:"type:decimal(3,2);default:0.0"`
	ReviewCount int       `json:"review_count" gorm:"default:0"`
	IsSale      bool      `json:"is_sale" gorm:"default:false"`
	SalePrice   *float64  `json:"sale_price,omitempty" gorm:"type:decimal(10,2)"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	IsActive    bool      `json:"is_active" gorm:"default:true;index"`
}

// TableName テーブル名を明示
func (Product) TableName() string {
	return "products"
}

// SearchQuery 検索条件
type SearchQuery struct {
	Keyword   string
	Category  string
	PriceMin  float64
	PriceMax  float64
	SortBy    string
	SortOrder string
	Page      int
	PageSize  int
}

// SearchResult 検索結果
type SearchResult struct {
	Products    []*Product `json:"products"`
	Total       int64      `json:"total"`
	CurrentPage int        `json:"current_page"`
	TotalPages  int        `json:"total_pages"`
	PageSize    int        `json:"page_size"`
}
