# Cart Service Dockerfile
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum
COPY services/cart-service/go.mod services/cart-service/go.sum ./services/cart-service/
COPY shared/go.mod ./shared/

WORKDIR /app/services/cart-service
RUN go mod download

# Copy shared module
WORKDIR /app
COPY shared ./shared

# Copy service code
COPY services/cart-service ./services/cart-service

WORKDIR /app/services/cart-service
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o cart-service .

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/services/cart-service/cart-service .

EXPOSE 8002

CMD ["./cart-service"]
