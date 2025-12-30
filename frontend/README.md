# EC Sample - Next.js フロントエンド

React + Next.js 14 + TailwindCSS で構築されたモダンなECサイトフロントエンド。

## 特徴

- **モダンな技術スタック**: Next.js 14 App Router、React 18、TypeScript
- **スタイリング**: TailwindCSS による効率的なUI構築
- **状態管理**: Zustand による軽量な状態管理
- **API統合**: マイクロサービスバックエンドとのシームレスな連携
- **責任分離**: Store、Components、Pages の明確な構造

## ディレクトリ構成

```
apps/web/
├── src/
│   ├── app/                 # Next.js App Router ページ
│   │   ├── layout.tsx       # ルートレイアウト
│   │   ├── page.tsx         # ホームページ
│   │   ├── products/        # 商品ページ
│   │   ├── cart/            # カートページ
│   │   ├── checkout/        # 決済ページ
│   │   ├── my-orders/       # 注文履歴
│   │   ├── orders/[id]/     # 注文詳細
│   │   └── my-points/       # ポイント管理
│   ├── components/          # React コンポーネント
│   │   ├── Header.tsx       # ヘッダー
│   │   ├── Footer.tsx       # フッター
│   │   ├── ProductCard.tsx  # 商品カード
│   │   └── Button.tsx       # ボタンコンポーネント
│   ├── lib/
│   │   └── api-client.ts    # API クライアント
│   ├── store/               # Zustand ストア
│   │   ├── cart.ts          # カート状態
│   │   ├── product.ts       # 商品状態
│   │   └── points.ts        # ポイント状態
│   ├── types/
│   │   └── api.ts           # API型定義
│   └── styles/
│       └── globals.css      # グローバルスタイル
├── package.json
├── next.config.js
├── tailwind.config.js
├── tsconfig.json
└── Dockerfile
```

## セットアップ

### ローカル開発

```bash
cd apps/web

# 依存関係をインストール
npm install

# 開発サーバーを起動
npm run dev
```

ブラウザで http://localhost:3000 を開く

### 環境変数

`.env.local` を作成（または `next.config.js` で設定）

```
NEXT_PUBLIC_API_BASE_URL=http://localhost:8001
NEXT_PUBLIC_CART_API_URL=http://localhost:8002
NEXT_PUBLIC_ORDER_API_URL=http://localhost:8003
NEXT_PUBLIC_POINT_API_URL=http://localhost:8004
NEXT_PUBLIC_PROMOTION_API_URL=http://localhost:8005
```

## 利用可能なコマンド

```bash
# 開発サーバー起動
npm run dev

# ビルド
npm run build

# 本番サーバー起動
npm start

# 型チェック
npm run type-check

# リンター実行
npm run lint
```

## 実装済みページ

### ホームページ (`/`)
- 新着商品表示
- カテゴリ分類
- 推奨商品紹介

### 商品一覧 (`/products`)
- 商品検索（キーワード）
- カテゴリフィルター
- ソート機能（新着・価格・評価）
- 商品カード表示

### 商品詳細 (`/products/[id]`)
- 詳細な商品情報
- 価格表示（セール価格対応）
- 在庫状況
- 数量選択
- カートに追加機能

### ショッピングカート (`/cart`)
- カートアイテム表示
- 数量変更
- アイテム削除
- 価格計算
- 注文確認画面へのリンク

### 決済 (`/checkout`)
- 配送先住所入力
- 支払い方法選択
- クーポンコード入力
- 注文確認
- 注文確定処理

### 注文履歴 (`/my-orders`)
- 過去注文一覧
- ステータス表示
- 注文詳細へのリンク

### 注文詳細 (`/orders/[id]`)
- 注文情報表示
- 注文商品一覧
- 配送情報
- 獲得ポイント表示
- 支払い方法表示

### ポイント管理 (`/my-points`)
- ポイント残高表示
- ポイント履歴
- トランザクション詳細
- ポイント利用ガイド

## API統合

### API クライアント

`src/lib/api-client.ts` に各マイクロサービスとのAPI統合を実装

- **ProductApi**: 商品検索・取得
- **CartApi**: カート操作
- **OrderApi**: 注文作成・取得
- **PointApi**: ポイント残高・トランザクション
- **PromotionApi**: クーポン・セール情報

### 型定義

`src/types/api.ts` で全API型定義を管理

## 状態管理

### Zustand ストア

軽量な状態管理ライブラリ Zustand を採用

```typescript
const { cart, addItem, removeItem } = useCartStore();
const { products, searchProducts } = useProductStore();
const { balance, fetchBalance } = usePointStore();
```

## UIコンポーネント

### Header
- ロゴ、ナビゲーション、検索、カートアイコン

### Footer
- リンク、会社情報、ポリシー

### ProductCard
- 商品画像、名前、評価、価格、アクション

### Button
複数のバリアントとサイズを備えた再利用可能コンポーネント

## スタイリング

TailwindCSS による効率的なレスポンシブデザイン実装

カラーパレット:
- Primary: Blue (`#3B82F6`)
- Secondary: Green (`#10B981`)
- Accent: Amber (`#F59E0B`)
- Danger: Red (`#EF4444`)

## デプロイ

### Docker

```bash
docker build -t ec-sample-web ./apps/web
docker run -p 3000:3000 ec-sample-web
```

### docker-compose

```bash
docker-compose up web
```


---

## 📚 参考リンク

- [プロジェクトルート README](../README.md)
- [API Schema & Type Definitions](../api-schema/README.md)
- [Next.js Documentation](https://nextjs.org/docs)
- [TailwindCSS Documentation](https://tailwindcss.com/docs)
- [Zustand Documentation](https://zustand-demo.pmnd.rs/)

**Happy Coding! 🚀**
