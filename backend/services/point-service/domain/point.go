package domain

import "time"

// TransactionType ポイントトランザクションタイプ
type TransactionType string

const (
	TransactionTypeEarn   TransactionType = "earn"
	TransactionTypeUse    TransactionType = "use"
	TransactionTypeExpire TransactionType = "expire"
	TransactionTypeRefund TransactionType = "refund"
)

// PointBalance ポイント残高
type PointBalance struct {
	UserID          string    `gorm:"primaryKey" json:"user_id"`
	TotalPoints     int       `json:"total_points"`
	AvailablePoints int       `json:"available_points"`
	ReservedPoints  int       `json:"reserved_points"`
	ExpiredPoints   int       `json:"expired_points"`
	LastUpdated     time.Time `json:"last_updated"`
}

// TableName テーブル名を指定
func (PointBalance) TableName() string {
	return "point_balances"
}

// PointTransaction ポイントトランザクション
type PointTransaction struct {
	ID             string          `gorm:"primaryKey" json:"id"`
	UserID         string          `json:"user_id" gorm:"index"`
	Type           TransactionType `json:"type"`
	Points         int             `json:"points"`
	Balance        int             `json:"balance"`
	Reason         string          `json:"reason"`
	RelatedOrderID string          `json:"related_order_id,omitempty"`
	ExpiresAt      *time.Time      `json:"expires_at,omitempty"`
	CreatedAt      time.Time       `json:"created_at" gorm:"index"`
}

// TableName テーブル名を指定
func (PointTransaction) TableName() string {
	return "point_transactions"
}

// EarnPointsRequest ポイント付与リクエスト
type EarnPointsRequest struct {
	Points         int        `json:"points"`
	Reason         string     `json:"reason"`
	RelatedOrderID string     `json:"related_order_id,omitempty"`
	ExpiresAt      *time.Time `json:"expires_at,omitempty"`
}

// UsePointsRequest ポイント使用リクエスト
type UsePointsRequest struct {
	Points         int    `json:"points"`
	Reason         string `json:"reason"`
	RelatedOrderID string `json:"related_order_id,omitempty"`
}
