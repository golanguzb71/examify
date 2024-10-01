# Stage 1: Build the binary
FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY . .

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o ielts-service ./cmd/main.go

# Stage 2: Create a minimal image
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/ielts-service .

COPY --from=builder /app/config/config.yaml ./config/

CMD ["./ielts-service"]
