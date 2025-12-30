# Product Service Dockerfile
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copy go.mod files
COPY services/product-service/go.mod services/product-service/go.sum ./services/product-service/
COPY shared/go.mod ./shared/

WORKDIR /app/services/product-service
RUN go mod download

# Copy shared module
WORKDIR /app
COPY shared ./shared

# Copy service code
COPY services/product-service ./services/product-service

WORKDIR /app/services/product-service
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o product-service .

FROM alpine:latest

RUN apk --no-cache add ca-certificates curl

WORKDIR /root/

COPY --from=builder /app/services/product-service/product-service .

EXPOSE 8001

CMD ["./product-service"]
