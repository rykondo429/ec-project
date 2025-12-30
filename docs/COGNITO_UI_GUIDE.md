# AWS Cognito OIDC認証UI - 完全ガイド

## 概要

EC-SampleプロジェクトにAWS Cognito OpenID Connect (OIDC)による認証UIを実装しました。

## 実装済み機能

### 1. 会員登録 (`/signup`)

**画面構成:**
- 名前入力フィールド（任意）
- メールアドレス入力
- パスワード入力（要件表示付き）
- 確認コード入力（2ステップ目）

**フロー:**
```
1. ユーザーが情報を入力 → signUp()実行
2. Cognitoがメールアドレスに確認コードを送信
3. ユーザーがコードを入力 → confirmSignUp()実行
4. 確認完了後、ログイン画面へリダイレクト
```

**ファイル:** [frontend/src/app/signup/page.tsx](../frontend/src/app/signup/page.tsx)

### 2. ログイン (`/login`)

**画面構成:**
- メールアドレス入力
- パスワード入力
- パスワードを忘れた場合のリンク
- 新規登録へのリンク
- テストアカウント情報表示

**フロー:**
```
1. ユーザーが認証情報を入力 → signIn()実行
2. Cognitoが認証処理
3. 成功時、IDトークン取得
4. AuthContextでユーザー情報を更新
5. リダイレクト先へ遷移（デフォルト: トップページ）
```

**ファイル:** [frontend/src/app/login/page.tsx](../frontend/src/app/login/page.tsx)

### 3. パスワードリセット (`/forgot-password`)

**画面構成:**
- ステップ1: メールアドレス入力
- ステップ2: 確認コード + 新パスワード入力

**フロー:**
```
1. メールアドレス入力 → forgotPassword()実行
2. Cognitoが確認コードをメール送信
3. コードと新パスワードを入力 → confirmPassword()実行
4. 成功後、ログイン画面へリダイレクト
```

**ファイル:** [frontend/src/app/forgot-password/page.tsx](../frontend/src/app/forgot-password/page.tsx)

### 4. マイページ (`/my-page`)

**画面構成:**
- ユーザー情報カード（メールアドレス、会員ID）
- メニューグリッド:
  - 注文履歴
  - ポイント残高
  - お気に入り
  - 設定
- 最近のアクティビティ（今後実装予定）
- ログアウトボタン

**アクセス制御:**
```typescript
// 未認証の場合、ログイン画面へリダイレクト
useEffect(() => {
  if (!loading && !user) {
    router.push('/login?redirect=/my-page');
  }
}, [user, loading, router]);
```

**ファイル:** [frontend/src/app/my-page/page.tsx](../frontend/src/app/my-page/page.tsx)

### 5. ヘッダー認証UI統合

**表示内容:**
- **未ログイン時:**
  - ログインボタン（`/login`へリンク）
  - 新規登録ボタン（`/signup`へリンク）

- **ログイン時:**
  - ユーザーメール表示（クリックで`/my-page`へ）
  - ログアウトボタン

**実装:**
```tsx
{!loading && (
  <>
    {user ? (
      <div className="flex items-center gap-2">
        <Link href="/my-page" className="...">
          <User size={16} />
          <span>{user.email}</span>
        </Link>
        <button onClick={signOut} className="...">
          <LogOut size={20} />
        </button>
      </div>
    ) : (
      <div className="flex items-center gap-2">
        <Link href="/login">ログイン</Link>
        <Link href="/signup">新規登録</Link>
      </div>
    )}
  </>
)}
```

**ファイル:** [frontend/src/components/Header.tsx](../frontend/src/components/Header.tsx)

## アーキテクチャ

### フロントエンド層

#### 1. Cognito SDK ラッパー (`lib/cognito.ts`)

```typescript
// 提供する関数
export const signIn = (params: SignInParams): Promise<CognitoUserSession>
export const signOut = (): void
export const getCurrentSession = (): Promise<CognitoUserSession>
export const getCurrentUser = (): Promise<AuthUser | null>
export const getIdToken = (): Promise<string | null>
export const signUp = (params: SignUpParams): Promise<any>
export const confirmSignUp = (email: string, code: string): Promise<any>
export const forgotPassword = (email: string): Promise<any>
export const confirmPassword = (email: string, code: string, newPassword: string): Promise<any>
```

**依存ライブラリ:** `amazon-cognito-identity-js`

#### 2. 認証コンテキスト (`contexts/AuthContext.tsx`)

```typescript
interface AuthContextType {
  user: AuthUser | null;           // 現在のユーザー情報
  loading: boolean;                // 初期化中フラグ
  signIn: (params: SignInParams) => Promise<void>;
  signOut: () => void;
  refreshUser: () => Promise<void>;  // セッション更新
}
```

**責務:**
- グローバルな認証状態管理
- ページリロード時のセッション復元
- ログイン/ログアウト処理の統一インターフェース

#### 3. API クライアント統合 (`lib/api-client.ts`)

```typescript
// Axiosインターセプターで自動トークン付与
this.client.interceptors.request.use(async (config) => {
  const token = await getIdToken();
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});
```

### バックエンド層

#### Cognito認証ミドルウェア (`backend/shared/auth/cognito.go`)

```go
// 機能:
// 1. JWKSキャッシング（パフォーマンス最適化）
// 2. トークン検証（RSA署名 + クレーム）
// 3. ユーザー情報をコンテキストに設定

func CognitoMiddleware(region, userPoolID string) echo.MiddlewareFunc {
  // Authorization ヘッダーからトークン抽出
  // JWKSでRSA公開鍵取得
  // 署名検証 + クレーム検証
  // c.Set("user_id", claims.Sub)
  // c.Set("email", claims.Email)
}
```

**検証項目:**
- Issuer: `https://cognito-idp.{region}.amazonaws.com/{userPoolID}`
- Audience: Client ID
- Expiration: 有効期限内
- Signature: RSA署名

## セットアップ手順

### 重要: LocalStack制限について

**LocalStack無料版ではCognito IDPサービスがサポートされていません。**

以下の2つの方法から選択してください：

### 1. ローカル開発環境（認証機能なし - 推奨）

最も簡単な方法です。認証UIは表示されますが、実際の認証処理はスキップされます。

```bash
# ステップ1: Dockerでインフラ起動
docker-compose up -d

# ステップ2: 開発モード設定
./scripts/setup-cognito-env.sh
# 選択: 1) 開発モード（認証機能無効）

# 設定内容:
# - NEXT_PUBLIC_AUTH_ENABLED=false
# - 認証チェックをすべてスキップ
# - すべてのページにアクセス可能

# ステップ3: フロントエンド起動
cd frontend
npm install
npm run dev

# ステップ4: ブラウザでアクセス
# http://localhost:3000
# ログイン画面は表示されますが、認証なしで全機能が利用可能
```

### 2. 本番環境（AWS Cognito）

```bash
# ステップ1: AWS Cognitoでユーザープール作成
aws cognito-idp create-user-pool \
  --pool-name ec-production-pool \
  --policies "PasswordPolicy={MinimumLength=8,RequireUppercase=true,RequireLowercase=true,RequireNumbers=true}" \
  --auto-verified-attributes email

# ステップ2: アプリクライアント作成
aws cognito-idp create-user-pool-client \
  --user-pool-id <USER_POOL_ID> \
  --client-name web-client \
  --explicit-auth-flows ALLOW_USER_PASSWORD_AUTH ALLOW_REFRESH_TOKEN_AUTH

# ステップ3: 環境変数設定（frontend/.env.production）
NEXT_PUBLIC_COGNITO_USER_POOL_ID=us-east-1_RealPoolID
NEXT_PUBLIC_COGNITO_CLIENT_ID=RealClientID
NEXT_PUBLIC_COGNITO_REGION=us-east-1

# ステップ4: バックエンド環境変数（.env.production）
COGNITO_REGION=us-east-1
COGNITO_USER_POOL_ID=us-east-1_RealPoolID
COGNITO_CLIENT_ID=RealClientID
```

## テスト方法

### 1. 会員登録フロー

```bash
# ブラウザで実行:
1. http://localhost:3000/signup にアクセス
2. 名前: Test User
   メール: test@example.com
   パスワード: TestPass123!
3. 「会員登録」ボタンをクリック
4. LocalStackの場合、コンソールログで確認コードを確認:
   docker logs ec-localstack | grep "Confirmation code"
5. 確認コード入力 → ログイン画面へ
```

### 2. ログインフロー

```bash
1. http://localhost:3000/login にアクセス
2. テストアカウントで入力:
   メール: testuser@example.com
   パスワード: TestPass123!
3. 「ログイン」ボタンをクリック
4. トップページへリダイレクト
5. ヘッダーにユーザー名が表示されることを確認
```

### 3. API呼び出しテスト

```bash
# ブラウザの開発者ツールで確認:
1. ログイン後、任意のAPIリクエストを実行（例: カート取得）
2. Networkタブで Request Headers を確認:
   Authorization: Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9...

# バックエンドログで検証結果を確認:
docker logs ec-cart-service | grep "user_id"
# 出力例: user_id=abc-123-def email=testuser@example.com
```

## トラブルシューティング

### 問題: ログイン後もユーザー情報が表示されない

**原因:** セッション復元の失敗

**解決策:**
```bash
# ブラウザの開発者ツールで確認:
1. Application → Local Storage → http://localhost:3000
2. 以下のキーが存在するか確認:
   - CognitoIdentityServiceProvider.{ClientID}.LastAuthUser
   - CognitoIdentityServiceProvider.{ClientID}.{username}.idToken

# 存在しない場合:
- ログアウト → 再ログイン
- LocalStorageをクリア → 再ログイン
```

### 問題: LocalStackで確認コードが届かない

**原因:** LocalStackはメール送信をエミュレートしない

**解決策:**
```bash
# LocalStackのログから確認コードを取得:
docker logs ec-localstack 2>&1 | grep -A5 "Confirmation code"

# または、開発環境では確認スキップ:
aws cognito-idp admin-confirm-sign-up \
  --endpoint-url http://localhost:4566 \
  --user-pool-id <USER_POOL_ID> \
  --username test@example.com
```

### 問題: トークン検証エラー（Backend）

**エラーメッセージ:** `Invalid token signature`

**原因:** LocalStack JWKSとの不一致

**解決策:**
```bash
# LocalStackのCognito JWKSエンドポイントを確認:
curl http://localhost:4566/us-east-1/<USER_POOL_ID>/.well-known/jwks.json

# バックエンド環境変数を再確認:
echo $COGNITO_REGION
echo $COGNITO_USER_POOL_ID

# 不一致の場合、docker-compose.ymlを修正して再起動:
docker-compose down
docker-compose up -d
```

## ファイル一覧

### フロントエンド
```
frontend/src/
├── app/
│   ├── login/page.tsx              # ログイン画面
│   ├── signup/page.tsx             # 会員登録画面
│   ├── forgot-password/page.tsx    # パスワードリセット
│   ├── my-page/page.tsx            # マイページ
│   └── layout.tsx                  # AuthProvider統合
├── components/
│   └── Header.tsx                  # 認証UI統合
├── contexts/
│   └── AuthContext.tsx             # 認証状態管理
└── lib/
    ├── cognito.ts                  # Cognito SDK ラッパー
    └── api-client.ts               # トークン自動付与
```

### バックエンド
```
backend/shared/auth/
├── cognito.go                      # Cognito認証ミドルウェア
└── config.go                       # 設定ヘルパー
```

### インフラ・スクリプト
```
scripts/
└── setup-cognito-env.sh            # 環境変数自動設定

init-aws.sh                         # LocalStack Cognito初期化
docker-compose.yml                  # Cognito環境変数定義
```

## 次のステップ

1. **MFA（多要素認証）の実装**
   - SMS/TOTPによる2段階認証
   - `amazon-cognito-identity-js`のMFA機能を使用

2. **ソーシャルログインの追加**
   - Google/Facebook/Appleログイン
   - CognitoのIdentity Providersを設定

3. **メール検証のカスタマイズ**
   - AWS SESでカスタムメールテンプレート
   - ブランディングされた確認メール

4. **セキュリティ強化**
   - リフレッシュトークンのローテーション
   - トークン失効時の自動再認証

## 参考資料

- [Amazon Cognito Developer Guide](https://docs.aws.amazon.com/cognito/latest/developerguide/)
- [amazon-cognito-identity-js GitHub](https://github.com/aws-amplify/amplify-js/tree/main/packages/amazon-cognito-identity-js)
- [OpenID Connect Specification](https://openid.net/specs/openid-connect-core-1_0.html)
