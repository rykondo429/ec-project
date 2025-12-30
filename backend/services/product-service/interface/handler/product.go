package handler

import (
	"net/http"

	"ec-sample/services/product-service/domain"
	"ec-sample/services/product-service/generated"
	"ec-sample/services/product-service/usecase"

	"github.com/labstack/echo/v4"
)

// ProductHandler HTTPハンドラー
type ProductHandler struct {
	uc *usecase.ProductUseCase
}

// NewProductHandler コンストラクタ
func NewProductHandler(uc *usecase.ProductUseCase) *ProductHandler {
	return &ProductHandler{uc: uc}
}

// GetProduct 商品詳細取得
// @Summary 商品詳細取得
// @Description 商品IDから商品情報を取得
// @Tags products
// @Produce json
// @Param id path string true "Product ID"
// @Success 200 {object} domain.Product
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/v1/products/{id} [get]
func (h *ProductHandler) GetProduct(c echo.Context, id string) error {
	product, err := h.uc.GetProductByID(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"error": err.Error(),
		})
	}
	return c.JSON(http.StatusOK, product)
}

// SearchProducts 商品検索
// @Summary 商品検索
// @Description キーワード・カテゴリ・価格範囲で商品検索
// @Tags products
// @Produce json
// @Param keyword query string false "Search keyword"
// @Param category query string false "Category"
// @Param price_min query number false "Minimum price"
// @Param price_max query number false "Maximum price"
// @Param sort_by query string false "Sort by (price, rating, created_at)"
// @Param sort_order query string false "Sort order (asc, desc)"
// @Param page query integer false "Page number" default(1)
// @Param page_size query integer false "Page size" default(20)
// @Success 200 {object} domain.SearchResult
// @Failure 400 {object} map[string]interface{}
// @Router /api/v1/products/search [get]
func (h *ProductHandler) SearchProducts(c echo.Context, params generated.SearchProductsParams) error {
	query := &domain.SearchQuery{
		Keyword:   "",
		Category:  "",
		SortBy:    "",
		SortOrder: "",
	}

	if params.Keyword != nil {
		query.Keyword = *params.Keyword
	}
	if params.Category != nil {
		query.Category = *params.Category
	}
	if params.SortBy != nil {
		query.SortBy = string(*params.SortBy)
	}
	if params.SortOrder != nil {
		query.SortOrder = string(*params.SortOrder)
	}
	if params.PriceMin != nil {
		query.PriceMin = float64(*params.PriceMin)
	}
	if params.PriceMax != nil {
		query.PriceMax = float64(*params.PriceMax)
	}
	if params.Page != nil {
		query.Page = *params.Page
	}
	if params.PageSize != nil {
		query.PageSize = *params.PageSize
	}

	result, err := h.uc.SearchProducts(c.Request().Context(), query)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, result)
}

// GetProductsByCategory カテゴリ別取得
// @Summary カテゴリ別商品取得
// @Description カテゴリから商品一覧を取得
// @Tags products
// @Produce json
// @Param category path string true "Category"
// @Param page query integer false "Page number" default(1)
// @Param page_size query integer false "Page size" default(20)
// @Success 200 {object} domain.SearchResult
// @Failure 400 {object} map[string]interface{}
// @Router /api/v1/products/category/{category} [get]
func (h *ProductHandler) GetProductsByCategory(c echo.Context, category string, params generated.GetProductsByCategoryParams) error {
	page := 1
	pageSize := 20

	if params.Page != nil && *params.Page > 0 {
		page = *params.Page
	}
	if params.PageSize != nil && *params.PageSize > 0 {
		pageSize = *params.PageSize
	}

	result, err := h.uc.GetProductsByCategory(c.Request().Context(), category, page, pageSize)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, result)
}

// GetFeaturedProducts おすすめ商品取得
// @Summary おすすめ商品取得
// @Description ピックアップ商品一覧を取得
// @Tags products
// @Produce json
// @Success 200 {array} domain.Product
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/products/featured [get]
func (h *ProductHandler) GetFeaturedProducts(c echo.Context) error {
	products, err := h.uc.GetFeaturedProducts(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		})
	}
	if products == nil {
		products = []*domain.Product{}
	}
	return c.JSON(http.StatusOK, products)
}

// GetCategories カテゴリー一覧取得
// @Summary カテゴリー一覧取得
// @Description 利用可能なカテゴリー一覧を取得
// @Tags products
// @Produce json
// @Success 200 {object} map[string][]string
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/products/categories [get]
func (h *ProductHandler) GetCategories(c echo.Context) error {
	categories, err := h.uc.GetCategories(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		})
	}
	if categories == nil {
		categories = []string{}
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"categories": categories,
	})
}

// SyncToElasticsearch MySQL→Elasticsearchの全商品同期
// @Summary Elasticsearch同期
// @Description MySQLの全商品をElasticsearchに同期
// @Tags admin
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/admin/sync-elasticsearch [post]
func (h *ProductHandler) SyncToElasticsearch(c echo.Context) error {
	count, err := h.uc.SyncAllProductsToElasticsearch(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error":   err.Error(),
			"synced":  count,
			"success": false,
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Elasticsearch sync completed successfully",
		"synced":  count,
		"success": true,
	})
}

// CreateProduct 商品作成
// @Summary 商品作成
// @Description 新しい商品を作成（MySQL + ES同期）
// @Tags products
// @Accept json
// @Produce json
// @Param product body domain.Product true "Product data"
// @Success 201 {object} domain.Product
// @Failure 400 {object} map[string]interface{}
// @Router /api/v1/products [post]
func (h *ProductHandler) CreateProduct(c echo.Context) error {
	var product domain.Product
	if err := c.Bind(&product); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "Invalid request body",
		})
	}

	if err := h.uc.CreateProduct(c.Request().Context(), &product); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, product)
}

// UpdateProduct 商品更新
// @Summary 商品更新
// @Description 商品情報を更新（MySQL + ES同期）
// @Tags products
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Param product body domain.Product true "Product data"
// @Success 200 {object} domain.Product
// @Failure 400 {object} map[string]interface{}
// @Router /api/v1/products/{id} [put]
func (h *ProductHandler) UpdateProduct(c echo.Context) error {
	id := c.Param("id")
	var product domain.Product
	if err := c.Bind(&product); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "Invalid request body",
		})
	}

	product.ID = id
	if err := h.uc.UpdateProduct(c.Request().Context(), &product); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, product)
}
