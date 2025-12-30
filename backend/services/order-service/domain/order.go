package domain

import "time"

// OrderStatus 注文ステータス
type OrderStatus string

const (
	OrderStatusPending    OrderStatus = "pending"
	OrderStatusProcessing OrderStatus = "processing"
	OrderStatusCompleted  OrderStatus = "completed"
	OrderStatusCancelled  OrderStatus = "cancelled"
)

// OrderItem 注文アイテム
type OrderItem struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	OrderID     string    `gorm:"type:varchar(36);index" json:"order_id"`
	ProductID   string    `json:"product_id"`
	ProductName string    `json:"product_name"`
	Price       float64   `json:"price"`
	Quantity    int       `json:"quantity"`
	Subtotal    float64   `json:"subtotal"`
	CreatedAt   time.Time `json:"created_at"`
}

// Order 注文ドメインモデル
type Order struct {
	ID              string      `gorm:"primaryKey" json:"id"`
	UserID          string      `json:"user_id" gorm:"index"`
	Items           []OrderItem `json:"items" gorm:"foreignKey:OrderID"`
	Subtotal        float64     `json:"subtotal"`
	ShippingCost    float64     `json:"shipping_cost"`
	Tax             float64     `json:"tax"`
	DiscountAmount  float64     `json:"discount_amount"`
	Total           float64     `json:"total"`
	Status          OrderStatus `json:"status" gorm:"index"`
	PaymentMethod   string      `json:"payment_method"`
	ShippingAddress string      `json:"shipping_address"`
	PointsUsed      int         `json:"points_used"`
	PointsEarned    int         `json:"points_earned"`
	CouponCode      string      `json:"coupon_code"`
	CreatedAt       time.Time   `json:"created_at" gorm:"index"`
	UpdatedAt       time.Time   `json:"updated_at"`
}

// TableName テーブル名を指定
func (Order) TableName() string {
	return "orders"
}

// TableName テーブル名を指定
func (OrderItem) TableName() string {
	return "order_items"
}

// CreateOrderRequest 注文作成リクエスト
type CreateOrderRequest struct {
	Items           []*OrderItem `json:"items"`
	PaymentMethod   string       `json:"payment_method"`
	ShippingAddress string       `json:"shipping_address"`
	PointsToUse     int          `json:"points_to_use"`
	CouponCode      string       `json:"coupon_code"`
}
