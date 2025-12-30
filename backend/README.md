# Backend - Go マイクロサービス

EC-Sample のバックエンドマイクロサービス群。Clean Architecture と DDD パターンに基づいた 5 つの独立したサービスで構成されています。

## 📋 目次

- [アーキテクチャ](#アーキテクチャ)
- [サービス一覧](#サービス一覧)
- [ディレクトリ構成](#ディレクトリ構成)
- [OpenAPIコード自動生成](#openapiコード自動生成)
- [共有ライブラリ](#共有ライブラリ)
- [開発ガイド](#開発ガイド)
- [API ドキュメント](#api-ドキュメント)

## 🏛 アーキテクチャ

### Clean Architecture

各サービスは以下の4層構造に従います：

```
┌─────────────────────────────────────────┐
│        Interface Layer (Handler)        │  ← HTTP リクエスト処理
├─────────────────────────────────────────┤
│         UseCase Layer (Business)        │  ← ビジネスロジック
├─────────────────────────────────────────┤
│      Domain Layer (Entity/Interface)    │  ← ドメインモデル
├─────────────────────────────────────────┤
│    Infrastructure Layer (Repository)    │  ← DB/外部サービス実装
└─────────────────────────────────────────┘
```

### 技術スタック

- **言語**: Go 1.21+
- **Web フレームワーク**: Echo v4
- **ORM**: GORM v2
- **ログ**: Uber Zap (構造化ログ)
- **認証**: JWT (golang-jwt/jwt v3)
- **AWS SDK**: AWS SDK Go v2

### 依存の方向

```
Interface → UseCase → Domain ← Infrastructure
```

- **Domain** は他の層に依存しない（純粋なビジネスロジック）
- **Infrastructure** は Domain のインターフェースを実装
- **UseCase** は Domain のインターフェースを通じて Infrastructure を利用
- **Interface** は UseCase を呼び出す

## 🔧 サービス一覧

### 1. Product Service (8001)

**役割**: 商品検索・管理

**技術**:
- Elasticsearch 8.10.2 (全文検索)
- Redis (キャッシング)
- MySQL (商品マスタ)

**主な機能**:
- 商品全文検索
- 商品一覧取得（ページネーション対応）
- 商品詳細取得
- カテゴリー別検索
- 価格範囲フィルタリング

**API**:
```bash
GET  /products          # 商品一覧
GET  /products/:id      # 商品詳細
GET  /products/search   # 検索
```

### 2. Cart Service (8002)

**役割**: ショッピングカート管理

**技術**:
- DynamoDB (カート永続化)
- ElastiCache/Redis (セッション管理)

**主な機能**:
- カート情報取得
- 商品追加・更新・削除
- カート内容のクリア
- 7日間自動有効期限

**API**:
```bash
GET    /carts          # カート取得
POST   /carts/items    # 商品追加
PUT    /carts/items/:productId  # 数量更新
DELETE /carts/items/:productId  # 商品削除
DELETE /carts          # カートクリア
```

### 3. Order Service (8003)

**役割**: 注文処理・管理

**技術**:
- MySQL (注文データ)
- GORM (ORM)
- AWS SQS (イベント送信)

**主な機能**:
- 注文作成
- 注文履歴取得
- 注文詳細確認
- 注文ステータス管理
- 決済連携（予定）

**API**:
```bash
POST /orders           # 注文作成
GET  /orders           # 注文履歴
GET  /orders/:id       # 注文詳細
```

### 4. Point Service (8004)

**役割**: ポイント管理

**技術**:
- MySQL (ポイント台帳)
- GORM (トランザクション管理)

**主な機能**:
- ポイント残高照会
- ポイント付与
- ポイント利用
- ポイント履歴取得
- 有効期限管理

**API**:
```bash
GET  /points/balance          # 残高照会
POST /points/earn             # ポイント付与
POST /points/use              # ポイント利用
GET  /points/transactions     # 履歴取得
```

### 5. Promotion Service (8005)

**役割**: プロモーション・割引管理

**技術**:
- MySQL (クーポン・セールデータ)
- GORM

**主な機能**:
- クーポン一覧取得
- クーポン詳細取得
- クーポン検証
- セール情報取得
- 割引計算

**API**:
```bash
GET  /coupons              # クーポン一覧
GET  /coupons/:code        # クーポン詳細
POST /coupons/validate     # クーポン検証
GET  /sales                # セール一覧
```

## 📁 ディレクトリ構成

### 全体構造

```
backend/
├── services/                # マイクロサービス
│   ├── product-service/
│   ├── cart-service/
│   ├── order-service/
│   ├── point-service/
│   └── promotion-service/
├── shared/                  # 共有ライブラリ
│   ├── auth/               # JWT 認証
│   ├── config/             # 設定管理
│   ├── errors/             # エラー定義
│   ├── logger/             # ロギング
│   └── go.mod
├── migrations/             # DB マイグレーション
│   └── mysql/
│       └── 01_init_schema.sql
└── README.md              # このファイル
```

### 各サービスの標準構成

すべてのサービスは同じディレクトリ構造に従います：

```
{service-name}/
├── main.go              # エントリポイント
├── domain/              # ドメイン層
│   ├── {model}.go      # ドメインモデル（struct定義）
│   └── repository.go   # Repositoryインターフェース定義
├── infrastructure/      # インフラ層
│   └── {db}/           # DB実装
│       └── {repo}.go   # Repository実装
├── interface/handler/   # インターフェース層
│   └── {model}.go      # HTTPハンドラ（Echoルート）
├── usecase/            # ユースケース層
│   └── {model}.go      # ビジネスロジック
├── generated/          # OpenAPIから自動生成 ✨
│   ├── openapi.go      # 型定義・インターフェース・ルート
│   └── README.md
├── openapi.yaml        # OpenAPI仕様
├── oapi-codegen.yaml   # コード生成設定
├── Makefile            # ビルドコマンド
├─ 🤖 OpenAPIコード自動生成

各サービスは**OpenAPI仕様からGoコードを自動生成**しています。

### 生成されるコード

`generated/openapi.go`に以下が含まれます：

1. **型定義（Models）** - リクエスト・レスポンスの構造体
2. **Echoサーバーインターフェース** - 実装すべきハンドラーメソッド
3. **ルート登録関数** - Echo用のルーティング設定
4. **OpenAPI仕様埋め込み** - 実行時にスキーマを参照可能

### 使用方法

```bash
# コード生成
cd backend/services/product-service
make generate

# ビルド（生成含む）
make build

# 実行（生成含む）
make run

# クリーンアップ
make clean
```

### ワークフロー

1. `openapi.yaml`を編集してAPI仕様を更新
2. `make generate`でコードを再生成
3. `interface/handler/*.go`で生成されたインターフェースを実装
4. `make build`でビルド確認

詳細は [CODE_GENERATION.md](CODE_GENERATION.md) を参照してください。

```
##─ tools.go            # Go generate用
├── go.mod
├── go.sum
└── Dockerfile
```

### 層の責務

#### Domain Layer
- **目的**: ビジネスルールとエンティティを定義
- **ファイル**: `domain/{model}.go`, `domain/repository.go`
- **依存**: なし（他の層に依存しない）

```go
// domain/order.go
type Order struct {
    ID        string
    UserID    string
    Items     []OrderItem
    Total     float64
    Status    string
    CreatedAt time.Time
}

// domain/repository.go
type OrderRepository interface {
    Create(ctx context.Context, order *Order) error
    FindByID(ctx context.Context, id string) (*Order, error)
    FindByUserID(ctx context.Context, userID string) ([]*Order, error)
}
```

#### Infrastructure Layer
- **目的**: 外部システムとの連携実装
- **ファイル**: `infrastructure/{db}/{repo}.go`
- **依存**: Domain Layer のインターフェースを実装

```go
// infrastructure/mysql/order_repository.go
type MySQLOrderRepository struct {
    db *gorm.DB
}

func (r *MySQLOrderRepository) Create(ctx context.Context, order *domain.Order) error {
    return r.db.WithContext(ctx).Create(order).Error
}
```

#### UseCase Layer
- **目的**: アプリケーション固有のビジネスロジック
- **ファイル**: `usecase/{model}.go`
- **依存**: Domain Layer のインターフェース経由で Infrastructure を利用

```go
// usecase/order.go
type OrderUseCase struct {
    repo domain.OrderRepository
}

func (uc *OrderUseCase) CreateOrder(ctx context.Context, req CreateOrderRequest) (*domain.Order, error) {
    // ビジネスロジック
    order := &domain.Order{...}
    return order, uc.repo.Create(ctx, order)
}
```

#### Interface Layer
- **目的**: HTTP リクエスト・レスポンス処理
- **ファイル**: `interface/handler/{model}.go`
- **依存**: UseCase Layer を呼び出す

```go
// interface/handler/order.go
type OrderHandler struct {
    usecase *usecase.OrderUseCase
}

func (h *OrderHandler) CreateOrder(c echo.Context) error {
    var req CreateOrderRequest
    if err := c.Bind(&req); err != nil {
        return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
    }
    
    userID := c.Get("user_id").(string)
    order, err := h.usecase.CreateOrder(c.Request().Context(), req)
    if err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
    }
    
    return c.JSON(http.StatusCreated, order)
}
```

## 🔗 共有ライブラリ

### shared/auth - JWT 認証

```go
// JWT トークン生成
token, err := auth.GenerateToken(userID, email)

// JWT トークン検証（middleware）
e.Use(auth.JWTMiddleware())
```

**機能**:
- JWT トークン生成
- トークン検証
- Echo middleware 統合
- Claims 抽出

### shared/config - 設定管理

```go
// 環境変数から設定を読み込み
cfg := config.Load()
fmt.Println(cfg.DatabaseDSN)
fmt.Println(cfg.RedisAddr)
```

**機能**:
- 環境変数読み込み
- デフォルト値設定
- 型安全な設定アクセス

### shared/errors - エラー定義

```go
// カスタムエラー
return errors.NotFound("order not found")
return errors.BadRequest("invalid input")
return errors.Unauthorized("token expired")
```

**機能**:
- 標準エラー型定義
- HTTP ステータスコード対応
- エラーレスポンス統一

### shared/logger - ロギング

```go
// 構造化ログ出力
logger := logger.NewLogger()
logger.Info("order created", zap.String("orderID", id))
logger.Error("failed to save", zap.Error(err))
```

**機能**:
- Zap 構造化ログ
- 開発/本番モード切り替え
- ログレベル設定

## 🚀 開発ガイド

### 環境構築

```bash
# 1. リポジトリクローン
git clone <repository-url>
cd ec-sample

# 2. 依存関係インストール（各サービス）
cd backend/services/product-service
go mod download

# 3. 共有ライブラリのインストール
cd backend/shared
go mod download
```

### ローカル開発

#### 個別サービスの起動

```bash
# Product Service
cd backend/services/product-service
go run main.go
# → http://localhost:8001

# Order Service
cd backend/services/order-service
go run main.go
# → http://localhost:8003
```

#### 環境変数設定

各サービスで必要な環境変数：

```bash
# 共通
export PORT=8001
export ENVIRONMENT=development
export JWT_SECRET_KEY=your-secret-key

# MySQL系サービス（order, point, promotion）
export MYSQL_DSN="user:pass@tcp(localhost:3306)/dbname?parseTime=true"

# Redis使用サービス（product, cart）
export REDIS_ADDR=localhost:6379

# Elasticsearch（product）
export ELASTICSEARCH_URL=http://localhost:9200

# DynamoDB（cart）
export AWS_REGION=us-east-1
export AWS_ENDPOINT_URL=http://localhost:4566
export AWS_ACCESS_KEY_ID=test
export AWS_SECRET_ACCESS_KEY=test
```

### ビルド

```bash
# 全サービスビルド（ルートから）
cd /path/to/ec-sample
bash build.sh

# 個別サービスビルド
cd backend/services/order-service
go build -o bin/main .
./bin/main
```

### テスト

```bash
# 単体テスト
cd backend/services/order-service
go test ./...

# カバレッジ付きテスト
go test -cover ./...

# 詳細カバレッジレポート
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### 依存関係管理

```bash
# 依存関係追加
go get github.com/some/package

# 不要な依存削除
go mod tidy

# ベンダリング（オプション）
go mod vendor
```

## 🔌 API ドキュメント

### OpenAPI 仕様

各サービスの OpenAPI 定義：

- [Product Service](services/product-service/openapi.yaml)
- [Cart Service](services/cart-service/openapi.yaml)
- [Order Service](services/order-service/openapi.yaml)
- [Point Service](services/point-service/openapi.yaml)
- [Promotion Service](services/promotion-service/openapi.yaml)

### API 使用例

詳細な API リクエスト例は [API_EXAMPLES.md](API_EXAMPLES.md) を参照してください。

### 認証

すべての API は JWT Bearer トークンによる認証をサポート：

```bash
curl -H "Authorization: Bearer <token>" \
     http://localhost:8001/products
```

未認証リクエストは `anonymous` ユーザーとして処理されます。

## 🐳 Docker 運用

### イメージビルド

```bash
# 個別サービス
docker build -t ec-product-service ./backend/services/product-service

# 全サービス（docker-compose）
docker-compose build
```

### コンテナ起動

```bash
# 全サービス起動
docker-compose up -d

# ログ確認
docker-compose logs -f ec-product-service

# サービス再起動
docker-compose restart ec-product-service
```

## 🛠 トラブルシューティング

### MySQL 接続エラー

```bash
# MySQL接続確認
mysql -h localhost -u ecuser -p ecsite

# DSN確認
echo $MYSQL_DSN
```

### Redis 接続エラー

```bash
# Redis接続確認
redis-cli ping

# ポート確認
netstat -an | grep 6379
```

### Elasticsearch 接続エラー

```bash
# Elasticsearch確認
curl http://localhost:9200/_cluster/health

# インデックス確認
curl http://localhost:9200/_cat/indices
```

### ポート競合

```bash
# 使用中ポート確認
lsof -i :8001

# プロセス終了
kill -9 <PID>
```

## 📚 参考資料

### 公式ドキュメント

- [Echo Framework](https://echo.labstack.com/)
- [GORM](https://gorm.io/)
- [Uber Zap](https://github.com/uber-go/zap)
- [AWS SDK for Go v2](https://aws.github.io/aws-sdk-go-v2/)

### アーキテクチャパターン

- [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [Domain-Driven Design](https://martinfowler.com/bliki/DomainDrivenDesign.html)
- [Repository Pattern](https://martinfowler.com/eaaCatalog/repository.html)

### プロジェクト内ドキュメント

- [プロジェクトルート README](../README.md)
- [API 使用例](API_EXAMPLES.md)
- [セットアップガイド](../SETUP.md)
- [フロントエンド README](../frontend/README.md)

## 🤝 開発ワークフロー

### 新しいエンドポイント追加

1. **Domain 定義**: `domain/{model}.go` にビジネスモデル追加
2. **Repository インターフェース**: `domain/repository.go` にメソッド追加
3. **Infrastructure 実装**: `infrastructure/{db}/{repo}.go` に実装追加
4. **UseCase 追加**: `usecase/{model}.go` にビジネスロジック追加
5. **Handler 追加**: `interface/handler/{model}.go` に HTTP ハンドラ追加
6. **ルート登録**: `main.go` にルート追加

### 新しいサービス追加

1. `backend/services/new-service/` ディレクトリ作成
2. 既存サービスの構造をテンプレートとしてコピー
3. `main.go` でサービス固有の設定
4. `docker-compose.yml` に新サービス定義追加
5. `build.sh` にビルドパス追加

## 📝 コーディング規約

### 命名規則

- **パッケージ**: 小文字、単数形（`order`, `product`）
- **インターフェース**: 名詞または動詞+er（`Repository`, `Handler`）
- **構造体**: PascalCase（`Order`, `OrderItem`）
- **関数/メソッド**: PascalCase（公開）、camelCase（非公開）
- **定数**: UPPER_SNAKE_CASE または PascalCase

### エラーハンドリング

```go
// 良い例
if err != nil {
    logger.Error("failed to create order", zap.Error(err))
    return nil, fmt.Errorf("create order: %w", err)
}

// 悪い例
if err != nil {
    panic(err)  // panic は使わない
}
```

### コンテキスト使用

```go
// タイムアウト設定
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

order, err := repo.FindByID(ctx, orderID)
```

### ログ出力

```go
// 構造化ログを使用
logger.Info("order created",
    zap.String("orderID", order.ID),
    zap.String("userID", order.UserID),
    zap.Float64("total", order.Total),
)
```

---

## 📚 参考リンク

- [プロジェクトルート README](../README.md)
- [API Schema & Type Definitions](../api-schema/README.md)
- [Echo Framework](https://echo.labstack.com/)
- [GORM](https://gorm.io/)
- [AWS SDK for Go v2](https://aws.github.io/aws-sdk-go-v2/)

**Happy Coding! 🚀**
