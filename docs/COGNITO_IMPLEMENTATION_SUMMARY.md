# AWS Cognito OpenID Connect 認証実装 - 完了サマリー

## 実装内容

AWS Cognito + OpenID Connect (OIDC) 認証を EC-Sample プロジェクトに完全統合しました。

## 追加・変更されたファイル

### バックエンド（Go）

**新規ファイル**:
- `backend/shared/auth/cognito.go` - Cognito認証ミドルウェア（JWKSキャッシュ、トークン検証）
- `backend/shared/auth/config.go` - Cognito設定ヘルパー関数

**更新ファイル**:
- `backend/services/cart-service/main.go` - Cognito認証統合
- `backend/services/product-service/main.go` - Cognito認証統合
- `backend/services/order-service/main.go` - Cognito認証統合
- `backend/services/point-service/main.go` - Cognito認証統合
- `backend/services/promotion-service/main.go` - Cognito認証統合

### フロントエンド（Next.js + TypeScript）

**新規ファイル**:
- `frontend/src/lib/cognito.ts` - Cognito認証ライブラリ（サインイン、サインアウト、トークン管理）
- `frontend/src/contexts/AuthContext.tsx` - 認証コンテキストプロバイダー
- `frontend/src/components/LoginForm.tsx` - ログインフォームコンポーネント

**更新ファイル**:
- `frontend/src/lib/api-client.ts` - IDトークン自動付与機能
- `frontend/package.json` - amazon-cognito-identity-js 追加

### インフラストラクチャ

**更新ファイル**:
- `init-aws.sh` - Cognito User Pool/Client作成、テストユーザー作成
- `docker-compose.yml` - 全サービスにCognito環境変数追加
- `.env.example` - Cognito設定項目追加
- `frontend/.env.example` - Cognito設定項目追加

### ドキュメント

**新規ファイル**:
- `docs/COGNITO_AUTH_GUIDE.md` - 詳細実装ガイド
- `docs/COGNITO_QUICKSTART.md` - 5分クイックスタート
- `.env.cognito` - Cognito環境変数テンプレート（参考用）

## 主要機能

### 1. **バックエンド認証ミドルウェア**

✅ OpenID Connect準拠のトークン検証
✅ JWKSキャッシング（パフォーマンス最適化）
✅ クレーム検証（issuer, audience, expiration）
✅ RSA署名検証
✅ アノニマスユーザーサポート（トークンなしアクセス）

### 2. **フロントエンド認証SDK**

✅ Cognito統合（amazon-cognito-identity-js使用）
✅ サインイン/サインアウト
✅ セッション管理
✅ トークン自動更新
✅ API呼び出し時の自動トークン付与

### 3. **LocalStack統合**

✅ ローカル開発環境でのCognitoエミュレート
✅ 自動User Pool/Client作成
✅ テストユーザー自動作成（`testuser@example.com` / `TestPass123!`）

### 4. **全サービス対応**

✅ Product Service (8001)
✅ Cart Service (8002)
✅ Order Service (8003)
✅ Point Service (8004)
✅ Promotion Service (8005)

## 使用方法

### クイックスタート

```bash
# 1. インフラ起動
docker-compose up -d

# 2. User Pool ID/Client IDを取得（コンソール出力を確認）

# 3. 環境変数設定（.envファイルを作成）

# 4. バックエンド起動
cd backend/services/cart-service && go run main.go

# 5. フロントエンド起動
cd frontend && npm install && npm run dev

# 6. ブラウザで http://localhost:3000 にアクセス
# テストユーザーでログイン: testuser@example.com / TestPass123!
```

詳細は `docs/COGNITO_QUICKSTART.md` を参照してください。

### APIリクエスト例

```bash
# トークン付きリクエスト
curl -H "Authorization: Bearer <ID_TOKEN>" \
  http://localhost:8002/api/v1/carts

# アノニマスアクセス（トークンなし）
curl http://localhost:8001/api/v1/products/search?keyword=test
```

### プログラムからの使用

**フロントエンド（TypeScript）**:
```typescript
import { signIn, getCurrentUser } from '@/lib/cognito';

// ログイン
await signIn({ 
  email: 'testuser@example.com', 
  password: 'TestPass123!' 
});

// ユーザー情報取得
const user = await getCurrentUser();
// { username: "testuser@example.com", email: "...", sub: "..." }
```

**バックエンド（Go）**:
```go
func (h *Handler) GetCart(c echo.Context) error {
    userID := c.Get("user_id").(string)  // Cognito sub
    email := c.Get("email").(string)
    
    cart, err := h.useCase.GetCart(ctx, userID)
    // ...
}
```

## 環境変数

### 必須設定（バックエンド）

```bash
COGNITO_USER_POOL_ID=us-east-1_XXXXXXXXX
COGNITO_CLIENT_ID=XXXXXXXXXXXXXXXXXXXXXXXXXX
COGNITO_REGION=us-east-1
COGNITO_JWKS_URL=http://localhost:4566/  # LocalStackの場合のみ
```

### 必須設定（フロントエンド）

```bash
NEXT_PUBLIC_COGNITO_USER_POOL_ID=us-east-1_XXXXXXXXX
NEXT_PUBLIC_COGNITO_CLIENT_ID=XXXXXXXXXXXXXXXXXXXXXXXXXX
NEXT_PUBLIC_COGNITO_REGION=us-east-1
```

## セキュリティ考慮事項

✅ トークン署名検証（RS256）
✅ クレーム検証（iss, aud, exp）
✅ JWKSキャッシュ（TTL: 1時間）
✅ HTTPS推奨（本番環境）
⚠️ LocalStorageトークン保存（HTTPOnly Cookie推奨）

## 本番環境への移行

1. AWS Cognito User Poolを作成
2. 環境変数を更新（`AWS_ENDPOINT_URL`削除、`COGNITO_JWKS_URL`削除）
3. HTTPSを有効化
4. CORS設定を制限（ワイルドカード削除）

詳細は `docs/COGNITO_AUTH_GUIDE.md` の「本番環境デプロイ」セクションを参照。

## トラブルシューティング

### よくある問題

1. **`invalid token` エラー**
   - JWKS URLを確認
   - LocalStackが起動しているか確認
   - User Pool ID/Client IDが正しいか確認

2. **`User Pool not configured` エラー**
   - 環境変数が設定されているか確認
   - Next.jsを再起動

3. **LocalStack Cognitoが動作しない**
   - `docker logs ec-localstack` でログ確認
   - `bash init-aws.sh` を手動実行

詳細は `docs/COGNITO_AUTH_GUIDE.md` のトラブルシューティングセクションを参照。

## テスト

### テストユーザー

デフォルトで以下のテストユーザーが作成されます：

- **メール**: testuser@example.com
- **パスワード**: TestPass123!

### 追加ユーザーの作成

```bash
aws cognito-idp admin-create-user \
  --endpoint-url http://localhost:4566 \
  --user-pool-id <USER_POOL_ID> \
  --username newuser@example.com \
  --user-attributes Name=email,Value=newuser@example.com \
  --temporary-password "TempPass123!" \
  --region us-east-1

aws cognito-idp admin-set-user-password \
  --endpoint-url http://localhost:4566 \
  --user-pool-id <USER_POOL_ID> \
  --username newuser@example.com \
  --password "NewPass123!" \
  --permanent \
  --region us-east-1
```

## 参考ドキュメント

- **クイックスタート**: `docs/COGNITO_QUICKSTART.md`
- **詳細ガイド**: `docs/COGNITO_AUTH_GUIDE.md`
- **環境変数テンプレート**: `.env.example`, `frontend/.env.example`

## 次のステップ

実装可能な追加機能：

- [ ] ユーザー登録フロー（サインアップ）
- [ ] パスワードリセット機能
- [ ] MFA（多要素認証）
- [ ] ソーシャルログイン（Google, Facebook）
- [ ] リフレッシュトークン自動更新
- [ ] HTTPOnly Cookieへのトークン保存

## サポート

問題が発生した場合：
1. `docs/COGNITO_AUTH_GUIDE.md` のトラブルシューティングを確認
2. `docker logs ec-localstack` でログ確認
3. GitHub Issuesで質問

---

**実装完了日**: 2025年12月30日
**対応バージョン**: Go 1.21+, Next.js 15.1+, LocalStack latest
