#!/bin/bash

# LocalStack 初期化スクリプト
# AWS リソースをセットアップ

echo "Initializing AWS resources..."

# Cognito User Pool 作成
echo "Creating Cognito User Pool..."
USER_POOL_ID=$(aws cognito-idp create-user-pool \
    --endpoint-url http://localhost:4566 \
    --pool-name ec-sample-user-pool \
    --policies "PasswordPolicy={MinimumLength=8,RequireUppercase=true,RequireLowercase=true,RequireNumbers=true,RequireSymbols=false}" \
    --auto-verified-attributes email \
    --username-attributes email \
    --region us-east-1 \
    --query 'UserPool.Id' \
    --output text 2>/dev/null || echo "")

if [ -z "$USER_POOL_ID" ]; then
    echo "User pool already exists or failed to create, getting existing pool ID..."
    USER_POOL_ID=$(aws cognito-idp list-user-pools \
        --endpoint-url http://localhost:4566 \
        --max-results 10 \
        --region us-east-1 \
        --query 'UserPools[?Name==`ec-sample-user-pool`].Id' \
        --output text)
fi

echo "User Pool ID: $USER_POOL_ID"

# Cognito User Pool Client 作成
echo "Creating Cognito User Pool Client..."
CLIENT_ID=$(aws cognito-idp create-user-pool-client \
    --endpoint-url http://localhost:4566 \
    --user-pool-id "$USER_POOL_ID" \
    --client-name ec-sample-web-client \
    --generate-secret \
    --explicit-auth-flows ALLOW_USER_PASSWORD_AUTH ALLOW_REFRESH_TOKEN_AUTH ALLOW_USER_SRP_AUTH \
    --region us-east-1 \
    --query 'UserPoolClient.ClientId' \
    --output text 2>/dev/null || echo "")

if [ -z "$CLIENT_ID" ]; then
    echo "Client already exists or failed to create, getting existing client ID..."
    CLIENT_ID=$(aws cognito-idp list-user-pool-clients \
        --endpoint-url http://localhost:4566 \
        --user-pool-id "$USER_POOL_ID" \
        --region us-east-1 \
        --query 'UserPoolClients[0].ClientId' \
        --output text)
fi

echo "Client ID: $CLIENT_ID"

# テストユーザー作成
echo "Creating test user..."
aws cognito-idp admin-create-user \
    --endpoint-url http://localhost:4566 \
    --user-pool-id "$USER_POOL_ID" \
    --username testuser@example.com \
    --user-attributes Name=email,Value=testuser@example.com Name=email_verified,Value=true \
    --temporary-password "TempPass123!" \
    --message-action SUPPRESS \
    --region us-east-1 2>/dev/null || echo "Test user already exists"

# テストユーザーのパスワードを永続化
aws cognito-idp admin-set-user-password \
    --endpoint-url http://localhost:4566 \
    --user-pool-id "$USER_POOL_ID" \
    --username testuser@example.com \
    --password "TestPass123!" \
    --permanent \
    --region us-east-1 2>/dev/null || echo "Password already set"

echo "Test User: testuser@example.com / TestPass123!"

# DynamoDB テーブル作成
aws dynamodb create-table \
    --endpoint-url http://localhost:4566 \
    --table-name carts \
    --attribute-definitions AttributeName=user_id,AttributeType=S \
    --key-schema AttributeName=user_id,KeyType=HASH \
    --billing-mode PAY_PER_REQUEST \
    --region us-east-1 || echo "Table 'carts' already exists"

# SQS キュー作成
aws sqs create-queue \
    --endpoint-url http://localhost:4566 \
    --queue-name order-events \
    --region us-east-1 || echo "Queue 'order-events' already exists"

aws sqs create-queue \
    --endpoint-url http://localhost:4566 \
    --queue-name point-events \
    --region us-east-1 || echo "Queue 'point-events' already exists"

# SNS トピック作成
aws sns create-topic \
    --endpoint-url http://localhost:4566 \
    --name order-notifications \
    --region us-east-1 || echo "Topic 'order-notifications' already exists"

aws sns create-topic \
    --endpoint-url http://localhost:4566 \
    --name point-notifications \
    --region us-east-1 || echo "Topic 'point-notifications' already exists"

echo ""
echo "==================================================================="
echo "AWS resources initialized successfully!"
echo "==================================================================="
echo "Cognito User Pool ID: $USER_POOL_ID"
echo "Cognito Client ID: $CLIENT_ID"
echo "Test User: testuser@example.com / TestPass123!"
echo "==================================================================="
echo ""
echo "Save these values to your .env file:"
echo "COGNITO_USER_POOL_ID=$USER_POOL_ID"
echo "COGNITO_CLIENT_ID=$CLIENT_ID"
echo "COGNITO_REGION=us-east-1"
echo "==================================================================="

