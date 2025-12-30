# Point Service Dockerfile
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copy go.mod files
COPY services/point-service/go.mod services/point-service/go.sum ./services/point-service/
COPY shared/go.mod ./shared/

WORKDIR /app/services/point-service
RUN go mod download

# Copy shared module
WORKDIR /app
COPY shared ./shared

# Copy service code
COPY services/point-service ./services/point-service

WORKDIR /app/services/point-service
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o point-service .

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/services/point-service/point-service .

EXPOSE 8004

CMD ["./point-service"]
