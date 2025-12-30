package handler

import (
	"net/http"

	"ec-sample/services/point-service/domain"
	"ec-sample/services/point-service/generated"
	"ec-sample/services/point-service/usecase"

	"github.com/labstack/echo/v4"
)

// PointHandler HTTPハンドラー (generated.ServerInterfaceを実装)
type PointHandler struct {
	uc *usecase.PointUseCase
}

// NewPointHandler コンストラクタ
func NewPointHandler(uc *usecase.PointUseCase) *PointHandler {
	return &PointHandler{uc: uc}
}

// 型チェック: generated.ServerInterfaceを実装していることを保証
var _ generated.ServerInterface = (*PointHandler)(nil)

// GetBalance ポイント残高取得
// @Summary ポイント残高取得
// @Tags points
// @Produce json
// @Success 200 {object} domain.PointBalance
// @Failure 500 {object} map[string]interface{}
// @Router /points/balance [get]
func (h *PointHandler) GetBalance(c echo.Context) error {
	userID := c.Get("user_id").(string)

	balance, err := h.uc.GetBalance(c.Request().Context(), userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, balance)
}

// EarnPoints ポイント付与
// @Summary ポイント付与
// @Tags points
// @Accept json
// @Produce json
// @Param request body domain.EarnPointsRequest true "Earn points request"
// @Success 200 {object} domain.PointBalance
// @Failure 400 {object} map[string]interface{}
// @Router /points/earn [post]
func (h *PointHandler) EarnPoints(c echo.Context) error {
	userID := c.Get("user_id").(string)

	var req domain.EarnPointsRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "invalid request body",
		})
	}

	balance, err := h.uc.EarnPoints(c.Request().Context(), userID, &req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, balance)
}

// UsePoints ポイント使用
// @Summary ポイント使用
// @Tags points
// @Accept json
// @Produce json
// @Param request body domain.UsePointsRequest true "Use points request"
// @Success 200 {object} domain.PointBalance
// @Failure 400 {object} map[string]interface{}
// @Router /points/use [post]
func (h *PointHandler) UsePoints(c echo.Context) error {
	userID := c.Get("user_id").(string)

	var req domain.UsePointsRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "invalid request body",
		})
	}

	balance, err := h.uc.UsePoints(c.Request().Context(), userID, &req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, balance)
}

// GetTransactions トランザクション履歴取得
// @Summary トランザクション履歴取得
// @Tags points
// @Produce json
// @Param page query integer false "Page number" default(1)
// @Param page_size query integer false "Page size" default(20)
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /points/transactions [get]
func (h *PointHandler) GetTransactions(c echo.Context, params generated.GetTransactionsParams) error {
	userID := c.Get("user_id").(string)

	page := 1
	pageSize := 20

	if params.Page != nil {
		page = *params.Page
	}
	if params.PageSize != nil {
		pageSize = *params.PageSize
	}

	transactions, total, err := h.uc.GetTransactions(c.Request().Context(), userID, page, pageSize)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
	}

	if transactions == nil {
		transactions = []*domain.PointTransaction{}
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"transactions": transactions,
		"total":        total,
		"current_page": page,
		"page_size":    pageSize,
	})
}
