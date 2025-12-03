FROM golang:1.25-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o user_service ./cmd/users_service

FROM alpine:3.19 AS runner
RUN apk add --no-cache ca-certificates && adduser -D -H appuser
WORKDIR /app
COPY --from=builder /app/user_service ./user_service
COPY --from=builder /app/.env ./
USER appuser
EXPOSE 8001
ENTRYPOINT ["./user_service"]
