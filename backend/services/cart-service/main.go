package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"ec-sample/services/cart-service/domain"
	"ec-sample/services/cart-service/generated"
	"ec-sample/services/cart-service/infrastructure/dynamodb"
	"ec-sample/services/cart-service/interface/handler"
	"ec-sample/services/cart-service/usecase"
	"ec-sample/shared/auth"

	"github.com/aws/aws-sdk-go-v2/config"
	awsdynamodb "github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
)

func main() {
	// ロガー初期化
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	// AWS Config 初期化
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		logger.Fatal("AWS config load failed", zap.Error(err))
	}

	// カスタムエンドポイント設定 (ローカル開発)
	if endpoint := os.Getenv("AWS_ENDPOINT_URL"); endpoint != "" {
		cfg.BaseEndpoint = &endpoint
	}

	// DynamoDB クライアント初期化
	dynamodbClient := awsdynamodb.NewFromConfig(cfg)

	// リポジトリ初期化
	var cartRepo domain.CartRepository
	cartRepo = dynamodb.NewDynamoDBRepository(dynamodbClient)

	// ユースケース初期化
	cartUC := usecase.NewCartUseCase(cartRepo)

	// ハンドラー初期化
	cartHandler := handler.NewCartHandler(cartUC)

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
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// DynamoDBへのPingテスト
		_, err := dynamodbClient.ListTables(ctx, &awsdynamodb.ListTablesInput{})
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"status": "unhealthy",
			})
		}

		return c.JSON(http.StatusOK, map[string]string{
			"status": "healthy",
		})
	})

	// API ルート (generated.RegisterHandlersを使用)
	api := e.Group("/api/v1")
	generated.RegisterHandlers(api, cartHandler)

	// OpenAPI スキーマ
	e.GET("/openapi.yaml", func(c echo.Context) error {
		return c.File("./openapi.yaml")
	})

	// サーバー起動
	port := getEnv("PORT", "8002")
	logger.Info(fmt.Sprintf("Cart Service starting on port %s", port))
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
