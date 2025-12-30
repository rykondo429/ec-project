#!/bin/bash

# 新商品作成テスト（MySQL + ES同期）

curl -X POST http://localhost:8001/api/v1/products \
  -H "Content-Type: application/json" \
  -d '{
    "id": "prod-test-001",
    "name": "テスト商品 - Nintendo Switch",
    "description": "人気のゲーム機。マリオ、ゼルダなど名作タイトルが遊べます。",
    "price": 32980,
    "stock": 50,
    "category": "electronics",
    "rating": 4.9,
    "review_count": 250,
    "is_sale": false,
    "is_active": true
  }'

echo ""
echo "---"
echo "MySQLから確認:"
sleep 1
curl -s "http://localhost:8001/api/v1/products/prod-test-001" | jq .

echo ""
echo "---"
echo "Elasticsearchから確認:"
sleep 1
curl -s "http://localhost:9200/products/_doc/prod-test-001" | jq ._source
