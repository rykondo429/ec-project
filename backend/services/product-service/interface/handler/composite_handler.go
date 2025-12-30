package handler

import (
	"ec-sample/services/product-service/generated"

	"github.com/labstack/echo/v4"
)

// CompositeHandler ProductとSKUの両方のハンドラーを統合
type CompositeHandler struct {
	*ProductHandler
	*SKUHandler
}

// NewCompositeHandler コンストラクタ
func NewCompositeHandler(productHandler *ProductHandler, skuHandler *SKUHandler) generated.ServerInterface {
	return &CompositeHandler{
		ProductHandler: productHandler,
		SKUHandler:     skuHandler,
	}
}

// SKUメソッドの委譲（OpenAPI互換のシグネチャ）
func (h *CompositeHandler) GetSKUsByProduct(c echo.Context, productId string) error {
	return h.SKUHandler.GetSKUsByProduct(c)
}

func (h *CompositeHandler) GetProductWithSKUs(c echo.Context, productId string) error {
	return h.SKUHandler.GetProductWithSKUs(c)
}

func (h *CompositeHandler) GetSKUByCode(c echo.Context, code string) error {
	return h.SKUHandler.GetSKUByCode(c)
}

func (h *CompositeHandler) GetLowStockSKUs(c echo.Context, params generated.GetLowStockSKUsParams) error {
	return h.SKUHandler.GetLowStockSKUs(c)
}

func (h *CompositeHandler) GetSKU(c echo.Context, id string) error {
	return h.SKUHandler.GetSKU(c)
}
