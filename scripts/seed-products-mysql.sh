#!/bin/bash

# MySQL商品データシードスクリプト
# Usage: ./scripts/seed-products-mysql.sh

MYSQL_HOST="${MYSQL_HOST:-127.0.0.1}"
MYSQL_USER="${MYSQL_USER:-ecuser}"
MYSQL_PASSWORD="${MYSQL_PASSWORD:-ecpassword}"
MYSQL_DATABASE="${MYSQL_DATABASE:-ecsite}"

echo "🚀 MySQLに商品データを投入します..."
echo "Host: $MYSQL_HOST"
echo "Database: $MYSQL_DATABASE"
echo ""

# MySQLに接続してデータ投入
mysql -h"$MYSQL_HOST" -u"$MYSQL_USER" -p"$MYSQL_PASSWORD" "$MYSQL_DATABASE" <<EOF

-- 既存データをクリア
TRUNCATE TABLE products;

-- サンプル商品データ投入
INSERT INTO products (id, name, description, price, stock, category, rating, review_count, is_sale, sale_price, is_active) VALUES
('prod-001', 'Nike Air Max 90', 'クラシックなデザインと優れたクッション性を備えたナイキの人気スニーカー。日常使いからスポーツシーンまで幅広く活躍します。', 12000.00, 50, 'shoes', 4.50, 128, false, NULL, true),
('prod-002', 'Adidas Ultraboost 22', '革新的なBoostフォームが生み出す最高のクッション性とエネルギーリターン。ランニングから街履きまで対応。', 18000.00, 30, 'shoes', 4.80, 95, true, 14400.00, true),
('prod-003', 'Sony WH-1000XM5 ワイヤレスノイズキャンセリングヘッドホン', '業界最高クラスのノイズキャンセリング性能。高音質と快適な装着感で長時間の使用も快適。', 45000.00, 20, 'electronics', 4.90, 203, false, NULL, true),
('prod-004', 'Apple AirPods Pro (第2世代)', 'アクティブノイズキャンセリングと空間オーディオ対応。Apple製品とのシームレスな連携が魅力。', 38000.00, 45, 'electronics', 4.70, 312, true, 33000.00, true),
('prod-005', 'Patagonia Better Sweater Jacket', '環境に配慮したリサイクル素材を使用。保温性と通気性を兼ね備えた定番フリース。', 15000.00, 60, 'clothing', 4.60, 87, false, NULL, true),
('prod-006', 'The North Face Nuptse Jacket', '高品質ダウンを使用した防寒性抜群のアイコニックなダウンジャケット。タウンユースにも最適。', 35000.00, 25, 'clothing', 4.80, 156, false, NULL, true),
('prod-007', 'Apple MacBook Air 13インチ M2チップ搭載', '驚異的なパフォーマンスと最大18時間のバッテリー駆動時間。薄型軽量ボディで持ち運びも快適。', 148000.00, 15, 'electronics', 4.90, 421, false, NULL, true),
('prod-008', 'Apple iPad Pro 11インチ M2チップ搭載', 'プロレベルのパフォーマンスと美しいLiquid Retinaディスプレイ。Apple Pencil対応。', 98000.00, 22, 'electronics', 4.70, 189, true, 88000.00, true),
('prod-009', 'New Balance 2002R', 'レトロなデザインと最新のクッション技術を融合。オールマイティに使えるライフスタイルシューズ。', 16000.00, 40, 'shoes', 4.60, 73, false, NULL, true),
('prod-010', 'ユニクロ ヒートテックウルトラウォームクルーネックT', '極暖インナー。通常のヒートテックの2.25倍暖かい。極寒の日も快適に過ごせます。', 1990.00, 150, 'clothing', 4.40, 542, true, 1490.00, true),
('prod-011', 'Bose QuietComfort Earbuds II', 'パーソナライズされたノイズキャンセリングとプレミアムサウンド。完全ワイヤレスの最上位モデル。', 35000.00, 18, 'electronics', 4.50, 134, false, NULL, true),
('prod-012', 'Converse Chuck Taylor All Star', '100年以上愛され続けるタイムレスなキャンバススニーカー。どんなスタイルにもマッチ。', 6500.00, 80, 'shoes', 4.30, 267, false, NULL, true);

EOF

if [ $? -eq 0 ]; then
    echo ""
    echo "✅ MySQLへのサンプル商品データ投入が完了しました！"
    echo ""
    echo "確認コマンド:"
    echo "  mysql -h$MYSQL_HOST -u$MYSQL_USER -p$MYSQL_PASSWORD $MYSQL_DATABASE -e 'SELECT id, name, price, stock FROM products;'"
    echo "  curl 'http://localhost:8001/api/v1/products/search?page=1&page_size=20'"
else
    echo ""
    echo "❌ データ投入に失敗しました"
    exit 1
fi
