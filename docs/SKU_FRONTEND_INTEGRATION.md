# SKUフロントエンド統合 - 実装ドキュメント

## 📋 概要

商品バリエーション（SKU）をフロントエンドに統合し、ユーザーがサイズや色を選択して購入できる機能を実装しました。

## 🎯 実装内容

### 1. API定義の更新

#### Product Service OpenAPI
- **新規エンドポイント追加**:
  - `GET /products/{product_id}/skus` - 商品のSKU一覧取得
  - `GET /products/{product_id}/with-skus` - 商品とSKU一括取得
  - `GET /skus/{id}` - SKU詳細取得
  - `GET /skus/code/{code}` - SKUコード検索
  - `GET /skus/low-stock` - 低在庫SKU一覧

- **新規スキーマ定義**:
  ```yaml
  SKU:
    type: object
    properties:
      id: string
      product_id: string
      sku_code: string
      name: string
      attributes:
        type: object
        additionalProperties: string  # color, size など
      price: number
      stock: integer
      is_active: boolean
  
  ProductWithSKUs:
    allOf:
      - $ref: '#/components/schemas/Product'
      - properties:
          skus:
            type: array
            items:
              $ref: '#/components/schemas/SKU'
  ```

#### Cart Service OpenAPI
- **CartItemとAddItemRequestを拡張**:
  ```yaml
  CartItem:
    properties:
      sku_id: string (optional)
      sku_code: string (optional)
      sku_attributes:
        type: object
        additionalProperties: string
  ```

### 2. TypeScript型生成

**実行コマンド**:
```bash
cd api-schema
npm run generate
```

**生成された型**:
- `types/product-service.ts` - SKU, ProductWithSKUs型
- `types/cart-service.ts` - SKU情報付きCartItem型

### 3. フロントエンド実装

#### APIクライアント拡張 ([lib/api-client.ts](../frontend/src/lib/api-client.ts))

```typescript
class ProductApi {
  async getProductWithSKUs(id: string): Promise<ProductWithSKUs> {
    const response = await this.client.get(`/api/v1/products/${id}/with-skus`);
    return response.data;
  }

  async getSKUsByProduct(productId: string): Promise<SKU[]> {
    const response = await this.client.get(`/api/v1/products/${productId}/skus`);
    return response.data;
  }
}
```

#### Zustand Store更新

**Product Store** ([store/product.ts](../frontend/src/store/product.ts)):
```typescript
interface ProductState {
  productWithSKUs: ProductWithSKUs | null;
  getProductWithSKUs: (id: string) => Promise<void>;
}
```

**Cart Store** ([store/cart.ts](../frontend/src/store/cart.ts)):
```typescript
interface CartState {
  addItem: (
    productId: string,
    productName: string,
    price: number,
    quantity: number,
    skuId?: string,              // 追加
    skuCode?: string,            // 追加
    skuAttributes?: Record<string, string> // 追加
  ) => Promise<void>;
}
```

#### 商品詳細ページUI ([app/products/[id]/page.tsx](../frontend/src/app/products/[id]/page.tsx))

**主な機能**:
1. **SKU選択UI**
   - 属性（color, size）ごとにボタン表示
   - 在庫なしSKUは無効化
   - 選択中のSKUを視覚的に表示

2. **動的価格・在庫表示**
   - 選択されたSKUの価格と在庫を表示
   - SKUごとに異なる価格設定可能

3. **カート追加**
   - SKU情報をカートAPIに送信
   - 同じ商品でもSKU IDが異なれば別アイテムとして追加

**UIコンポーネント構成**:
```tsx
// 属性選択ボタン
{attributeKeys.map(attrKey => (
  <div key={attrKey}>
    <label>{attrKey === 'size' ? 'サイズ' : 'カラー'}</label>
    <div className="flex gap-2">
      {options.map(option => (
        <button
          onClick={() => handleAttributeSelect(attrKey, option)}
          disabled={!isAvailable}
          className={isSelected ? 'border-primary bg-primary' : ''}
        >
          {option}
        </button>
      ))}
    </div>
  </div>
))}

// 選択中のSKU表示
{selectedSKU && (
  <div className="bg-gray-50 p-3">
    <p>選択中: {selectedSKU.name}</p>
    <p>SKUコード: {selectedSKU.sku_code}</p>
  </div>
)}
```

### 4. バックエンド実装

#### Cart Service Domain更新 ([cart-service/domain/cart.go](../backend/services/cart-service/domain/cart.go))

```go
type CartItem struct {
    ProductID     string            `json:"product_id"`
    ProductName   string            `json:"product_name"`
    SKUID         string            `json:"sku_id,omitempty"`
    SKUCode       string            `json:"sku_code,omitempty"`
    SKUAttributes map[string]string `json:"sku_attributes,omitempty"`
    Price         float64           `json:"price"`
    Quantity      int               `json:"quantity"`
    Subtotal      float64           `json:"subtotal"`
    AddedAt       time.Time         `json:"added_at"`
}

type Cart struct {
    UserID       string      `json:"user_id" dynamodbav:"user_id"`
    Items        []*CartItem `json:"items" dynamodbav:"items"`
    Total        float64     `json:"total" dynamodbav:"total"`
    ItemCount    int         `json:"item_count" dynamodbav:"item_count"`
    LastModified time.Time   `json:"last_modified" dynamodbav:"last_modified"`
    ExpiresAt    time.Time   `json:"expires_at" dynamodbav:"expires_at"`
}
```

**重要**: DynamoDB用の`dynamodbav`タグを追加してマーシャリングを正しく動作させる

#### Cart UseCase更新 ([cart-service/usecase/cart.go](../backend/services/cart-service/usecase/cart.go))

**SKU IDベースのアイテム識別**:
```go
// 既存アイテム確認（同じ商品IDとSKU IDの組み合わせ）
itemKey := req.ProductID
if req.SKUID != "" {
    itemKey = req.ProductID + "_" + req.SKUID
}

itemFound := false
for i, item := range cart.Items {
    existingKey := item.ProductID
    if item.SKUID != "" {
        existingKey = item.ProductID + "_" + item.SKUID
    }

    if existingKey == itemKey {
        // 数量更新
        cart.Items[i].Quantity += req.Quantity
        cart.Items[i].Subtotal = math.Round(float64(cart.Items[i].Quantity)*req.Price*100) / 100
        itemFound = true
        break
    }
}

if !itemFound {
    // 新規アイテム追加
    newItem := &domain.CartItem{
        ProductID:     req.ProductID,
        ProductName:   req.ProductName,
        SKUID:         req.SKUID,
        SKUCode:       req.SKUCode,
        SKUAttributes: req.SKUAttributes,
        Price:         req.Price,
        Quantity:      req.Quantity,
        Subtotal:      math.Round(float64(req.Quantity)*req.Price*100) / 100,
        AddedAt:       time.Now(),
    }
    cart.Items = append(cart.Items, newItem)
}
```

## 🧪 動作テスト結果

### テスト1: 商品とSKU一括取得

**リクエスト**:
```bash
curl "http://localhost:8001/api/v1/products/prod-001/with-skus"
```

**レスポンス**:
```json
{
  "id": "prod-001",
  "name": "Nike Air Max 90",
  "price": 12000,
  "skus": [
    {
      "id": "sku-001",
      "sku_code": "NIKE-AM90-RED-S",
      "name": "Nike Air Max 90 - Red / S",
      "attributes": {"color": "Red", "size": "S"},
      "price": 12000,
      "stock": 20
    },
    {
      "id": "sku-002",
      "sku_code": "NIKE-AM90-RED-M",
      "attributes": {"color": "Red", "size": "M"},
      "stock": 20
    },
    ... (全6 SKUs)
  ]
}
```

### テスト2: SKU付きカート追加

**リクエスト**:
```bash
curl -X POST "http://localhost:8002/api/v1/carts/items" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer test-token" \
  -d '{
    "product_id": "prod-001",
    "product_name": "Nike Air Max 90 - Red / M",
    "sku_id": "sku-002",
    "sku_code": "NIKE-AM90-RED-M",
    "sku_attributes": {"color": "Red", "size": "M"},
    "price": 12000,
    "quantity": 2
  }'
```

**レスポンス**:
```json
{
  "user_id": "test-user-123",
  "items": [
    {
      "product_id": "prod-001",
      "product_name": "Nike Air Max 90 - Red / M",
      "sku_id": "sku-002",
      "sku_code": "NIKE-AM90-RED-M",
      "sku_attributes": {"color": "Red", "size": "M"},
      "price": 12000,
      "quantity": 2,
      "subtotal": 24000
    }
  ],
  "total": 24000,
  "item_count": 2
}
```

### テスト3: 異なるSKUの追加

**リクエスト** (Black / Lサイズを追加):
```bash
curl -X POST "http://localhost:8002/api/v1/carts/items" \
  -d '{
    "product_id": "prod-001",
    "sku_id": "sku-006",
    "sku_code": "NIKE-AM90-BLK-L",
    "sku_attributes": {"color": "Black", "size": "L"},
    "price": 12000,
    "quantity": 1
  }'
```

**結果**: 
- 同じ商品ID (`prod-001`) でも別アイテムとして追加 ✅
- カート内に2つのアイテム（Red/M と Black/L）が存在 ✅

### テスト4: 同じSKUの再追加

**リクエスト** (Red / M を再度追加):
```bash
curl -X POST "http://localhost:8002/api/v1/carts/items" \
  -d '{
    "product_id": "prod-001",
    "sku_id": "sku-002",
    "sku_code": "NIKE-AM90-RED-M",
    "price": 12000,
    "quantity": 1
  }'
```

**結果**:
- Red / M の数量が 2 → 3 に増加 ✅
- Black / L は変更なし ✅

**最終カート内容**:
```json
{
  "items": [
    {
      "sku_code": "NIKE-AM90-RED-M",
      "quantity": 3,
      "subtotal": 36000
    },
    {
      "sku_code": "NIKE-AM90-BLK-L",
      "quantity": 1,
      "subtotal": 12000
    }
  ],
  "total": 48000,
  "item_count": 4
}
```

## 🎨 UIデザインパターン

### SKU選択ボタン

```tsx
// 選択中: プライマリカラーの背景
className="border-2 border-primary bg-primary text-white"

// 未選択: グレーボーダー
className="border-2 border-gray-300 hover:border-primary"

// 在庫なし: 無効化＋取り消し線
className="opacity-40 cursor-not-allowed line-through"
```

### 属性表示名の日本語化

```typescript
const displayName = attrKey === 'size' ? 'サイズ' : 
                   attrKey === 'color' ? 'カラー' : attrKey;
```

### チェックマーク表示

```tsx
{isSelected && (
  <Check size={16} className="absolute top-1 right-1" />
)}
```

## 🔧 起動方法

### 開発環境セットアップ

```bash
# 1. LocalStackを起動（DynamoDB）
cd /path/to/ec-sample
docker compose up -d localstack

# 2. Product Serviceを起動
cd backend/services/product-service
PORT=8001 go run main.go

# 3. Cart Serviceを起動
cd backend/services/cart-service
AWS_ENDPOINT_URL=http://localhost:4566 \
AWS_REGION=us-east-1 \
AWS_ACCESS_KEY_ID=test \
AWS_SECRET_ACCESS_KEY=test \
PORT=8002 \
go run main.go

# 4. Frontendを起動（Node.js 18以上が必要）
cd frontend
npm run dev
```

### 動作確認URL

- フロントエンド: http://localhost:3000
- Product Service: http://localhost:8001/api/v1/products/prod-001/with-skus
- Cart Service: http://localhost:8002/api/v1/carts

## 📝 実装のポイント

### 1. SKU識別キーの設計

**課題**: 同じ商品でも異なるSKUは別アイテムとして扱う必要がある

**解決策**: `product_id + "_" + sku_id` を複合キーとして使用

```go
itemKey := req.ProductID
if req.SKUID != "" {
    itemKey = req.ProductID + "_" + req.SKUID
}
```

### 2. DynamoDB属性マッピング

**課題**: Go構造体のフィールド名とDynamoDB属性名の不一致

**解決策**: `dynamodbav` タグを明示的に追加

```go
type Cart struct {
    UserID string `json:"user_id" dynamodbav:"user_id"`
}
```

### 3. フロントエンドの状態管理

**課題**: 複数の属性（色、サイズ）の選択状態を管理

**解決策**: 属性キーと値のマップを使用

```typescript
const [selectedAttributes, setSelectedAttributes] = 
  React.useState<Record<string, string>>({});

// 選択時
setSelectedAttributes(prev => ({
  ...prev,
  [attrKey]: value
}));
```

### 4. 在庫切れSKUの表示

**課題**: 選択不可のSKUを視覚的に区別

**解決策**: 在庫チェック＋スタイル変更

```typescript
const isAvailable = matchingSKU && matchingSKU.stock > 0;

<button
  disabled={!isAvailable}
  className={!isAvailable && 'opacity-40 line-through'}
>
```

## 🚀 今後の拡張予定

### Phase 1: UI/UX改善
- [ ] SKUごとの商品画像表示
- [ ] 在庫数リアルタイム更新
- [ ] カートページでのSKU表示改善

### Phase 2: 高度な機能
- [ ] SKUごとのレビュー・評価
- [ ] おすすめSKU表示（人気の組み合わせ）
- [ ] 在庫切れ通知機能

### Phase 3: 分析機能
- [ ] SKU別売上ランキング
- [ ] 属性別購入傾向分析
- [ ] 在庫最適化アラート

## 📚 関連ドキュメント

- [SKU管理機能](SKU_MANAGEMENT.md) - バックエンド実装詳細
- [API Examples](../API_EXAMPLES.md) - 全エンドポイントの使用例
- [OpenAPI Specifications](../api-schema/openapi/) - API仕様定義
