# API Schema Code Generation Makefile
# 
# このディレクトリはGoモジュールではありません。generated/内のGoファイルで
# コンパイルエラーが表示されても問題ありません。
# 生成されたコードは `make sync-to-services` で各サービスにコピーされ、
# そこで適切なgo.modと依存関係と共にビルドされます。

.PHONY: install
install:
	@echo "Installing oapi-codegen..."
	@go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@v1.16.2
	@echo "✅ oapi-codegen installed"

.PHONY: generate
generate: generate-go generate-ts

.PHONY: generate-go
generate-go: install
	@echo "Generating Go server code from OpenAPI specs..."
	@echo "⚠️  Note: generated/*.go files will show compile errors here (no go.mod)"
	@echo "    This is expected! They will be compiled in each service directory."
	@mkdir -p generated/cart-service
	@mkdir -p generated/order-service
	@mkdir -p generated/point-service
	@mkdir -p generated/product-service
	@mkdir -p generated/promotion-service
	@$$(go env GOPATH)/bin/oapi-codegen -config oapi-codegen-configs/cart-service.yaml openapi/cart-service.yaml
	@$$(go env GOPATH)/bin/oapi-codegen -config oapi-codegen-configs/order-service.yaml openapi/order-service.yaml
	@$$(go env GOPATH)/bin/oapi-codegen -config oapi-codegen-configs/point-service.yaml openapi/point-service.yaml
	@$$(go env GOPATH)/bin/oapi-codegen -config oapi-codegen-configs/product-service.yaml openapi/product-service.yaml
	@$$(go env GOPATH)/bin/oapi-codegen -config oapi-codegen-configs/promotion-service.yaml openapi/promotion-service.yaml
	@echo "✅ Go code generated successfully"

.PHONY: generate-ts
generate-ts:
	@echo "Generating TypeScript types from OpenAPI specs..."
	@npm run generate:types
	@echo "✅ TypeScript types generated successfully"

.PHONY: sync-to-services
sync-to-services: generate-go
	@echo "Syncing generated code to service directories..."
	@cp generated/cart-service/openapi.go ../backend/services/cart-service/generated/
	@cp generated/order-service/openapi.go ../backend/services/order-service/generated/
	@cp generated/point-service/openapi.go ../backend/services/point-service/generated/
	@cp generated/product-service/openapi.go ../backend/services/product-service/generated/
	@cp generated/promotion-service/openapi.go ../backend/services/promotion-service/generated/
	@echo "✅ Generated code synced to all services"

.PHONY: clean
clean:
	@echo "Cleaning generated files..."
	@rm -rf generated/
	@echo "✅ Cleaned generated files"

.PHONY: validate
validate:
	@echo "Validating OpenAPI specifications..."
	@npm run validate
	@echo "✅ Validation complete"
