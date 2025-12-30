# Promotion Service Dockerfile
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copy go.mod files
COPY services/promotion-service/go.mod services/promotion-service/go.sum ./services/promotion-service/
COPY shared/go.mod ./shared/

WORKDIR /app/services/promotion-service
RUN go mod download

# Copy shared module
WORKDIR /app
COPY shared ./shared

# Copy service code
COPY services/promotion-service ./services/promotion-service

WORKDIR /app/services/promotion-service
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o promotion-service .

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/services/promotion-service/promotion-service .

EXPOSE 8005

CMD ["./promotion-service"]
