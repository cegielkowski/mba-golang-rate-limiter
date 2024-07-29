# Etapa de construção
FROM golang:1.18-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o rate-limiter ./cmd/main.go

# Etapa de execução
FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/rate-limiter .
ENTRYPOINT ["./rate-limiter"]