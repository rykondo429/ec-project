package handler

import (
	"net/http"

	"ec-sample/services/order-service/domain"
	"ec-sample/services/order-service/generated"
	"ec-sample/services/order-service/usecase"

	"github.com/labstack/echo/v4"
)

// OrderHandler HTTPハンドラー (generated.ServerInterfaceを実装)
type OrderHandler struct {
	uc *usecase.OrderUseCase
}

// NewOrderHandler コンストラクタ
func NewOrderHandler(uc *usecase.OrderUseCase) *OrderHandler {
	return &OrderHandler{uc: uc}
}

// 型チェック: generated.ServerInterfaceを実装していることを保証
var _ generated.ServerInterface = (*OrderHandler)(nil)

// GetOrder 注文取得
// @Summary 注文取得
// @Tags orders
// @Produce json
// @Param id path string true "Order ID"
// @Success 200 {object} domain.Order
// @Failure 404 {object} map[string]interface{}
// @Router /orders/{id} [get]
func (h *OrderHandler) GetOrder(c echo.Context, id string) error {
	order, err := h.uc.GetOrder(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"error": err.Error(),
		})
	}
	return c.JSON(http.StatusOK, order)
}

// CreateOrder 注文作成
// @Summary 注文作成
// @Tags orders
// @Accept json
// @Produce json
// @Param request body domain.CreateOrderRequest true "Create order request"
// @Success 201 {object} domain.Order
// @Failure 400 {object} map[string]interface{}
// @Router /orders [post]
func (h *OrderHandler) CreateOrder(c echo.Context) error {
	userID := c.Get("user_id").(string)

	var req domain.CreateOrderRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "invalid request body",
		})
	}

	order, err := h.uc.CreateOrder(c.Request().Context(), userID, &req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, order)
}

// GetUserOrders ユーザーの注文履歴取得
// @Summary 注文履歴取得
// @Tags orders
// @Produce json
// @Param page query integer false "Page number" default(1)
// @Param page_size query integer false "Page size" default(20)
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /orders [get]
func (h *OrderHandler) GetUserOrders(c echo.Context, params generated.GetUserOrdersParams) error {
	userID := c.Get("user_id").(string)

	page := 1
	pageSize := 20

	if params.Page != nil && *params.Page > 0 {
		page = *params.Page
	}
	if params.PageSize != nil && *params.PageSize > 0 {
		pageSize = *params.PageSize
	}

	orders, total, err := h.uc.GetUserOrders(c.Request().Context(), userID, page, pageSize)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
	}

	if orders == nil {
		orders = []*domain.Order{}
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"orders":       orders,
		"total":        total,
		"current_page": page,
		"page_size":    pageSize,
	})
}

// UpdateOrderStatus 注文ステータス更新
// @Summary 注文ステータス更新
// @Tags orders
// @Accept json
// @Produce json
// @Param id path string true "Order ID"
// @Param status query string true "New status"
// @Success 200 {object} domain.Order
// @Failure 400 {object} map[string]interface{}
// @Router /orders/{id}/status [put]
func (h *OrderHandler) UpdateOrderStatus(c echo.Context, id string, params generated.UpdateOrderStatusParams) error {
	status := domain.OrderStatus(params.Status)

	order, err := h.uc.UpdateOrderStatus(c.Request().Context(), id, status)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, order)
}
