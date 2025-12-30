# EC-Sample - マイクロサービス E コマースプラットフォーム

Go + Next.js で構築された、Clean Architecture ベースのマイクロサービス型 EC サイトのサンプル実装

[![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![Next.js](https://img.shields.io/badge/Next.js-14-000000?style=flat&logo=next.js)](https://nextjs.org/)
[![TypeScript](https://img.shields.io/badge/TypeScript-5.3-3178C6?style=flat&logo=typescript)](https://www.typescriptlang.org/)
[![Docker](https://img.shields.io/badge/Docker-Compose-2496ED?style=flat&logo=docker)](https://www.docker.com/)
[![CI](https://github.com/YOUR_USERNAME/ec-sample/workflows/CI/badge.svg)](https://github.com/YOUR_USERNAME/ec-sample/actions)

## 📋 目次

- [特徴](#特徴)
- [アーキテクチャ](#アーキテクチャ)
- [ディレクトリ構成](#ディレクトリ構成)
- [クイックスタート](#クイックスタート)
- [開発ガイド](#開発ガイド)
- [CI/CD](#cicd)
- [API ドキュメント](#api-ドキュメント)
- [ドキュメント](#ドキュメント)

## ✨ 特徴

- 🏗 **Clean Architecture**: ドメイン駆動設計に基づいた階層構造
- 🔌 **マイクロサービス**: 5つの独立したバックエンドサービス
- ⚛️ **モダンフロントエンド**: Next.js 14 App Router + TailwindCSS
- 🐳 **Docker対応**: ワンコマンドで開発環境構築
- 🔐 **JWT認証**: セキュアな認証・認可
- 📊 **複数DB対応**: MySQL, DynamoDB, Elasticsearch, Redis
- 🎨 **型安全**: TypeScript による完全な型定義
- 📝 **OpenAPI**: 各サービスの API 仕様を定義
- 🚀 **CI/CD**: GitHub Actions による自動テスト・ビルド

## 🏛 アーキテクチャ

### 技術スタック

#### バックエンド
- **言語**: Go 1.21+
- **フレームワーク**: Echo v4
- **ORM**: GORM
- **ログ**: Uber Zap (構造化ログ)
- **認証**: JWT (golang-jwt/jwt v3)

#### フロントエンド
- **フレームワーク**: Next.js 14 (App Router)
- **UI**: React 18 + TypeScript 5.3
- **スタイリング**: TailwindCSS 3.3
- **状態管理**: Zustand 4.4
- **HTTP クライアント**: Axios 1.6

#### データストア
- **MySQL 8.0**: 商品・注文・ポイント・プロモーション管理（主データソース）
- **DynamoDB**: カート管理 (LocalStack)
- **Elasticsearch 8.10**: 商品検索インデックス（MySQLから同期）
- **Kibana 8.10**: Elasticsearch可視化・管理ツール
- **Redis**: キャッシュレイヤー (ElastiCache API)

#### インフラ
- **Docker & Docker Compose**: コンテナオーケストレーション
- **LocalStack**: AWS サービスのローカルエミュレーション

## マイクロサービス

### 1. Product Service (商品管理)
- 商品データ管理（MySQL）
- 商品検索（MySQL FULLTEXT + Elasticsearch同期）
- 商品一覧取得
- 商品詳細取得
- 在庫管理
- ポート: 8001

### 2. Cart Service (カート管理)
- カート情報の取得・更新
- DynamoDB + ElastiCache を活用
- ポート: 8002

### 3. Order Service (注文管理)
- 注文作成・確認
- 注文履歴取得
- イベント駆動で他サービスと連携
- ポート: 8003

### 4. Point Service (ポイント管理)
- ポイント残高取得
- ポイント付与・使用
- トランザクション管理
- ポート: 8004

### 5. Promotion Service (プロモーション管理)
- クーポン・セール情報管理
- 割引計算
- ポート: 8005

## 📁 ディレクトリ構成

```
ec-sample/
├── frontend/                    # Next.js フロントエンド (Port 3000)
│   ├── src/
│   │   ├── app/                # Next.js App Router ページ
│   │   ├── components/         # React コンポーネント
│   │   ├── lib/                # API クライアント
│   │   ├── store/              # Zustand 状態管理
│   │   ├── types/              # TypeScript 型定義
│   │   └── styles/             # TailwindCSS スタイル
│   ├── package.json
│   ├── next.config.js
│   ├── tailwind.config.js
│   └── Dockerfile
│
├── backend/                     # Go マイクロサービス
│   ├── services/               # 5つの独立したサービス
│   │   ├── product-service/    # 商品検索 (8001 - Elasticsearch)
│   │   ├── cart-service/       # カート管理 (8002 - DynamoDB)
│   │   ├── order-service/      # 注文管理 (8003 - MySQL)
│   │   ├── point-service/      # ポイント管理 (8004 - MySQL)
│   │   └── promotion-service/  # プロモーション (8005 - MySQL)
│   ├── shared/                 # 共有ライブラリ
│   │   ├── auth/               # JWT 認証
│   │   ├── config/             # 設定管理
│   │   ├── errors/             # エラー定義
│   │   └── logger/             # Zap ロギング
│   └── migrations/             # DB マイグレーション
│       └── mysql/
│           └── 01_init_schema.sql
│
├── api-schema/                 # API 仕様・スキーマ
│   ├── openapi/                # OpenAPI YAML 定義
│   └── types/                  # 共有型定義
│
├── docs/                       # プロジェクトドキュメント
│
├── docker-compose.yml          # インフラストラクチャ定義
├── build.sh                    # ビルドスクリプト
├── init-aws.sh                 # LocalStack 初期化
└── README.md                   # このファイル
```

### 各サービスの標準構造

各バックエンドサービスは同じ Clean Architecture パターンに従います：

```
backend/services/{service-name}/
├── main.go              # エントリポイント、DI設定
├── domain/              # ドメイン層
│   ├── {model}.go      # ビジネスモデル
│   └── repository.go   # Repository インターフェース
├── infrastructure/      # インフラ層
│   └── {db}/           # DB 実装 (GORM, DynamoDB SDK, etc.)
│       └── {repo}.go
├── interface/handler/   # インターフェース層
│   └── {model}.go      # HTTP ハンドラ (Echo)
├── usecase/            # ユースケース層
│   └── {model}.go      # ビジネスロジック
├── openapi.yaml        # OpenAPI 仕様
├── go.mod
└── Dockerfile
```

## 🚀 クイックスタート

### 前提条件

- Docker & Docker Compose
- Go 1.23+
- Node.js 20+ (Dockerのみの場合は不要)

### 環境変数の設定

`.env`ファイルは既に開発用設定で作成済みです。本番環境では以下の手順で設定：

```bash
# 本番環境用のテンプレートをコピー
cp .env.production.example .env.production

# 本番用の値を設定（秘密鍵、DB接続情報など）
vi .env.production

# 本番環境で起動（.envではなく.env.productionを使用）
docker-compose --env-file .env.production up -d
```

### Docker開発環境（推奨）

**ローカルにNode.jsをインストールせずに、Dockerだけで開発できます。**

```bash
# すべてのサービスを起動（データベース + バックエンド + フロントエンド）
docker-compose up -d

# ログ確認
docker-compose logs -f

# 停止
docker-compose down
```

**開発環境と本番環境の違い:**

| 環境 | 構成ファイル | 用途 | 特徴 |
|------|------------|------|------|
| 開発 | `docker-compose.yml` + `.env` | フルスタック開発 | LocalStack使用、全サービスビルド |
| 開発（軽量） | `docker-compose.dev.yml` | フロントエンド開発 | インフラのみ起動、ホットリロード |
| 本番 | `docker-compose.yml` + `.env.production` | デプロイ | 実際のAWSサービス、最適化ビルド |

**アクセス:**
- フロントエンド: http://localhost:3000
- Product API: http://localhost:8001
- Cart API: http://localhost:8002
- Order API: http://localhost:8003
- Point API: http://localhost:8004
- Promotion API: http://localhost:8005

### ローカル開発（Go + Node.js）

#### 1. インフラストラクチャの起動

```bash
# MySQL, Redis, Elasticsearch, LocalStack を起動
docker-compose up -d mysql redis elasticsearch localstack

# ヘルスチェック
docker-compose ps
```

#### 2. サンプルデータの投入

**すべてのデータはMySQLで管理されます**：

```bash
# 商品データの投入（12件のサンプル商品）
./scripts/seed-products-mysql.sh

# 確認
mysql -h127.0.0.1 -uecuser -pecpassword ecsite -e "SELECT id, name, price, stock FROM products;"
curl 'http://localhost:8001/api/v1/products/search?page=1&page_size=10'
```

投入される商品:
- Nike Air Max 90
- Adidas Ultraboost 22 (セール中)
- Sony WH-1000XM5
- Apple AirPods Pro (セール中)
- Patagonia フリースジャケット
- The North Face ダウンジャケット
- MacBook Air M2
- iPad Pro (セール中)
- New Balance 2002R
- Uniqlo ヒートテック (セール中)
- Bose QuietComfort Earbuds
- Converse Chuck Taylor

> **重要**: Product ServiceはMySQLを主データソースとして使用します。Elasticsearchは検索インデックスとしてのみ使用し、データはMySQLから同期されます。

#### 3. バックエンドサービスの起動

各サービスは`.env`ファイルから環境変数を読み込みます：

```bash
# Product Service
cd backend/services/product-service && go run main.go  # :8001

# Cart Service
cd backend/services/cart-service && go run main.go     # :8002

# Order Service
cd backend/services/order-service && go run main.go    # :8003

# Point Service
cd backend/services/point-service && go run main.go    # :8004

# Promotion Service
cd backend/services/promotion-service && go run main.go # :8005
```

#### 4. フロントエンドの起動

```bash
cd frontend
npm install
npm run dev  # http://localhost:3000
```

## 🔧 開発ガイド

### Docker Compose設定の使い分け

**`docker-compose.yml`（デフォルト）:**
- 全サービスをビルド＆起動（バックエンド5個 + フロントエンド + インフラ）
- 本番同等の環境でテスト可能
- `.env`ファイルから環境変数を自動読み込み

**`docker-compose.dev.yml`（開発専用）:**
- インフラ（MySQL, Redis, Elasticsearch, LocalStack）のみ
- フロントエンドはホットリロード対応（`Dockerfile.dev`使用）
- バックエンドはローカルで`go run`して開発（高速リロード）

```bash
# パターン1: 全てDocker（本番同等環境）
docker-compose up -d

# パターン2: フロントエンドのみDocker開発モード
docker-compose -f docker-compose.dev.yml up web

# パターン3: インフラのみDocker、アプリはローカル
docker-compose -f docker-compose.dev.yml up -d mysql redis elasticsearch
cd backend/services/product-service && go run main.go
cd frontend && npm run dev
```

### ビルド

```bash
# 全サービスビルド（バックエンド + フロントエンド）
bash build.sh

# 個別サービスビルド
cd backend/services/order-service
go build -o bin/main .

# フロントエンドビルド
cd frontend
npm run build
```

### テスト

```bash
# バックエンド
cd backend/services/order-service
go test ./...

# フロントエンド
cd frontend
npm test
```

### Docker 運用

```bash
# イメージビルド（全体）
docker-compose build

# 個別イメージビルド
docker build -t ec-product-service ./backend/services/product-service

# コンテナ起動
docker-compose up -d

# ログ確認
docker-compose logs -f ec-order-service

# 停止
docker-compose down
```

### 本番環境へのデプロイ

```bash
# 1. 本番用環境変数を設定
cp .env.production.example .env.production
vi .env.production  # 実際の値を設定（DB_PASSWORD, JWT_SECRET_KEY など）

# 2. 本番用設定でビルド＆起動
docker-compose --env-file .env.production build
docker-compose --env-file .env.production up -d

# 3. ログ確認
docker-compose logs -f

# 4. ヘルスチェック
curl http://localhost:8001/health
curl http://localhost:3000
```

**本番環境での注意点:**
- `.env.production`は**絶対にGitにコミットしない**（`.gitignore`に設定済み）
- `JWT_SECRET_KEY`は強力なランダム文字列を使用（`openssl rand -base64 32`）
- `AWS_ENDPOINT_URL`は空にする（LocalStack不使用）
- データベースは実際のRDS/ElastiCache/OpenSearchを使用
- 環境変数はAWS Secrets ManagerやKubernetesのSecretsで管理推奨

## 🔐 認証フロー

### AWS Cognito OIDC認証

本プロジェクトは**AWS Cognito OpenID Connect (OIDC)** による認証を実装しています。

#### 認証機能

✅ **会員登録** (`/signup`)
- メールアドレス + パスワードで新規登録
- メール確認コードによる本人確認
- パスワード要件: 8文字以上、大文字・小文字・数字を含む

✅ **ログイン** (`/login`)
- メール + パスワード認証
- JWT IDトークンの自動取得・保存
- リダイレクトパラメータ対応

✅ **パスワードリセット** (`/forgot-password`)
- メールアドレスで確認コード送信
- コード確認後に新パスワード設定

✅ **マイページ** (`/my-page`)
- 認証済みユーザーのみアクセス可能
- ユーザー情報表示
- 注文履歴・ポイント・お気に入りへのリンク

✅ **ヘッダー統合**
- 未ログイン時: ログイン・新規登録ボタン表示
- ログイン時: ユーザー名表示、ログアウトボタン

#### ローカル開発環境セットアップ

**⚠️ 重要: LocalStack無料版ではCognitoがサポートされていません**

開発環境では以下の2つの方法があります：

**方法1: 認証機能なし（推奨・最速）**

```bash
# 1. Docker起動
docker-compose up -d

# 2. 開発モード（認証無効）で環境変数設定
./scripts/setup-cognito-env.sh
# 選択: 1) 開発モード（認証機能無効）

# 3. フロントエンド起動
cd frontend && npm run dev

# 4. ブラウザでアクセス
open http://localhost:3000
# 認証なしで全機能が利用可能
```

**方法2: AWS Cognito本番環境を使用**

```bash
# 1. AWSコンソールでCognitoユーザープール作成
# https://console.aws.amazon.com/cognito/

# 2. 環境変数設定スクリプト実行
./scripts/setup-cognito-env.sh
# 選択: 2) AWS Cognito本番環境
# User Pool ID、Client IDを入力

# 3. バックエンド環境変数も設定
# docker-compose.ymlまたは.envファイルに追加:
# COGNITO_REGION=us-east-1
# COGNITO_USER_POOL_ID=us-east-1_xxxxxxxxx
# COGNITO_CLIENT_ID=xxxxxxxxxxxxxx

# 4. サービス起動
docker-compose up -d
cd frontend && npm run dev
```

#### 認証フロー詳細

**開発モード（認証無効）:**
- すべてのAPIリクエストは認証チェックをスキップ
- ログイン/ログアウト機能は無効
- 商品閲覧、カート、注文などすべての機能にアクセス可能

**本番モード（AWS Cognito）:**

1. **トークン管理**
   - IDトークン（JWT）を`localStorage`に保存
   - API呼び出し時に`Authorization: Bearer <id_token>`ヘッダーを自動付与
   - `frontend/src/lib/api-client.ts`のAxiosインターセプターで実装

2. **バックエンド検証**
   - `backend/shared/auth/cognito.go`でJWKS検証
   - RSA署名検証 + クレーム検証（issuer, audience, expiration）
   - 検証成功時に`user_id`, `email`をコンテキストに設定

3. **セッション管理**
   - `frontend/src/contexts/AuthContext.tsx`で状態管理
   - ページリロード時に自動セッション復元
   - トークン期限切れ時は自動ログアウト

4. **未認証アクセス**
   - トークンなしのリクエストは「アノニマスユーザー」として処理
   - 商品閲覧などの公開機能は認証不要

#### 詳細ドキュメント

- [Cognito認証ガイド](docs/COGNITO_AUTH_GUIDE.md) - 実装の詳細解説
- [クイックスタート](docs/COGNITO_QUICKSTART.md) - 5分で始める認証機能

## 📚 API ドキュメント

各サービスの OpenAPI 仕様:
- Product Service: `backend/services/product-service/openapi.yaml`
- Cart Service: `backend/services/cart-service/openapi.yaml`
- Order Service: `backend/services/order-service/openapi.yaml`
- Point Service: `backend/services/point-service/openapi.yaml`
- Promotion Service: `backend/services/promotion-service/openapi.yaml`

### API使用例

#### 認証

すべての保護されたエンドポイントには、Authorization ヘッダーに JWT トークンを付与：

```bash
curl -H "Authorization: Bearer YOUR_JWT_TOKEN" http://localhost:8001/api/v1/...
```

#### Product Service

```bash
# 商品検索
curl "http://localhost:8001/api/v1/products/search?keyword=スニーカー&page=1&page_size=20"

# 商品詳細
curl "http://localhost:8001/api/v1/products/prod-123"

# カテゴリ別
curl "http://localhost:8001/api/v1/products/category/shoes"
```

#### Cart Service

```bash
# カート取得
curl -H "Authorization: Bearer YOUR_TOKEN" http://localhost:8002/api/v1/carts

# 商品追加
curl -X POST -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{"product_id":"prod-123","product_name":"Nike Air","price":12000,"quantity":2}' \
  http://localhost:8002/api/v1/carts/items
```

#### Order Service

```bash
# 注文作成
curl -X POST -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "items":[{"product_id":"prod-123","product_name":"Nike Air","price":12000,"quantity":2,"subtotal":24000}],
    "payment_method":"credit_card",
    "shipping_address":"東京都渋谷区...",
    "points_to_use":0
  }' \
  http://localhost:8003/api/v1/orders

# 注文履歴
curl -H "Authorization: Bearer YOUR_TOKEN" http://localhost:8003/api/v1/orders
```

#### Point Service

```bash
# 残高照会
curl -H "Authorization: Bearer YOUR_TOKEN" http://localhost:8004/api/v1/points/balance

# ポイント履歴
curl -H "Authorization: Bearer YOUR_TOKEN" http://localhost:8004/api/v1/points/transactions
```

#### Promotion Service

```bash
# クーポン検証
curl -X POST -H "Content-Type: application/json" \
  -d '{"code":"SUMMER2024","order_total":25000,"user_id":"user-123"}' \
  http://localhost:8005/api/v1/promotions/validate-coupon

# アクティブなクーポン一覧
curl http://localhost:8005/api/v1/promotions/coupons
```

## 🚀 CI/CD

このプロジェクトはGitHub Actionsを使用した包括的なCI/CDパイプラインを実装しています。

### ワークフロー

- **Backend CI**: Go サービスのビルド・テスト・Lint ([backend-ci.yml](.github/workflows/backend-ci.yml))
- **Frontend CI**: Next.js のビルド・テスト・型チェック ([frontend-ci.yml](.github/workflows/frontend-ci.yml))
- **OpenAPI Validation**: API仕様の検証と同期確認 ([openapi-validation.yml](.github/workflows/openapi-validation.yml))
- **Docker Build**: イメージビルドとセキュリティスキャン ([docker-build.yml](.github/workflows/docker-build.yml))

### ローカルでのCI実行

```bash
# Backend Lint
golangci-lint run ./backend/services/...

# Backend Tests
cd backend/services/cart-service
go test -v -race ./...

# Frontend Lint
cd frontend
npm run lint
npx tsc --noEmit

# OpenAPI検証
swagger-cli validate api-schema/openapi/cart-service.yaml
```

詳細は [CI/CD Documentation](.github/CI_README.md) を参照してください。

## 🌟 主な機能

### フロントエンド (8ページ)

| ページ | 機能 | URL |
|--------|------|-----|
| ホーム | 注目商品・カテゴリー表示 | `/` |
| 商品一覧 | 検索・フィルタ・ソート | `/products` |
| 商品詳細 | 詳細情報・カート追加 | `/products/[id]` |
| カート | 数量変更・削除 | `/cart` |
| チェックアウト | 注文確定・クーポン適用 | `/checkout` |
| 注文履歴 | 過去の注文一覧 | `/my-orders` |
| 注文詳細 | 注文状況確認 | `/orders/[id]` |
| ポイント管理 | 残高・履歴表示 | `/my-points` |

### バックエンド機能

#### Product Service (8001)
- Elasticsearch 全文検索
- Redis キャッシング
- 商品一覧・詳細取得

#### Cart Service (8002)
- DynamoDB カート永続化
- ElastiCache セッション管理
- 7日間自動削除

#### Order Service (8003)
- MySQL トランザクション管理
- SQS イベント送信
- 注文履歴管理

#### Point Service (8004)
- ポイント残高管理
- トランザクション履歴
- 有効期限管理

#### Promotion Service (8005)
- クーポン管理
- セール適用
- 割引計算

## 📖 ドキュメント

各ディレクトリに詳細なREADMEがあります：
- [backend/README.md](backend/README.md) - バックエンドサービス詳細
- [frontend/README.md](frontend/README.md) - フロントエンド詳細
- [api-schema/README.md](api-schema/README.md) - API仕様・型定義
- [.github/copilot-instructions.md](.github/copilot-instructions.md) - AI エージェント向けガイド

## 🛠 トラブルシューティング

### サービスが起動しない

```bash
# ヘルスチェック
docker-compose ps

# ログ確認
docker-compose logs <service-name>

# 特定のサービスを再起動
docker-compose restart <service-name>
```

### データベース接続エラー

```bash
# MySQL接続確認
docker exec ec-mysql mysql -uroot -prootpassword -e "SELECT 1"

# Redis接続確認
docker exec ec-redis redis-cli ping

# Elasticsearch接続確認
curl http://localhost:9200/_cluster/health

# Kibana起動確認（ブラウザでアクセス）
# http://localhost:5601
```

### MySQL初期化エラー

```bash
# ボリュームをリセット
docker-compose down -v
docker-compose up -d

# ログ確認
docker logs ec-mysql
```

### ポート競合エラー

以下のポートが使用されていないか確認：
- 3000 (フロントエンド)
- 3306 (MySQL)
- 5601 (Kibana)
- 6379 (Redis)
- 8001-8005 (バックエンドサービス)
- 9200 (Elasticsearch)
- 4566 (LocalStack)

### Elasticsearch起動失敗

```bash
# メモリ不足の可能性
docker logs ec-elasticsearch

# Linuxの場合、vm.max_map_countを増やす
sudo sysctl -w vm.max_map_count=262144
```

### ビルドエラー

```bash
# Go依存関係更新
cd backend/services/<service>
go mod tidy

# フロントエンド依存関係更新
cd frontend
npm install

# Dockerイメージを完全再ビルド
docker-compose build --no-cache
```

### ホットリロードが効かない（Windows）

```bash
# docker-compose.dev.ymlを使用
docker-compose -f docker-compose.dev.yml up web
```

## 🤝 開発ワークフロー

### 新しいエンドポイント追加

1. `domain/{model}.go` - ビジネスモデル定義
2. `infrastructure/{db}/{repo}.go` - DB 実装
3. `usecase/{model}.go` - ビジネスロジック
4. `interface/handler/{model}.go` - HTTP ハンドラ
5. `main.go` - ルート登録

### 新しいサービス追加

1. `backend/services/new-service/` を作成
2. 既存サービスの構造をテンプレートとしてコピー
3. `docker-compose.yml` に定義追加
4. `build.sh` にサービスパス追加
5. `api-schema/openapi/new-service.yaml` を作成

## 🌍 環境変数管理のベストプラクティス

### ローカル開発
```bash
.env                    # Gitにコミット済み（開発用デフォルト値）
```

### 本番環境
```bash
.env.production         # Gitにコミット禁止（.gitignoreに設定済み）
.env.production.example # テンプレート（Gitにコミット可）
```

### Docker Composeでの使用

```bash
# ローカル開発（.envを自動読み込み）
docker-compose up -d

# 本番環境（明示的に.env.productionを指定）
docker-compose --env-file .env.production up -d

# ステージング環境
docker-compose --env-file .env.staging up -d
```

### セキュリティ推奨事項

1. **本番環境の秘密情報は環境変数で管理**
   - AWS Secrets Manager
   - Kubernetes Secrets
   - GitHub Actions Secrets

2. **開発環境の`.env`には本物の秘密情報を含めない**
   - ダミー値やLocalStack用の値のみ
   - `development-secret-key`のような明示的な名前

3. **環境変数の優先順位を理解**
   ```
   シェル環境変数 > .env > docker-compose.ymlのデフォルト値
   ```

## 📝 ライセンス

このプロジェクトはサンプル実装です。

## 🙏 謝辞

このプロジェクトは、マイクロサービスアーキテクチャと Clean Architecture のベストプラクティスを示すために作成されました。

ブラウザで以下にアクセス:
- フロントエンド: http://localhost:3000
- Product API: http://localhost:8001
- Cart API: http://localhost:8002
- Order API: http://localhost:8003
- Point API: http://localhost:8004
- Promotion API: http://localhost:8005

### 自動ビルド

```bash
# すべてのサービスをビルド
bash build.sh

# Dockerイメージもビルド
bash build.sh --docker
```

## OpenAPI

各サービスのOpenAPI定義は `services/{service}/openapi.yaml` に配置
