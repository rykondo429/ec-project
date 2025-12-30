# LocalStack Cognito制限と回避策

## 問題

LocalStack無料版ではAWS Cognito IDPサービスがサポートされていません。

```bash
$ ./scripts/setup-cognito-env.sh

An error occurred (InternalFailure) when calling the ListUserPools operation: 
The API for service 'cognito-idp' is either not included in your current license 
plan or has not yet been emulated by LocalStack.
```

## 解決策

### オプション1: 開発モード（認証無効）- 推奨

最も簡単な方法です。認証UIは実装されていますが、実際の認証処理はスキップされます。

**メリット:**
- 設定不要、即座に開発開始可能
- AWS アカウント不要
- すべてのページにアクセス可能

**デメリット:**
- 認証フローのテストができない
- 本番環境と異なる動作

**セットアップ:**

```bash
# 1. 環境変数設定
./scripts/setup-cognito-env.sh
# 選択: 1) 開発モード（認証機能無効）

# 2. 起動
docker-compose up -d
cd frontend && npm run dev

# 3. アクセス
open http://localhost:3000
```

**動作:**
- `NEXT_PUBLIC_AUTH_ENABLED=false`が設定される
- AuthContextは常に`user=null`を返す
- ログイン画面にアクセスするとエラーメッセージ表示
- すべてのAPIリクエストは認証チェックをスキップ

### オプション2: AWS Cognito本番環境

実際のAWS Cognitoを使用します。

**メリット:**
- 完全な認証フローのテスト可能
- 本番環境と同じ動作

**デメリット:**
- AWSアカウントが必要
- 設定手順が増える
- AWS利用料金が発生する可能性（無料枠内で収まる）

**セットアップ:**

```bash
# 1. AWS CognitoでUser Pool作成
aws cognito-idp create-user-pool \
  --pool-name ec-sample-dev \
  --policies "PasswordPolicy={MinimumLength=8,RequireUppercase=true,RequireLowercase=true,RequireNumbers=true}" \
  --auto-verified-attributes email \
  --username-attributes email

# 出力からUserPoolIdを保存
# 例: us-east-1_AbCdEfGhI

# 2. App Client作成
aws cognito-idp create-user-pool-client \
  --user-pool-id us-east-1_AbCdEfGhI \
  --client-name web-client \
  --explicit-auth-flows ALLOW_USER_PASSWORD_AUTH ALLOW_REFRESH_TOKEN_AUTH ALLOW_USER_SRP_AUTH

# 出力からClientIdを保存
# 例: 1a2b3c4d5e6f7g8h9i0j

# 3. 環境変数設定
./scripts/setup-cognito-env.sh
# 選択: 2) AWS Cognito本番環境
# User Pool ID: us-east-1_AbCdEfGhI
# Client ID: 1a2b3c4d5e6f7g8h9i0j
# Region: us-east-1

# 4. バックエンド環境変数設定
# .envファイルに追加:
cat >> .env << EOF
COGNITO_REGION=us-east-1
COGNITO_USER_POOL_ID=us-east-1_AbCdEfGhI
COGNITO_CLIENT_ID=1a2b3c4d5e6f7g8h9i0j
EOF

# 5. 起動
docker-compose down
docker-compose up -d
cd frontend && npm run dev

# 6. テストユーザー作成
aws cognito-idp admin-create-user \
  --user-pool-id us-east-1_AbCdEfGhI \
  --username testuser@example.com \
  --user-attributes Name=email,Value=testuser@example.com Name=email_verified,Value=true \
  --temporary-password "TempPass123!" \
  --message-action SUPPRESS

# 7. パスワード永続化
aws cognito-idp admin-set-user-password \
  --user-pool-id us-east-1_AbCdEfGhI \
  --username testuser@example.com \
  --password "TestPass123!" \
  --permanent
```

### オプション3: LocalStack Pro（有料）

LocalStack Proライセンスを購入すると、Cognitoがサポートされます。

**料金:** $35/月〜

**参考:** https://localstack.cloud/pricing

## 推奨フロー

開発段階に応じて使い分けることをお勧めします：

1. **初期開発（UI実装）**
   - オプション1: 開発モード（認証無効）
   - 認証UIの見た目や動作を確認

2. **認証フローテスト**
   - オプション2: AWS Cognito本番環境
   - 実際のログイン/ログアウトをテスト

3. **本番デプロイ**
   - オプション2: AWS Cognito本番環境
   - 本番用のUser Poolを作成

## トラブルシューティング

### Q: 開発モードでログイン画面にアクセスするとエラーが出る

**A:** 正常な動作です。`NEXT_PUBLIC_AUTH_ENABLED=false`の場合、ログイン処理はエラーを返します。

エラーメッセージ:
```
認証機能は現在無効です。.env.localでNEXT_PUBLIC_AUTH_ENABLED=trueを設定してください。
```

開発モードでは認証画面を使用せず、直接各ページにアクセスしてください。

### Q: AWS Cognitoでサインアップしたユーザーがログインできない

**A:** メールアドレスの確認が必要です。

```bash
# 確認コードをスキップ（開発用）
aws cognito-idp admin-confirm-sign-up \
  --user-pool-id <USER_POOL_ID> \
  --username <email>
```

本番環境では、Cognitoがメールで確認コードを送信します。

## 参考リンク

- [LocalStack Coverage](https://docs.localstack.cloud/references/coverage/)
- [AWS Cognito Documentation](https://docs.aws.amazon.com/cognito/)
- [EC-Sample Cognito UI Guide](./COGNITO_UI_GUIDE.md)
