FROM golang:1.24-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o product_service ./cmd/Product_Service

FROM alpine:3.19 AS runner
RUN apk add --no-cache ca-certificates && adduser -D -H appuser
WORKDIR /app
COPY --from=builder /app/product_service ./product_service
COPY --from=builder /app/.env ./
USER appuser
EXPOSE 8082
ENTRYPOINT ["./product_service"]
