package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"ec-sample/services/promotion-service/domain"
	"ec-sample/services/promotion-service/generated"
	"ec-sample/services/promotion-service/infrastructure/mysql"
	"ec-sample/services/promotion-service/interface/handler"
	"ec-sample/services/promotion-service/usecase"
	"ec-sample/shared/auth"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
	mysqldriver "gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	// ロガー初期化
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	// MySQL 接続（GORM）
	dsn := os.Getenv("MYSQL_DSN")
	if dsn == "" {
		dsn = "ecuser:ecpassword@tcp(localhost:3306)/ecsite?parseTime=true"
	}

	db, err := gorm.Open(mysqldriver.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Fatal("MySQL connection failed", zap.Error(err))
	}

	// DB接続設定
	sqlDB, err := db.DB()
	if err != nil {
		logger.Fatal("Failed to get DB instance", zap.Error(err))
	}
	sqlDB.SetMaxOpenConns(10)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)

	// マイグレーション実行
	if err := db.AutoMigrate(&domain.Coupon{}, &domain.Sale{}); err != nil {
		logger.Fatal("Migration failed", zap.Error(err))
	}

	// リポジトリ初期化
	var couponRepo domain.CouponRepository
	var saleRepo domain.SaleRepository
	couponRepo = mysql.NewMySQLCouponRepository(db)
	saleRepo = mysql.NewMySQLSaleRepository(db)

	// ユースケース初期化
	promotionUC := usecase.NewPromotionUseCase(couponRepo, saleRepo)

	// ハンドラー初期化
	promotionHandler := handler.NewPromotionHandler(promotionUC)

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

		sqlDB, err := db.DB()
		if err != nil {
			return c.JSON(400, map[string]string{
				"status": "unhealthy",
			})
		}

		if err := sqlDB.PingContext(ctx); err != nil {
			return c.JSON(400, map[string]string{
				"status": "unhealthy",
			})
		}

		return c.JSON(200, map[string]string{
			"status": "healthy",
		})
	})

	// API ルート (generated.RegisterHandlersを使用)
	api := e.Group("/api/v1")
	generated.RegisterHandlers(api, promotionHandler)

	// OpenAPI スキーマ
	e.GET("/openapi.yaml", func(c echo.Context) error {
		return c.File("./openapi.yaml")
	})

	// サーバー起動
	port := getEnv("PORT", "8005")
	logger.Info(fmt.Sprintf("Promotion Service starting on port %s", port))
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
