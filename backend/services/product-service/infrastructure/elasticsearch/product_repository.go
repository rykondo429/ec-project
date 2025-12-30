package elasticsearch

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"strings"

	"ec-sample/services/product-service/domain"

	"github.com/opensearch-project/opensearch-go/v2"
	"github.com/opensearch-project/opensearch-go/v2/opensearchapi"
)

const (
	productIndexName = "products"
)

// ElasticsearchRepository Elasticsearch実装
type ElasticsearchRepository struct {
	client *opensearch.Client
}

// NewElasticsearchRepository コンストラクタ
func NewElasticsearchRepository(client *opensearch.Client) domain.ProductRepository {
	return &ElasticsearchRepository{client: client}
}

// GetByID IDで商品を取得
func (r *ElasticsearchRepository) GetByID(ctx context.Context, id string) (*domain.Product, error) {
	req := opensearchapi.GetRequest{
		Index:      productIndexName,
		DocumentID: id,
	}

	res, err := req.Do(ctx, r.client)
	if err != nil {
		return nil, fmt.Errorf("failed to get product: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode == 404 {
		return nil, fmt.Errorf("product not found")
	}

	var result struct {
		Source *domain.Product `json:"_source"`
	}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return result.Source, nil
}

// Search 検索
func (r *ElasticsearchRepository) Search(ctx context.Context, query *domain.SearchQuery) (*domain.SearchResult, error) {
	// クエリ構築
	searchBody := buildSearchQuery(query)

	searchReq := opensearchapi.SearchRequest{
		Index: []string{productIndexName},
		Body:  strings.NewReader(searchBody),
	}

	res, err := searchReq.Do(ctx, r.client)
	if err != nil {
		return nil, fmt.Errorf("search failed: %w", err)
	}
	defer res.Body.Close()

	var searchResult struct {
		Hits struct {
			Total struct {
				Value int64 `json:"value"`
			} `json:"total"`
			Hits []struct {
				Source *domain.Product `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}

	if err := json.NewDecoder(res.Body).Decode(&searchResult); err != nil {
		return nil, fmt.Errorf("failed to decode search result: %w", err)
	}

	products := make([]*domain.Product, len(searchResult.Hits.Hits))
	for i, hit := range searchResult.Hits.Hits {
		products[i] = hit.Source
	}

	total := searchResult.Hits.Total.Value
	totalPages := int(math.Ceil(float64(total) / float64(query.PageSize)))

	return &domain.SearchResult{
		Products:    products,
		Total:       total,
		CurrentPage: query.Page,
		TotalPages:  totalPages,
		PageSize:    query.PageSize,
	}, nil
}

// GetByCategory カテゴリで取得
func (r *ElasticsearchRepository) GetByCategory(ctx context.Context, category string, page, pageSize int) (*domain.SearchResult, error) {
	from := (page - 1) * pageSize

	query := map[string]interface{}{
		"from": from,
		"size": pageSize,
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": []map[string]interface{}{
					{
						"match": map[string]interface{}{
							"category": category,
						},
					},
					{
						"match": map[string]interface{}{
							"is_active": true,
						},
					},
				},
			},
		},
	}

	queryBytes, _ := json.Marshal(query)

	searchReq := opensearchapi.SearchRequest{
		Index: []string{productIndexName},
		Body:  strings.NewReader(string(queryBytes)),
	}

	res, err := searchReq.Do(ctx, r.client)
	if err != nil {
		return nil, fmt.Errorf("search failed: %w", err)
	}
	defer res.Body.Close()

	var searchResult struct {
		Hits struct {
			Total struct {
				Value int64 `json:"value"`
			} `json:"total"`
			Hits []struct {
				Source *domain.Product `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}

	if err := json.NewDecoder(res.Body).Decode(&searchResult); err != nil {
		return nil, fmt.Errorf("failed to decode search result: %w", err)
	}

	products := make([]*domain.Product, len(searchResult.Hits.Hits))
	for i, hit := range searchResult.Hits.Hits {
		products[i] = hit.Source
	}

	total := searchResult.Hits.Total.Value
	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))

	return &domain.SearchResult{
		Products:    products,
		Total:       total,
		CurrentPage: page,
		TotalPages:  totalPages,
		PageSize:    pageSize,
	}, nil
}

// GetCategories カテゴリー一覧取得（Elasticsearchのaggregation使用）
func (r *ElasticsearchRepository) GetCategories(ctx context.Context) ([]string, error) {
	query := map[string]interface{}{
		"size": 0,
		"aggs": map[string]interface{}{
			"categories": map[string]interface{}{
				"terms": map[string]interface{}{
					"field": "category.keyword",
					"size":  100,
				},
			},
		},
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": []map[string]interface{}{
					{
						"term": map[string]interface{}{
							"is_active": true,
						},
					},
					{
						"exists": map[string]interface{}{
							"field": "category",
						},
					},
				},
			},
		},
	}

	queryBytes, _ := json.Marshal(query)

	searchReq := opensearchapi.SearchRequest{
		Index: []string{productIndexName},
		Body:  strings.NewReader(string(queryBytes)),
	}

	res, err := searchReq.Do(ctx, r.client)
	if err != nil {
		return nil, fmt.Errorf("failed to get categories: %w", err)
	}
	defer res.Body.Close()

	var result struct {
		Aggregations struct {
			Categories struct {
				Buckets []struct {
					Key string `json:"key"`
				} `json:"buckets"`
			} `json:"categories"`
		} `json:"aggregations"`
	}

	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode categories result: %w", err)
	}

	categories := make([]string, 0, len(result.Aggregations.Categories.Buckets))
	for _, bucket := range result.Aggregations.Categories.Buckets {
		if bucket.Key != "" {
			categories = append(categories, bucket.Key)
		}
	}

	return categories, nil
}

// GetFeatured おすすめ商品取得
func (r *ElasticsearchRepository) GetFeatured(ctx context.Context) ([]*domain.Product, error) {
	query := map[string]interface{}{
		"size": 10,
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": []map[string]interface{}{
					{
						"match": map[string]interface{}{
							"is_active": true,
						},
					},
				},
			},
		},
		"sort": []map[string]interface{}{
			{
				"rating": map[string]string{
					"order": "desc",
				},
			},
		},
	}

	queryBytes, _ := json.Marshal(query)

	searchReq := opensearchapi.SearchRequest{
		Index: []string{productIndexName},
		Body:  strings.NewReader(string(queryBytes)),
	}

	res, err := searchReq.Do(ctx, r.client)
	if err != nil {
		return nil, fmt.Errorf("search failed: %w", err)
	}
	defer res.Body.Close()

	var searchResult struct {
		Hits struct {
			Hits []struct {
				Source *domain.Product `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}

	if err := json.NewDecoder(res.Body).Decode(&searchResult); err != nil {
		return nil, fmt.Errorf("failed to decode search result: %w", err)
	}

	products := make([]*domain.Product, len(searchResult.Hits.Hits))
	for i, hit := range searchResult.Hits.Hits {
		products[i] = hit.Source
	}

	return products, nil
}

// Create 商品を作成
func (r *ElasticsearchRepository) Create(ctx context.Context, product *domain.Product) error {
	body, _ := json.Marshal(product)

	indexReq := opensearchapi.IndexRequest{
		Index:      productIndexName,
		DocumentID: product.ID,
		Body:       strings.NewReader(string(body)),
	}

	_, err := indexReq.Do(ctx, r.client)
	return err
}

// Update 商品を更新
func (r *ElasticsearchRepository) Update(ctx context.Context, product *domain.Product) error {
	return r.Create(ctx, product)
}

// buildSearchQuery 検索クエリ構築
func buildSearchQuery(query *domain.SearchQuery) string {
	from := (query.Page - 1) * query.PageSize

	must := []map[string]interface{}{
		{
			"match": map[string]interface{}{
				"is_active": true,
			},
		},
	}

	// キーワード検索
	if query.Keyword != "" {
		must = append(must, map[string]interface{}{
			"multi_match": map[string]interface{}{
				"query":  query.Keyword,
				"fields": []string{"name", "description"},
			},
		})
	}

	// カテゴリ検索
	if query.Category != "" {
		must = append(must, map[string]interface{}{
			"match": map[string]interface{}{
				"category": query.Category,
			},
		})
	}

	// 価格範囲
	if query.PriceMin > 0 || query.PriceMax > 0 {
		rangeFilter := map[string]interface{}{}
		if query.PriceMin > 0 {
			rangeFilter["gte"] = query.PriceMin
		}
		if query.PriceMax > 0 {
			rangeFilter["lte"] = query.PriceMax
		}
		must = append(must, map[string]interface{}{
			"range": map[string]interface{}{
				"price": rangeFilter,
			},
		})
	}

	// ソート
	sort := []map[string]interface{}{}
	if query.SortBy != "" {
		order := "asc"
		if query.SortOrder == "desc" {
			order = "desc"
		}
		sort = append(sort, map[string]interface{}{
			query.SortBy: map[string]string{
				"order": order,
			},
		})
	} else {
		sort = append(sort, map[string]interface{}{
			"_score": map[string]string{
				"order": "desc",
			},
		})
	}

	searchBody := map[string]interface{}{
		"from": from,
		"size": query.PageSize,
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": must,
			},
		},
		"sort": sort,
	}

	body, _ := json.Marshal(searchBody)
	return string(body)
}

// GetAll 全商品取得（MySQLから同期する際のダミー実装）
func (r *ElasticsearchRepository) GetAll(ctx context.Context) ([]*domain.Product, error) {
	// Elasticsearchでは全件取得は非推奨のため、大量データを扱う場合はscroll APIを使用
	// ここではシンプルに最大10000件まで取得
	searchBody := `{
		"query": {"match_all": {}},
		"size": 10000
	}`

	searchReq := opensearchapi.SearchRequest{
		Index: []string{productIndexName},
		Body:  strings.NewReader(searchBody),
	}

	res, err := searchReq.Do(ctx, r.client)
	if err != nil {
		return nil, fmt.Errorf("search failed: %w", err)
	}
	defer res.Body.Close()

	var searchResult struct {
		Hits struct {
			Hits []struct {
				Source *domain.Product `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}

	if err := json.NewDecoder(res.Body).Decode(&searchResult); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	products := make([]*domain.Product, len(searchResult.Hits.Hits))
	for i, hit := range searchResult.Hits.Hits {
		products[i] = hit.Source
	}

	return products, nil
}
