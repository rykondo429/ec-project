# API Schema & Type Definitions

**フロントエンドとバックエンドの両方で使用する、統合OpenAPI定義とコード生成の管理ディレクトリです。**

## 📁 ディレクトリ構成

```
api-schema/
├── openapi/                         # OpenAPI YAML定義（単一の真実の源）
│   ├── product-service.yaml        # Product Service API
│   ├── cart-service.yaml           # Cart Service API
│   ├── order-service.yaml          # Order Service API
│   ├── point-service.yaml          # Point Service API
│   └── promotion-service.yaml      # Promotion Service API
├── oapi-codegen-configs/            # Go コード生成設定
│   ├── cart-service.yaml
│   ├── order-service.yaml
│   ├── point-service.yaml
│   ├── product-service.yaml
│   └── promotion-service.yaml
├── generated/                       # 自動生成Goコード（gitignore推奨）
│   ├── cart-service/
│   ├── order-service/
│   ├── point-service/
│   ├── product-service/
│   └── promotion-service/
├── types/                           # 自動生成TypeScript型（手動編集禁止）
│   ├── product-service.ts
│   ├── cart-service.ts
│   ├── order-service.ts
│   ├── point-service.ts
│   ├── promotion-service.ts
│   └── index.ts                    # 統合エクスポート
├── Makefile                         # Go コード生成用
├── package.json                     # TypeScript型生成用
└── README.md
```

## ⚠️ 重要な注意事項

### OpenAPI定義の管理

**`api-schema/openapi/`が唯一の真実の源（Single Source of Truth）です。**

- ✅ `api-schema/openapi/*.yaml` - **ここを編集**
- ❌ `backend/services/*/openapi.yaml` - **廃止（削除推奨）**

### Generated Goコードについて

**このディレクトリの`generated/`内のGoファイルは、コンパイルエラーを無視して問題ありません。**

理由：
1. **api-schemaはGoモジュールではない** - TypeScript/Node.jsプロジェクトです
2. **go.modがない** - Goの依存関係管理がありません
3. **生成コードの最終的な配置先は各サービス** - `make sync-to-services`で`backend/services/*/generated/`にコピーされます
4. **各サービスでビルドされる** - 各サービスのgo.modで依存関係が管理されます

VSCode設定（`.vscode/settings.json`）で`api-schema/generated`はGoplsの対象外に設定済みです。

### TypeScript型について

**`types/`ディレクトリのファイルは自動生成されます。手動編集しないでください。**

型を更新するには：
```bash
cd api-schema
npm run generate
```

## 🚀 クイックスタート

### 1. 依存関係のインストール

```bash
cd api-schema
npm install
```

### 2. コード生成

#### すべてのコード生成（TypeScript + Go）

```bash
# フロントエンド用TypeScript型とバックエンド用Goコードを両方生成
npm run generate
# または
make generate
```

#### TypeScript型のみ生成（フロントエンド用）

```bash
# 全サービスの型を生成
npm run generate:types

# 個別サービスの型を生成
npm run generate:product
npm run generate:cart
npm run generate:order
npm run generate:point
npm run generate:promotion
```

#### Goサーバーコードのみ生成（バックエンド用）

```bash
# すべてのサービスのEcho serverコードを生成
make generate-go

# 生成後に各サービスディレクトリに同期
make sync-to-services
```

### 3. OpenAPI仕様の検証

```bash
# 全サービスの検証
npm run validate

# 個別サービスの検証
npm run validate:product
```

## 📋 OpenAPI定義

### Product Service (8001)

**エンドポイント:**
- `GET /products/{id}` - 商品詳細取得
- `GET /products/search` - 商品検索
- `GET /products/category/{category}` - カテゴリ別取得
- `GET /products/featured` - おすすめ商品

**主要なスキーマ:**
- `Product` - 商品情報
- `SearchResponse` - 検索結果
- `Category` - カテゴリ情報

### Cart Service (8002)

**エンドポイント:**
- `GET /carts` - カート取得
- `POST /carts/items` - 商品追加
- `PUT /carts/items/{productId}` - 数量更新
- `DELETE /carts/items/{productId}` - 商品削除
- `DELETE /carts` - カートクリア

**主要なスキーマ:**
- `Cart` - カート情報
- `CartItem` - カート商品
- `AddToCartRequest` - 追加リクエスト

### Order Service (8003)

**エンドポイント:**
- `POST /orders` - 注文作成
- `GET /orders` - 注文履歴
- `GET /orders/{id}` - 注文詳細

**主要なスキーマ:**
- `Order` - 注文情報
- `OrderItem` - 注文明細
- `CreateOrderRequest` - 注文作成リクエスト

### Point Service (8004)

**エンドポイント:**
- `GET /points/balance` - 残高照会
- `POST /points/earn` - ポイント付与
- `POST /points/use` - ポイント利用
- `GET /points/transactions` - 履歴取得

**主要なスキーマ:**
- `PointBalance` - ポイント残高
- `PointTransaction` - ポイント履歴
- `EarnPointsRequest` - 付与リクエスト

### Promotion Service (8005)

**エンドポイント:**
- `GET /coupons` - クーポン一覧
- `GET /coupons/{code}` - クーポン詳細
- `POST /coupons/validate` - クーポン検証
- `GET /sales` - セール一覧

**主要なスキーマ:**
- `Coupon` - クーポン情報
- `Sale` - セール情報
- `ValidateCouponRequest` - 検証リクエスト

## 🔧 TypeScript型の使用方法

### フロントエンドでの使用

```typescript
// api-schema/types からインポート
import type { components } from '../api-schema/types/product-service';

type Product = components['schemas']['Product'];
type SearchResponse = components['schemas']['SearchResponse'];

// 使用例
const product: Product = {
  id: 'prod-123',
  name: 'スニーカー',
  price: 5980,
  category: 'shoes',
  // ...
};

// API レスポンスの型付け
async function searchProducts(keyword: string): Promise<SearchResponse> {
  const response = await fetch(`/api/products/search?keyword=${keyword}`);
  return response.json();
}
```

### パス定義の使用

```typescript
import type { paths } from '../api-schema/types/product-service';

// GET /products/{id} のレスポンス型
type GetProductResponse = paths['/products/{id}']['get']['responses']['200']['content']['application/json'];

// GET /products/search のパラメータ型
type SearchParams = paths['/products/search']['get']['parameters']['query'];
```

### 統合インポート

```typescript
// すべてのサービス型を一括インポート
import * as ProductAPI from '../api-schema/types/product-service';
import * as CartAPI from '../api-schema/types/cart-service';
import * as OrderAPI from '../api-schema/types/order-service';
import * as PointAPI from '../api-schema/types/point-service';
import * as PromotionAPI from '../api-schema/types/promotion-service';
```

## 📝 OpenAPI定義の更新フロー（重要！）

### 統合管理の原則

**`api-schema/openapi/` が唯一の真実の源（Single Source of Truth）です。**

OpenAPI定義は以下のディレクトリで一元管理されています：
- ✅ `api-schema/openapi/*.yaml` - **ここを編集**
- ❌ `backend/services/*/openapi.yaml` - **廃止（削除推奨）**

### 標準的な更新フロー

#### 1. OpenAPI定義を編集

```bash
# api-schema/openapi/ の YAML ファイルを直接編集
cd api-schema
vim openapi/product-service.yaml  # または任意のエディタ
```

#### 2. コード生成（フロントエンド＋バックエンド）

```bash
# TypeScript型とGoサーバーコードを両方生成
cd api-schema
npm run generate
# または
make generate

# 生成されたGoコードを各サービスに同期
make sync-to-services
```

#### 3. 各サービスで生成されたコードを確認

```bash
# 自動的に以下にコピーされます
# backend/services/product-service/generated/openapi.go
# backend/services/cart-service/generated/openapi.go
# など

# 各サービスをリビルド
cd backend/services/product-service
make build
```

#### 4. フロントエンドで最新の型を使用

```typescript
// api-schema/types/ から最新の型が利用可能
import type { components } from '../api-schema/types/product-service';
```

### 個別サービスからの生成（従来互換）

各サービスディレクトリで `make generate` を実行すると、自動的に統合OpenAPI定義から生成されます：

```bash
cd backend/services/cart-service
make generate
# → api-schema/openapi/cart-service.yaml を参照して生成
```

## 🔧 Goサーバーコードの使用方法

### 生成されるコード

`oapi-codegen` により以下が生成されます：

```go
// generated/openapi.go

// サーバーインターフェース（実装が必要）
type ServerInterface interface {
    GetProducts(ctx echo.Context) error
    CreateOrder(ctx echo.Context) error
    // ...
}

// リクエスト/レスポンス型
type Product struct {
    ID       string  `json:"id"`
    Name     string  `json:"name"`
    Price    float64 `json:"price"`
    // ...
}

// Echo ルート登録関数
func RegisterHandlers(e *echo.Echo, si ServerInterface)
```

### サービスでの使用例

```go
// main.go
import (
    "github.com/labstack/echo/v4"
    "ec-sample/services/product-service/generated"
    "ec-sample/services/product-service/interface/handler"
)

func main() {
    e := echo.New()
    
    // ハンドラー実装
    h := handler.NewProductHandler(/* ... */)
    
    // 自動生成されたルート登録
    generated.RegisterHandlers(e, h)
    
    e.Start(":8001")
}
```

```go
// interface/handler/product_handler.go
// ServerInterface を実装
func (h *ProductHandler) GetProducts(c echo.Context) error {
    // 自動生成された型を使用
    products := []generated.Product{
        {ID: "prod-1", Name: "商品A", Price: 1000},
    }
    return c.JSON(200, products)
}
```

## 🔧 TypeScript型の使用方法
- 軽量で高速

**インストール:**
```bash
npm install -D openapi-typescript
```

**使用例:**
```bash
npx openapi-typescript openapi/product-service.yaml -o types/product-service.ts
```

### Swagger UI（オプション）

OpenAPI仕様をブラウザで確認できます。

**起動方法:**
```bash
# Docker で Swagger UI を起動
docker run -p 8080:8080 -e SWAGGER_JSON=/openapi/product-service.yaml -v $(pwd)/openapi:/openapi swaggerapi/swagger-ui
```

ブラウザで http://localhost:8080 にアクセス

### OpenAPI Validator（オプション）

OpenAPI仕様の検証ツール。

**インストール:**
```bash
npm install -D @ibm/openapi-validator
```

**使用例:**
```bash
npx openapi-validator openapi/product-service.yaml
```

## 🔄 自動生成ワークフロー

### CI/CD での統合

```yaml
# .github/workflows/codegen.yml の例
name: Code Generation
on:
  push:
    paths:
      - 'api-schema/openapi/**'
      
jobs:
  generate:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Setup Node.js
        uses: actions/setup-node@v3
        with:
          node-version: '18'
          
      - name: Setup Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.21'
          
      - name: Install dependencies
        working-directory: api-schema
        run: npm install
        
      - name: Generate all code
        working-directory: api-schema
        run: |
          npm run generate:types
          make generate-go
          make sync-to-services
          
      - name: Commit generated code
        run: |
          git config --global user.name "GitHub Actions"
          git config --global user.email "actions@github.com"
          git add .
          git commit -m "chore: regenerate types and server code" || echo "No changes"
          git push
```

### Pre-commit Hook（推奨）

```bash
# .git/hooks/pre-commit
#!/bin/bash
cd api-schema

# OpenAPI定義が変更された場合のみ実行
if git diff --cached --name-only | grep -q "^api-schema/openapi/"; then
  echo "OpenAPI定義が変更されました。コードを再生成します..."
  
  # 型とサーバーコードを生成
  npm run generate:types
  make generate-go
  make sync-to-services
  
  # 変更を追加
  git add types/ generated/ ../backend/services/*/generated/
fi
```

## 🛠 開発ツール

### openapi-typescript

OpenAPI仕様からTypeScript型を生成します。

**特徴:**
- 厳密な型定義
- paths、components、operationsをすべてサポート
- 高速な生成

### oapi-codegen

OpenAPI仕様からGoサーバーコードを生成します。

**特徴:**
- Echo/Chi/Gin など複数フレームワーク対応
- モデル型とサーバーインターフェースの自動生成
- リクエスト/レスポンスバリデーション

**設定例:**
```yaml
# oapi-codegen-configs/cart-service.yaml
package: generated
generate:
  echo-server: true    # Echoサーバーコード生成
  models: true         # モデル型生成
output: ../generated/cart-service/openapi.go
```

### Swagger UI

OpenAPI仕様をブラウザで確認できます。

**起動方法:**
```bash
# Docker で Swagger UI を起動
docker run -p 8080:8080 -e SWAGGER_JSON=/openapi/product-service.yaml -v $(pwd)/openapi:/openapi swaggerapi/swagger-ui
```

ブラウザで http://localhost:8080 にアクセス

### OpenAPI Validator（オプション）

OpenAPI仕様の検証ツール。

**インストール:**
```bash
npm install -D @ibm/openapi-validator
```

**使用例:**
```bash
npx openapi-validator openapi/product-service.yaml
```

## 🚨 トラブルシューティング

### 生成されたGoコードのコンパイルエラー

```bash
# パッケージ名の不一致
# → oapi-codegen-configs/*.yaml の package を確認

# import path の問題
# → go.mod のモジュール名を確認

# 再生成
cd api-schema
make clean
make generate-go
make sync-to-services
```

### TypeScript型が更新されない

```bash
# node_modules をクリーンアップ
cd api-schema
rm -rf node_modules package-lock.json
npm install
npm run generate:types
```

### Makefile で permission denied

```bash
# 実行権限を付与
chmod +x api-schema/Makefile
```

## 📚 参考リンク

- [OpenAPI Specification](https://swagger.io/specification/)
- [openapi-typescript](https://github.com/drwpow/openapi-typescript)
- [oapi-codegen](https://github.com/deepmap/oapi-codegen)
- [Echo Framework](https://echo.labstack.com/)

## 🔄 自動生成ワークフロー（旧バージョン）

### CI/CD統合（推奨）

```yaml
# .github/workflows/generate-types.yml
name: Generate API Types

on:
  push:
    paths:
      - 'api-schema/openapi/**'

jobs:
  generate:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-node@v3
        with:
          node-version: '18'
      - run: cd api-schema && npm install
      - run: cd api-schema && npm run generate
      - run: git diff --exit-code types/ || (git add types/ && git commit -m "chore: regenerate API types")
```

### Git Hooks（オプション）

```bash
# .git/hooks/pre-commit
#!/bin/sh
cd api-schema
npm run generate
git add types/
```

## 📚 参考資料

### OpenAPI仕様

- [OpenAPI Specification](https://swagger.io/specification/)
- [OpenAPI Guide](https://swagger.io/docs/specification/about/)

### ツール

- [openapi-typescript](https://github.com/drwpow/openapi-typescript)
- [Swagger Editor](https://editor.swagger.io/)
- [Swagger UI](https://swagger.io/tools/swagger-ui/)

### プロジェクト内ドキュメント

- [バックエンド README](../backend/README.md)
- [フロントエンド README](../frontend/README.md)
- [API 使用例](../API_EXAMPLES.md)

## 🤝 貢献ガイドライン

### OpenAPI定義の編集

1. **api-schema/openapi/** の YAML ファイルを直接編集
2. **make generate** または **npm run generate** で両方のコードを生成
3. **npm run validate** で検証
4. **make sync-to-services** でGoコードを各サービスに配布
5. プルリクエスト作成

### 型定義の確認

生成された型は以下に保存されます：
- **types/** - TypeScript型（手動編集禁止）
- **generated/** - Goサーバーコード（手動編集禁止）

## ⚠️ 注意事項

- **types/** と **generated/** ディレクトリのファイルは自動生成されます。手動編集しないでください。
- **api-schema/openapi/** が唯一の真実の源です。各サービスの openapi.yaml は廃止されました。
- OpenAPI定義の変更時は必ず両方のコード（TypeScript + Go）を再生成してください。
- バックエンドとフロントエンドで型の整合性を保つため、定期的に同期してください。
- **generated/** ディレクトリを .gitignore に追加することを推奨します（各サービスで生成可能なため）。

## 🗂 .gitignore 推奨設定

```gitignore
# api-schema/.gitignore
node_modules/
generated/
*.log
```

## 🔗 関連リンク

- [Product Service OpenAPI](openapi/product-service.yaml)
- [Cart Service OpenAPI](openapi/cart-service.yaml)
- [Order Service OpenAPI](openapi/order-service.yaml)
- [Point Service OpenAPI](openapi/point-service.yaml)
- [Promotion Service OpenAPI](openapi/promotion-service.yaml)

---

**Happy API Development! 🚀**

