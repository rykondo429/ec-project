package handler

import (
	"net/http"
	"strconv"

	"ec-sample/services/product-service/domain"
	"ec-sample/services/product-service/usecase"

	"github.com/labstack/echo/v4"
)

// SKUHandler SKU HTTPハンドラー
type SKUHandler struct {
	uc *usecase.SKUUseCase
}

// NewSKUHandler コンストラクタ
func NewSKUHandler(uc *usecase.SKUUseCase) *SKUHandler {
	return &SKUHandler{uc: uc}
}

// CreateSKU SKU作成
func (h *SKUHandler) CreateSKU(c echo.Context) error {
	var req domain.SKUCreateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "Invalid request body",
		})
	}

	sku, err := h.uc.CreateSKU(c.Request().Context(), &req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, sku)
}

// GetSKU SKU取得
func (h *SKUHandler) GetSKU(c echo.Context) error {
	id := c.Param("id")
	sku, err := h.uc.GetSKUByID(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"error": err.Error(),
		})
	}
	return c.JSON(http.StatusOK, sku)
}

// GetSKUByCode SKUコードでSKU取得
func (h *SKUHandler) GetSKUByCode(c echo.Context) error {
	code := c.Param("code")
	sku, err := h.uc.GetSKUBySKUCode(c.Request().Context(), code)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"error": err.Error(),
		})
	}
	return c.JSON(http.StatusOK, sku)
}

// GetSKUsByProduct 商品のSKU一覧取得
func (h *SKUHandler) GetSKUsByProduct(c echo.Context) error {
	productID := c.Param("product_id")
	activeOnly := c.QueryParam("active_only") != "false"

	skus, err := h.uc.GetSKUsByProductID(c.Request().Context(), productID, activeOnly)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
	}

	if skus == nil {
		skus = []*domain.SKU{}
	}
	return c.JSON(http.StatusOK, skus)
}

// GetProductWithSKUs 商品とSKU一覧を取得
func (h *SKUHandler) GetProductWithSKUs(c echo.Context) error {
	productID := c.Param("product_id")

	result, err := h.uc.GetProductWithSKUs(c.Request().Context(), productID)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, result)
}

// UpdateSKU SKU更新
func (h *SKUHandler) UpdateSKU(c echo.Context) error {
	id := c.Param("id")
	var req domain.SKUUpdateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "Invalid request body",
		})
	}

	sku, err := h.uc.UpdateSKU(c.Request().Context(), id, &req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, sku)
}

// DeleteSKU SKU削除
func (h *SKUHandler) DeleteSKU(c echo.Context) error {
	id := c.Param("id")
	if err := h.uc.DeleteSKU(c.Request().Context(), id); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
	}
	return c.NoContent(http.StatusNoContent)
}

// AdjustStock 在庫調整
func (h *SKUHandler) AdjustStock(c echo.Context) error {
	var req domain.SKUStockAdjustment
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "Invalid request body",
		})
	}

	sku, err := h.uc.AdjustStock(c.Request().Context(), &req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, sku)
}

// GetLowStockSKUs 低在庫SKU一覧取得
func (h *SKUHandler) GetLowStockSKUs(c echo.Context) error {
	threshold := 10
	if t := c.QueryParam("threshold"); t != "" {
		if val, err := strconv.Atoi(t); err == nil {
			threshold = val
		}
	}

	skus, err := h.uc.GetLowStockSKUs(c.Request().Context(), threshold)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
	}

	if skus == nil {
		skus = []*domain.SKU{}
	}
	return c.JSON(http.StatusOK, skus)
}

// CreateBulkSKUs 一括SKU作成
func (h *SKUHandler) CreateBulkSKUs(c echo.Context) error {
	productID := c.Param("product_id")
	var requests []*domain.SKUCreateRequest
	if err := c.Bind(&requests); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "Invalid request body",
		})
	}

	skus, err := h.uc.CreateBulkSKUs(c.Request().Context(), productID, requests)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, skus)
}
