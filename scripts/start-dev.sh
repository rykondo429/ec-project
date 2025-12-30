#!/bin/bash
# フロントエンド開発環境 クイックスタートスクリプト

set -e

echo "🚀 EC Sample フロントエンド開発環境を起動します..."
echo ""

# Dockerがインストールされているか確認
if ! command -v docker &> /dev/null; then
    echo "❌ Dockerがインストールされていません"
    echo "   https://www.docker.com/products/docker-desktop からインストールしてください"
    exit 1
fi

# Docker Composeのバージョンを確認
if docker compose version &> /dev/null; then
    COMPOSE_CMD="docker compose"
elif docker-compose --version &> /dev/null; then
    COMPOSE_CMD="docker-compose"
else
    echo "❌ Docker Composeが利用できません"
    exit 1
fi

echo "✅ Docker環境を確認しました"
echo ""

# オプション表示
echo "起動モードを選択してください:"
echo "  1) フロントエンドのみ（推奨 - 最も軽量）"
echo "  2) インフラ + フロントエンド（MySQL, Redis等も起動）"
echo "  3) 全サービス起動（バックエンドAPI含む）"
echo ""
read -p "番号を入力 [1-3]: " choice

case $choice in
  1)
    echo ""
    echo "📦 フロントエンドのみを起動します..."
    $COMPOSE_CMD -f docker-compose.dev.yml up web
    ;;
  2)
    echo ""
    echo "📦 インフラ + フロントエンドを起動します..."
    $COMPOSE_CMD -f docker-compose.dev.yml up
    ;;
  3)
    echo ""
    echo "📦 全サービスを起動します（時間がかかる場合があります）..."
    $COMPOSE_CMD up
    ;;
  *)
    echo "❌ 無効な選択です"
    exit 1
    ;;
esac
