package usecase

import (
	"context"
	"fmt"
	"math"
	"time"

	"ec-sample/services/cart-service/domain"
)

// CartUseCase カートユースケース
type CartUseCase struct {
	repo domain.CartRepository
}

// NewCartUseCase コンストラクタ
func NewCartUseCase(repo domain.CartRepository) *CartUseCase {
	return &CartUseCase{repo: repo}
}

// GetCart カート取得
func (uc *CartUseCase) GetCart(ctx context.Context, userID string) (*domain.Cart, error) {
	if userID == "" {
		return nil, fmt.Errorf("user_id is required")
	}

	cart, err := uc.repo.GetCart(ctx, userID)
	if err != nil {
		return nil, err
	}

	// カート有効期限チェック
	if cart != nil && cart.ExpiresAt.Before(time.Now()) {
		// 期限切れカートはクリア
		_ = uc.repo.ClearCart(ctx, userID)
		return &domain.Cart{
			UserID:       userID,
			Items:        []*domain.CartItem{},
			Total:        0,
			ItemCount:    0,
			LastModified: time.Now(),
			ExpiresAt:    time.Now().Add(7 * 24 * time.Hour),
		}, nil
	}

	return cart, nil
}

// AddItem アイテム追加
func (uc *CartUseCase) AddItem(ctx context.Context, userID string, req *domain.AddItemRequest) (*domain.Cart, error) {
	if userID == "" {
		return nil, fmt.Errorf("user_id is required")
	}
	if req.ProductID == "" {
		return nil, fmt.Errorf("product_id is required")
	}
	if req.Quantity <= 0 {
		return nil, fmt.Errorf("quantity must be greater than 0")
	}
	if req.Price < 0 {
		return nil, fmt.Errorf("price must be non-negative")
	}

	cart, err := uc.repo.GetCart(ctx, userID)
	if err != nil {
		return nil, err
	}

	// 新規カート作成
	if cart == nil {
		cart = &domain.Cart{
			UserID:       userID,
			Items:        []*domain.CartItem{},
			Total:        0,
			ItemCount:    0,
			LastModified: time.Now(),
			ExpiresAt:    time.Now().Add(7 * 24 * time.Hour),
		}
	}

	// 既存アイテム確認（同じ商品IDとSKU IDの組み合わせ）
	itemKey := req.ProductID
	if req.SKUID != "" {
		itemKey = req.ProductID + "_" + req.SKUID
	}

	itemFound := false
	for i, item := range cart.Items {
		existingKey := item.ProductID
		if item.SKUID != "" {
			existingKey = item.ProductID + "_" + item.SKUID
		}

		if existingKey == itemKey {
			// 数量更新
			cart.Items[i].Quantity += req.Quantity
			cart.Items[i].Subtotal = math.Round(float64(cart.Items[i].Quantity)*req.Price*100) / 100
			cart.Items[i].AddedAt = time.Now()
			itemFound = true
			break
		}
	}

	// 新規アイテム追加
	if !itemFound {
		newItem := &domain.CartItem{
			ProductID:     req.ProductID,
			ProductName:   req.ProductName,
			SKUID:         req.SKUID,
			SKUCode:       req.SKUCode,
			SKUAttributes: req.SKUAttributes,
			Price:         req.Price,
			Quantity:      req.Quantity,
			Subtotal:      math.Round(float64(req.Quantity)*req.Price*100) / 100,
			AddedAt:       time.Now(),
		}
		cart.Items = append(cart.Items, newItem)
	}

	// 合計計算
	uc.calculateCartTotal(cart)

	// 保存
	if err := uc.repo.SaveCart(ctx, cart); err != nil {
		return nil, err
	}

	return cart, nil
}

// UpdateItem アイテム更新
func (uc *CartUseCase) UpdateItem(ctx context.Context, userID, productID string, quantity int) (*domain.Cart, error) {
	if userID == "" {
		return nil, fmt.Errorf("user_id is required")
	}
	if productID == "" {
		return nil, fmt.Errorf("product_id is required")
	}

	if quantity < 0 {
		return nil, fmt.Errorf("quantity must be non-negative")
	}

	cart, err := uc.repo.GetCart(ctx, userID)
	if err != nil {
		return nil, err
	}

	if cart == nil {
		return nil, fmt.Errorf("cart not found")
	}

	// 数量が0の場合は削除
	if quantity == 0 {
		return uc.RemoveItem(ctx, userID, productID)
	}

	// アイテム更新
	found := false
	for i, item := range cart.Items {
		if item.ProductID == productID {
			cart.Items[i].Quantity = quantity
			cart.Items[i].Subtotal = math.Round(float64(quantity)*item.Price*100) / 100
			cart.Items[i].AddedAt = time.Now()
			found = true
			break
		}
	}

	if !found {
		return nil, fmt.Errorf("product not found in cart")
	}

	// 合計計算
	uc.calculateCartTotal(cart)

	// 保存
	if err := uc.repo.SaveCart(ctx, cart); err != nil {
		return nil, err
	}

	return cart, nil
}

// RemoveItem アイテム削除
func (uc *CartUseCase) RemoveItem(ctx context.Context, userID, productID string) (*domain.Cart, error) {
	if userID == "" {
		return nil, fmt.Errorf("user_id is required")
	}
	if productID == "" {
		return nil, fmt.Errorf("product_id is required")
	}

	cart, err := uc.repo.GetCart(ctx, userID)
	if err != nil {
		return nil, err
	}

	if cart == nil {
		return nil, fmt.Errorf("cart not found")
	}

	// アイテム削除
	for i, item := range cart.Items {
		if item.ProductID == productID {
			cart.Items = append(cart.Items[:i], cart.Items[i+1:]...)
			break
		}
	}

	// 合計計算
	uc.calculateCartTotal(cart)

	// 保存
	if err := uc.repo.SaveCart(ctx, cart); err != nil {
		return nil, err
	}

	return cart, nil
}

// ClearCart カートクリア
func (uc *CartUseCase) ClearCart(ctx context.Context, userID string) error {
	if userID == "" {
		return fmt.Errorf("user_id is required")
	}

	return uc.repo.ClearCart(ctx, userID)
}

// calculateCartTotal カート合計計算
func (uc *CartUseCase) calculateCartTotal(cart *domain.Cart) {
	total := 0.0
	itemCount := 0

	for _, item := range cart.Items {
		total += item.Subtotal
		itemCount += item.Quantity
	}

	cart.Total = math.Round(total*100) / 100
	cart.ItemCount = itemCount
	cart.LastModified = time.Now()
}
