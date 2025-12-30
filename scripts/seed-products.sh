#!/bin/bash

# Elasticsearch商品データシードスクリプト
# Usage: ./scripts/seed-products.sh

ES_URL="${ELASTICSEARCH_URL:-http://localhost:9200}"
INDEX_NAME="products"

echo "🚀 Elasticsearchに商品データを投入します..."
echo "URL: $ES_URL"
echo "Index: $INDEX_NAME"
echo ""

# インデックスが存在するか確認
if curl -s -o /dev/null -w "%{http_code}" "$ES_URL/$INDEX_NAME" | grep -q "200"; then
    echo "⚠️  既存のインデックスを削除します..."
    curl -X DELETE "$ES_URL/$INDEX_NAME"
    echo ""
fi

# インデックス作成（マッピング定義）
echo "📝 インデックスを作成します..."
curl -X PUT "$ES_URL/$INDEX_NAME" -H 'Content-Type: application/json' -d '{
  "settings": {
    "number_of_shards": 1,
    "number_of_replicas": 0,
    "analysis": {
      "analyzer": {
        "japanese_analyzer": {
          "type": "custom",
          "tokenizer": "standard",
          "filter": ["lowercase"]
        }
      }
    }
  },
  "mappings": {
    "properties": {
      "id": { "type": "keyword" },
      "name": { 
        "type": "text",
        "analyzer": "japanese_analyzer",
        "fields": {
          "keyword": { "type": "keyword" }
        }
      },
      "description": { 
        "type": "text",
        "analyzer": "japanese_analyzer"
      },
      "price": { "type": "double" },
      "stock": { "type": "integer" },
      "category": { "type": "keyword" },
      "images": { "type": "keyword" },
      "rating": { "type": "double" },
      "review_count": { "type": "integer" },
      "is_sale": { "type": "boolean" },
      "sale_price": { "type": "double" },
      "created_at": { "type": "date" },
      "updated_at": { "type": "date" },
      "is_active": { "type": "boolean" }
    }
  }
}'
echo ""
echo ""

# サンプル商品データ投入
echo "📦 サンプル商品を投入します..."

# 1. Nike Air Max
curl -X POST "$ES_URL/$INDEX_NAME/_doc/prod-001" -H 'Content-Type: application/json' -d '{
  "id": "prod-001",
  "name": "Nike Air Max 90",
  "description": "クラシックなデザインと優れたクッション性を備えたナイキの人気スニーカー。日常使いからスポーツシーンまで幅広く活躍します。",
  "price": 12000,
  "stock": 50,
  "category": "shoes",
  "images": ["/images/products/nike-air-max-90.jpg"],
  "rating": 4.5,
  "review_count": 128,
  "is_sale": false,
  "is_active": true,
  "created_at": "2024-12-01T00:00:00Z",
  "updated_at": "2024-12-01T00:00:00Z"
}'

# 2. Adidas Ultraboost
curl -X POST "$ES_URL/$INDEX_NAME/_doc/prod-002" -H 'Content-Type: application/json' -d '{
  "id": "prod-002",
  "name": "Adidas Ultraboost 22",
  "description": "革新的なBoostフォームが生み出す最高のクッション性とエネルギーリターン。ランニングから街履きまで対応。",
  "price": 18000,
  "stock": 30,
  "category": "shoes",
  "images": ["/images/products/adidas-ultraboost.jpg"],
  "rating": 4.8,
  "review_count": 95,
  "is_sale": true,
  "sale_price": 14400,
  "is_active": true,
  "created_at": "2024-12-01T00:00:00Z",
  "updated_at": "2024-12-15T00:00:00Z"
}'

# 3. Sony WH-1000XM5
curl -X POST "$ES_URL/$INDEX_NAME/_doc/prod-003" -H 'Content-Type: application/json' -d '{
  "id": "prod-003",
  "name": "Sony WH-1000XM5 ワイヤレスノイズキャンセリングヘッドホン",
  "description": "業界最高クラスのノイズキャンセリング性能。高音質と快適な装着感で長時間の使用も快適。",
  "price": 45000,
  "stock": 20,
  "category": "electronics",
  "images": ["/images/products/sony-wh1000xm5.jpg"],
  "rating": 4.9,
  "review_count": 203,
  "is_sale": false,
  "is_active": true,
  "created_at": "2024-11-15T00:00:00Z",
  "updated_at": "2024-11-15T00:00:00Z"
}'

# 4. Apple AirPods Pro
curl -X POST "$ES_URL/$INDEX_NAME/_doc/prod-004" -H 'Content-Type: application/json' -d '{
  "id": "prod-004",
  "name": "Apple AirPods Pro (第2世代)",
  "description": "アクティブノイズキャンセリングと空間オーディオ対応。Apple製品とのシームレスな連携が魅力。",
  "price": 38000,
  "stock": 45,
  "category": "electronics",
  "images": ["/images/products/airpods-pro.jpg"],
  "rating": 4.7,
  "review_count": 312,
  "is_sale": true,
  "sale_price": 33000,
  "is_active": true,
  "created_at": "2024-12-05T00:00:00Z",
  "updated_at": "2024-12-20T00:00:00Z"
}'

# 5. Patagonia フリースジャケット
curl -X POST "$ES_URL/$INDEX_NAME/_doc/prod-005" -H 'Content-Type: application/json' -d '{
  "id": "prod-005",
  "name": "Patagonia Better Sweater Jacket",
  "description": "環境に配慮したリサイクル素材を使用。保温性と通気性を兼ね備えた定番フリース。",
  "price": 15000,
  "stock": 60,
  "category": "clothing",
  "images": ["/images/products/patagonia-fleece.jpg"],
  "rating": 4.6,
  "review_count": 87,
  "is_sale": false,
  "is_active": true,
  "created_at": "2024-11-01T00:00:00Z",
  "updated_at": "2024-11-01T00:00:00Z"
}'

# 6. The North Face ダウンジャケット
curl -X POST "$ES_URL/$INDEX_NAME/_doc/prod-006" -H 'Content-Type: application/json' -d '{
  "id": "prod-006",
  "name": "The North Face Nuptse Jacket",
  "description": "高品質ダウンを使用した防寒性抜群のアイコニックなダウンジャケット。タウンユースにも最適。",
  "price": 35000,
  "stock": 25,
  "category": "clothing",
  "images": ["/images/products/tnf-nuptse.jpg"],
  "rating": 4.8,
  "review_count": 156,
  "is_sale": false,
  "is_active": true,
  "created_at": "2024-10-15T00:00:00Z",
  "updated_at": "2024-10-15T00:00:00Z"
}'

# 7. MacBook Air M2
curl -X POST "$ES_URL/$INDEX_NAME/_doc/prod-007" -H 'Content-Type: application/json' -d '{
  "id": "prod-007",
  "name": "Apple MacBook Air 13インチ M2チップ搭載",
  "description": "驚異的なパフォーマンスと最大18時間のバッテリー駆動時間。薄型軽量ボディで持ち運びも快適。",
  "price": 148000,
  "stock": 15,
  "category": "electronics",
  "images": ["/images/products/macbook-air-m2.jpg"],
  "rating": 4.9,
  "review_count": 421,
  "is_sale": false,
  "is_active": true,
  "created_at": "2024-09-01T00:00:00Z",
  "updated_at": "2024-09-01T00:00:00Z"
}'

# 8. iPad Pro
curl -X POST "$ES_URL/$INDEX_NAME/_doc/prod-008" -H 'Content-Type: application/json' -d '{
  "id": "prod-008",
  "name": "Apple iPad Pro 11インチ M2チップ搭載",
  "description": "プロレベルのパフォーマンスと美しいLiquid Retinaディスプレイ。Apple Pencil対応。",
  "price": 98000,
  "stock": 22,
  "category": "electronics",
  "images": ["/images/products/ipad-pro.jpg"],
  "rating": 4.7,
  "review_count": 189,
  "is_sale": true,
  "sale_price": 88000,
  "is_active": true,
  "created_at": "2024-11-20T00:00:00Z",
  "updated_at": "2024-12-10T00:00:00Z"
}'

# 9. New Balance 2002R
curl -X POST "$ES_URL/$INDEX_NAME/_doc/prod-009" -H 'Content-Type: application/json' -d '{
  "id": "prod-009",
  "name": "New Balance 2002R",
  "description": "レトロなデザインと最新のクッション技術を融合。オールマイティに使えるライフスタイルシューズ。",
  "price": 16000,
  "stock": 40,
  "category": "shoes",
  "images": ["/images/products/nb-2002r.jpg"],
  "rating": 4.6,
  "review_count": 73,
  "is_sale": false,
  "is_active": true,
  "created_at": "2024-12-10T00:00:00Z",
  "updated_at": "2024-12-10T00:00:00Z"
}'

# 10. Uniqlo ヒートテックウルトラウォーム
curl -X POST "$ES_URL/$INDEX_NAME/_doc/prod-010" -H 'Content-Type: application/json' -d '{
  "id": "prod-010",
  "name": "ユニクロ ヒートテックウルトラウォームクルーネックT",
  "description": "極暖インナー。通常のヒートテックの2.25倍暖かい。極寒の日も快適に過ごせます。",
  "price": 1990,
  "stock": 150,
  "category": "clothing",
  "images": ["/images/products/heattech-ultra.jpg"],
  "rating": 4.4,
  "review_count": 542,
  "is_sale": true,
  "sale_price": 1490,
  "is_active": true,
  "created_at": "2024-11-01T00:00:00Z",
  "updated_at": "2024-12-01T00:00:00Z"
}'

# 11. Bose QuietComfort Earbuds
curl -X POST "$ES_URL/$INDEX_NAME/_doc/prod-011" -H 'Content-Type: application/json' -d '{
  "id": "prod-011",
  "name": "Bose QuietComfort Earbuds II",
  "description": "パーソナライズされたノイズキャンセリングとプレミアムサウンド。完全ワイヤレスの最上位モデル。",
  "price": 35000,
  "stock": 18,
  "category": "electronics",
  "images": ["/images/products/bose-qc-earbuds.jpg"],
  "rating": 4.5,
  "review_count": 134,
  "is_sale": false,
  "is_active": true,
  "created_at": "2024-10-20T00:00:00Z",
  "updated_at": "2024-10-20T00:00:00Z"
}'

# 12. Converse Chuck Taylor
curl -X POST "$ES_URL/$INDEX_NAME/_doc/prod-012" -H 'Content-Type: application/json' -d '{
  "id": "prod-012",
  "name": "Converse Chuck Taylor All Star",
  "description": "100年以上愛され続けるタイムレスなキャンバススニーカー。どんなスタイルにもマッチ。",
  "price": 6500,
  "stock": 80,
  "category": "shoes",
  "images": ["/images/products/converse-chuck.jpg"],
  "rating": 4.3,
  "review_count": 267,
  "is_sale": false,
  "is_active": true,
  "created_at": "2024-09-15T00:00:00Z",
  "updated_at": "2024-09-15T00:00:00Z"
}'

echo ""
echo ""
echo "✅ サンプル商品データの投入が完了しました！"
echo ""
echo "確認コマンド:"
echo "  curl 'http://localhost:9200/products/_search?size=20&pretty'"
echo "  curl 'http://localhost:8001/api/v1/products/search?page=1&page_size=20'"
