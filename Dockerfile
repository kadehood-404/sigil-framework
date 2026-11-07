# --- Build Stage ---
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy module files and dependencies
COPY go.mod ./
RUN go mod download || true

# Copy source code files
COPY main.go ./

# Explicitly build main.go to avoid collision with server.go
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o sigil-proxy main.go

# --- Runtime Stage ---
FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /root/
COPY --from=builder /app/sigil-proxy .

EXPOSE 11435
CMD ["./sigil-proxy"]