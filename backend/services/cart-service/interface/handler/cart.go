package handler

import (
	"net/http"

	"ec-sample/services/cart-service/domain"
	"ec-sample/services/cart-service/generated"
	"ec-sample/services/cart-service/usecase"

	"github.com/labstack/echo/v4"
)

// CartHandler HTTPハンドラー (generated.ServerInterfaceを実装)
type CartHandler struct {
	uc *usecase.CartUseCase
}

// NewCartHandler コンストラクタ
func NewCartHandler(uc *usecase.CartUseCase) *CartHandler {
	return &CartHandler{uc: uc}
}

// 型チェック: generated.ServerInterfaceを実装していることを保証
var _ generated.ServerInterface = (*CartHandler)(nil)

// GetCart カート取得
// @Summary カート取得
// @Tags carts
// @Produce json
// @Security BearerAuth
// @Success 200 {object} domain.Cart
// @Failure 400 {object} map[string]interface{}
// @Router /carts [get]
func (h *CartHandler) GetCart(c echo.Context) error {
	userID := c.Get("user_id").(string)

	cart, err := h.uc.GetCart(c.Request().Context(), userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		})
	}

	if cart == nil {
		return c.JSON(http.StatusOK, domain.Cart{
			UserID:    userID,
			Items:     []*domain.CartItem{},
			Total:     0,
			ItemCount: 0,
		})
	}

	return c.JSON(http.StatusOK, cart)
}

// AddItem アイテム追加
// @Summary カートにアイテムを追加
// @Tags carts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body domain.AddItemRequest true "Add item request"
// @Success 200 {object} domain.Cart
// @Failure 400 {object} map[string]interface{}
// @Router /carts/items [post]
func (h *CartHandler) AddItem(c echo.Context) error {
	userID := c.Get("user_id").(string)

	var req domain.AddItemRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "invalid request body",
		})
	}

	cart, err := h.uc.AddItem(c.Request().Context(), userID, &req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, cart)
}

// UpdateItem アイテム更新
// @Summary カートアイテムを更新
// @Tags carts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param product_id path string true "Product ID"
// @Param request body domain.UpdateItemRequest true "Update item request"
// @Success 200 {object} domain.Cart
// @Failure 400 {object} map[string]interface{}
// @Router /carts/items/{product_id} [put]
func (h *CartHandler) UpdateItem(c echo.Context, productId string) error {
	userID := c.Get("user_id").(string)

	var req domain.UpdateItemRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "invalid request body",
		})
	}

	cart, err := h.uc.UpdateItem(c.Request().Context(), userID, productId, req.Quantity)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, cart)
}

// RemoveItem アイテム削除
// @Summary カートからアイテムを削除
// @Tags carts
// @Produce json
// @Security BearerAuth
// @Param product_id path string true "Product ID"
// @Success 200 {object} domain.Cart
// @Failure 400 {object} map[string]interface{}
// @Router /carts/items/{product_id} [delete]
func (h *CartHandler) RemoveItem(c echo.Context, productId string) error {
	userID := c.Get("user_id").(string)

	cart, err := h.uc.RemoveItem(c.Request().Context(), userID, productId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, cart)
}

// ClearCart カートをクリア
// @Summary カートをクリア
// @Tags carts
// @Produce json
// @Security BearerAuth
// @Success 204
// @Failure 400 {object} map[string]interface{}
// @Router /carts [delete]
func (h *CartHandler) ClearCart(c echo.Context) error {
	userID := c.Get("user_id").(string)

	err := h.uc.ClearCart(c.Request().Context(), userID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
	}

	return c.NoContent(http.StatusNoContent)
}
