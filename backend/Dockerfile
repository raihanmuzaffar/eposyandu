# Build Stage
FROM golang:1.24-alpine AS builder
WORKDIR /app

# Set toolchain agar otomatis mengunduh jika butuh versi spesifik
ENV GOTOOLCHAIN=auto

# Mengunduh dependency terlebih dahulu
COPY go.mod go.sum ./
RUN go mod download

# Copy seluruh source code
COPY . .

# Build binary Go secara terkompilasi statis
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o main .

# Final Run Stage
FROM alpine:latest
WORKDIR /root/

RUN apk --no-cache add ca-certificates tzdata
COPY --from=builder /app/main .

EXPOSE 3000

CMD ["./main"]