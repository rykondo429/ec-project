#!/bin/bash

# EC Sample Auto Build Script
# 全てのマイクロサービスとフロントエンドを自動ビルドする

set -e

# プロジェクトルートディレクトリを保存
PROJECT_ROOT=$(pwd)

SERVICES=(
    "backend/services/product-service"
    "backend/services/order-service"
    "backend/services/cart-service"
    "backend/services/point-service"
    "backend/services/promotion-service"
)

echo "=========================================="
echo "EC Sample Auto Build Starting"
echo "=========================================="

for SERVICE in "${SERVICES[@]}"; do
    echo ""
    echo "Building: $SERVICE"
    cd "$PROJECT_ROOT/$SERVICE"
    
    # go mod tidy
    echo "  - Running go mod tidy..."
    go mod tidy
    
    # go build
    echo "  - Building..."
    go build -o bin/main .
    
    if [ $? -eq 0 ]; then
        echo "  ✓ Build successful"
    else
        echo "  ✗ Build failed"
        exit 1
    fi
done

# Frontend build (use Docker)
echo ""
echo "Building: frontend (Next.js with Docker)"
cd "$PROJECT_ROOT"

# Check if Docker is available
if command -v docker &> /dev/null; then
    echo "  - Building Docker image for frontend..."
    docker compose build web
    
    if [ $? -eq 0 ]; then
        echo "  ✓ Frontend Docker build successful"
    else
        echo "  ✗ Frontend Docker build failed"
        echo "  Note: Node.js 18+ required. Use Docker instead: docker compose build web"
        # Frontend build failure is warning only (backend still works)
    fi
else
    echo "  ⚠️  Docker not found - skipping frontend build"
    echo "  Install Docker or use: cd frontend && npm run build (requires Node.js 18+)"
fi

echo ""
echo "=========================================="
echo "Backend services built successfully!"
echo "=========================================="

# Docker Compose full build (optional)
if [ "$1" == "--all" ]; then
    echo ""
    echo "Building all Docker images..."
    docker compose build
    echo "  ✓ All Docker images built"
fi
