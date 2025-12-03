FROM golang:1.24-alpine AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o gateway ./cmd/gateway

FROM alpine:3.19
RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=builder /app/gateway .
COPY --from=builder /app/configs ./configs

USER nobody

CMD ["./gateway"]
