package domain

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

// SKU (Stock Keeping Unit) - 商品バリエーション管理
type SKU struct {
	ID         string    `json:"id" gorm:"primaryKey;type:varchar(36)"`
	ProductID  string    `json:"product_id" gorm:"type:varchar(36);not null;index:idx_product_sku"`
	SKUCode    string    `json:"sku_code" gorm:"type:varchar(100);uniqueIndex;not null"` // 例: NIKE-AIR-RED-M
	Name       string    `json:"name" gorm:"type:varchar(255);not null"`                 // 例: Nike Air Max 90 - Red / M
	Attributes SKUAttrs  `json:"attributes" gorm:"type:json"`                            // サイズ、色など
	Price      float64   `json:"price" gorm:"type:decimal(10,2);not null"`               // SKU固有価格（オプション）
	Stock      int       `json:"stock" gorm:"not null;default:0"`
	IsActive   bool      `json:"is_active" gorm:"default:true;index"`
	CreatedAt  time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt  time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// SKUAttrs SKU属性（サイズ、色、その他のバリエーション）
type SKUAttrs map[string]string // 例: {"size": "M", "color": "Red"}

// Scan implements sql.Scanner interface for SKUAttrs
func (a *SKUAttrs) Scan(value interface{}) error {
	if value == nil {
		*a = make(SKUAttrs)
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to unmarshal SKUAttrs value: %v", value)
	}

	result := make(SKUAttrs)
	if err := json.Unmarshal(bytes, &result); err != nil {
		return err
	}

	*a = result
	return nil
}

// Value implements driver.Valuer interface for SKUAttrs
func (a SKUAttrs) Value() (driver.Value, error) {
	if a == nil {
		return nil, nil
	}
	return json.Marshal(a)
}

// TableName テーブル名を明示
func (SKU) TableName() string {
	return "skus"
}

// SKUCreateRequest SKU作成リクエスト
type SKUCreateRequest struct {
	ProductID  string   `json:"product_id" binding:"required"`
	SKUCode    string   `json:"sku_code" binding:"required"`
	Name       string   `json:"name" binding:"required"`
	Attributes SKUAttrs `json:"attributes"`
	Price      float64  `json:"price"`
	Stock      int      `json:"stock"`
}

// SKUUpdateRequest SKU更新リクエスト
type SKUUpdateRequest struct {
	Name       *string   `json:"name"`
	Attributes *SKUAttrs `json:"attributes"`
	Price      *float64  `json:"price"`
	Stock      *int      `json:"stock"`
	IsActive   *bool     `json:"is_active"`
}

// SKUStockAdjustment 在庫調整
type SKUStockAdjustment struct {
	SKUID    string `json:"sku_id" binding:"required"`
	Quantity int    `json:"quantity" binding:"required"` // 正数=増加、負数=減少
	Reason   string `json:"reason"`                      // 調整理由
}

// ProductWithSKUs 商品とそのSKU一覧
type ProductWithSKUs struct {
	*Product
	SKUs []*SKU `json:"skus"`
}
