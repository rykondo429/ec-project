# EC-Sample Copilot Instructions

## Architecture Overview

**Go-based microservices e-commerce platform** using Clean Architecture and DDD. Five independent services communicate via REST APIs and event-driven messaging (SQS/SNS).

### Service Map
Located in `backend/services/` - each service runs independently on its own port:
- **Product Service** (8001): Elasticsearch full-text search, Redis caching, MySQL product master
- **Cart Service** (8002): DynamoDB cart persistence, ElastiCache/Redis session management
- **Order Service** (8003): MySQL + GORM for transactional orders, SQS event publishing
- **Point Service** (8004): MySQL point ledger with transaction management
- **Promotion Service** (8005): MySQL coupon/discount calculation engine

### Technology Stack
- **Backend**: Go 1.21+, Echo v4, GORM v2, AWS SDK v2
- **Databases**: MySQL 8.0, DynamoDB (LocalStack), Elasticsearch 8.10.2, Redis 7
- **Auth**: JWT Bearer tokens ([backend/shared/auth/jwt.go](backend/shared/auth/jwt.go)) - fallback to "anonymous" if missing
- **AWS Emulation**: LocalStack for DynamoDB, SQS, SNS, Cognito, S3 (initialized via [init-aws.sh](init-aws.sh))
- **Frontend**: Next.js 14 App Router, TypeScript 5.3, TailwindCSS 3.3, Zustand 4.4
- **Logging**: Uber Zap structured JSON logging - NEVER use `fmt.Print`
- **API Schema**: OpenAPI 3.0 specs → oapi-codegen → generated Go server code

## Code Organization Pattern (MANDATORY)

**Every service MUST follow this identical 4-layer structure** in `backend/services/{service-name}/`:

```
backend/services/{service-name}/
├── main.go                      # Echo setup, DI, middleware chain, route registration
├── domain/
│   ├── {model}.go              # Pure Go structs + JSON tags (Order, Cart, etc.)
│   └── repository.go           # Repository interface definitions ONLY
├── infrastructure/
│   └── {db}/{repo}_repository.go  # Concrete repository implementations
├── interface/handler/
│   └── {model}_handler.go      # Echo HTTP handlers (extract context, return JSON)
├── usecase/
│   └── {model}_usecase.go      # Business logic orchestration
├── generated/
│   └── openapi.go              # Auto-generated from openapi.yaml (DO NOT EDIT)
├── openapi.yaml                # OpenAPI 3.0 spec (source of truth)
├── oapi-codegen.yaml           # Code generation config
├── Makefile                    # `make generate` → regenerate code
└── Dockerfile
```

### Dependency Flow (Critical!)
```
Interface → UseCase → Domain ← Infrastructure
```
- **Domain** = Zero external dependencies (pure business entities)
- **Infrastructure** implements Domain interfaces (`type CartRepository interface {...}`)
- **UseCase** depends ONLY on Domain interfaces, never concrete implementations
- **Interface** calls UseCase methods, handles HTTP concerns only

### Layer Responsibilities

**Domain Layer** ([example](backend/services/cart-service/domain/repository.go)):
```go
// Pure interface definitions - NO implementations
type CartRepository interface {
    GetCart(ctx context.Context, userID string) (*Cart, error)
    SaveCart(ctx context.Context, cart *Cart) error
}
```

**Infrastructure Layer** ([example](backend/services/cart-service/infrastructure/dynamodb/cart_repository.go)):
```go
// Concrete implementation using DynamoDB SDK
func (r *DynamoDBRepository) GetCart(ctx context.Context, userID string) (*domain.Cart, error) {
    // AWS SDK calls here
}
```

**UseCase Layer** (business logic only):
```go
type CartUseCase struct {
    repo domain.CartRepository  // Interface, not concrete type
}

func (uc *CartUseCase) GetCart(ctx context.Context, userID string) (*domain.Cart, error) {
    // Validation, business rules, then delegate to repo
    return uc.repo.GetCart(ctx, userID)
}
```

**Handler Layer** ([example](backend/services/cart-service/interface/handler/cart_handler.go)):
```go
func (h *CartHandler) GetCart(c echo.Context) error {
    userID := c.Get("user_id").(string)  // Extract from middleware
    cart, err := h.useCase.GetCart(c.Request().Context(), userID)
    if err != nil {
        return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
    }
    return c.JSON(http.StatusOK, cart)
}
```

## Critical Workflows

### Local Development Setup (Complete Sequence)
```bash
# 1. Start ALL infrastructure (blocks until healthy)
docker-compose up -d
# Wait for: MySQL, Redis, Elasticsearch, LocalStack healthchecks

# 2. Verify services are healthy
docker-compose ps
# All services should show "healthy" status

# 3. Initialize LocalStack AWS resources (auto-run via docker-compose)
# Creates: DynamoDB tables, SQS queues, SNS topics
# See init-aws.sh for details

# 4. Start individual Go services (each in separate terminal)
cd backend/services/product-service && go run main.go  # :8001
cd backend/services/cart-service && go run main.go     # :8002
cd backend/services/order-service && go run main.go    # :8003
cd backend/services/point-service && go run main.go    # :8004
cd backend/services/promotion-service && go run main.go # :8005

# 5. Start frontend (separate terminal)
cd frontend && npm install && npm run dev  # :3000
```

### Build & Test
```bash
# Build all services (uses build.sh orchestrator)
bash build.sh
# Runs: go mod tidy → go build -o bin/main for each service
# Creates bin/main executables in each service directory

# Test individual service (from service directory)
cd backend/services/order-service
go test ./...  # Tests all packages

# Regenerate OpenAPI code (when openapi.yaml changes)
cd backend/services/cart-service
make generate  # Runs oapi-codegen → generates generated/openapi.go

# Docker build (individual service)
docker build -t ec-order-service ./backend/services/order-service

# Docker compose full rebuild
docker-compose build
```

### Testing API Endpoints (Quick Reference)
```bash
# Product search
curl "http://localhost:8001/api/v1/products/search?keyword=shoes&page=1&page_size=20"

# Get cart (requires auth header)
curl -H "Authorization: Bearer YOUR_TOKEN" http://localhost:8002/api/v1/carts

# Add to cart
curl -X POST -H "Content-Type: application/json" -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{"product_id":"prod-123","product_name":"Nike Air","price":12000,"quantity":2}' \
  http://localhost:8002/api/v1/carts/items

# See API_EXAMPLES.md for complete examples
```

## Key Integration Patterns

### Database Connection Patterns (Critical - Each Service Different!)

**MySQL (GORM) - Order/Point/Promotion Services**:
```go
// main.go pattern (see backend/services/order-service/main.go)
dsn := os.Getenv("MYSQL_DSN")  // user:pass@tcp(host:3306)/dbname?parseTime=true
db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
sqlDB, _ := db.DB()
sqlDB.SetMaxOpenConns(10)
sqlDB.SetMaxIdleConns(5)
sqlDB.SetConnMaxLifetime(time.Hour)
// Auto-migrate models on startup:
db.AutoMigrate(&domain.Order{}, &domain.OrderItem{})
```

**DynamoDB (AWS SDK v2) - Cart Service**:
```go
// main.go pattern (see backend/services/cart-service/main.go)
cfg, _ := config.LoadDefaultConfig(context.Background())
if endpoint := os.Getenv("AWS_ENDPOINT_URL"); endpoint != "" {
    cfg.BaseEndpoint = &endpoint  // LocalStack support
}
dynamoClient := awsdynamodb.NewFromConfig(cfg)
// Table must exist in LocalStack (see init-aws.sh)
```

**Elasticsearch - Product Service**:
```go
// main.go pattern (see backend/services/product-service/main.go)
esURL := os.Getenv("ELASTICSEARCH_URL")  // http://elasticsearch:9200
cfg := opensearchgo.Config{Addresses: []string{esURL}}
client, _ := opensearchgo.NewClient(cfg)
```

**Redis - All Services (caching)**:
```go
rdb := redis.NewClient(&redis.Options{
    Addr: os.Getenv("REDIS_ADDR"),  // redis:6379
})
```

### Authentication & Authorization Flow

**Middleware Setup** ([backend/shared/auth/jwt.go](backend/shared/auth/jwt.go)):
```go
// In main.go of each service:
e.Use(auth.JWTMiddleware(os.Getenv("JWT_SECRET_KEY")))
// Header: "Authorization: Bearer <token>"
// Sets: c.Set("user_id", claims.UserID), c.Set("email", claims.Email)
// Fallback: Sets "anonymous" if no token present
```

**Handler Usage**:
```go
func (h *Handler) GetCart(c echo.Context) error {
    userID := c.Get("user_id").(string)  // ALWAYS available after middleware
    // Use for user-scoped operations
}
```

### Frontend API Integration Pattern

**API Client Setup** ([frontend/src/lib/api-client.ts](frontend/src/lib/api-client.ts)):
```typescript
// Separate client per service (different base URLs)
const cartApi = new CartApi('http://localhost:8002');
const productApi = new ProductApi('http://localhost:8001');

// Auto-attach JWT from localStorage
this.client.interceptors.request.use((config) => {
  const token = localStorage.getItem('auth_token');
  if (token) config.headers.Authorization = `Bearer ${token}`;
  return config;
});
```

**Zustand State Management** ([frontend/src/store/cart.ts](frontend/src/store/cart.ts)):
```typescript
// Per-domain stores, not global
import { create } from 'zustand';
const useCartStore = create((set) => ({
  cart: null,
  addItem: async (item) => { /* call CartApi */ },
}));
```

## Environment Variables

**Shared (all services)**:
```
PORT=8001                       # Service port
ENVIRONMENT=development
JWT_SECRET_KEY=your-secret-key
MYSQL_DSN=ecuser:ecpassword@tcp(mysql:3306)/ecsite
REDIS_ADDR=redis:6379
ELASTICSEARCH_URL=http://elasticsearch:9200
AWS_REGION=us-east-1
AWS_ENDPOINT_URL=http://localstack:4566  # For local development
AWS_ACCESS_KEY_ID=test
AWS_SECRET_ACCESS_KEY=test
```

See `.env.example` and `docker-compose.yml` for complete list.

## Project Conventions

1. **No cross-service database access**: Each service owns its data; use REST APIs for inter-service communication
2. **Immediate context cancellation**: Always use `context.WithTimeout()` for external calls
3. **Package organization**: Separate by architectural layer (domain, infrastructure, interface, usecase)
4. **Logging**: Use zap logger, not fmt.Print
5. **Repository pattern mandatory**: All data access via repository interfaces
6. **HTTP methods**: GET (query), POST (create), PUT (full update), DELETE (remove)

## Common Tasks

### Adding a New Endpoint to Existing Service
1. Add handler method to `interface/handler/{model}.go`
2. Register route in `main.go`: `e.GET("/path", handler.Method)`
3. Add business logic to `usecase/{model}.go` if needed
4. Test with: `curl -H "Authorization: Bearer test" http://localhost:PORT/path`

### Adding a New Service
1. Create `backend/services/new-service/` with standard structure
2. Copy pattern from similar service (cart-service for DynamoDB, order-service for MySQL)
3. Update `docker-compose.yml` with new service definition
4. Update `build.sh` to include new service
5. Add OpenAPI spec to `api-schema/openapi/new-service.yaml`

### Debugging Service Issues
```bash
# Check service logs
docker logs ec-cart-service
docker logs ec-order-service

# Test database connectivity
docker exec ec-mysql mysql -uroot -p${DB_ROOT_PASSWORD} -e "SELECT 1"
docker exec ec-redis redis-cli ping

# View LocalStack logs
docker logs ec-localstack
```

## OpenAPI Code Generation Workflow

**Critical**: Services use oapi-codegen to auto-generate Echo server boilerplate from OpenAPI specs.

**Files per service**:
- `openapi.yaml` - Source of truth (edit this)
- `oapi-codegen.yaml` - Generation config (package: generated, echo-server: true)
- `generated/openapi.go` - Auto-generated (NEVER edit manually)
- `Makefile` - Run `make generate` after openapi.yaml changes

**Pattern** ([backend/services/order-service/Makefile](backend/services/order-service/Makefile)):
```makefile
generate:
	go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@v1.16.2
	oapi-codegen -config oapi-codegen.yaml openapi.yaml > generated/openapi.go
```

**When to regenerate**:
1. After changing openapi.yaml (routes, models, parameters)
2. After pulling changes that modify openapi.yaml
3. Before building for production

**Frontend Types** ([api-schema/](api-schema/)):
```bash
cd api-schema
npm run generate  # Creates TypeScript types from all service OpenAPI specs
# Generates: types/{service-name}.ts from openapi/{service-name}.yaml
```

## Important Files Reference

- [build.sh](build.sh) - Service build orchestration
- [docker-compose.yml](docker-compose.yml) - Infrastructure definition
- [backend/shared/auth/jwt.go](../../backend/shared/auth/jwt.go) - JWT middleware
- [backend/shared/config/config.go](backend/shared/config/config.go) - Config loader
- [backend/migrations/mysql/01_init_schema.sql](backend/migrations/mysql/01_init_schema.sql) - DB schema
