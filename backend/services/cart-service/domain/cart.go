package domain

import "time"

// CartItem カートアイテム
type CartItem struct {
	ProductID     string            `json:"product_id"`
	ProductName   string            `json:"product_name"`
	SKUID         string            `json:"sku_id,omitempty"`
	SKUCode       string            `json:"sku_code,omitempty"`
	SKUAttributes map[string]string `json:"sku_attributes,omitempty"`
	Price         float64           `json:"price"`
	Quantity      int               `json:"quantity"`
	Subtotal      float64           `json:"subtotal"`
	AddedAt       time.Time         `json:"added_at"`
}

// Cart カートドメインモデル
type Cart struct {
	UserID       string      `json:"user_id" dynamodbav:"user_id"`
	Items        []*CartItem `json:"items" dynamodbav:"items"`
	Total        float64     `json:"total" dynamodbav:"total"`
	ItemCount    int         `json:"item_count" dynamodbav:"item_count"`
	LastModified time.Time   `json:"last_modified" dynamodbav:"last_modified"`
	ExpiresAt    time.Time   `json:"expires_at" dynamodbav:"expires_at"`
}

// AddItemRequest アイテム追加リクエスト
type AddItemRequest struct {
	ProductID     string            `json:"product_id"`
	ProductName   string            `json:"product_name"`
	SKUID         string            `json:"sku_id,omitempty"`
	SKUCode       string            `json:"sku_code,omitempty"`
	SKUAttributes map[string]string `json:"sku_attributes,omitempty"`
	Price         float64           `json:"price"`
	Quantity      int               `json:"quantity"`
}

// UpdateItemRequest アイテム更新リクエスト
type UpdateItemRequest struct {
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
}
