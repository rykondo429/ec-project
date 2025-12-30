package usecase

import (
	"context"
	"fmt"
	"math"
	"time"

	"ec-sample/services/order-service/domain"

	"github.com/google/uuid"
)

const (
	shippingCost = 500.0
	taxRate      = 0.1
	pointsPerYen = 1
)

// OrderUseCase 注文ユースケース
type OrderUseCase struct {
	repo domain.OrderRepository
}

// NewOrderUseCase コンストラクタ
func NewOrderUseCase(repo domain.OrderRepository) *OrderUseCase {
	return &OrderUseCase{repo: repo}
}

// CreateOrder 注文作成
func (uc *OrderUseCase) CreateOrder(ctx context.Context, userID string, req *domain.CreateOrderRequest) (*domain.Order, error) {
	if userID == "" {
		return nil, fmt.Errorf("user_id is required")
	}
	if len(req.Items) == 0 {
		return nil, fmt.Errorf("items are required")
	}

	// 小計計算
	subtotal := 0.0
	for _, item := range req.Items {
		subtotal += item.Subtotal
	}

	// 割引計算 (クーポンは簡易実装)
	discountAmount := 0.0
	if req.CouponCode != "" {
		// 実装時: クーポンサービスと連携して割引額を取得
		discountAmount = subtotal * 0.1 // 例: 10% 割引
	}

	// ポイント使用による割引
	pointsDiscount := float64(req.PointsToUse) / 100.0 // 100ポイント = 1円

	// 税計算
	taxableAmount := subtotal - discountAmount - pointsDiscount
	tax := math.Round(taxableAmount*taxRate*100) / 100

	// 合計
	total := subtotal - discountAmount - pointsDiscount + shippingCost + tax

	// ポイント付与計算
	pointsEarned := int(subtotal * pointsPerYen)

	// OrderItem の型変換（ポインタスライス → 値スライス）
	items := make([]domain.OrderItem, len(req.Items))
	for i, item := range req.Items {
		items[i] = *item
	}

	order := &domain.Order{
		ID:              uuid.New().String(),
		UserID:          userID,
		Items:           items,
		Subtotal:        subtotal,
		ShippingCost:    shippingCost,
		Tax:             tax,
		DiscountAmount:  discountAmount,
		Total:           math.Round(total*100) / 100,
		Status:          domain.OrderStatusPending,
		PaymentMethod:   req.PaymentMethod,
		ShippingAddress: req.ShippingAddress,
		PointsUsed:      req.PointsToUse,
		PointsEarned:    pointsEarned,
		CouponCode:      req.CouponCode,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	// 注文保存
	if err := uc.repo.Create(ctx, order); err != nil {
		return nil, err
	}

	return order, nil
}

// GetOrder 注文取得
func (uc *OrderUseCase) GetOrder(ctx context.Context, id string) (*domain.Order, error) {
	if id == "" {
		return nil, fmt.Errorf("order_id is required")
	}

	return uc.repo.GetByID(ctx, id)
}

// GetUserOrders ユーザーの注文履歴取得
func (uc *OrderUseCase) GetUserOrders(ctx context.Context, userID string, page, pageSize int) ([]*domain.Order, int64, error) {
	if userID == "" {
		return nil, 0, fmt.Errorf("user_id is required")
	}

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}

	return uc.repo.GetByUserID(ctx, userID, page, pageSize)
}

// UpdateOrderStatus 注文ステータス更新
func (uc *OrderUseCase) UpdateOrderStatus(ctx context.Context, id string, status domain.OrderStatus) (*domain.Order, error) {
	if id == "" {
		return nil, fmt.Errorf("order_id is required")
	}

	// ステータス検証
	validStatuses := map[domain.OrderStatus]bool{
		domain.OrderStatusPending:    true,
		domain.OrderStatusProcessing: true,
		domain.OrderStatusCompleted:  true,
		domain.OrderStatusCancelled:  true,
	}

	if !validStatuses[status] {
		return nil, fmt.Errorf("invalid status: %s", status)
	}

	if err := uc.repo.UpdateStatus(ctx, id, status); err != nil {
		return nil, err
	}

	return uc.repo.GetByID(ctx, id)
}
