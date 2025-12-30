package usecase

import (
	"context"
	"fmt"
	"time"

	"ec-sample/services/point-service/domain"

	"github.com/google/uuid"
)

// PointUseCase ポイントユースケース
type PointUseCase struct {
	repo domain.PointRepository
}

// NewPointUseCase コンストラクタ
func NewPointUseCase(repo domain.PointRepository) *PointUseCase {
	return &PointUseCase{repo: repo}
}

// GetBalance ポイント残高取得
func (uc *PointUseCase) GetBalance(ctx context.Context, userID string) (*domain.PointBalance, error) {
	if userID == "" {
		return nil, fmt.Errorf("user_id is required")
	}

	balance, err := uc.repo.GetBalance(ctx, userID)
	if err != nil {
		return nil, err
	}

	// 未登録の場合は初期化
	if balance == nil {
		balance = &domain.PointBalance{
			UserID:          userID,
			TotalPoints:     0,
			AvailablePoints: 0,
			ReservedPoints:  0,
			ExpiredPoints:   0,
			LastUpdated:     time.Now(),
		}
		_ = uc.repo.SaveBalance(ctx, balance)
	}

	return balance, nil
}

// EarnPoints ポイント付与
func (uc *PointUseCase) EarnPoints(ctx context.Context, userID string, req *domain.EarnPointsRequest) (*domain.PointBalance, error) {
	if userID == "" {
		return nil, fmt.Errorf("user_id is required")
	}
	if req.Points <= 0 {
		return nil, fmt.Errorf("points must be greater than 0")
	}
	if req.Reason == "" {
		return nil, fmt.Errorf("reason is required")
	}

	// 残高取得
	balance, err := uc.GetBalance(ctx, userID)
	if err != nil {
		return nil, err
	}

	// 残高更新
	balance.TotalPoints += req.Points
	balance.AvailablePoints += req.Points
	balance.LastUpdated = time.Now()

	// トランザクション作成
	transaction := &domain.PointTransaction{
		ID:             uuid.New().String(),
		UserID:         userID,
		Type:           domain.TransactionTypeEarn,
		Points:         req.Points,
		Balance:        balance.AvailablePoints,
		Reason:         req.Reason,
		RelatedOrderID: req.RelatedOrderID,
		ExpiresAt:      req.ExpiresAt,
		CreatedAt:      time.Now(),
	}

	// トランザクション保存
	if err := uc.repo.CreateTransaction(ctx, transaction); err != nil {
		return nil, err
	}

	// 残高保存
	if err := uc.repo.SaveBalance(ctx, balance); err != nil {
		return nil, err
	}

	return balance, nil
}

// UsePoints ポイント使用
func (uc *PointUseCase) UsePoints(ctx context.Context, userID string, req *domain.UsePointsRequest) (*domain.PointBalance, error) {
	if userID == "" {
		return nil, fmt.Errorf("user_id is required")
	}
	if req.Points <= 0 {
		return nil, fmt.Errorf("points must be greater than 0")
	}
	if req.Reason == "" {
		return nil, fmt.Errorf("reason is required")
	}

	// 残高取得
	balance, err := uc.GetBalance(ctx, userID)
	if err != nil {
		return nil, err
	}

	// ポイント不足チェック
	if balance.AvailablePoints < req.Points {
		return nil, fmt.Errorf("insufficient points")
	}

	// 残高更新
	balance.AvailablePoints -= req.Points
	balance.ReservedPoints += req.Points
	balance.LastUpdated = time.Now()

	// トランザクション作成
	transaction := &domain.PointTransaction{
		ID:             uuid.New().String(),
		UserID:         userID,
		Type:           domain.TransactionTypeUse,
		Points:         req.Points,
		Balance:        balance.AvailablePoints,
		Reason:         req.Reason,
		RelatedOrderID: req.RelatedOrderID,
		CreatedAt:      time.Now(),
	}

	// トランザクション保存
	if err := uc.repo.CreateTransaction(ctx, transaction); err != nil {
		return nil, err
	}

	// 残高保存
	if err := uc.repo.SaveBalance(ctx, balance); err != nil {
		return nil, err
	}

	return balance, nil
}

// GetTransactions トランザクション履歴取得
func (uc *PointUseCase) GetTransactions(ctx context.Context, userID string, page, pageSize int) ([]*domain.PointTransaction, int64, error) {
	if userID == "" {
		return nil, 0, fmt.Errorf("user_id is required")
	}

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}

	return uc.repo.GetTransactions(ctx, userID, page, pageSize)
}
