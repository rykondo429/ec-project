# GitHub Actions セットアップガイド

## バッジの設定

README.mdのCIバッジを有効にするには、以下の手順を実施してください：

1. GitHubリポジトリのURLに合わせて、README.mdの以下の行を更新：

```markdown
[![CI](https://github.com/YOUR_USERNAME/ec-sample/workflows/CI/badge.svg)](https://github.com/YOUR_USERNAME/ec-sample/actions)
```

`YOUR_USERNAME` を実際のGitHubユーザー名またはOrganization名に置き換えてください。

例：
```markdown
[![CI](https://github.com/johndoe/ec-sample/workflows/CI/badge.svg)](https://github.com/johndoe/ec-sample/actions)
```

## 初回push後の確認

1. リポジトリにコードをpushすると、GitHub Actionsが自動的に実行されます
2. GitHubリポジトリの「Actions」タブでワークフローの実行状況を確認できます
3. 各ワークフローのステータス（成功/失敗）がバッジに反映されます

## Codecov（オプション）

カバレッジレポートをCodecovにアップロードする場合：

1. https://codecov.io/ でアカウントを作成
2. リポジトリを連携
3. GitHubリポジトリの Settings > Secrets and variables > Actions で以下を設定：
   - `CODECOV_TOKEN`: Codecovから取得したトークン

## GitHub Secretsの設定（本番環境用）

本番環境のデプロイに必要なシークレットを設定：

1. GitHubリポジトリの Settings > Secrets and variables > Actions
2. 以下のシークレットを追加（必要に応じて）：
   - `AWS_ACCESS_KEY_ID`
   - `AWS_SECRET_ACCESS_KEY`
   - `DOCKER_USERNAME`
   - `DOCKER_PASSWORD`
   - その他のシークレット情報

## トラブルシューティング

### ワークフローが実行されない

- `.github/workflows/` ディレクトリがリポジトリのルートに存在することを確認
- YAMLファイルのシンタックスエラーを確認
- リポジトリのSettings > Actions > General で「Allow all actions」が有効になっていることを確認

### テストが失敗する

- ローカルで `go test ./...` を実行してテストが通ることを確認
- 依存関係が正しく定義されているか確認（go.mod）
- テストに必要な環境変数が設定されているか確認

### Lintエラー

- ローカルで `golangci-lint run` を実行してエラーを修正
- `.golangci.yml` の設定を確認
