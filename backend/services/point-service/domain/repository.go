package domain

import (
	"context"
)

// PointRepository ポイントリポジトリインターフェース
type PointRepository interface {
	GetBalance(ctx context.Context, userID string) (*PointBalance, error)
	SaveBalance(ctx context.Context, balance *PointBalance) error
	CreateTransaction(ctx context.Context, transaction *PointTransaction) error
	GetTransactions(ctx context.Context, userID string, page, pageSize int) ([]*PointTransaction, int64, error)
	GetTransactionByID(ctx context.Context, id string) (*PointTransaction, error)
}
