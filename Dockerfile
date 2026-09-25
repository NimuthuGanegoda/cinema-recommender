# ==========================================
# Multi-Stage Apple-Caliber Docker Build
# ==========================================

# Stage 1: Build binary
FROM golang:1.22-alpine AS builder

WORKDIR /src

# Cache dependencies
COPY go.mod ./
RUN go mod download

# Copy source tree
COPY . .

# Compile optimized static binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-s -w -X main.version=2.0.0" \
    -o /app/cinema-recommender ./cmd/api

# Stage 2: Minimal runtime image
FROM alpine:3.20

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/cinema-recommender /app/cinema-recommender

EXPOSE 8080

ENTRYPOINT ["/app/cinema-recommender", "-server", "-port", ":8080"]
