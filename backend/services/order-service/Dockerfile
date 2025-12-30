# Order Service Dockerfile
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copy go.mod files
COPY services/order-service/go.mod services/order-service/go.sum ./services/order-service/
COPY shared/go.mod ./shared/

WORKDIR /app/services/order-service
RUN go mod download

# Copy shared module
WORKDIR /app
COPY shared ./shared

# Copy service code
COPY services/order-service ./services/order-service

WORKDIR /app/services/order-service
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o order-service .

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/services/order-service/order-service .

EXPOSE 8003

CMD ["./order-service"]
