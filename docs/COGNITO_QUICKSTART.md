# Cognito認証 - クイックスタートガイド

このガイドに従って、AWS Cognito認証を5分で起動できます。

## 前提条件

- Docker & Docker Compose インストール済み
- Go 1.21+ インストール済み
- Node.js 18+ インストール済み
- AWS CLI インストール済み

## ステップ 1: インフラストラクチャの起動

```bash
# プロジェクトルートで実行
docker-compose up -d

# 全サービスがhealthyになるまで待機（約30秒）
docker-compose ps
```

## ステップ 2: Cognito User Pool IDとClient IDを取得

init-aws.shが自動実行され、以下のような出力が表示されます：

```
===================================================================
AWS resources initialized successfully!
===================================================================
Cognito User Pool ID: us-east-1_12345ABCD
Cognito Client ID: abcdefg1234567890
Test User: testuser@example.com / TestPass123!
===================================================================
```

**重要**: この値をコピーしてください！

## ステップ 3: 環境変数の設定

### バックエンド

プロジェクトルートに `.env` ファイルを作成：

```bash
cat > .env << EOF
# Database
DB_HOST=mysql
DB_PORT=3306
DB_NAME=ecsite
DB_USER=ecuser
DB_PASSWORD=ecpassword
DB_ROOT_PASSWORD=rootpassword

# Redis
REDIS_ADDR=redis:6379

# Elasticsearch
ELASTICSEARCH_URL=http://elasticsearch:9200

# AWS
AWS_REGION=us-east-1
AWS_ACCESS_KEY_ID=test
AWS_SECRET_ACCESS_KEY=test
AWS_ENDPOINT_URL=http://localstack:4566

# Cognito (ステップ2で取得した値を設定)
COGNITO_USER_POOL_ID=us-east-1_12345ABCD
COGNITO_CLIENT_ID=abcdefg1234567890
COGNITO_REGION=us-east-1
COGNITO_JWKS_URL=http://localhost:4566/

# JWT (レガシー、必要に応じて)
JWT_SECRET_KEY=your-secret-key-here
EOF
```

### フロントエンド

```bash
cat > frontend/.env.local << EOF
# API URLs
NEXT_PUBLIC_API_BASE_URL=http://localhost:8001
NEXT_PUBLIC_CART_API_URL=http://localhost:8002
NEXT_PUBLIC_ORDER_API_URL=http://localhost:8003
NEXT_PUBLIC_POINT_API_URL=http://localhost:8004
NEXT_PUBLIC_PROMOTION_API_URL=http://localhost:8005

# Cognito (ステップ2で取得した値を設定)
NEXT_PUBLIC_COGNITO_USER_POOL_ID=us-east-1_12345ABCD
NEXT_PUBLIC_COGNITO_CLIENT_ID=abcdefg1234567890
NEXT_PUBLIC_COGNITO_REGION=us-east-1
EOF
```

## ステップ 4: バックエンドサービスの起動

別々のターミナルで各サービスを起動（または tmux/screen を使用）:

```bash
# Terminal 1 - Product Service
cd backend/services/product-service && go run main.go

# Terminal 2 - Cart Service
cd backend/services/cart-service && go run main.go

# Terminal 3 - Order Service
cd backend/services/order-service && go run main.go

# Terminal 4 - Point Service
cd backend/services/point-service && go run main.go

# Terminal 5 - Promotion Service
cd backend/services/promotion-service && go run main.go
```

**または、すべてをバックグラウンドで起動**:

```bash
cd backend/services/product-service && go run main.go > /tmp/product.log 2>&1 &
cd backend/services/cart-service && go run main.go > /tmp/cart.log 2>&1 &
cd backend/services/order-service && go run main.go > /tmp/order.log 2>&1 &
cd backend/services/point-service && go run main.go > /tmp/point.log 2>&1 &
cd backend/services/promotion-service && go run main.go > /tmp/promotion.log 2>&1 &
```

## ステップ 5: フロントエンドの起動

```bash
cd frontend

# 依存関係のインストール（初回のみ）
npm install

# 開発サーバー起動
npm run dev
```

## ステップ 6: 動作確認

### 1. ブラウザでアクセス

http://localhost:3000 を開く

### 2. ログインページにアクセス

ログインコンポーネントを使用するか、直接認証をテスト：

```typescript
// ブラウザのコンソールで実行
import { signIn } from '@/lib/cognito';

await signIn({
  email: 'testuser@example.com',
  password: 'TestPass123!'
});
```

### 3. APIテスト（curlで確認）

```bash
# まず、トークンを取得（簡易テスト用スクリプト）
# 注：実際の環境では、フロントエンドから取得してください

# カートAPIをテスト（トークンなし - anonymousユーザー）
curl http://localhost:8002/api/v1/carts

# 製品検索をテスト
curl "http://localhost:8001/api/v1/products/search?keyword=test&page=1&page_size=10"
```

## トラブルシューティング

### 問題: サービスが起動しない

```bash
# ログを確認
docker-compose logs localstack
docker-compose logs mysql

# サービスを再起動
docker-compose restart
```

### 問題: Cognito User Pool IDが見つからない

```bash
# LocalStackのログを確認
docker logs ec-localstack | grep "User Pool ID"

# 手動でUser Poolを作成
bash init-aws.sh
```

### 問題: トークン検証エラー

1. 環境変数が正しく設定されているか確認
2. LocalStackが正しく起動しているか確認
3. COGNITO_JWKS_URLが`http://localhost:4566/`に設定されているか確認

### 問題: フロントエンドでCognito not configured

```bash
# 環境変数を確認
cd frontend
cat .env.local

# Next.jsを再起動
npm run dev
```

## 次のステップ

- [詳細なガイド](./COGNITO_AUTH_GUIDE.md)を参照
- ユーザー登録機能の実装
- パスワードリセット機能の実装
- ソーシャルログインの追加

## テストユーザー

デフォルトで以下のテストユーザーが作成されています：

- **メール**: `testuser@example.com`
- **パスワード**: `TestPass123!`

## よくある質問

**Q: 本番環境ではどうすればいい？**
A: `AWS_ENDPOINT_URL`を削除し、実際のAWS Cognitoの認証情報を設定してください。

**Q: 複数のユーザーを作成するには？**
A: AWS CLIを使用：
```bash
aws cognito-idp admin-create-user \
  --endpoint-url http://localhost:4566 \
  --user-pool-id <USER_POOL_ID> \
  --username newuser@example.com \
  --user-attributes Name=email,Value=newuser@example.com \
  --temporary-password "TempPass123!" \
  --region us-east-1
```

**Q: トークンの有効期限は？**
A: デフォルトでは1時間です。リフレッシュトークンを使用して更新できます。

## サポート

問題が解決しない場合は、以下を確認してください：
- [詳細ガイド](./COGNITO_AUTH_GUIDE.md)
- [LocalStack Documentation](https://docs.localstack.cloud/)
- プロジェクトのIssuesページ
