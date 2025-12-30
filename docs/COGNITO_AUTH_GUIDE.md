# AWS Cognito OpenID Connect 認証実装ガイド

このドキュメントでは、EC-SampleプロジェクトにおけるAWS Cognito OpenID Connect認証の実装と使用方法について説明します。

## 概要

このプロジェクトでは、AWS Cognitoを使用したOpenID Connect (OIDC) 認証を実装しています。ローカル開発環境ではLocalStackを使用してCognitoをエミュレートします。

## アーキテクチャ

### バックエンド認証フロー

1. クライアントがCognitoからIDトークンを取得
2. APIリクエスト時に`Authorization: Bearer <ID_TOKEN>`ヘッダーを付与
3. バックエンドサービスがトークンを検証
   - JWKSエンドポイントから公開鍵を取得
   - トークンの署名を検証
   - クレーム（issuer, audience, expiration）を検証
4. 検証成功後、ユーザー情報をリクエストコンテキストに設定

### 実装コンポーネント

#### バックエンド（Go）

- **`backend/shared/auth/cognito.go`**: Cognito認証ミドルウェア
- **`backend/shared/auth/config.go`**: Cognito設定ヘルパー
- **各サービスの`main.go`**: 認証ミドルウェアの統合

#### フロントエンド（Next.js + TypeScript）

- **`frontend/src/lib/cognito.ts`**: Cognito認証ライブラリ
- **`frontend/src/contexts/AuthContext.tsx`**: 認証コンテキスト
- **`frontend/src/components/LoginForm.tsx`**: ログインフォーム
- **`frontend/src/lib/api-client.ts`**: APIクライアント（トークン自動付与）

## セットアップ手順

### 1. LocalStack環境の起動

```bash
# Docker Composeでインフラを起動
docker-compose up -d

# サービスのヘルスチェック
docker-compose ps
```

### 2. Cognito User Poolの作成

LocalStackが起動すると、`init-aws.sh`スクリプトが自動的に実行され、以下を作成します：

- Cognito User Pool
- User Pool Client
- テストユーザー (`testuser@example.com` / `TestPass123!`)

スクリプトの出力から、以下の値を取得してください：

```
Cognito User Pool ID: us-east-1_XXXXXXXXX
Cognito Client ID: XXXXXXXXXXXXXXXXXXXXXXXXXX
```

### 3. 環境変数の設定

#### バックエンド（.envファイル）

プロジェクトルートに`.env`ファイルを作成し、以下を設定：

```bash
# Cognito Configuration
COGNITO_USER_POOL_ID=us-east-1_XXXXXXXXX  # init-aws.shの出力から取得
COGNITO_CLIENT_ID=XXXXXXXXXXXXXXXXXXXXXXXXXX  # init-aws.shの出力から取得
COGNITO_REGION=us-east-1
COGNITO_JWKS_URL=http://localhost:4566/

# その他の環境変数（.env.exampleを参照）
DB_HOST=mysql
DB_PORT=3306
# ... (その他の設定)
```

#### フロントエンド（frontend/.env.local）

```bash
# API URLs
NEXT_PUBLIC_API_BASE_URL=http://localhost:8001
NEXT_PUBLIC_CART_API_URL=http://localhost:8002
NEXT_PUBLIC_ORDER_API_URL=http://localhost:8003
NEXT_PUBLIC_POINT_API_URL=http://localhost:8004
NEXT_PUBLIC_PROMOTION_API_URL=http://localhost:8005

# Cognito Configuration
NEXT_PUBLIC_COGNITO_USER_POOL_ID=us-east-1_XXXXXXXXX
NEXT_PUBLIC_COGNITO_CLIENT_ID=XXXXXXXXXXXXXXXXXXXXXXXXXX
NEXT_PUBLIC_COGNITO_REGION=us-east-1
```

### 4. バックエンドサービスの起動

各サービスを個別のターミナルで起動：

```bash
# Product Service
cd backend/services/product-service
go run main.go

# Cart Service
cd backend/services/cart-service
go run main.go

# Order Service
cd backend/services/order-service
go run main.go

# Point Service
cd backend/services/point-service
go run main.go

# Promotion Service
cd backend/services/promotion-service
go run main.go
```

### 5. フロントエンドの起動

```bash
cd frontend

# 依存関係のインストール（初回のみ）
npm install
# Cognito Identity SDKのインストール
npm install amazon-cognito-identity-js

# 開発サーバー起動
npm run dev
```

フロントエンドは http://localhost:3000 で起動します。

## 使用方法

### フロントエンドでのログイン

1. ブラウザで http://localhost:3000/login にアクセス
2. テストアカウントでログイン：
   - メール: `testuser@example.com`
   - パスワード: `TestPass123!`
3. ログイン成功後、IDトークンが自動的に取得され、APIリクエストに付与されます

### プログラムからの認証

#### サインイン

```typescript
import { signIn } from '@/lib/cognito';

try {
  const session = await signIn({
    email: 'testuser@example.com',
    password: 'TestPass123!',
  });
  console.log('Signed in:', session.getIdToken().getJwtToken());
} catch (error) {
  console.error('Sign in error:', error);
}
```

#### 現在のユーザー取得

```typescript
import { getCurrentUser } from '@/lib/cognito';

const user = await getCurrentUser();
if (user) {
  console.log('Current user:', user);
  // { username: "testuser@example.com", email: "testuser@example.com", sub: "..." }
}
```

#### サインアウト

```typescript
import { signOut } from '@/lib/cognito';

signOut();
```

### バックエンドでのユーザー情報取得

認証ミドルウェアが自動的にユーザー情報をコンテキストに設定します：

```go
func (h *Handler) GetCart(c echo.Context) error {
    userID := c.Get("user_id").(string)  // Cognitoのsub（ユーザーID）
    email := c.Get("email").(string)      // ユーザーのメールアドレス
    
    // ユーザー固有の処理
    cart, err := h.useCase.GetCart(c.Request().Context(), userID)
    // ...
}
```

### APIリクエストのテスト

#### curlでのテスト

```bash
# 1. トークンを取得（LocalStackの場合、簡易的な方法）
TOKEN="<ID_TOKEN>"  # フロントエンドのログイン後、開発者ツールで取得

# 2. 認証付きAPIリクエスト
curl -H "Authorization: Bearer $TOKEN" \
  http://localhost:8002/api/v1/carts
```

## トラブルシューティング

### トークン検証エラー

**問題**: `invalid token` エラーが発生する

**解決策**:
1. JWKS URLが正しいか確認
2. LocalStackが起動しているか確認
3. User Pool IDとClient IDが正しいか確認

```bash
# LocalStackのログを確認
docker logs ec-localstack

# Cognitoのエンドポイントをテスト
curl http://localhost:4566/_aws/cognito-idp/
```

### User Pool not configured エラー

**問題**: フロントエンドで `Cognito User Pool not configured` エラー

**解決策**:
1. 環境変数が正しく設定されているか確認
   ```bash
   echo $NEXT_PUBLIC_COGNITO_USER_POOL_ID
   echo $NEXT_PUBLIC_COGNITO_CLIENT_ID
   ```
2. `.env.local`ファイルがfrontendディレクトリに存在するか確認
3. Next.jsを再起動

### LocalStack Cognitoの制限事項

LocalStackのCognito実装には以下の制限があります：

- JWKSエンドポイントが本番環境と異なる場合がある
- 一部の高度な機能（MFA、カスタム認証フローなど）が未サポート
- トークンの有効期限が短い

**本番環境への移行時の注意**:
- `COGNITO_JWKS_URL`環境変数を削除（自動構築される）
- `AWS_ENDPOINT_URL`を削除
- 本物のAWS認証情報を設定

## 本番環境デプロイ

### AWS Cognito User Poolの作成

```bash
# AWS CLIで本番用User Poolを作成
aws cognito-idp create-user-pool \
  --pool-name ec-sample-production \
  --policies "PasswordPolicy={MinimumLength=8,RequireUppercase=true,RequireLowercase=true,RequireNumbers=true}" \
  --auto-verified-attributes email \
  --region us-east-1

# User Pool Clientの作成
aws cognito-idp create-user-pool-client \
  --user-pool-id <USER_POOL_ID> \
  --client-name ec-sample-web \
  --explicit-auth-flows ALLOW_USER_PASSWORD_AUTH ALLOW_REFRESH_TOKEN_AUTH \
  --region us-east-1
```

### 環境変数の更新

```bash
# 本番環境変数（.env.production）
COGNITO_USER_POOL_ID=<本番のUser Pool ID>
COGNITO_CLIENT_ID=<本番のClient ID>
COGNITO_REGION=us-east-1
# COGNITO_JWKS_URLは設定不要（自動構築）
# AWS_ENDPOINT_URLは削除

# AWSクレデンシャルは環境に応じて設定
AWS_REGION=us-east-1
AWS_ACCESS_KEY_ID=<本番キー>
AWS_SECRET_ACCESS_KEY=<本番シークレット>
```

## セキュリティベストプラクティス

1. **トークンの保存**
   - LocalStorageではなくHTTPOnly Cookieの使用を検討
   - 機密情報をトークンに含めない

2. **HTTPS必須**
   - 本番環境では必ずHTTPSを使用
   - トークンの送信は常に暗号化

3. **トークンのリフレッシュ**
   - アクセストークンの有効期限を短く設定
   - リフレッシュトークンを使用して自動更新

4. **CORS設定**
   - 本番環境では具体的なオリジンを指定
   - ワイルドカード（`*`）は使用しない

## 参考資料

- [AWS Cognito Documentation](https://docs.aws.amazon.com/cognito/)
- [OpenID Connect Specification](https://openid.net/specs/openid-connect-core-1_0.html)
- [LocalStack Cognito](https://docs.localstack.cloud/user-guide/aws/cognito/)
- [amazon-cognito-identity-js](https://github.com/aws-amplify/amplify-js/tree/main/packages/amazon-cognito-identity-js)
