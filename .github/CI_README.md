# CI/CD Documentation

## Overview

このプロジェクトは GitHub Actions を使用して包括的なCI/CDパイプラインを実装しています。

## ワークフロー一覧

### 1. CI (統合ワークフロー)
**ファイル**: [.github/workflows/ci.yml](.github/workflows/ci.yml)

全てのCIワークフローを統合し、プルリクエスト時に自動実行されます。

### 2. Backend CI
**ファイル**: [.github/workflows/backend-ci.yml](.github/workflows/backend-ci.yml)

Go製マイクロサービスのビルド、テスト、Lintを実行します。

**実行内容**:
- 各サービス（cart, order, point, product, promotion）のテスト実行
- `go vet` による静的解析
- Race detectorを有効にしたテスト
- カバレッジレポート生成
- golangci-lint による包括的なLint
- OpenAPIコード生成の検証

**トリガー**:
- `backend/**` 配下のファイル変更時
- main, developブランチへのpush/PR

### 3. Frontend CI
**ファイル**: [.github/workflows/frontend-ci.yml](.github/workflows/frontend-ci.yml)

Next.js フロントエンドのビルド、テスト、型チェックを実行します。

**実行内容**:
- ESLint による静的解析
- TypeScript型チェック
- Next.js ビルド検証
- テストの実行（設定されている場合）
- API型定義の同期確認

**トリガー**:
- `frontend/**`, `api-schema/**` 配下のファイル変更時
- main, developブランチへのpush/PR

### 4. OpenAPI Validation
**ファイル**: [.github/workflows/openapi-validation.yml](.github/workflows/openapi-validation.yml)

OpenAPI仕様の検証と同期確認を行います。

**実行内容**:
- 各サービスのOpenAPIスキーマ検証（Swagger CLI使用）
- api-schema とサービスディレクトリ間のスキーマ同期確認
- Spectral による OpenAPI Linting

**トリガー**:
- `api-schema/openapi/**`, `backend/services/*/openapi.yaml` 変更時
- main, developブランチへのpush/PR

### 5. Docker Build
**ファイル**: [.github/workflows/docker-build.yml](.github/workflows/docker-build.yml)

Dockerイメージのビルドとセキュリティスキャンを実行します。

**実行内容**:
- 全サービスのDockerイメージビルド
- フロントエンドのDockerイメージビルド
- docker-compose設定ファイルの検証
- Trivy によるセキュリティスキャン（脆弱性検出）

**トリガー**:
- `backend/**`, `frontend/**`, `docker-compose.yml` 変更時
- main, developブランチへのpush/PR

## ローカルでの実行方法

### Backend Lint
```bash
# golangci-lintのインストール
brew install golangci-lint  # macOS
# or
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# 各サービスでLint実行
cd backend/services/cart-service
golangci-lint run

# または全サービス一括
for service in backend/services/*/; do
  echo "Linting $service"
  (cd "$service" && golangci-lint run)
done
```

### Backend Tests
```bash
# 各サービスでテスト実行
cd backend/services/order-service
go test -v -race ./...

# カバレッジ付き
go test -v -race -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Frontend Lint & Build
```bash
cd frontend

# Lint
npm run lint

# 型チェック
npx tsc --noEmit

# ビルド
npm run build
```

### OpenAPI Validation
```bash
# Swagger CLIのインストール
npm install -g @apidevtools/swagger-cli

# 検証
swagger-cli validate api-schema/openapi/cart-service.yaml

# Spectralのインストール
npm install -g @stoplight/spectral-cli

# Lint
spectral lint api-schema/openapi/cart-service.yaml
```

### Docker Build
```bash
# 各サービスのビルド
docker build -t ec-cart-service ./backend/services/cart-service

# docker-compose検証
docker compose config

# セキュリティスキャン（Trivyインストール必要）
brew install trivy  # macOS
trivy image ec-cart-service:latest
```

## Dependabot

依存関係の自動更新を設定済みです。

**設定ファイル**: [.github/dependabot.yml](.github/dependabot.yml)

**更新対象**:
- Go modules（各サービス + shared）
- npm packages（frontend, api-schema）
- Docker images
- GitHub Actions

**更新頻度**: 毎週

## golangci-lint設定

**設定ファイル**: [.golangci.yml](.golangci.yml)

有効なLinter:
- `errcheck`: エラーチェック漏れ検出
- `gosimple`: コード簡略化提案
- `govet`: 静的解析
- `staticcheck`: 詳細な静的解析
- `gofmt`, `goimports`: フォーマット・import整理
- `gosec`: セキュリティチェック
- `revive`: 包括的なLint
- その他多数

## トラブルシューティング

### CIが失敗する場合

1. **OpenAPI生成エラー**
   ```bash
   cd backend/services/cart-service
   make generate
   git add generated/
   ```

2. **golangci-lint エラー**
   ```bash
   golangci-lint run --fix
   ```

3. **型チェックエラー（Frontend）**
   ```bash
   cd frontend
   npm run lint -- --fix
   ```

4. **Docker build失敗**
   ```bash
   # キャッシュクリア
   docker builder prune -a
   docker build --no-cache -t service-name .
   ```

### よくあるエラー

**"Generated files are out of sync"**
- OpenAPI仕様変更後に `make generate` を実行してコミットしてください

**"Schemas are out of sync"**
- `api-schema/openapi/` と `backend/services/*/openapi.yaml` が一致していることを確認してください

**"Go module checksum mismatch"**
- `go mod tidy` を実行してください

## ベストプラクティス

1. **PRを作成する前に**:
   - ローカルでLintとテストを実行
   - OpenAPI変更時は必ず `make generate` を実行
   - `go mod tidy` でモジュール整理

2. **コミット前のチェックリスト**:
   - [ ] `golangci-lint run` 成功
   - [ ] `go test ./...` 成功
   - [ ] `npm run lint` 成功（frontend）
   - [ ] `npm run build` 成功（frontend）

3. **PR作成時**:
   - テンプレートに従って記入
   - 適切なラベルを付与
   - 関連Issueをリンク

## 参考リンク

- [GitHub Actions ドキュメント](https://docs.github.com/en/actions)
- [golangci-lint](https://golangci-lint.run/)
- [Trivy](https://github.com/aquasecurity/trivy)
- [Spectral](https://stoplight.io/open-source/spectral)
