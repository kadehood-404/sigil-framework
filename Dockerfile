# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy dependency manifests first for efficient caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Compile a statically linked Go binary
RUN CGO_ENABLED=0 GOOS=linux go build -o sigil-proxy .

# Final minimal runtime stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /root/
COPY --from=builder /app/sigil-proxy .

EXPOSE 11435
CMD ["./sigil-proxy"]