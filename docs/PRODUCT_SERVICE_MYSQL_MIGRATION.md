# Product Service - MySQLマスターデータ管理への移行

## 📋 変更概要

Product ServiceをElasticsearch単独からMySQL主データソース + Elasticsearch検索インデックスの構成に変更しました。

## 🔄 変更内容

### 1. アーキテクチャ変更

**変更前**:
```
Product Service → Elasticsearch (データ + 検索)
```

**変更後**:
```
Product Service → MySQL (マスターデータ)
                → Elasticsearch (検索インデックス - 今後実装)
```

### 2. 新規追加ファイル

#### バックエンド
- `backend/services/product-service/infrastructure/mysql/product_repository.go`
  - MySQLリポジトリ実装
  - FULLTEXT検索対応
  - CRUD操作
  - ページネーション、フィルタリング、ソート

#### スクリプト
- `scripts/seed-products-mysql.sh`
  - MySQLへのサンプルデータ投入スクリプト
  - 12件の商品データ
  - カテゴリ: shoes, electronics, clothing

### 3. 変更ファイル

#### `backend/services/product-service/main.go`
- MySQL接続の追加
- `gorm.io/driver/mysql`を使用
- コネクションプール設定
- ヘルスチェックをMySQLに変更

#### `backend/services/product-service/domain/product.go`
- GORMタグの追加
- `TableName()` メソッド実装
- MySQL スキーマとの完全一致

#### `backend/services/product-service/go.mod`
- 依存関係追加:
  - `gorm.io/driver/mysql v1.5.2`
  - `gorm.io/gorm v1.25.5`

#### `README.md`
- データストアの説明を更新
- サンプルデータ投入手順を更新
- MySQLが主データソースであることを明記

## 📊 データベーススキーマ

```sql
CREATE TABLE products (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description LONGTEXT,
    price DECIMAL(10, 2) NOT NULL,
    stock INT NOT NULL DEFAULT 0,
    category VARCHAR(100),
    rating DECIMAL(3, 2) DEFAULT 0.0,
    review_count INT DEFAULT 0,
    is_sale BOOLEAN DEFAULT FALSE,
    sale_price DECIMAL(10, 2),
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_category (category),
    INDEX idx_is_active (is_active),
    FULLTEXT INDEX ft_name_description (name, description)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4;
```

## 🚀 セットアップ手順

### 1. MySQLマイグレーション実行

```bash
# スキーマ作成（初回のみ）
mysql -h127.0.0.1 -uecuser -pecpassword ecsite < backend/migrations/mysql/01_init_schema.sql
```

### 2. サンプルデータ投入

```bash
# 商品データ投入（12件）
./scripts/seed-products-mysql.sh
```

### 3. Product Service起動

```bash
# ローカル実行
cd backend/services/product-service
go run main.go

# または Docker
docker-compose up -d product-service
```

### 4. 動作確認

```bash
# ヘルスチェック
curl http://localhost:8001/health

# 商品一覧取得
curl 'http://localhost:8001/api/v1/products/search?page=1&page_size=20'

# キーワード検索（FULLTEXT）
curl 'http://localhost:8001/api/v1/products/search?keyword=nike'

# カテゴリフィルタ
curl 'http://localhost:8001/api/v1/products/search?category=shoes'

# 価格範囲
curl 'http://localhost:8001/api/v1/products/search?price_min=10000&price_max=50000'
```

## 📈 投入データ統計

```
総商品数: 12件
総在庫数: 555個
平均価格: 39,040円
カテゴリ: shoes (4), electronics (5), clothing (3)
```

### サンプル商品一覧

| ID | 商品名 | 価格 | 在庫 | カテゴリ | セール |
|----|--------|------|------|----------|--------|
| prod-001 | Nike Air Max 90 | 12,000円 | 50 | shoes | - |
| prod-002 | Adidas Ultraboost 22 | 18,000円 | 30 | shoes | ✓ 14,400円 |
| prod-003 | Sony WH-1000XM5 | 45,000円 | 20 | electronics | - |
| prod-004 | Apple AirPods Pro | 38,000円 | 45 | electronics | ✓ 33,000円 |
| prod-005 | Patagonia Fleece | 15,000円 | 60 | clothing | - |
| prod-006 | The North Face Nuptse | 35,000円 | 25 | clothing | - |
| prod-007 | MacBook Air M2 | 148,000円 | 15 | electronics | - |
| prod-008 | iPad Pro | 98,000円 | 22 | electronics | ✓ 88,000円 |
| prod-009 | New Balance 2002R | 16,000円 | 40 | shoes | - |
| prod-010 | ユニクロ ヒートテック | 1,990円 | 150 | clothing | ✓ 1,490円 |
| prod-011 | Bose QC Earbuds | 35,000円 | 18 | electronics | - |
| prod-012 | Converse Chuck Taylor | 6,500円 | 80 | shoes | - |

## 🔍 検索機能

### 実装済み

1. **FULLTEXT検索** - MySQLの全文検索インデックス
   - `name`と`description`フィールドを対象
   - 日本語対応

2. **フィルタリング**
   - カテゴリ
   - 価格範囲（最小〜最大）
   - アクティブ状態

3. **ソート**
   - 価格順（昇順/降順）
   - 評価順
   - 作成日順
   - スコア順（デフォルト）

4. **ページネーション**
   - ページ番号
   - ページサイズ
   - 総件数・総ページ数計算

### 今後の拡張予定

- [ ] Elasticsearch同期機能
- [ ] より高度な検索クエリ
- [ ] ファセット検索
- [ ] サジェスト機能
- [ ] 商品レコメンド

## ⚠️ 重要な注意点

### データの一貫性

**MySQLがマスターデータ**
- すべてのCRUD操作はMySQLに対して実行
- Elasticsearchは検索用のインデックスとして使用
- データの真実の源泉(Source of Truth)はMySQL

### 在庫管理

現在の実装:
- 商品レベルの在庫（`stock`カラム）
- SKU管理は未実装
- バリエーション（サイズ、色など）は未対応

今後の拡張:
- SKUテーブルの追加
- バリエーション管理
- 在庫引当ロジック

## 🔧 トラブルシューティング

### MySQLに接続できない

```bash
# コンテナ起動確認
docker-compose ps mysql

# ログ確認
docker-compose logs mysql

# 手動接続テスト
mysql -h127.0.0.1 -uecuser -pecpassword ecsite -e "SELECT 1"
```

### テーブルが存在しない

```bash
# マイグレーション実行
mysql -h127.0.0.1 -uecuser -pecpassword ecsite < backend/migrations/mysql/01_init_schema.sql

# テーブル確認
mysql -h127.0.0.1 -uecuser -pecpassword ecsite -e "SHOW TABLES"
```

### データが空

```bash
# サンプルデータ投入
./scripts/seed-products-mysql.sh

# 件数確認
mysql -h127.0.0.1 -uecuser -pecpassword ecsite -e "SELECT COUNT(*) FROM products"
```

## 📝 関連ドキュメント

- [README.md](../README.md) - プロジェクト全体のドキュメント
- [backend/README.md](../backend/README.md) - バックエンドアーキテクチャ
- [backend/migrations/mysql/01_init_schema.sql](../backend/migrations/mysql/01_init_schema.sql) - MySQLスキーマ定義
