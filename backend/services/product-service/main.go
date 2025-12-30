package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"ec-sample/services/product-service/domain"
	"ec-sample/services/product-service/generated"
	"ec-sample/services/product-service/infrastructure/elasticsearch"
	mysqlrepo "ec-sample/services/product-service/infrastructure/mysql"
	redisrepo "ec-sample/services/product-service/infrastructure/redis"
	"ec-sample/services/product-service/interface/handler"
	"ec-sample/services/product-service/usecase"
	"ec-sample/shared/auth"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	opensearchgo "github.com/opensearch-project/opensearch-go/v2"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	// ロガー初期化
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	// MySQL接続初期化
	dsn := getEnv("MYSQL_DSN", "ecuser:ecpassword@tcp(localhost:3306)/ecsite?parseTime=true&loc=Local")
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Fatal("MySQL接続失敗", zap.Error(err))
	}

	sqlDB, _ := db.DB()
	sqlDB.SetMaxOpenConns(10)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(time.Hour)

	logger.Info("MySQL接続成功")

	// Redisクライアント初期化
	redisAddr := getEnv("REDIS_ADDR", "localhost:6379")
	redisClient := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: getEnv("REDIS_PASSWORD", ""),
		DB:       0,
	})

	// Redis接続確認
	ctx := context.Background()
	if err := redisClient.Ping(ctx).Err(); err != nil {
		logger.Warn("Redis接続失敗 - キャッシュなしで続行", zap.Error(err))
	} else {
		logger.Info("Redis接続成功")
	}

	// Elasticsearchクライアント初期化
	esConfig := opensearchgo.Config{
		Addresses: []string{
			getEnv("ELASTICSEARCH_URL", "http://localhost:9200"),
		},
		Transport: http.DefaultTransport,
	}

	esClient, err := opensearchgo.NewClient(esConfig)
	if err != nil {
		logger.Fatal("Elasticsearch接続失敗", zap.Error(err))
	}

	// リポジトリ初期化（MySQLを主データソースとして使用）
	var productRepo domain.ProductRepository
	productRepo = mysqlrepo.NewMySQLRepository(db)

	// SKUリポジトリ初期化
	skuRepo := mysqlrepo.NewSKURepository(db)

	// Elasticsearch用のリポジトリ（同期用に保持）
	esRepo := elasticsearch.NewElasticsearchRepository(esClient)

	// Redisキャッシュリポジトリ初期化
	cacheRepo := redisrepo.NewCacheRepository(redisClient)

	// ユースケース初期化（MySQLとES両方を渡す）
	productUC := usecase.NewProductUseCase(productRepo, esRepo, cacheRepo, logger)
	skuUC := usecase.NewSKUUseCase(skuRepo, productRepo, logger)

	// ハンドラー初期化
	productHandler := handler.NewProductHandler(productUC)
	skuHandler := handler.NewSKUHandler(skuUC)
	compositeHandler := handler.NewCompositeHandler(productHandler, skuHandler)

	// Echo インスタンス初期化
	e := echo.New()

	// ミドルウェア
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE},
	}))

	// Cognito 認証ミドルウェア
	cognitoConfig := auth.GetCognitoConfig()
	e.Use(auth.CognitoMiddleware(cognitoConfig, logger))

	// ヘルスチェック
	e.GET("/health", func(c echo.Context) error {
		// MySQL ヘルスチェック
		sqlDB, err := db.DB()
		if err != nil || sqlDB.Ping() != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"status": "unhealthy",
				"error":  "MySQL connection failed",
			})
		}

		// Redis ヘルスチェック（オプショナル）
		redisStatus := "connected"
		if err := redisClient.Ping(c.Request().Context()).Err(); err != nil {
			redisStatus = "disconnected"
		}

		return c.JSON(http.StatusOK, map[string]interface{}{
			"status": "healthy",
			"redis":  redisStatus,
		})
	})

	// API ルート
	api := e.Group("/api/v1")

	// OpenAPI定義のエンドポイント (generated.RegisterHandlersを使用)
	// GET /products/search, GET /products/featured, GET /products/categories, GET /products/category/{category}, GET /products/{id}
	// GET /products/{product_id}/skus, GET /products/{product_id}/with-skus
	// GET /skus/{id}, GET /skus/code/{code}, GET /skus/low-stock
	generated.RegisterHandlers(api, compositeHandler)

	// 追加のエンドポイント（OpenAPIに未定義）
	// 商品管理エンドポイント
	api.POST("/products", productHandler.CreateProduct)
	api.PUT("/products/:id", productHandler.UpdateProduct)

	// SKU管理エンドポイント
	api.POST("/products/:product_id/skus/bulk", skuHandler.CreateBulkSKUs)
	api.POST("/skus", skuHandler.CreateSKU)
	api.POST("/skus/stock/adjust", skuHandler.AdjustStock)
	api.PUT("/skus/:id", skuHandler.UpdateSKU)
	api.DELETE("/skus/:id", skuHandler.DeleteSKU)

	// 管理者エンドポイント（Elasticsearch同期）
	admin := e.Group("/api/v1/admin")
	admin.POST("/sync-elasticsearch", productHandler.SyncToElasticsearch)

	// OpenAPI スキーマ
	e.GET("/openapi.yaml", func(c echo.Context) error {
		return c.File("./openapi.yaml")
	})

	// サーバー起動
	port := getEnv("PORT", "8001")
	logger.Info(fmt.Sprintf("Product Service starting on port %s", port))
	if err := e.Start(":" + port); err != nil {
		logger.Fatal("Server start failed", zap.Error(err))
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
