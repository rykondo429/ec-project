package domain

import "time"

// CouponStatus クーポンステータス
type CouponStatus string

const (
	CouponStatusActive    CouponStatus = "active"
	CouponStatusExpired   CouponStatus = "expired"
	CouponStatusUsed      CouponStatus = "used"
	CouponStatusSuspended CouponStatus = "suspended"
)

// Coupon クーポンドメインモデル
type Coupon struct {
	ID                   string       `gorm:"primaryKey" json:"id"`
	Code                 string       `gorm:"type:varchar(100);uniqueIndex" json:"code"`
	Status               CouponStatus `json:"status" gorm:"index"`
	DiscountType         string       `json:"discount_type"` // percentage, fixed_amount
	DiscountValue        float64      `json:"discount_value"`
	MaxDiscount          float64      `json:"max_discount,omitempty"`
	MinOrderAmount       float64      `json:"min_order_amount,omitempty"`
	StartDate            time.Time    `json:"start_date"`
	EndDate              time.Time    `json:"end_date"`
	MaxUsageCount        int          `json:"max_usage_count"`
	CurrentUsageCount    int          `json:"current_usage_count"`
	MaxPerUser           int          `json:"max_per_user"`
	ApplicableCategories string       `json:"applicable_categories,omitempty" gorm:"type:json"` // JSON カラムとして保存
	Description          string       `json:"description"`
	CreatedAt            time.Time    `json:"created_at" gorm:"index"`
}

// TableName テーブル名を指定
func (Coupon) TableName() string {
	return "coupons"
}

// Sale セール情報
type Sale struct {
	ID                   string    `gorm:"primaryKey" json:"id"`
	Name                 string    `json:"name"`
	Description          string    `json:"description"`
	StartDate            time.Time `json:"start_date"`
	EndDate              time.Time `json:"end_date" gorm:"index"`
	DiscountType         string    `json:"discount_type"` // percentage, fixed_amount
	DiscountValue        float64   `json:"discount_value"`
	MaxDiscount          float64   `json:"max_discount,omitempty"`
	ApplicableCategories string    `json:"applicable_categories" gorm:"type:json"` // JSON カラムとして保存
	ApplicableProducts   string    `json:"applicable_products" gorm:"type:json"`   // JSON カラムとして保存
	CreatedAt            time.Time `json:"created_at" gorm:"index"`
}

// TableName テーブル名を指定
func (Sale) TableName() string {
	return "sales"
}

// ValidateCouponRequest クーポン検証リクエスト
type ValidateCouponRequest struct {
	Code       string  `json:"code"`
	OrderTotal float64 `json:"order_total"`
	UserID     string  `json:"user_id"`
	Items      []struct {
		ProductID string `json:"product_id"`
		Category  string `json:"category"`
	} `json:"items"`
}

// CouponValidationResult クーポン検証結果
type CouponValidationResult struct {
	Valid          bool    `json:"valid"`
	DiscountAmount float64 `json:"discount_amount"`
	Message        string  `json:"message"`
	Coupon         *Coupon `json:"coupon,omitempty"`
}
