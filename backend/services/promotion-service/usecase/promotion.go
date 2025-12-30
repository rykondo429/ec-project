package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"time"

	"ec-sample/services/promotion-service/domain"
)

// PromotionUseCase プロモーションユースケース
type PromotionUseCase struct {
	couponRepo domain.CouponRepository
	saleRepo   domain.SaleRepository
}

// NewPromotionUseCase コンストラクタ
func NewPromotionUseCase(couponRepo domain.CouponRepository, saleRepo domain.SaleRepository) *PromotionUseCase {
	return &PromotionUseCase{
		couponRepo: couponRepo,
		saleRepo:   saleRepo,
	}
}

// ValidateCoupon クーポン検証
func (uc *PromotionUseCase) ValidateCoupon(ctx context.Context, req *domain.ValidateCouponRequest) (*domain.CouponValidationResult, error) {
	if req.Code == "" {
		return &domain.CouponValidationResult{
			Valid:   false,
			Message: "coupon code is required",
		}, nil
	}

	// クーポン取得
	coupon, err := uc.couponRepo.GetByCode(ctx, req.Code)
	if err != nil || coupon == nil {
		return &domain.CouponValidationResult{
			Valid:   false,
			Message: "coupon not found",
		}, nil
	}

	// ステータスチェック
	if coupon.Status != domain.CouponStatusActive {
		return &domain.CouponValidationResult{
			Valid:   false,
			Message: "coupon is not active",
			Coupon:  coupon,
		}, nil
	}

	// 有効期限チェック
	now := time.Now()
	if now.Before(coupon.StartDate) || now.After(coupon.EndDate) {
		return &domain.CouponValidationResult{
			Valid:   false,
			Message: "coupon has expired",
			Coupon:  coupon,
		}, nil
	}

	// 使用回数チェック
	if coupon.MaxUsageCount > 0 && coupon.CurrentUsageCount >= coupon.MaxUsageCount {
		return &domain.CouponValidationResult{
			Valid:   false,
			Message: "coupon usage limit reached",
			Coupon:  coupon,
		}, nil
	}

	// 最低注文金額チェック
	if coupon.MinOrderAmount > 0 && req.OrderTotal < coupon.MinOrderAmount {
		return &domain.CouponValidationResult{
			Valid:   false,
			Message: fmt.Sprintf("minimum order amount is %.2f", coupon.MinOrderAmount),
			Coupon:  coupon,
		}, nil
	}

	// 適用カテゴリチェック
	if len(coupon.ApplicableCategories) > 0 {
		var categories []string
		if err := json.Unmarshal([]byte(coupon.ApplicableCategories), &categories); err == nil {
			applicable := false
			for _, cat := range categories {
				for _, item := range req.Items {
					if item.Category == cat {
						applicable = true
						break
					}
				}
				if applicable {
					break
				}
			}
			if !applicable {
				return &domain.CouponValidationResult{
					Valid:   false,
					Message: "coupon not applicable to these items",
					Coupon:  coupon,
				}, nil
			}
		}
	}

	// 割引額計算
	discountAmount := 0.0
	if coupon.DiscountType == "percentage" {
		discountAmount = math.Round(req.OrderTotal*coupon.DiscountValue/100*100) / 100
		if coupon.MaxDiscount > 0 && discountAmount > coupon.MaxDiscount {
			discountAmount = coupon.MaxDiscount
		}
	} else if coupon.DiscountType == "fixed_amount" {
		discountAmount = coupon.DiscountValue
		if discountAmount > req.OrderTotal {
			discountAmount = req.OrderTotal
		}
	}

	return &domain.CouponValidationResult{
		Valid:          true,
		DiscountAmount: discountAmount,
		Message:        "coupon is valid",
		Coupon:         coupon,
	}, nil
}

// GetActiveCoupons アクティブなクーポン取得
func (uc *PromotionUseCase) GetActiveCoupons(ctx context.Context) ([]*domain.Coupon, error) {
	return uc.couponRepo.GetActiveCoupons(ctx)
}

// GetActiveSales アクティブなセール取得
func (uc *PromotionUseCase) GetActiveSales(ctx context.Context) ([]*domain.Sale, error) {
	return uc.saleRepo.GetActiveSales(ctx)
}

// GetSalesByCategory カテゴリ別セール取得
func (uc *PromotionUseCase) GetSalesByCategory(ctx context.Context, category string) ([]*domain.Sale, error) {
	if category == "" {
		return nil, fmt.Errorf("category is required")
	}
	return uc.saleRepo.GetSalesByCategory(ctx, category)
}

// GetSalesByProduct 商品別セール取得
func (uc *PromotionUseCase) GetSalesByProduct(ctx context.Context, productID string) ([]*domain.Sale, error) {
	if productID == "" {
		return nil, fmt.Errorf("product_id is required")
	}
	return uc.saleRepo.GetSalesByProduct(ctx, productID)
}

// CalculateDiscount 割引計算
func (uc *PromotionUseCase) CalculateDiscount(sale *domain.Sale, basePrice float64) float64 {
	if sale == nil {
		return 0
	}

	discountAmount := 0.0
	if sale.DiscountType == "percentage" {
		discountAmount = math.Round(basePrice*sale.DiscountValue/100*100) / 100
		if sale.MaxDiscount > 0 && discountAmount > sale.MaxDiscount {
			discountAmount = sale.MaxDiscount
		}
	} else if sale.DiscountType == "fixed_amount" {
		discountAmount = sale.DiscountValue
		if discountAmount > basePrice {
			discountAmount = basePrice
		}
	}

	return discountAmount
}
