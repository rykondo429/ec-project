package mysql

import (
	"context"
	"fmt"

	"ec-sample/services/point-service/domain"

	"gorm.io/gorm"
)

// MySQLRepository MySQL実装
type MySQLRepository struct {
	db *gorm.DB
}

// NewMySQLRepository コンストラクタ
func NewMySQLRepository(db *gorm.DB) domain.PointRepository {
	return &MySQLRepository{db: db}
}

// GetBalance 残高取得
func (r *MySQLRepository) GetBalance(ctx context.Context, userID string) (*domain.PointBalance, error) {
	var balance domain.PointBalance
	if err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		First(&balance).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get balance: %w", err)
	}
	return &balance, nil
}

// SaveBalance 残高保存
func (r *MySQLRepository) SaveBalance(ctx context.Context, balance *domain.PointBalance) error {
	if err := r.db.WithContext(ctx).
		Save(balance).Error; err != nil {
		return fmt.Errorf("failed to save balance: %w", err)
	}
	return nil
}

// CreateTransaction トランザクション作成
func (r *MySQLRepository) CreateTransaction(ctx context.Context, transaction *domain.PointTransaction) error {
	if err := r.db.WithContext(ctx).Create(transaction).Error; err != nil {
		return fmt.Errorf("failed to create transaction: %w", err)
	}
	return nil
}

// GetTransactions トランザクション取得
func (r *MySQLRepository) GetTransactions(ctx context.Context, userID string, page, pageSize int) ([]*domain.PointTransaction, int64, error) {
	var transactions []*domain.PointTransaction
	var total int64

	// 総件数取得
	if err := r.db.WithContext(ctx).
		Model(&domain.PointTransaction{}).
		Where("user_id = ?", userID).
		Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to get count: %w", err)
	}

	offset := (page - 1) * pageSize

	// トランザクション取得
	if err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&transactions).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to get transactions: %w", err)
	}

	return transactions, total, nil
}

// GetTransactionByID ID で取得
func (r *MySQLRepository) GetTransactionByID(ctx context.Context, id string) (*domain.PointTransaction, error) {
	var transaction domain.PointTransaction
	if err := r.db.WithContext(ctx).
		Where("id = ?", id).
		First(&transaction).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("transaction not found")
		}
		return nil, fmt.Errorf("failed to get transaction: %w", err)
	}
	return &transaction, nil
}
