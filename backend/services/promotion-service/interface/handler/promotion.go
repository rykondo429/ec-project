package handler

import (
	"net/http"

	"ec-sample/services/promotion-service/domain"
	"ec-sample/services/promotion-service/generated"
	"ec-sample/services/promotion-service/usecase"

	"github.com/labstack/echo/v4"
)

// PromotionHandler HTTPハンドラー (generated.ServerInterfaceを実装)
type PromotionHandler struct {
	uc *usecase.PromotionUseCase
}

// NewPromotionHandler コンストラクタ
func NewPromotionHandler(uc *usecase.PromotionUseCase) *PromotionHandler {
	return &PromotionHandler{uc: uc}
}

// 型チェック: generated.ServerInterfaceを実装していることを保証
var _ generated.ServerInterface = (*PromotionHandler)(nil)

// ValidateCoupon クーポン検証
// @Summary クーポン検証
// @Tags promotions
// @Accept json
// @Produce json
// @Param request body domain.ValidateCouponRequest true "Validate coupon request"
// @Success 200 {object} domain.CouponValidationResult
// @Router /promotions/validate-coupon [post]
func (h *PromotionHandler) ValidateCoupon(c echo.Context) error {
	var req domain.ValidateCouponRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "invalid request body",
		})
	}

	result, err := h.uc.ValidateCoupon(c.Request().Context(), &req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, result)
}

// GetActiveCoupons アクティブなクーポン取得
// @Summary アクティブなクーポン取得
// @Tags promotions
// @Produce json
// @Success 200 {array} domain.Coupon
// @Router /promotions/coupons [get]
func (h *PromotionHandler) GetActiveCoupons(c echo.Context) error {
	coupons, err := h.uc.GetActiveCoupons(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		})
	}

	if coupons == nil {
		coupons = []*domain.Coupon{}
	}

	return c.JSON(http.StatusOK, coupons)
}

// GetActiveSales アクティブなセール取得
// @Summary アクティブなセール取得
// @Tags promotions
// @Produce json
// @Success 200 {array} domain.Sale
// @Router /promotions/sales [get]
func (h *PromotionHandler) GetActiveSales(c echo.Context) error {
	sales, err := h.uc.GetActiveSales(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		})
	}

	if sales == nil {
		sales = []*domain.Sale{}
	}

	return c.JSON(http.StatusOK, sales)
}

// GetSalesByCategory カテゴリ別セール取得
// @Summary カテゴリ別セール取得
// @Tags promotions
// @Produce json
// @Param category path string true "Category"
// @Success 200 {array} domain.Sale
// @Router /promotions/sales/category/{category} [get]
func (h *PromotionHandler) GetSalesByCategory(c echo.Context, category string) error {
	sales, err := h.uc.GetSalesByCategory(c.Request().Context(), category)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
	}

	if sales == nil {
		sales = []*domain.Sale{}
	}

	return c.JSON(http.StatusOK, sales)
}

// GetSalesByProduct 商品別セール取得
// @Summary 商品別セール取得
// @Tags promotions
// @Produce json
// @Param product_id path string true "Product ID"
// @Success 200 {array} domain.Sale
// @Router /promotions/sales/product/{product_id} [get]
func (h *PromotionHandler) GetSalesByProduct(c echo.Context, productId string) error {
	sales, err := h.uc.GetSalesByProduct(c.Request().Context(), productId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
	}

	if sales == nil {
		sales = []*domain.Sale{}
	}

	return c.JSON(http.StatusOK, sales)
}
