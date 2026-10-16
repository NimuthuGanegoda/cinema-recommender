# Go major tag tracks supported security releases; pin a tested digest in your deployment.
FROM golang:1-alpine AS builder
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /cinebite ./cmd/api

FROM alpine:3
RUN apk --no-cache add ca-certificates tzdata && addgroup -S cinebite && adduser -S -G cinebite -u 10001 cinebite
WORKDIR /app
COPY --from=builder /cinebite /app/cinebite
USER 10001:10001
ENV PORT=:8080 TZ=Asia/Colombo
EXPOSE 8080
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s CMD wget -qO- http://127.0.0.1:8080/health || exit 1
ENTRYPOINT ["/app/cinebite", "-server"]
