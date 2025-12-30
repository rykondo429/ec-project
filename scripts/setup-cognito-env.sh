#!/bin/bash

# Cognito環境変数設定スクリプト
# LocalStack無料版ではCognitoがサポートされていないため、
# 本番AWS Cognitoまたは開発モード（認証無効）で使用します

echo "==================================================================="
echo "EC-Sample Cognito環境変数設定"
echo "==================================================================="
echo ""
echo "⚠️  注意: LocalStack無料版ではCognitoがサポートされていません"
echo ""
echo "設定モードを選択してください:"
echo "  1) 開発モード（認証機能無効）- 推奨"
echo "  2) AWS Cognito本番環境（手動設定）"
echo ""
read -p "選択 [1/2]: " MODE

if [ "$MODE" = "2" ]; then
    echo ""
    echo "AWS Cognito本番環境の設定値を入力してください:"
    echo ""
    read -p "User Pool ID: " USER_POOL_ID
    read -p "Client ID: " CLIENT_ID
    read -p "Region [us-east-1]: " REGION
    REGION=${REGION:-us-east-1}
    AUTH_ENABLED=true
else
    echo ""
    echo "開発モード（認証無効）で設定します..."
    USER_POOL_ID=""
    CLIENT_ID=""
    REGION="us-east-1"
    AUTH_ENABLED=false
fi

# frontend/.env.local ファイルを生成
cat > frontend/.env.local <<EOF
# Frontend API Base URLs
NEXT_PUBLIC_API_BASE_URL=http://localhost:8001
NEXT_PUBLIC_CART_API_URL=http://localhost:8002
NEXT_PUBLIC_ORDER_API_URL=http://localhost:8003
NEXT_PUBLIC_POINT_API_URL=http://localhost:8004
NEXT_PUBLIC_PROMOTION_API_URL=http://localhost:8005

# Authentication
NEXT_PUBLIC_AUTH_ENABLED=${AUTH_ENABLED}
NEXT_PUBLIC_COGNITO_USER_POOL_ID=${USER_POOL_ID}
NEXT_PUBLIC_COGNITO_CLIENT_ID=${CLIENT_ID}
NEXT_PUBLIC_COGNITO_REGION=${REGION}
EOF

echo ""
echo "✅ frontend/.env.local ファイルを生成しました"
echo ""

if [ "$MODE" = "2" ]; then
    echo "📝 AWS Cognito設定:"
    echo "   User Pool ID: ${USER_POOL_ID}"
    echo "   Client ID: ${CLIENT_ID}"
    echo "   Region: ${REGION}"
    echo ""
    echo "⚠️  バックエンドサービスの環境変数も設定してください:"
    echo "   docker-compose.ymlまたは.envファイルに以下を追加:"
    echo "   COGNITO_REGION=${REGION}"
    echo "   COGNITO_USER_POOL_ID=${USER_POOL_ID}"
    echo "   COGNITO_CLIENT_ID=${CLIENT_ID}"
else
    echo "📝 開発モード設定完了"
    echo "   認証機能: 無効"
    echo "   すべてのAPIリクエストは認証なしで動作します"
fi

echo ""
echo "🚀 次のステップ:"
echo "   1. フロントエンドを起動: cd frontend && npm run dev"
echo "   2. ブラウザで http://localhost:3000 にアクセス"
if [ "$MODE" = "2" ]; then
    echo "   3. 右上の「ログイン」ボタンからログイン"
else
    echo "   3. 認証なしで全機能を利用できます"
fi
echo ""
