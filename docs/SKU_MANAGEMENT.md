# SKU管理機能 - 実装ドキュメント

## 📋 概要

商品バリエーション（サイズ、色など）ごとの在庫管理を実現するSKU（Stock Keeping Unit）管理機能を実装しました。

## 🎯 主な機能

### 1. SKU基本管理
- ✅ SKU作成・更新・削除（論理削除）
- ✅ SKUコード検索
- ✅ 商品ごとのSKU一覧取得
- ✅ SKU一括作成

### 2. 在庫管理
- ✅ 在庫調整（増減）
- ✅ 低在庫SKU検索
- ✅ トランザクション制御による在庫整合性保証

### 3. バリエーション管理
- ✅ JSON型での柔軟な属性管理
- ✅ サイズ・色・その他カスタム属性対応

## 🗂 データベーススキーマ

### SKUsテーブル

```sql
CREATE TABLE skus (
    id VARCHAR(36) PRIMARY KEY,
    product_id VARCHAR(36) NOT NULL,           -- 商品ID（外部キー）
    sku_code VARCHAR(100) NOT NULL UNIQUE,     -- SKUコード
    name VARCHAR(255) NOT NULL,                -- SKU名
    attributes JSON,                           -- バリエーション属性
    price DECIMAL(10, 2) NOT NULL DEFAULT 0,   -- SKU価格
    stock INT NOT NULL DEFAULT 0,              -- 在庫数
    is_active BOOLEAN DEFAULT TRUE,            -- アクティブフラグ
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    INDEX idx_product_id (product_id),
    INDEX idx_sku_code (sku_code),
    INDEX idx_is_active (is_active),
    INDEX idx_product_sku (product_id, is_active),
    FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE
);
```

### SKU在庫履歴テーブル（将来拡張用）

```sql
CREATE TABLE sku_stock_history (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    sku_id VARCHAR(36) NOT NULL,
    quantity_change INT NOT NULL,              -- 在庫変動量
    stock_before INT NOT NULL,                 -- 変動前在庫
    stock_after INT NOT NULL,                  -- 変動後在庫
    reason VARCHAR(255),                       -- 理由
    created_by VARCHAR(100),                   -- 実行者
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    INDEX idx_sku_id (sku_id),
    INDEX idx_created_at (created_at),
    FOREIGN KEY (sku_id) REFERENCES skus(id) ON DELETE CASCADE
);
```

## 📁 実装ファイル

### ドメイン層

**[domain/sku.go](backend/services/product-service/domain/sku.go)**
- SKUモデル定義
- SKUAttrs（JSON型カスタムScanner/Valuer実装）
- リクエスト/レスポンス型定義

```go
type SKU struct {
    ID         string
    ProductID  string
    SKUCode    string
    Name       string
    Attributes SKUAttrs  // map[string]string
    Price      float64
    Stock      int
    IsActive   bool
    CreatedAt  time.Time
    UpdatedAt  time.Time
}

// JSON型のカスタムScanner実装
func (a *SKUAttrs) Scan(value interface{}) error
func (a SKUAttrs) Value() (driver.Value, error)
```

**[domain/sku_repository.go](backend/services/product-service/domain/sku_repository.go)**
- SKUリポジトリインターフェース定義

### インフラストラクチャ層

**[infrastructure/mysql/sku_repository.go](backend/services/product-service/infrastructure/mysql/sku_repository.go)**
- MySQL実装（全12メソッド）
- トランザクション制御
- 在庫調整ロジック

### ユースケース層

**[usecase/sku.go](backend/services/product-service/usecase/sku.go)**
- ビジネスロジック
- バリデーション
- ログ出力

### インターフェース層

**[interface/handler/sku.go](backend/services/product-service/interface/handler/sku.go)**
- HTTPハンドラー
- リクエスト/レスポンス変換

## 🚀 API エンドポイント

### SKU管理

| メソッド | エンドポイント | 説明 |
|---------|--------------|------|
| `POST` | `/api/v1/skus` | SKU作成 |
| `GET` | `/api/v1/skus/:id` | SKU取得 |
| `PUT` | `/api/v1/skus/:id` | SKU更新 |
| `DELETE` | `/api/v1/skus/:id` | SKU削除 |
| `GET` | `/api/v1/skus/code/:code` | SKUコード検索 |

### 商品-SKU連携

| メソッド | エンドポイント | 説明 |
|---------|--------------|------|
| `GET` | `/api/v1/products/:product_id/skus` | 商品のSKU一覧 |
| `GET` | `/api/v1/products/:product_id/with-skus` | 商品+SKU一括取得 |
| `POST` | `/api/v1/products/:product_id/skus/bulk` | SKU一括作成 |

### 在庫管理

| メソッド | エンドポイント | 説明 |
|---------|--------------|------|
| `POST` | `/api/v1/skus/stock/adjust` | 在庫調整 |
| `GET` | `/api/v1/skus/low-stock` | 低在庫SKU一覧 |

## 💡 使用例

### 1. 商品のSKU一覧取得

```bash
curl -s "http://localhost:8001/api/v1/products/prod-001/skus" | jq .

# レスポンス例
[
  {
    "id": "sku-001",
    "product_id": "prod-001",
    "sku_code": "NIKE-AM90-RED-S",
    "name": "Nike Air Max 90 - Red / S",
    "attributes": {
      "color": "Red",
      "size": "S"
    },
    "price": 12000,
    "stock": 15,
    "is_active": true
  },
  ...
]
```

### 2. 商品とSKUをまとめて取得

```bash
curl -s "http://localhost:8001/api/v1/products/prod-001/with-skus" | jq .

# レスポンス例
{
  "id": "prod-001",
  "name": "Nike Air Max 90",
  "price": 12000,
  "skus": [
    {
      "sku_code": "NIKE-AM90-BLK-L",
      "stock": 20,
      "attributes": {"color": "Black", "size": "L"}
    },
    ...
  ]
}
```

### 3. SKU作成

```bash
curl -X POST "http://localhost:8001/api/v1/skus" \
  -H "Content-Type: application/json" \
  -d '{
    "product_id": "prod-001",
    "sku_code": "NIKE-AM90-WHT-M",
    "name": "Nike Air Max 90 - White / M",
    "attributes": {
      "color": "White",
      "size": "M"
    },
    "price": 12000,
    "stock": 30
  }'
```

### 4. SKU一括作成

```bash
curl -X POST "http://localhost:8001/api/v1/products/prod-003/skus/bulk" \
  -H "Content-Type: application/json" \
  -d '[
    {
      "sku_code": "SONY-WH1000XM5-BLK",
      "name": "Sony WH-1000XM5 - Black",
      "attributes": {"color": "Black"},
      "price": 45000,
      "stock": 10
    },
    {
      "sku_code": "SONY-WH1000XM5-SLV",
      "name": "Sony WH-1000XM5 - Silver",
      "attributes": {"color": "Silver"},
      "price": 45000,
      "stock": 10
    }
  ]'
```

### 5. 在庫調整

```bash
# 在庫減少（販売）
curl -X POST "http://localhost:8001/api/v1/skus/stock/adjust" \
  -H "Content-Type: application/json" \
  -d '{
    "sku_id": "sku-001",
    "quantity": -5,
    "reason": "販売"
  }'

# 在庫追加（入荷）
curl -X POST "http://localhost:8001/api/v1/skus/stock/adjust" \
  -H "Content-Type: application/json" \
  -d '{
    "sku_id": "sku-001",
    "quantity": 10,
    "reason": "入荷"
  }'
```

### 6. 低在庫SKU検索

```bash
# 在庫15以下のSKUを取得
curl "http://localhost:8001/api/v1/skus/low-stock?threshold=15" | jq .

# レスポンス例
[
  {
    "id": "sku-004",
    "sku_code": "NIKE-AM90-BLK-S",
    "stock": 10,
    "attributes": {"color": "Black", "size": "S"}
  },
  ...
]
```

### 7. SKU更新

```bash
curl -X PUT "http://localhost:8001/api/v1/skus/sku-001" \
  -H "Content-Type: application/json" \
  -d '{
    "stock": 50,
    "price": 11500
  }'
```

### 8. SKU削除（論理削除）

```bash
curl -X DELETE "http://localhost:8001/api/v1/skus/sku-001"
```

## 📊 サンプルデータ

マイグレーションで以下のサンプルデータが投入されます：

### Nike Air Max 90（6 SKUs）
- Red / S, M, L
- Black / S, M, L

### Adidas Ultraboost 22（3 SKUs）
- White / S, M, L

### ユニクロ ヒートテック（5 SKUs）
- Black / S, M, L
- White / S, M

**合計**: 14 SKUs、総在庫: 300個

```bash
# 確認コマンド
docker exec ec-mysql mysql -uecuser -pecpassword ecsite \
  -e "SELECT COUNT(*) as total, SUM(stock) as total_stock FROM skus"
```

## 🔧 技術的特徴

### 1. JSON型のカスタムScanner実装

MySQLのJSON型をGo構造体にマッピングするため、`sql.Scanner`と`driver.Valuer`インターフェースを実装：

```go
type SKUAttrs map[string]string

func (a *SKUAttrs) Scan(value interface{}) error {
    bytes, _ := value.([]byte)
    result := make(SKUAttrs)
    json.Unmarshal(bytes, &result)
    *a = result
    return nil
}

func (a SKUAttrs) Value() (driver.Value, error) {
    return json.Marshal(a)
}
```

### 2. トランザクション制御による在庫整合性

```go
func (r *SKURepository) AdjustStock(ctx context.Context, skuID string, quantity int) error {
    return r.db.Transaction(func(tx *gorm.DB) error {
        // 1. 現在在庫取得（行ロック）
        var sku SKU
        tx.Where("id = ?", skuID).First(&sku)
        
        // 2. 在庫計算
        newStock := sku.Stock + quantity
        if newStock < 0 {
            return fmt.Errorf("insufficient stock")
        }
        
        // 3. 更新
        tx.Model(&SKU{}).Where("id = ?", skuID).Update("stock", newStock)
        return nil
    })
}
```

### 3. ルーティング順序の最適化

具体的なルート（`/products/:id/skus`）を動的ルート（`/products/:id`）より先に定義：

```go
// 具体的なルートを先に
api.GET("/products/:product_id/with-skus", skuHandler.GetProductWithSKUs)
api.GET("/products/:product_id/skus", skuHandler.GetSKUsByProduct)

// 動的ルートは後
api.GET("/products/:id", productHandler.GetProduct)
```

## 🎯 今後の拡張予定

### Phase 1: 在庫履歴
- [ ] `sku_stock_history`テーブルの活用
- [ ] 在庫変動の追跡・監査

### Phase 2: 高度な在庫管理
- [ ] 引当在庫（予約在庫）管理
- [ ] 安全在庫アラート
- [ ] 自動発注機能

### Phase 3: バリエーション拡張
- [ ] 画像のSKUごと管理
- [ ] SKU固有の説明文
- [ ] 複数価格帯（会員価格など）

### Phase 4: 統計・分析
- [ ] SKUごとの売上分析
- [ ] 人気バリエーション特定
- [ ] 在庫回転率計算

## 📝 運用Tips

### 在庫不足エラーへの対応

```bash
# エラー例
{
  "error": "failed to adjust stock: insufficient stock: current=5, requested=-10"
}

# 対処: 現在在庫を確認
curl "http://localhost:8001/api/v1/skus/sku-001" | jq '.stock'
```

### 低在庫SKUの定期監視

```bash
# 在庫10以下のSKUを毎日チェック
curl "http://localhost:8001/api/v1/skus/low-stock?threshold=10" \
  | jq '.[] | {code: .sku_code, stock: .stock, product: .product_id}'
```

### SKUコードの命名規則

推奨フォーマット: `{ブランド}-{モデル}-{色}-{サイズ}`

例:
- `NIKE-AM90-RED-M`
- `ADIDAS-UB22-WHT-L`
- `UNIQLO-HT-BLK-S`

## 🐛 トラブルシューティング

### JSON Scanエラー

**症状**: `unsupported Scan, storing driver.Value type []uint8 into type *domain.SKUAttrs`

**原因**: SKUAttrsのScanメソッド未実装

**解決**: domain/sku.goにScanner/Valuerインターフェース実装済み

### ルーティング404エラー

**症状**: `/products/{id}/skus`が404を返す

**原因**: ルーティング順序の問題

**解決**: 具体的なルートを動的ルートより先に定義

## 📚 関連ドキュメント

- [Product Service MySQL移行](PRODUCT_SERVICE_MYSQL_MIGRATION.md)
- [Elasticsearch同期機能](ELASTICSEARCH_SYNC.md)
- [API Examples](../API_EXAMPLES.md)
