# Build stage
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o bot .

# Runtime stage
FROM alpine:3.21
RUN apk add --no-cache ca-certificates tzdata
# Set timezone to Asia/Taipei for local use
ENV TZ=Asia/Taipei
WORKDIR /app
COPY --from=builder /app/bot .
# Using named volumes for persistent data and project configs
VOLUME ["/data", "/projects"]
# Add a non-root user for better security
RUN adduser -D -u 1001 botuser && chown -R botuser /app /data /projects 2>/dev/null || true
USER botuser
CMD ["./bot"]
